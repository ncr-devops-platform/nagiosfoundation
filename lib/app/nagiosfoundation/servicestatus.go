package nagiosfoundation

import (
	"errors"
	"fmt"
	"strings"
)

const serviceCheckName = "CheckService"

type getServiceInfoFunc func(string) (string, string, string, int, error)

type serviceInfo struct {
	// The name of the service to process.
	desiredName string

	// The state of the service to match.
	desiredState string

	// The user of the service to match.
	desiredUser string

	// User only wants current state
	currentStateWanted bool

	actualName      string
	actualStateText string
	actualStateNbr  int
	actualUser      string

	getServiceInfo getServiceInfoFunc
}

// Returns the actual name of the service resulting from the service query.
func (i *serviceInfo) ActualName() string {
	return i.actualName
}

// Returns the actual state text of the service resulting from the service query.
func (i *serviceInfo) ActualStateText() string {
	return i.actualStateText
}

// Returns the actual state number of the service resulting from the service query.
func (i *serviceInfo) ActualStateNbr() int {
	return i.actualStateNbr
}

// Returns the actual user of the service resulting from the service query.
func (i *serviceInfo) ActualUser() string {
	return i.actualUser
}

// Checks for a match against the actual name of the service. The comparison
// is case insensitive.
func (i *serviceInfo) IsName(name string) bool {
	return strings.EqualFold(i.ActualName(), name)
}

// Checks for a match against the actual state of the service. The comparison
// is case insensitive.
func (i *serviceInfo) IsState(state string) bool {
	return strings.EqualFold(i.ActualStateText(), state)
}

// Checks for a match against the actual user of the service. The comparison
// is case insensitive.
func (i *serviceInfo) IsUser(user string) bool {
	return strings.EqualFold(i.ActualUser(), user)
}

// Executes the OS constrained function to retrieve information about a service.
// This information is derived differently in Windows and Linux and must execute
// an OS constrained method named getInfoOsConstrained().
func (i *serviceInfo) GetInfo() error {
	var err error

	if i.getServiceInfo == nil {
		return errors.New("No get service info handler declared")
	}

	i.actualName, i.actualUser, i.actualStateText, i.actualStateNbr, err = i.getServiceInfo(i.desiredName)

	return err
}

// Process the desired service info against the actual service info and return
// check text and a return code.
func (i *serviceInfo) ProcessInfo() (string, int) {
	var (
		checkInfo, nagiosInfo string
		retcode               int
	)

	if !i.IsName(i.desiredName) {
		if i.currentStateWanted == true {
			checkInfo = fmt.Sprintf("%s does not exist", i.desiredName)
			nagiosInfo = fmt.Sprintf("service_state=255 service_name=%s", i.desiredName)
			retcode = 0
		} else {
			checkInfo = fmt.Sprintf("%s does not exist", i.desiredName)
			retcode = 2
		}
	} else if i.currentStateWanted == true {
		checkInfo = fmt.Sprintf("%s is in a %s state", i.desiredName, i.ActualStateText())
		nagiosInfo = fmt.Sprintf("service_state=%d service_name=%s", i.ActualStateNbr(), i.desiredName)
		retcode = 0
	} else if i.desiredState != "" && i.desiredUser != "" {
		if i.IsState(i.desiredState) && i.IsUser(i.desiredUser) {
			checkInfo = fmt.Sprintf("%s in a %s state and started by user %s",
				i.ActualName(), i.ActualStateText(), i.ActualUser())
			retcode = 0
		} else {
			checkInfo = fmt.Sprintf("%s either not in a %s state or not started by user %s",
				i.ActualName(), i.desiredState, i.desiredUser)
			retcode = 2
		}
	} else if i.desiredState != "" {
		if i.IsState(i.desiredState) {
			checkInfo = fmt.Sprintf("%s in a %s state",
				i.ActualName(), i.ActualStateText())
			retcode = 0
		} else {
			checkInfo = fmt.Sprintf("%s not in a %s state",
				i.ActualName(), i.desiredState)
			retcode = 2
		}
	} else if i.desiredUser != "" {
		if i.IsUser(i.desiredUser) {
			checkInfo = fmt.Sprintf("%s started by user %s",
				i.ActualName(), i.ActualUser())
			retcode = 0
		} else {
			checkInfo = fmt.Sprintf("%s not started by user %s",
				i.ActualName(), i.desiredUser)
			retcode = 2
		}
	} else {
		checkInfo = fmt.Sprintf("%s in a %s state and started by user %s",
			i.ActualName(), i.ActualStateText(), i.ActualUser())
		retcode = 0
	}

	var responseStateText, actualInfo string

	if retcode == 0 {
		responseStateText = statusTextOK
		actualInfo = ""
	} else {
		responseStateText = statusTextCritical
		actualInfo = fmt.Sprintf(" (Name: %s, State: %s, User: %s)",
			i.ActualName(), i.ActualStateText(), i.ActualUser())
	}

	msg, _ := resultMessage(serviceCheckName, responseStateText, checkInfo+actualInfo, nagiosInfo)

	return msg, retcode
}

// CheckService checks a service based on name, state,
// user, and manager
func CheckService(name, state, user string, currentStateWanted bool, manager string) (string, int) {
	return checkServiceOsConstrained(name, state, user, currentStateWanted, manager)
}
