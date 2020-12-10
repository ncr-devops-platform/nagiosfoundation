package nagiosfoundation

import (
	"errors"

	"github.com/ncr-devops-platform/nagiosfoundation/lib/pkg/memory"
	"github.com/ncr-devops-platform/nagiosfoundation/lib/pkg/nagiosformatters"
)

// CheckProcessMemoryWithHandler determines the percentage of
// memory used by a process and emits a critical response if it's over
// the critical argument, a warning response if it's over the
// warning argument, and good response otherwise.
func CheckProcessMemoryWithHandler(warning, critical int, processName, metricName string, memoryHandler func(string) (float64, error)) (string, int) {
	const checkName = "CheckProcessMemory"

	var msg string
	var retcode int
	var usedMemoryPercentage float64
	var err error

	if memoryHandler == nil {
		err = errors.New("No GetProcessMemoryPercentage service")
	} else {
		usedMemoryPercentage, err = memoryHandler(processName)
	}

	if err != nil {
		msg, _ = resultMessage(checkName, statusTextCritical, err.Error())
		retcode = 2
	} else {
		msg, retcode = nagiosformatters.GreaterFormatNagiosCheck(checkName, usedMemoryPercentage, float64(warning), float64(critical), metricName)
	}

	return msg, retcode
}

// CheckProcessMemory executes CheckProcessMemoryWithHandler(),
// passing it the OS constrained GetProcessMemoryPercentage() function, prints
// the returned message and exits with the returned exit code.
//
// Returns are those of CheckProcessMemoryWithHandler()
func CheckProcessMemory(warning, critical int, processName, metricName string) (string, int) {
	return CheckProcessMemoryWithHandler(warning, critical, processName, metricName, memory.GetProcessMemoryPercentage)
}
