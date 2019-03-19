package nagiosfoundation

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func getPidNameWithHandler(readFile func(string) ([]byte, error), pid int) (string, error) {
	procFile := fmt.Sprintf("/proc/%d/stat", pid)
	procDataBytes, err := readFile(procFile)
	if err != nil {
		return "", err
	}

	procData := string(procDataBytes)

	procNameStart := strings.IndexRune(procData, '(') + 1
	procNameEnd := strings.IndexRune(procData, ')')

	if procNameStart >= procNameEnd {
		return "", errors.New("Could not parse process name")
	}

	procName := procData[procNameStart:procNameEnd]

	return procName, nil
}

func getPidName(pid int) (string, error) {
	return getPidNameWithHandler(ioutil.ReadFile, pid)
}

type processByNameHandlers struct {
	open       func(string) (*os.File, error)
	close      func(*os.File) error
	readDir    func(*os.File, int) ([]os.FileInfo, error)
	getPidName func(readFile func(string) ([]byte, error), pid int) (string, error)
	readFile   func(string) ([]byte, error)
}

func getProcessesByNameWithHandlers(svc processByNameHandlers, name string) ([]os.FileInfo, error) {
	var errorReturn error
	matchingEntries := make([]os.FileInfo, 0)

	dir, err := svc.open("/proc")
	if err != nil {
		matchingEntries = nil
		errorReturn = err
	}

	defer svc.close(dir)

	var procEntries []os.FileInfo
	if errorReturn == nil {
		procEntries, err = svc.readDir(dir, 0)

		if err != nil {
			matchingEntries = nil
			errorReturn = err
		}
	}

	if errorReturn == nil {
		for _, procEntry := range procEntries {
			// Skip entries that aren't directories
			if !procEntry.IsDir() {
				continue
			}

			// Skip entries that aren't numbers
			pid, err := strconv.Atoi(procEntry.Name())
			if err != nil {
				continue
			}

			if procName, _ := svc.getPidName(svc.readFile, pid); procName == name {
				matchingEntries = append(matchingEntries, procEntry)
			}
		}
	}

	return matchingEntries, errorReturn
}

func getProcessesByName(name string) ([]os.FileInfo, error) {
	svc := processByNameHandlers{
		open: os.Open,
		close: func(f *os.File) error {
			return f.Close()
		},
		readDir: func(f *os.File, entries int) ([]os.FileInfo, error) {
			return f.Readdir(entries)
		},
		getPidName: getPidNameWithHandler,
		readFile:   ioutil.ReadFile,
	}

	return getProcessesByNameWithHandlers(svc, name)
}

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

// GetHelpProcess returns help
func GetHelpProcess() string {
	return `Perform a check for a process by name to determine if the process
is running or not running. The default is to check for a running process.

The --name (-n) option is always required.
` + getHelpOsConstrained()
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

// checkProcessWithService provides a way to inject a custom
// service for interrogating the OS for the named process.
// This is mainly used for testing but can also be used for any
// application wishing to override the normal interrogations.
func checkProcessWithService(name string, checkType string, processService ProcessService) (string, int) {
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

// checkProcessCmd will interrogate the OS for details on
// a named process. The details of the interrogation
// depend on the check type.
func checkProcessCmd(name string, checkType string, checkProcess func(string, string, ProcessService) (string, int), processService ProcessService) (string, int) {
	var invalidCmdMsg string
	var msg string
	var retcode int

	checkType = strings.ToLower(checkType)

	if name == "" {
		invalidCmdMsg = invalidCmdMsg +
			"A process name must be specified with the --name (-n) option."
	} else if checkType != "running" && checkType != "notrunning" {
		invalidCmdMsg = invalidCmdMsg +
			fmt.Sprintf("Invalid check type (%s). Only \"running\" and \"notrunning\" are supported.",
				checkType)
	}

	if invalidCmdMsg != "" {
		msg = fmt.Sprintf("CheckProcess CRITICAL - %s", invalidCmdMsg)
		retcode = 2
	} else {
		msg, retcode = checkProcess(name, checkType, processService)
	}

	return msg, retcode
}

// CheckProcess checks for you
func CheckProcess(name string, checkType string) {
	msg, retcode := checkProcessCmd(name, checkType, checkProcessWithService, new(processHandler))

	if retcode >= 0 {
		fmt.Println(msg)
	}

	os.Exit(retcode)
}
