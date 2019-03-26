// +build windows

package nagiosfoundation

import (
	"errors"
	"fmt"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"

	"github.com/StackExchange/wmi"
)

func getInfoWmi(name string) (string, string, string, error) {
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

func getStateText(state svc.State) string {
	var txtState string

	switch state {
	case windows.SERVICE_STOPPED:
		txtState = "Stopped"
	case windows.SERVICE_START_PENDING:
		txtState = "Start Pending"
	case windows.SERVICE_STOP_PENDING:
		txtState = "Stop Pending"
	case windows.SERVICE_RUNNING:
		txtState = "Running"
	case windows.SERVICE_CONTINUE_PENDING:
		txtState = "Continue Pending"
	case windows.SERVICE_PAUSE_PENDING:
		txtState = "Pause Pending"
	case windows.SERVICE_PAUSED:
		txtState = "Paused"
	default:
		txtState = "Unknown"
	}

	return txtState
}

func getInfoSvcMgr(name string) (string, string, string, error) {
	var serviceName string
	var serviceStartName string
	var serviceState string

	mgrPtr, err := mgr.Connect()
	if err != nil {
		err = errors.New("Connect to Service Manager failed: " + err.Error())
	}

	var service *mgr.Service
	if err == nil {
		service, err = mgrPtr.OpenService(name)
		if err == nil {
			serviceName = service.Name
		} else {
			// Error is valid - service doesn't exist which is
			// what is being checked.
			err = nil
		}
	}

	if service != nil {
		var config mgr.Config
		if err == nil {
			config, err = service.Config()

			if err == nil {
				serviceStartName = config.ServiceStartName
			} else {
				err = errors.New("Getting service configuration failed: " + err.Error())
			}
		}

		var status svc.Status
		if err == nil {
			status, err = service.Query()

			if err == nil {
				serviceState = getStateText(status.State)
			} else {
				err = errors.New("Query service failed: " + err.Error())
			}
		}
	}

	return serviceName, serviceStartName, serviceState, err
}

func checkServiceOsConstrained(name string, state string, user string, manager string) (string, int) {
	managers := make(map[string]getServiceInfoFunc)
	managers["wmi"] = getInfoWmi
	managers["svcmgr"] = getInfoSvcMgr

	if manager == "" {
		manager = "wmi"
	}

	var msg string
	var retcode int

	if _, ok := managers[manager]; !ok {
		managersList := ""
		for key, _ := range managers {
			if managersList != "" {
				managersList = managersList + ", "
			}
			managersList = managersList + "\"" + key + "\""
		}

		msg = fmt.Sprintf("%s CRITICAL - Service manager \"%s\" not valid. Valid managers are %s.", serviceCheckName, manager, managersList)
		retcode = 2
	} else {
		i := serviceInfo{
			desiredName:    name,
			desiredState:   state,
			desiredUser:    user,
			getServiceInfo: managers[manager],
		}

		err := i.GetInfo()

		if err != nil {
			msg = fmt.Sprintf("%s CRITICAL - %s", serviceCheckName, err)
			retcode = 2
		} else {
			msg, retcode = i.ProcessInfo()
		}
	}

	return msg, retcode
}
