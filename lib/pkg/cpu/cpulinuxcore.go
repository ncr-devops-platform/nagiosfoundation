//go:build !windows
// +build !windows

package cpu

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/ncr-devops-platform/nagiosfoundation/lib/pkg/process"
)

func getNonWhiteSpaceLines(input []string) []string {
	output := []string{}
	for _, line := range input {
		if len(strings.TrimSpace(line)) <= 0 {
			continue
		}

		output = append(output, line)
	}

	return output
}

type sampleStartEnd struct {
	Start int
	End   int
}

// getSampleRanges returns 'pidstat' sample blocks' indices based on given lines
// Given lines MUST NOT contain empty lines.
// Range returned is inclusive for start index and exclusive for end index - for slicing support of given lines.
// Example of data returned:
// 0: 3 -> 5
// 1: 8 -> 10
// Above example means that there were 2 samples found in 'nonEmptyLines' array. Both blocks contain 2 lines about CPU usage.
// First one starts at line 3 and ends at line 5 (line 3 and 4 contain pid stats).
// Second starts at line 8 and ends at line 10 (line 8 and 9 contain pid stats).
func getSampleRanges(nonEmptyLines []string) map[int]*sampleStartEnd {
	sampleRanges := map[int]*sampleStartEnd{}
	currentSample := -1
	for i := 0; i < len(nonEmptyLines); i++ {
		line := nonEmptyLines[i]
		if strings.HasPrefix(line, "#") {
			if currentSample == -1 {
				currentSample++
				sampleRanges[currentSample] = &sampleStartEnd{
					Start: i + 1,
				}
				continue
			}

			sampleRanges[currentSample].End = i
			currentSample++
			sampleRanges[currentSample] = &sampleStartEnd{
				Start: i + 1,
			}
		}
	}

	if sampleRanges[currentSample].End == 0 {
		sampleRanges[currentSample].End = len(nonEmptyLines)
	}

	// safety check
	deleteValues := []int{}
	for key, sample := range sampleRanges {
		// special case for unexpected sample end, e.g. double "# Time UID TGID ..." line or no samples at all
		if sample.End <= sample.Start {
			deleteValues = append(deleteValues, key)
		}
	}

	// remove invalid entries
	for _, key := range deleteValues {
		delete(sampleRanges, key)
	}

	return sampleRanges
}

type pidStatLine struct {
	TID  int
	CPU  float64
	Core int
}

// parsePIDStatLine parses single PID stat line. See expected format below
func parsePIDStatLine(line string) (*pidStatLine, error) {
	// When 10 Fields: Time, UID, TGID, TID, %usr, %system, %guest, %CPU, CPU, Command
	// When 11 Fields: Time, UID, TGID, TID, %usr, %system, %guest, %wait, %CPU, CPU, Command
	fmt.Println(line)
	lineSplit := strings.Fields(line)
	if len(lineSplit) != 10 && len(lineSplit) != 11 {
		return nil, fmt.Errorf("Unexpected PID stat line. Expected at least 10 or 11 elements, actual %v", len(lineSplit))
	}

	cpuValuesColumnIndex := 7
	cpuCoreValuesColumnIndex := 8

	if len(lineSplit) == 11 {
		cpuValuesColumnIndex = 8
		cpuCoreValuesColumnIndex = 9
	}

	tidStr := lineSplit[3]
	var tid int
	var err error
	if tidStr == "-" {
		tid = 0. // Assign a default value of 0 for TID when '-' is encountered
	} else {
		tid, err = strconv.Atoi(tidStr)
		if err != nil {
			return nil, err
		}
	}

	cpuValue, err := strconv.ParseFloat(lineSplit[cpuValuesColumnIndex], 64)
	if err != nil {
		return nil, err
	}
	cpuCore, err := strconv.Atoi(lineSplit[cpuCoreValuesColumnIndex])
	if err != nil {
		return nil, err
	}

	return &pidStatLine{
		TID:  tid,
		CPU:  cpuValue,
		Core: cpuCore,
	}, nil
}

// parsePIDStatSample attempts to parse given lines to PID Stat Lines.
// Each line is expected to adhere to format defined in parsePIDStatLine
func parsePIDStatSample(sampleLines []string) ([]*pidStatLine, error) {
	result := []*pidStatLine{}
	for _, line := range sampleLines {
		pidStatLine, err := parsePIDStatLine(line)
		if err != nil {
			return nil, err
		}

		result = append(result, pidStatLine)
	}

	return result, nil
}

// calculateCPUCoreUsage calculates per-core utilization from a given PID Stat sample.
// TID == 0 will be ignored (this is primary process that is reported alongside tasks of a primary process;
// TGID of TID 0 will be equal to requested PID)
// Utilization of primary process will be reported with a task of a process (TID will be matching PID).
func calculateCPUCoreUsage(pidStatLines []*pidStatLine) map[int]float64 {
	result := map[int]float64{}
	for _, entry := range pidStatLines {
		// exclude 0 utilization to exclude unused cores
		if entry.TID == 0 || entry.CPU <= 0.0 {
			continue
		}
		if _, ok := result[entry.Core]; !ok {
			result[entry.Core] = 0.0
		}
		result[entry.Core] += entry.CPU
	}

	return result
}

func calculateCPUCoreAverage(input map[int]map[int]float64) map[int]float64 {
	result := map[int]float64{}
	samplesUsed := map[int]int{}
	for _, cpuCoreUsage := range input {
		for core, utilization := range cpuCoreUsage {
			if _, ok := result[core]; !ok {
				result[core] = 0.0
				samplesUsed[core] = 0
			}

			result[core] += utilization
			samplesUsed[core]++
		}
	}

	for core, value := range result {
		result[core] = value / float64(samplesUsed[core])
	}

	return result
}

func getCPUCoreUsage(output string, shellErr error) (float64, error) {
	if shellErr != nil {
		return 0.0, fmt.Errorf("%v. Full shell output: %v", shellErr, output)
	}

	split := strings.Split(output, "\n")
	split = getNonWhiteSpaceLines(split)
	sampleRanges := getSampleRanges(split)
	sampleCoreUsage := map[int]map[int]float64{}
	for instance, sampleRange := range sampleRanges {
		parsedSample, err := parsePIDStatSample(split[sampleRange.Start:sampleRange.End])
		if err != nil {
			return 0.0, err
		}
		sampleCoreUsage[instance] = calculateCPUCoreUsage(parsedSample)
	}

	cpuAverage := calculateCPUCoreAverage(sampleCoreUsage)
	var max float64 = 0.0
	for _, value := range cpuAverage {
		max = math.Max(max, value)
	}

	return max, nil
}

func getProcessCoreCPULoad(processInfo []process.GeneralInfo) (float64, error) {
	var max float64 = 0
	for _, entry := range processInfo {
		command := fmt.Sprintf("pidstat -p %v -th 1 5", entry.PID)
		out, err := execBash(command)
		maxCoreUsage, err := getCPUCoreUsage(string(out), err)
		if err != nil {
			return 0, err
		}

		max = math.Max(max, maxCoreUsage)
	}

	return max, nil
}
