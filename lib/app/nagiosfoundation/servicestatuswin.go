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

func getStateNbrFromText(state string) int {
	var nbrState int

	switch state {
	case "Stopped":
		nbrState = 6
	case "Start Pending":
		nbrState = 2
	case "Stop Pending":
		nbrState = 5
	case "Running":
		nbrState = 0
	case "Continue Pending":
		nbrState = 4
	case "Pause Pending":
		nbrState = 3
	case "Paused":
		nbrState = 1
	default:
		nbrState = 7
	}

	return nbrState
}

func getInfoWmi(name string) (string, string, string, int, error) {
	type win32_Service struct {
		Name      string
		State     string
		StartName string
	}

	var dst []win32_Service
	var actualName string
	var actualUser string
	var actualStateText string
	var actualStateNbr int

	w := fmt.Sprintf("where name = '%v'", name)

	query := wmi.CreateQuery(&dst, w)

	err := wmi.Query(query, &dst)

	if err == nil && len(dst) >= 1 {
		actualName = dst[0].Name
		actualUser = dst[0].StartName
		actualStateText = dst[0].State
		actualStateNbr = getStateNbrFromText(actualStateText)
	}

	return actualName, actualUser, actualStateText, actualStateNbr, err
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

func getStateNbr(state svc.State) int {
	var nbrState int

	switch state {
	case windows.SERVICE_STOPPED:
		nbrState = 6
	case windows.SERVICE_START_PENDING:
		nbrState = 2
	case windows.SERVICE_STOP_PENDING:
		nbrState = 5
	case windows.SERVICE_RUNNING:
		nbrState = 0
	case windows.SERVICE_CONTINUE_PENDING:
		nbrState = 4
	case windows.SERVICE_PAUSE_PENDING:
		nbrState = 3
	case windows.SERVICE_PAUSED:
		nbrState = 1
	default:
		nbrState = 7
	}

	return nbrState
}

func getInfoSvcMgr(name string) (string, string, string, int, error) {
	var serviceName string
	var serviceStartName string
	var serviceStateText string
	var serviceStateNbr int

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
				serviceStateText = getStateText(status.State)
				serviceStateNbr = getStateNbr(status.State)
			} else {
				err = errors.New("Query service failed: " + err.Error())
			}
		}
	}

	return serviceName, serviceStartName, serviceStateText, serviceStateNbr, err
}

func checkServiceOsConstrained(name string, state string, user string, currentStateWanted bool, manager string) (string, int) {
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
			desiredName:        name,
			desiredState:       state,
			desiredUser:        user,
			currentStateWanted: currentStateWanted,
			getServiceInfo:     managers[manager],
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
