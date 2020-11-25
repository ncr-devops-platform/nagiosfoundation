// +build windows

package cpu

import (
	"fmt"
	"regexp"
	"runtime"
	"time"

	"github.com/StackExchange/wmi"
)

type Win32_PerfFormattedData_PerfOS_Processor struct {
	PercentProcessorTime uint64
}

func getCPULoadOsConstrained() (float64, error) {
	var dst []Win32_PerfFormattedData_PerfOS_Processor
	q := wmi.CreateQuery(&dst, "WHERE name='_Total'")
	err := wmi.Query(q, &dst)
	if err != nil {
		return 0.0, err
	}
	return (float64)(dst[0].PercentProcessorTime), nil
}

// See WMI class: Win32_PerfFormattedData_PerfProc_Process
// Must match the original class name (case insensitive)
type win32_PerfFormattedData_PerfProc_Process struct {
	Name                 string
	PercentProcessorTime uint64
}

var getWin32Processes = func(query string, dst interface{}) error {
	return wmi.Query(query, dst)
}

var getCPUCount = runtime.NumCPU

func getProcessCPULoadOsConstrained(processName string, _ bool) (float64, error) {
	cpuValues := []float64{}
	for i := 0; i < 3; i++ {

		var dst []win32_PerfFormattedData_PerfProc_Process
		q := wmi.CreateQuery(&dst, "")
		err := getWin32Processes(q, &dst)
		if err != nil {
			return 0.0, err
		}

		var percentage uint64
		for _, entry := range dst {
			pattern := fmt.Sprintf("^%v#[0-9]+$", processName)
			match, err := regexp.MatchString(pattern, entry.Name)
			if err != nil {
				return 0.0, err
			}

			if entry.Name == processName || match {
				// PercentProcessorTime is reported as a value in range [0, CORES * 100]
				percentage += entry.PercentProcessorTime
			}
		}
		cpuValues = append(cpuValues, float64(percentage))
		time.Sleep(1 * time.Second)
	}

	result := average(cpuValues) / float64(getCPUCount())
	return result, nil
}
