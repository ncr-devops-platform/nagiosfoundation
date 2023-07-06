// +build !windows

package nagiosfoundation

import (
	"fmt"
	"os/exec"
	"strings"
)

// Not Used
//
// To be implemented in the future to make this Linux module
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

func systemdServiceTest(serviceName string, currentStateWanted bool) (string, int) {
	cmd := exec.Command("systemctl", "check", serviceName)
	out, err := cmd.CombinedOutput()
	state := strings.TrimSpace(string(out))

	var retcode int
	var info string
	var serviceState int

	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			info = fmt.Sprintf("%s not in a running state", serviceName)
			retcode = 2
		} else {
			info = fmt.Sprintf("Failed to execute systemctl. %s Status unknown: %v", serviceName, err)
			retcode = 2
		}

		serviceState = 0
	} else {
		info = fmt.Sprintf("%s in a running state", serviceName)
		serviceState = 1
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

	msg := fmt.Sprintf("%s %s - %s%s", serviceCheckName, responseStateText, info, actualInfo)

	if currentStateWanted {
		msg = msg + fmt.Sprintf(" | service_state=%d service_name=%s",
			serviceState, serviceName)
		retcode = 0
	}

	return msg, retcode
}

func checkServiceOsConstrained(name string, state string, user string, currentStateWanted bool, manager string) (string, int) {
	var msg string
	var retcode int

	switch manager {
	case "systemd":
		msg, retcode = systemdServiceTest(name, currentStateWanted)
	default:
		msg = fmt.Sprintf("%s CRITICAL - %s is not a valid service manager.", serviceCheckName, manager)
		retcode = 2
	}

	return msg, retcode
}
