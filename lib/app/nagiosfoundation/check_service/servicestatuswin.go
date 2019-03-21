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
