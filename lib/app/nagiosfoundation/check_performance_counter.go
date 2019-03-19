package nagiosfoundation

import (
	"errors"
	"fmt"

	"github.com/ncr-devops-platform/nagiosfoundation/lib/pkg/nagiosformatters"
	"github.com/ncr-devops-platform/nagiosfoundation/lib/pkg/perfcounters"
)

// GetHelpPerformanceCounter returns a string containing help for the
// performance counter functionality.
func GetHelpPerformanceCounter() string {
	return `The performance counter check is Windows only. It retrieves a Windows Performance Counter
(--counter_name) and compares it to --critical and --warning then outputs an appropriate response
based on the check. Many flags make this check quite configurable.

The defaults for this check have the --critical and --warning values set to 0, and the counter value
retrieved is compared to be lesser than those values. Generally a counter value will be > 0, causing
this check to generally emit an OK response when using these defaults.`
}

// CheckPerformanceCounterWithHandler fetches a performance counter
// specified with the -counter_name flag. It then performs checks against
// the value based on the threshold test specified along with the
// warning and critical thresholds.
//
// Returns are a message stating the results of the check and a return
// value from the check.
func CheckPerformanceCounterWithHandler(warning, critical float64, greaterThan bool, pollingAttempts, pollingDelay int, metricName, counterName string, perfCounterHandler func(string, int, int) (perfcounters.PerformanceCounter, error)) (string, int) {
	var msg string
	var retcode int
	var counter perfcounters.PerformanceCounter
	var err error

	if perfCounterHandler == nil {
		err = errors.New("No ReadPerformanceCounter() service")
	} else {
		counter, err = perfCounterHandler(counterName, pollingAttempts, pollingDelay)
	}

	if err == nil {
		if greaterThan {
			msg, retcode = nagiosformatters.GreaterFormatNagiosCheck(counterName, counter.Value, warning, critical, metricName)
		} else {
			msg, retcode = nagiosformatters.LesserFormatNagiosCheck(counterName, counter.Value, warning, critical, metricName)
		}
	} else {
		msg = fmt.Sprintf("%s CRITICAL - %s", counterName, err)
		retcode = 2
	}

	return msg, retcode
}
