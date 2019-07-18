// +build windows linux darwin

package nagiosfoundation

import (
	"errors"

	"github.com/ncr-devops-platform/nagiosfoundation/lib/pkg/cpu"
	"github.com/ncr-devops-platform/nagiosfoundation/lib/pkg/nagiosformatters"
)

// CheckCPUWithHandler gets the CPU load then emits a critical response
// if it's above the critical argument, a warning if it's above
// warning argument and good response for everything else.
//
// Returns are a response message and response code.
func CheckCPUWithHandler(warning, critical int, metricName string, cpuHandler func() (float64, error)) (string, int) {
	const checkName = "CheckAVGCPULoad"

	var msg string
	var retcode int
	var value float64
	var err error

	if cpuHandler == nil {
		err = errors.New("No GetCPULoad() service")
	} else {
		value, err = cpuHandler()
	}

	if err == nil {
		msg, retcode = nagiosformatters.GreaterFormatNagiosCheck(checkName, value, float64(warning), float64(critical), metricName)
	} else {
		msg, _ = resultMessage(checkName, statusTextCritical, err.Error())
		retcode = 2
	}

	return msg, retcode
}

// CheckCPU executes CheckCPUWithHandler(), passing it the OS
// constrained GetCPULoad() function, prints the returned message
// and exits with the returned exit code.
func CheckCPU(warning, critical int, metricName string) (string, int) {
	return CheckCPUWithHandler(warning, critical, metricName, cpu.GetCPULoad)
}
