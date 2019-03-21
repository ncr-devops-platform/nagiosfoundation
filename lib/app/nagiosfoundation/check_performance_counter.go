package nagiosfoundation

import (
	"errors"
	"fmt"

	"github.com/ncr-devops-platform/nagiosfoundation/lib/pkg/nagiosformatters"
	"github.com/ncr-devops-platform/nagiosfoundation/lib/pkg/perfcounters"
)

// CheckPerformanceCounterWithHandler fetches a performance counter
// specified with the counterName parameter. It then performs checks
// against the value based on the threshold test specified along with
// the warning and critical thresholds.
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
