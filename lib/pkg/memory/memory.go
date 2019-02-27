// +build !windows

package memory

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

func GetFreeMemory() (float64, error) {
	out, err := exec.Command("free").Output()
	if err != nil {
		if glog.V(1) {
			glog.Errorf("Error running free: %v", err)
		}

		return 0.0, err
	}

	lines := strings.Split(string(out), "\n")

	free, err := strconv.ParseFloat(strings.Fields(lines[2])[3], 64)
	if err != nil {
		if glog.V(1) {
			glog.Errorf("Error parsing free output to float64: %v", err)
		}

		return 0.0, err
	}

	return free, nil
}
