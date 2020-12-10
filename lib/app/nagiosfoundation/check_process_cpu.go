package nagiosfoundation

import (
	"errors"

	"github.com/ncr-devops-platform/nagiosfoundation/lib/pkg/cpu"

	"github.com/ncr-devops-platform/nagiosfoundation/lib/pkg/nagiosformatters"
)

// CheckProcessCPUWithHandler gets the CPU load of a process then emits a critical response
// if it's above the critical argument, a warning if it's above
// warning argument and good response for everything else.
//
// Returns are a response message and response code.
func CheckProcessCPUWithHandler(warning, critical int, processName, metricName string, perCoreCalculation bool,
	processCPUCoreHandler func(string, bool) (float64, error)) (string, int) {
	const checkName = "CheckProcessCPULoad"

	var msg string
	var retcode int
	var value float64
	var err error

	if processCPUCoreHandler == nil {
		err = errors.New("No GetProcessCPULoad() service")
	} else {
		value, err = processCPUCoreHandler(processName, perCoreCalculation)
	}

	if err == nil {
		msg, retcode = nagiosformatters.GreaterFormatNagiosCheck(checkName, value, float64(warning), float64(critical), metricName)
	} else {
		msg, _ = resultMessage(checkName, statusTextCritical, err.Error())
		retcode = 2
	}

	return msg, retcode
}

// CheckProcessCPU executes CheckProcessCPUWithHandler(), passing it the OS
// constrained GetProcessCPULoad() function, prints the returned message
// and exits with the returned exit code.
func CheckProcessCPU(warning, critical int, processName, metricName string, perCoreCalculation bool) (string, int) {
	return CheckProcessCPUWithHandler(warning, critical, processName, metricName, perCoreCalculation, cpu.GetProcessCPULoad)
}
