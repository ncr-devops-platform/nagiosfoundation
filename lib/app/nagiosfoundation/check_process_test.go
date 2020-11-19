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
	_, retcode = checkProcessWithService(testProcessGoodName, "running", "metric", new(testProcessHandler))
	if retcode != statusCodeOK {
		t.Errorf("Running check with running process failed with retcode %d", retcode)
	}

	// Not running check with running process
	_, retcode = checkProcessWithService(testProcessGoodName, "notrunning", "metric", new(testProcessHandler))
	if retcode != statusCodeCritical {
		t.Errorf("Not running check with running process failed with retcode %d", retcode)
	}

	// Running check with not running process
	_, retcode = checkProcessWithService(testProcessBadName, "running", "metric", new(testProcessHandler))
	if retcode != statusCodeCritical {
		t.Errorf("Running check with not running process failed with retcode %d", retcode)
	}

	// Not running check with not running process
	_, retcode = checkProcessWithService(testProcessBadName, "notrunning", "metric", new(testProcessHandler))
	if retcode != statusCodeOK {
		t.Errorf("Not running check with not running process failed with retcode %d", retcode)
	}

	// Invalid check type
	_, retcode = checkProcessWithService(testProcessGoodName, "", "metric", new(testProcessHandler))
	if retcode != statusCodeCritical {
		t.Errorf("Invalid check type not detected with retcode %d", retcode)
	}

	testMsg := "Test Message"
	testCheckProcess := func(name, checkType, metricName string, processService ProcessService) (string, int) {
		return testMsg, statusCodeOK
	}

	_, retcode = checkProcessCmd("dummyprocess", "running", "metric", testCheckProcess, new(testProcessHandler))

	if retcode != statusCodeOK {
		t.Error("valid check process test should have returned OK")
	}

	_, retcode = checkProcessCmd("", "dummytype", "metric", testCheckProcess, new(testProcessHandler))

	if retcode != statusCodeCritical {
		t.Error("check process with no -name should return CRITICAL")
	}

	_, retcode = checkProcessCmd("", "", "metric", testCheckProcess, new(testProcessHandler))

	if retcode != statusCodeCritical {
		t.Error("check process test with no parameters should have returned CRITICAL")
	}

	_, retcode = checkProcessCmd("dummyprocess", "badtype", "metric", testCheckProcess, new(testProcessHandler))

	if retcode != statusCodeCritical {
		t.Error("check process test with invalid type should have returned CRITICAL")
	}
}
