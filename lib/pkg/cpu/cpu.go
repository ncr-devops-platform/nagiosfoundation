package cpu

import (
	"errors"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

func getStats(getStatsData func() (string, error)) ([]float64, error) {
	if getStatsData == nil {
		return []float64{}, errors.New("No stats data handler given")
	}

	statsData, err := getStatsData()

	if err != nil {
		return []float64{}, err
	}

	line := strings.Split(statsData, "\n")[0]
	fields := strings.Fields(line)

	if len(fields) < 3 {
		return []float64{}, errors.New("CPU data not found")
	}

	stats := fields[1:]
	result := make([]float64, len(stats))

	for i := range stats {
		result[i], err = strconv.ParseFloat(stats[i], 64)
		if err != nil {
			return result, err
		}
	}

	return result, nil
}

func getCPULoadLinuxWithHandler(getStatsData func() (string, error)) (float64, error) {
	var sleep = 1
	var usage, totalDiff float64

	beforeStats, err := getStats(getStatsData)
	if err != nil {
		return usage, err
	}

	time.Sleep(time.Duration(sleep) * time.Second)

	afterStats, err := getStats(getStatsData)
	if err != nil {
		return usage, err
	}

	diffStats := make([]float64, len(beforeStats))
	for i := range beforeStats {
		diffStats[i] = afterStats[i] - beforeStats[i]
		totalDiff += diffStats[i]
	}

	usage = 100.0 * (totalDiff - diffStats[3]) / totalDiff
	return usage, nil
}

func getStatsDataService() (string, error) {
	var statsData string

	contents, err := ioutil.ReadFile("/proc/stat")

	if err == nil {
		statsData = string(contents)
	}

	return statsData, err
}

func getCPULoadLinux() (float64, error) {
	return getCPULoadLinuxWithHandler(getStatsDataService)
}

// GetCPULoad returns the current CPU load as a percentage.
func GetCPULoad() (float64, error) {
	return getCPULoadOsConstrained()
}
