package nagiosfoundation

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

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

func showHelp() {
	fmt.Printf(
		`check_process -name <process name> [ other options ]
  Perform various checks for a process. These checks depend on the -check-type
  flag which defaults to "running". The -name option is always required.

	-name <process name>: Required. The name of the process to check
	-type <check type>: Defaults to "running". Supported types are "running"
	  "notrunning".
`)

	showHelpOsConstrained()
}

func checkRunning(processCheck ProcessCheck, invert bool) (string, int) {
	var msg string
	var retcode int
	var responseStateText string
	var checkInfo string

	result := processCheck.IsProcessRunning()
	if result != invert {
		retcode = 0
		responseStateText = "OK"
	} else {
		retcode = 2
		responseStateText = "CRITICAL"
	}

	if result == true {
		checkInfo = ""
	} else {
		checkInfo = "not "
	}

	msg = fmt.Sprintf("CheckProcess %s - Process %s is %srunning", responseStateText, processCheck.ProcessName, checkInfo)

	return msg, retcode
}

// CheckProcessWithService provides a way to inject a custom
// service for interrogating the OS for the named process.
// This is mainly used for testing but can also be used for any
// application wishing to override the normal interrogations.
func CheckProcessWithService(name string, checkType string, processService ProcessService) (string, int) {
	pc := ProcessCheck{
		ProcessName:         name,
		ProcessCheckHandler: processService,
	}

	var msg string
	var retcode int

	switch checkType {
	case "running":
		msg, retcode = checkRunning(pc, false)
	case "notrunning":
		msg, retcode = checkRunning(pc, true)
	default:
		msg = fmt.Sprintf("Invalid check type: %s", checkType)
		retcode = 3
	}

	return msg, retcode
}

// CheckProcess will interrogate the OS for details on
// a named process. The details of the interrogation
// depend on the check type.
func CheckProcess(name string, checkType string) (string, int) {
	msg, retcode := CheckProcessWithService(name, checkType, new(processHandler))

	return msg, retcode
}

// CheckProcessFlags provides an entry point with no
// parameters. Instead it relies on command line flags for
// the appropriate parameters and calls CheckProcess.
func CheckProcessFlags() {
	namePtr := flag.String("name", "", "process name")
	checkTypePtr := flag.String("type", "running", "type of check (currently only \"running\" is supported")
	flag.Parse()

	*checkTypePtr = strings.ToLower(*checkTypePtr)

	invalidCmdMsg := ""

	if *namePtr == "" {
		invalidCmdMsg = invalidCmdMsg +
			"A process name must be specified with the -name option.\n"
	}

	if *checkTypePtr != "running" && *checkTypePtr != "notrunning" {
		invalidCmdMsg = invalidCmdMsg +
			fmt.Sprintf("Invalid check type (%s). Only \"running\" and \"notrunning\" are supported.\n",
				*checkTypePtr)
	}

	if invalidCmdMsg != "" {
		fmt.Printf("%s\n", invalidCmdMsg)
		showHelp()
	} else {
		msg, retcode := CheckProcess(*namePtr, *checkTypePtr)

		fmt.Println(msg)
		os.Exit(retcode)
	}
}
