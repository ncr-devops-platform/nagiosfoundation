// +build windows

package memory

import (
	"github.com/ncr-devops-platform/nagiosfoundation/lib/pkg/perfcounters"
)

func GetFreeMemory() (float64, error) {
	counter, err := perfcounters.ReadPerformanceCounter("\\Memory\\Available MBytes", 2, 1)
	if err == nil {
		return counter.Value, nil
	}
	return 0.0, err
}
