package nagiosfoundation

import (
	"errors"
	"fmt"

	"github.com/ncr-devops-platform/nagiosfoundation/lib/pkg/memory"
	"github.com/ncr-devops-platform/nagiosfoundation/lib/pkg/nagiosformatters"
)

// CheckAvailableMemoryWithHandler determines the percentage of
// memory used and emits a critical response if it's over
// the critical argument, a warning response if it's over the
// warning argument, and good response otherwise.
func CheckAvailableMemoryWithHandler(warning, critical int, metricName string, memoryHandler func() uint64) (string, int) {
	const checkName = "CheckAvailableMemory"

	var msg string
	var retcode int
	var usedMemoryPercentage uint64
	var err error

	if memoryHandler == nil {
		err = errors.New("No used memory percentage service")
	} else {
		usedMemoryPercentage = memoryHandler()
	}

	if err != nil || usedMemoryPercentage == 0 {
		msg = fmt.Sprintf("%s CRITICAL - %s", checkName, err)
		retcode = 2
	} else {
		msg, retcode = nagiosformatters.GreaterFormatNagiosCheck(checkName, float64(usedMemoryPercentage), float64(warning), float64(critical), metricName)
	}

	return msg, retcode
}

// CheckAvailableMemory executes CheckAvailableMemoryWithHandler(),
// passing it the OS constranted GetFreeMemory() function, prints
// the returned message and exits with the returned exit code.
//
// Returns are those of CheckAvailableMemoryWithHandler()
func CheckAvailableMemory(warning, critical int, metricName string) (string, int) {
	return CheckAvailableMemoryWithHandler(warning, critical, metricName, memory.GetUsedMemoryPercentage)
}
