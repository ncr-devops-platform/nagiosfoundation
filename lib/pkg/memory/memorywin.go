// +build windows

package memory

import (
	"github.com/StackExchange/wmi"
	"github.com/ncr-devops-platform/nagiosfoundation/lib/pkg/perfcounters"
	memlib "github.com/pbnjay/memory"
)

func getFreeMemoryOsConstrained() uint64 {
	counter, err := perfcounters.ReadPerformanceCounter("\\Memory\\Available Bytes", 2, 1)

	var memoryAvailable uint64

	if err != nil {
		memoryAvailable = uint64(0)
	} else {
		memoryAvailable = uint64(counter.Value)
	}

	return uint64(memoryAvailable)
}

// See WMI class: Win32_Process
// Must match the original class name (case insensitive)
type win32_process struct {
	Name      string
	ProcessID uint32

	// ParentProcessId may refer to a process reusing process identifier (in case parent has already terminated)
	ParentProcessID uint32

	// WorkingSetSize field seems to be the most accurate representation of currently used RAM
	WorkingSetSize uint64
}

type processInfo struct {
	Process           win32_process
	UsedInCalculation bool
}

var getTotalMemory = memlib.TotalMemory
var getWin32Processes = wmi.Query

func getMemoryUsedByPIDAndItsChildren(data map[uint32]*processInfo, pid uint32) uint64 {
	rootProcess, isRootFound := data[pid]
	if !isRootFound || rootProcess.UsedInCalculation {
		return 0
	}

	result := rootProcess.Process.WorkingSetSize
	rootProcess.UsedInCalculation = true

	for processPID, process := range data {
		if process.Process.ParentProcessID == pid {
			result += getMemoryUsedByPIDAndItsChildren(data, processPID)
		}
	}

	return result
}

func getMemoryUsedByProcessNameAndItsChildren(data map[uint32]*processInfo, name string) uint64 {
	var result uint64 = 0
	for pid, process := range data {
		if process.Process.Name == name {
			result += getMemoryUsedByPIDAndItsChildren(data, pid)
		}
	}

	return result
}

func getProcessMemoryPercentageOsContrained(processName string) (float64, error) {
	memoryTotal := getTotalMemory()

	// query all processes
	var dst []win32_process
	q := wmi.CreateQuery(&dst, "")
	err := getWin32Processes(q, &dst)
	if err != nil {
		return 0.0, err
	}

	pidDataMap := map[uint32]*processInfo{}
	for _, entry := range dst {
		pidDataMap[entry.ProcessID] = &processInfo{
			Process:           entry,
			UsedInCalculation: false,
		}
	}

	// calculate consumed memory, including children of matched processes
	nameToMatch := processName + ".exe"
	memoryUsedTotal := getMemoryUsedByProcessNameAndItsChildren(pidDataMap, nameToMatch)
	percentage := float64(memoryUsedTotal) / float64(memoryTotal) * 100
	return percentage, nil
}
