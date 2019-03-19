// +build windows

package memory

import (
	"github.com/ncr-devops-platform/nagiosfoundation/lib/pkg/perfcounters"
)

func getFreeMemoryOsConstrained() uint64 {
	counter, err := perfcounters.ReadPerformanceCounter("\\Memory\\Available Bytes", 2, 1)

	var memoryAvailable uint64

	if err != nil {
		memoryAvailable = uint64(0)
	} else {
		memoryAvailable = uint64(counter.Value)
	}

	return uint64(memoryAvailable)
}
