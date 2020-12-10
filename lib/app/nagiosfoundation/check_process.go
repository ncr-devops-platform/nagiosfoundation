package nagiosfoundation

import (
	"fmt"
	"strconv"
	"strings"
)

const checkProcessName = "CheckProcess"

// ProcessService is an interface required by ProcessCheck.
//
// The given a process name, the method IsProcessRunning()
// must return true if the named process is running, otherwise
// false. Note the code will be different for each OS.
type ProcessService interface {
	IsProcessRunning(string) bool
}

type processHandler struct{}

func (p processHandler) IsProcessRunning(name string) bool {
	return isProcessRunningOsConstrained(name)
}

// ProcessCheck is used to encapsulate a named process
// along with the methods used to get information about
// that process. Currently the only check is for the named
// process running.
type ProcessCheck struct {
	ProcessName string

	ProcessCheckHandler ProcessService
}

// IsProcessRunning interrogates the OS for the named
// process to check if it's running. Note this function
// calls IsProcessRunning in the injected service and
// in this implementation will ultimately call an OS
// constrained function.
func (p ProcessCheck) IsProcessRunning() bool {
	return p.ProcessCheckHandler.IsProcessRunning(p.ProcessName)
}

func checkRunning(processCheck ProcessCheck, metricName string, invert bool) (string, int) {
	var msg string
	var retcode int
	var responseStateText string
	var checkInfo string

	result := processCheck.IsProcessRunning()
	if result != invert {
		retcode = statusCodeOK
		responseStateText = statusTextOK
	} else {
		retcode = statusCodeCritical
		responseStateText = statusTextCritical
	}

	nagiosOutput := metricName + "="
	if result == true {
		checkInfo = ""
		nagiosOutput = nagiosOutput + strconv.Itoa(statusCodeOK)
	} else {
		checkInfo = "not "
		nagiosOutput = nagiosOutput + strconv.Itoa(statusCodeCritical)
	}

	msg, _ = resultMessage(checkProcessName, responseStateText,
		fmt.Sprintf("Process %s is %srunning", processCheck.ProcessName, checkInfo),
		nagiosOutput)

	return msg, retcode
}

// checkProcessWithService provides a way to inject a custom
// service for interrogating the OS for the named process.
// This is mainly used for testing but can also be used for any
// application wishing to override the normal interrogations.
func checkProcessWithService(name, checkType, metricName string, processService ProcessService) (string, int) {
	pc := ProcessCheck{
		ProcessName:         name,
		ProcessCheckHandler: processService,
	}

	var msg string
	var retcode int

	switch checkType {
	case "running":
		msg, retcode = checkRunning(pc, metricName, false)
	case "notrunning":
		msg, retcode = checkRunning(pc, metricName, true)
	default:
		msg = fmt.Sprintf("Invalid check type: %s", checkType)
		retcode = statusCodeCritical
	}

	return msg, retcode
}

// checkProcessCmd will interrogate the OS for details on
// a named process. The details of the interrogation
// depend on the check type.
func checkProcessCmd(name, checkType, metricName string, checkProcess func(string, string, string, ProcessService) (string, int), processService ProcessService) (string, int) {
	var invalidParametersMsg string
	var msg string
	var retcode int

	checkType = strings.ToLower(checkType)

	if name == "" {
		invalidParametersMsg = invalidParametersMsg +
			"A process name must be specified."
	} else if checkType != "running" && checkType != "notrunning" {
		invalidParametersMsg = invalidParametersMsg +
			fmt.Sprintf("Invalid check type (%s). Only \"running\" and \"notrunning\" are supported.",
				checkType)
	}

	if invalidParametersMsg != "" {
		msg, _ = resultMessage(checkProcessName, statusTextCritical, invalidParametersMsg)
		retcode = statusCodeCritical
	} else {
		msg, retcode = checkProcess(name, checkType, metricName, processService)
	}

	return msg, retcode
}

// CheckProcess finds a process by name to determine
// if it is running or not running.
func CheckProcess(name, checkType, metricName string) (string, int) {
	return checkProcessCmd(name, checkType, metricName, checkProcessWithService, new(processHandler))
}
