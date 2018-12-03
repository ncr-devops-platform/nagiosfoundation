// +build !windows

package memory

import (
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
)

func GetAvgFreeMemory(pollingAttempts int, pollingDelay int) (float64, error) {
	var total float64
	total = 0
	for i := 0; i < pollingAttempts; i++ {
		freeMemory, err := pollForFreeMemory()
		if err != nil {
			glog.Errorf("Could not aggregate free memory result %v", err)
			return 0.0, err
		}
		total += freeMemory
		time.Sleep(time.Duration(pollingDelay) * time.Second)
	}
	return (total / float64(pollingAttempts)), nil
}

func pollForFreeMemory() (float64, error) {
	out, err := exec.Command("free").Output()
	if err != nil {
		glog.Errorf("Error running free: %v", err)
		return 0.0, err
	}

	lines := strings.Split(string(out), "\n")

	free, err := strconv.ParseFloat(strings.Fields(lines[2])[3], 64)
	if err != nil {
		glog.Errorf("Error parsing free output to float64: %v", err)
		return 0.0, err
	}

	return free, nil
}
