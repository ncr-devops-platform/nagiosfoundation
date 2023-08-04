package nagiosfoundation

import "testing"

func TestActualIs(t *testing.T) {
	var goodName = "goodName"
	var goodState = "goodState"
	var goodUser = "goodUser"

	var badName = "badName"
	var badState = "badState"
	var badUser = "badUser"

	var matchName = "GOODNAME"
	var matchStateText = "GOODSTATE"
	var matchStateNbr = 0
	var matchUser = "GOODUSER"

	si := serviceInfo{
		getServiceInfo: func(n string) (string, string, string, int, error) {
			return goodName, goodUser, goodState, 0, nil
		},
	}

	var err error
	err = si.GetInfo()
	if err != nil {
		t.Errorf("GetInfo() returned error but was fed good data")
	}

	var actualResult string

	actualResult = si.ActualName()
	if actualResult != goodName {
		t.Errorf("ActualName (%s) does not match desiredName (%s)", actualResult, goodName)
	}

	actualResult = si.ActualStateText()
	if actualResult != goodState {
		t.Errorf("ActualState (%s) does not match desiredState (%s)", actualResult, goodState)
	}

	actualResult = si.ActualUser()
	if actualResult != goodUser {
		t.Errorf("ActualUser (%s) does not match desiredUser (%s)", actualResult, goodUser)
	}

	actualResult = si.ActualName()
	if actualResult != goodName {
		t.Errorf("ActualName (%s) does not match desiredName (%s)", actualResult, goodName)
	}

	var isResult bool
	isResult = si.IsName(matchName)
	if !isResult {
		t.Errorf("IsName(%s) does not match actualName (%s)", matchName, si.ActualName())
	}

	isResult = si.IsState(matchStateText)
	if !isResult {
		t.Errorf("IsState(%s) does not match actualState (%s)", matchStateText, si.ActualStateText())
	}

	if si.ActualStateNbr() != matchStateNbr {
		isResult = false
	} else {
		isResult = true
	}
	if !isResult {
		t.Errorf("IsState(%d) does not match actualState (%d)", matchStateNbr, si.ActualStateNbr())
	}

	isResult = si.IsUser(matchUser)
	if !isResult {
		t.Errorf("IsUser(%s) does not match actualUser (%s)", matchUser, si.ActualUser())
	}

	si.desiredName = goodName
	si.desiredState = goodState
	si.desiredUser = goodUser

	var msg string
	var retcode int

	// All good check
	msg, retcode = si.ProcessInfo()
	if retcode != 0 {
		t.Errorf("ProcessInfo() failed on good data with retcode %d, msg %s", retcode, msg)
	}

	// Check all with bad name
	si.desiredName = badName
	msg, retcode = si.ProcessInfo()
	if retcode != 2 {
		t.Errorf("ProcessInfo() failed on bad name with retcode %d, msg %s", retcode, msg)
	}

	// Check all with bad state
	si.desiredName = goodName
	si.desiredState = badState
	msg, retcode = si.ProcessInfo()
	if retcode != 2 {
		t.Errorf("ProcessInfo() failed on bad state with retcode %d, msg %s", retcode, msg)
	}

	// Check all with bad user
	si.desiredState = goodState
	si.desiredUser = badUser
	msg, retcode = si.ProcessInfo()
	if retcode != 2 {
		t.Errorf("ProcessInfo() failed on bad user with retcode %d, msg %s", retcode, msg)
	}

	// Check good state only
	si.desiredUser = ""
	msg, retcode = si.ProcessInfo()
	if retcode != 0 {
		t.Errorf("ProcessInfo() failed on blank user with retcode %d, msg %s", retcode, msg)
	}

	// Check bad state only
	si.desiredState = badState
	msg, retcode = si.ProcessInfo()
	if retcode != 2 {
		t.Errorf("ProcessInfo() failed on bad state with retcode %d, msg %s", retcode, msg)
	}

	// Check good user only
	si.desiredUser = goodUser
	si.desiredState = ""
	msg, retcode = si.ProcessInfo()
	if retcode != 0 {
		t.Errorf("ProcessInfo() failed on blank state with retcode %d, msg %s", retcode, msg)
	}

	// Check bad user only
	si.desiredUser = badUser
	msg, retcode = si.ProcessInfo()
	if retcode != 2 {
		t.Errorf("ProcessInfo() failed on bad user with retcode %d, msg %s", retcode, msg)
	}

	// Get service info only
	si.desiredState = ""
	si.desiredUser = ""
	msg, retcode = si.ProcessInfo()
	if retcode != 0 {
		t.Errorf("ProcessInfo() failed returning service info only retcode %d, msg %s", retcode, msg)
	}

	si.currentStateWanted = true
	msg, retcode = si.ProcessInfo()
	if retcode != 0 {
		t.Errorf("ProcessInfo() failed when fetching current state")
	}

	si.desiredName = badName
	msg, retcode = si.ProcessInfo()
	if retcode != 0 {
		t.Errorf("ProcessInfo() failed when fetching current state for unknown service")
	}

	si.getServiceInfo = nil
	err = si.GetInfo()
	if err == nil {
		t.Errorf("GetInfo() returned no error but had a nil handler")
	}
}
