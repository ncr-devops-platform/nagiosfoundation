// +build windows

package memory

import (
	"github.com/jkerry/nagiosfoundation/lib/pkg/perfcounters"
)

func GetAvgFreeMemory(pollingAttempts int, pollingDelay int) (float64, error) {
	counter, err := perfcounters.ReadPerformanceCounter("\\Memory\\Available MBytes", pollingAttempts, pollingDelay)
	if err == nil {
		return counter.Value, nil
	}
	return 0.0, err
}
