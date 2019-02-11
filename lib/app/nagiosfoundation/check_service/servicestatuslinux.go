// +build !windows

package nagiosfoundation

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Not Used
//
// To be impemented in the future to make this Linux module
// compliant with the serviceInfo structure like the Windows
// version is.
//
// When implemented, the name parameter will contain the name of the
// service for which to get info and returns will be the actual name,
// state, and user of the service being checked. Note the service
// manager is not passed in. When implemented, this will probably
// need to be added as a parameter but ignored in the Windows build
// or assigned to a local global in CheckServiceOsConstrained() then
// accessed here.
//
// See the Windows version of GetInfoOsConstrained() for an example.
func getInfoOsConstrained(name string) (string, string, string, error) {
	return "", "", "", nil
}

func systemdServiceTest(serviceName string) (string, int) {
	cmd := exec.Command("systemctl", "check", serviceName)
	out, err := cmd.CombinedOutput()
	state := strings.TrimSpace(string(out))

	var retcode int
	var info string

	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			info = fmt.Sprintf("%s not in a running state", serviceName)
			retcode = 2
		} else {
			info = fmt.Sprintf("failed to execute systemctl. %s Status unknown: %v", serviceName, err)
			retcode = 3
		}
	} else {
		info = fmt.Sprintf("%s in a running state", serviceName)
		retcode = 0
	}

	var responseStateText string
	var actualInfo string

	if retcode == 0 {
		responseStateText = "OK"
	} else {
		responseStateText = "CRITICAL"
		actualInfo = fmt.Sprintf(" (State: %s)", state)
	}

	msg := fmt.Sprintf("CheckService %s - %s%s", responseStateText, info, actualInfo)
	return msg, retcode
}

func checkServiceOsConstrained(name string, state string, user string, manager string) {
	if manager == "" {
		fmt.Fprintln(os.Stderr, "A service manager must be specified. Currently only systemd is supported.")
		os.Exit(0)
	}
	if state != "" || user != "" {
		fmt.Fprintln(os.Stderr, "Linux functionality is limited to only checking for a running state.")
		fmt.Fprintln(os.Stderr, "The -state and -user options are not used but at least one of the options was given.")
		os.Exit(0)
	}

	var msg string
	var retcode int

	switch manager {
	case "systemd":
		msg, retcode = systemdServiceTest(name)
	default:
		msg = fmt.Sprintf("CheckServiceRunning CRITICAL - %s not in a valid service manager.", manager)
		retcode = 3
	}

	fmt.Println(msg)
	os.Exit(retcode)
}

// Show help specific to the operating system.
func showHelpOsConstrained() {
	fmt.Printf(`    -manager <service manager>: Required. The service manager executed for the check.
	  systemd is the only supported service manager

    The only check done is for a running state. Both the -name and -manager options must
    be specified the the service is only checked to see if it is running.
`)
}
