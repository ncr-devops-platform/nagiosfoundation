package nagiosfoundation

import (
	"fmt"
	"testing"
)

const testProcessGoodName = "goodName"
const testProcessBadName = "badName"

type testProcessHandler struct{}

func (p testProcessHandler) IsProcessRunning(name string) bool {
	retval := false

	if name == testProcessGoodName {
		retval = true
	}

	return retval
}
func TestCheckProcess(t *testing.T) {
	fmt.Println("TestCheckProcess()")

	pc := ProcessCheck{
		ProcessName:         testProcessGoodName,
		ProcessCheckHandler: new(testProcessHandler),
	}

	if pc.IsProcessRunning() != true {
		t.Error("isProcessRunning() failed")
	}

	pc.ProcessName = testProcessBadName

	if pc.IsProcessRunning() != false {
		t.Error("isProcessRunning() failed")
	}

	var retcode int
	// Running check with running process
	_, retcode = CheckProcessWithService(testProcessGoodName, "running", new(testProcessHandler))
	if retcode != 0 {
		t.Errorf("Running check with running process failed with retcode %d", retcode)
	}

	// Not running check with running process
	_, retcode = CheckProcessWithService(testProcessGoodName, "notrunning", new(testProcessHandler))
	if retcode != 2 {
		t.Errorf("Not running check with running process failed with retcode %d", retcode)
	}

	// Running check with not running process
	_, retcode = CheckProcessWithService(testProcessBadName, "running", new(testProcessHandler))
	if retcode != 2 {
		t.Errorf("Running check with not running process failed with retcode %d", retcode)
	}

	// Not running check with not running process
	_, retcode = CheckProcessWithService(testProcessBadName, "notrunning", new(testProcessHandler))
	if retcode != 0 {
		t.Errorf("Not running check with not running process failed with retcode %d", retcode)
	}

	// Invalid check type
	_, retcode = CheckProcessWithService(testProcessGoodName, "", new(testProcessHandler))
	if retcode != 3 {
		t.Errorf("Invalid check type not detected with retcode %d", retcode)
	}

	showHelp()
}
