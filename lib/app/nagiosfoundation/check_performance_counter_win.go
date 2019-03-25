package nagiosfoundation

import (
	"github.com/ncr-devops-platform/nagiosfoundation/lib/pkg/perfcounters"
)

// CheckPerformanceCounter executes CheckPerformanceCounterWitHandler(),
// passing it the OS constrained ReadPerformanceCounter() function, prints
// the returned message and exits with the returned exit code.
func CheckPerformanceCounter(warning, critical float64, greaterThan bool, pollingAttempts, pollingDelay int, metricName, counterName string) (string, int) {
	return CheckPerformanceCounterWithHandler(warning,
		critical, greaterThan, pollingAttempts, pollingDelay,
		metricName, counterName, perfcounters.ReadPerformanceCounter)
}
