package memory

import (
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	m "github.com/pbnjay/memory"
)

func getMemInfoEntryFromFile(filename string, memInfoEntry string) uint64 {
	var memoryInfo uint64

	reader, err := os.Open(filename)

	if err == nil {
		memoryInfo = getMemInfoEntryFromReader(reader, memInfoEntry)
	}

	return memoryInfo
}

func getMemInfoEntryFromReader(reader io.Reader, memInfoEntry string) uint64 {
	var memoryInfo uint64

	memData, err := ioutil.ReadAll(reader)

	if err == nil {
		memDataSplit := strings.Split(string(memData), "\n")
		exp := regexp.MustCompile("^" + memInfoEntry + ":.*\\s(\\d*).*\\skB$")

		for _, entry := range memDataSplit {
			if strings.HasPrefix(entry, memInfoEntry+":") {
				err = nil
				m := exp.FindStringSubmatch(entry)

				if m != nil && len(m) == 2 {
					if memoryInfo, err = strconv.ParseUint(m[1], 10, 64); err != nil {
						memoryInfo = 0
					}
				}

				break
			}
		}
	}

	return memoryInfo * 1024
}

func getFreeMemoryWithHandler(freeMemory func() uint64) uint64 {
	return freeMemory()
}

// GetFreeMemory returns the amount of available memory.
func GetFreeMemory() uint64 {
	return getFreeMemoryWithHandler(getFreeMemoryOsConstrained)
}

func getTotalMemoryWithHandler(totalMemory func() uint64) uint64 {
	return totalMemory()
}

// GetTotalMemory returns the total amount of memory.
func GetTotalMemory() uint64 {
	return getTotalMemoryWithHandler(m.TotalMemory)
}

func getUsedMemoryWithHandlers(totalMemory func() uint64, freeMemory func() uint64) uint64 {
	var used uint64

	total := getTotalMemoryWithHandler(totalMemory)
	free := getFreeMemoryWithHandler(freeMemory)

	if total != 0 && free != 0 {
		used = total - free
	}

	return used
}

// GetUsedMemory returns the amount of used memory.
func GetUsedMemory() uint64 {
	return getUsedMemoryWithHandlers(m.TotalMemory, getFreeMemoryOsConstrained)
}

func getUsedMemoryPercentageWithHandlers(totalMemory func() uint64, freeMemory func() uint64) uint64 {
	var percentageUsed uint64

	used := getUsedMemoryWithHandlers(totalMemory, freeMemory)
	total := getTotalMemoryWithHandler(totalMemory)

	if used != 0 && total != 0 {
		percentageUsed = uint64(float64(used) / float64(total) * 100)
	}

	return percentageUsed
}

// GetUsedMemoryPercentage returns the amount of used memory
// as a percentage.
func GetUsedMemoryPercentage() uint64 {
	return getUsedMemoryPercentageWithHandlers(m.TotalMemory, getFreeMemoryOsConstrained)
}

// GetProcessMemoryPercentage returns the percentage of used memory of a process against total system memory.
func GetProcessMemoryPercentage(processName string) (float64, error) {
	return getProcessMemoryPercentageOsContrained(processName)
}
