// +build !windows

package cpu

import (
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

func GetCPULoad() (float64, error) {
	var sleep = 1
	var usage, totalDiff float64

	beforeStats, err := getStats()
	if err != nil {
		return usage, err
	}

	time.Sleep(time.Duration(sleep) * time.Second)

	afterStats, err := getStats()
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

func getStats() ([]float64, error) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return []float64{}, err
	}

	line := strings.Split(string(contents), "\n")[0]
	stats := strings.Fields(line)[1:]

	result := make([]float64, len(stats))
	for i := range stats {
		result[i], err = strconv.ParseFloat(stats[i], 64)
		if err != nil {
			return result, err
		}
	}

	return result, nil
}
