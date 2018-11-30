// +build !windows

package nagiosfoundation

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func CheckServiceRunning() {
	var serviceName = flag.String("service_name", "", "the name of the service to check")
	var serviceManager = flag.String("service_manager", "", "the name of local service manager. Allowed options are: systemd")
	flag.Parse()
	var msg string
	var retcode int
	switch *serviceManager {
	case "systemd":
		systemdMsg, systemdCode := systemdServiceTest(*serviceName)
		msg = systemdMsg
		retcode = systemdCode
	default:
		msg = fmt.Sprintf("CheckServiceRunning CRITICAL - %s not in a valid service manager.", *serviceManager)
		retcode = 3
	}
	fmt.Println(msg)
	os.Exit(retcode)
}

func systemdServiceTest(serviceName string) (string, int) {
	cmd := exec.Command("systemctl", "check", serviceName)
	out, err := cmd.CombinedOutput()
	var msg string
	var retcode int
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			msg = fmt.Sprintf("CheckServiceRunning CRITICAL - %s not in a running state: %s", serviceName, out)
			retcode = 2
		} else {
			msg = fmt.Sprintf("CheckServiceRunning CRITICAL - failed to execute systemctl. %s Status unknown: %v", serviceName, err)
			retcode = 3
		}
	} else {
		msg = fmt.Sprintf("CheckServiceRunning OK - %s in a running state", serviceName)
		retcode = 0
	}
	return msg, retcode
}
