// +build windows

package nagiosfoundation

import (
	"fmt"
	"log"
	"os"

	"github.com/StackExchange/wmi"
)

func getInfoOsConstrained(name string) (string, string, string, error) {
	type win32_Service struct {
		Name      string
		State     string
		StartName string
	}

	var dst []win32_Service
	var actualName string
	var actualUser string
	var actualState string

	w := fmt.Sprintf("where name = '%v'", name)

	query := wmi.CreateQuery(&dst, w)

	err := wmi.Query(query, &dst)

	if err == nil && len(dst) >= 1 {
		actualName = dst[0].Name
		actualUser = dst[0].StartName
		actualState = dst[0].State
	}

	return actualName, actualUser, actualState, err
}

func checkServiceOsConstrained(name string, state string, user string, manager string) {
	if manager != "" {
		fmt.Fprintln(os.Stderr, "Service manager", manager, "set but is not used in Windows")
		os.Exit(0)
	}

	i := serviceInfo{
		desiredName:    name,
		desiredState:   state,
		desiredUser:    user,
		getServiceInfo: func(n string) (string, string, string, error) { return getInfoOsConstrained(n) },
	}

	err := i.GetInfo()

	if err != nil {
		log.Fatal(err)
	}

	msg, retcode := i.ProcessInfo()

	fmt.Println(msg)
	os.Exit(retcode)
}

// Show help specific to the operating system.
func showHelpOsConstrained() {
	fmt.Printf(`    -state <service state>: State to check. Examples: running, stopped
    -user <user name>: User to check. Example: "NT AUTHORITY\LocalService"

    Some examples:
      check_service.exe -name audiosrv
        Checks for the service to exist and shows the service state and user.
      check_service.exe -name audiosrv -state running
        Checks for the service in the running state.
      check_service.exe -name audiosrv -state running -user "NT AUTHORITY\LocalService"
        Checks for the service in the running state and running as user.
      check_service.exe -name audiosrv -user "NT AUTHORITY\LocalService"
        Checks for the service to exist and would be run as user.
`)
}
