// +build !windows

package memory

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/ncr-devops-platform/nagiosfoundation/lib/pkg/process"
)

// getFreeMemoryOsConstrained returns the amount of available memory.
func getFreeMemoryOsConstrained() uint64 {
	return getMemInfoEntryFromFile("/proc/meminfo", "MemAvailable")
}

var getProcessesByName = process.GetProcessesByName
var execBash = func(command string) ([]byte, error) {
	return exec.Command("bash", "-c", command).CombinedOutput()
}

func handlePSOutput(output string, err error) (float64, error) {
	if err != nil {
		return 0.0, fmt.Errorf("%v. Full shell output: %v", err, output)
	}
	split := strings.Split(output, "\n")
	if len(split) < 2 {
		return 0.0, fmt.Errorf("Invalid output from 'ps' command. Expected 1 header and 1 value line")
	}
	// second line should contain memory consumption
	value, err := strconv.ParseFloat(strings.TrimSpace(split[1]), 64)
	if err != nil {
		return 0.0, err
	}

	return value, nil
}

func getProcessMemoryPercentageOsContrained(processName string) (float64, error) {
	processInfo, err := getProcessesByName(processName)
	if err != nil {
		return 0.0, err
	}

	var memoryUsedTotal float64 = 0
	for _, entry := range processInfo {
		command := fmt.Sprintf("ps -p %v -o %%mem", entry.PID)
		out, err := execBash(command)
		memoryUsed, err := handlePSOutput(string(out), err)
		if err != nil {
			return 0.0, err
		}

		memoryUsedTotal += memoryUsed
	}

	return memoryUsedTotal, nil
}
