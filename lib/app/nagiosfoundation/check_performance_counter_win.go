package nagiosfoundation

import (
	"fmt"
	"os"

	"github.com/ncr-devops-platform/nagiosfoundation/lib/pkg/perfcounters"
)

// CheckPerformanceCounter executes CheckPerformanceCounterWitHandler(),
// passing it the OS constrained ReadPerformanceCounter() function, prints
// the returned message and exits with the returned exit code.
func CheckPerformanceCounter(warning, critical float64, greaterThan bool, pollingAttempts, pollingDelay int, metricName, counterName string) {
	msg, retval := CheckPerformanceCounterWithHandler(warning,
		critical, greaterThan, pollingAttempts, pollingDelay,
		metricName, counterName, perfcounters.ReadPerformanceCounter)

	fmt.Println(msg)
	os.Exit(retval)
}
