// +build !windows

package cpu

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/ncr-devops-platform/nagiosfoundation/lib/pkg/process"
)

func getCPULoadOsConstrained() (float64, error) {
	return getCPULoadLinux()
}

var getProcessesByName = process.GetProcessesByName
var getCPUCount = runtime.NumCPU

var execBash = func(command string) ([]byte, error) {
	return exec.Command("bash", "-c", command).CombinedOutput()
}

func parseTopSamples(lines []string) (float64, error) {
	var result float64 = 0
	// Fields: PID, USER, PR, NI, VIRT, RES, SHR, S, %CPU, %MEM, TIME+, COMMAND
	for _, line := range lines {
		lineSplit := strings.Fields(line)
		if len(lineSplit) != 12 {
			return 0.0, fmt.Errorf("Unexpected Top stat line. Expected 12 elements, actual %v", len(lineSplit))
		}

		cpuValue, err := strconv.ParseFloat(lineSplit[8], 64)
		if err != nil {
			return 0.0, err
		}

		result += cpuValue
	}

	return result, nil
}
func handleTopOutput(output string, err error) (float64, error) {
	if err != nil {
		return 0.0, fmt.Errorf("%v. Full shell output: %v", err, output)
	}

	split := strings.Split(output, "\n")
	totalLines := len(split)

	cpuSamples := []float64{}
	ignoredFirstSample := false
	for i := 0; i < totalLines; i++ {
		line := split[i]

		// read till next PID command
		if !strings.Contains(line, "PID") {
			continue
		}

		// ignore first sample since on older machines it will report PID lifetime stats
		if !ignoredFirstSample {
			ignoredFirstSample = true
			continue
		}

		sampleLines := []string{}
		for idx := i + 1; idx < totalLines; idx++ {
			line := split[idx]
			if len(line) <= 0 {
				// reached end of >sample< output
				i = idx
				break
			}
			sampleLines = append(sampleLines, split[idx])
		}

		// either reached end of sample output or end of output
		cpuSum, err := parseTopSamples(sampleLines)
		if err != nil {
			return 0.0, err
		}
		cpuSamples = append(cpuSamples, cpuSum)
	}

	result := average(cpuSamples) / float64(getCPUCount())
	return result, nil
}

func getProcessCPULoad(processInfo []process.GeneralInfo) (float64, error) {
	if len(processInfo) < 1 {
		return 0.0, nil
	}

	pidsString := ""
	for _, entry := range processInfo {
		pidsString = fmt.Sprintf("%v,%v", pidsString, entry.PID)
	}
	pidsString = pidsString[1:]

	command := fmt.Sprintf("top -b -n 4 -d 1 -p %v", pidsString)
	out, err := execBash(command)
	max, err := handleTopOutput(string(out), err)
	return max, err
}

// returns highest single core CPU load of a given process based on multiple samples from pidstat
// tasks of a process will be considered as well
func getProcessCPULoadOsConstrained(processName string, perCoreCalculation bool) (float64, error) {
	processInfo, err := getProcessesByName(processName)
	if err != nil {
		return 0, err
	}

	if perCoreCalculation {
		return getProcessCoreCPULoad(processInfo)
	}

	return getProcessCPULoad(processInfo)
}
