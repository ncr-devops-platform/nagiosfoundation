// +build windows

package cpu

import (
	"github.com/ncr-devops-platform/nagiosfoundation/lib/pkg/perfcounters"
)

func getCPULoadOsConstrained() (float64, error) {
	counter, err := perfcounters.ReadPerformanceCounter("\\Processor(_Total)\\% Processor Time", 2, 1)
	if err == nil {
		return counter.Value, nil
	}
	return 0.0, err
}
