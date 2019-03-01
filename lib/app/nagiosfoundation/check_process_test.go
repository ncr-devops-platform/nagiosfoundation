package nagiosfoundation

import (
	"flag"
	"fmt"
	"os"
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

	// Test flags
	// Save args and flagset for restoration
	savedArgs := os.Args
	savedFlagCommandLine := flag.CommandLine
	pgmName := "TestCheckProcess"
	testMsg := "Test Message"
	testCheckProcess := func(name string, checkType string, processService ProcessService) (string, int) {
		return testMsg, 0
	}

	// Reset the default flag set
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{pgmName, "-name", "dummyprocess"}
	_, retcode = CheckProcessFlags(testCheckProcess, new(testProcessHandler))

	if retcode != 0 {
		t.Error("valid check process test should have returned OK")
	}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{pgmName}
	_, retcode = CheckProcessFlags(testCheckProcess, new(testProcessHandler))

	if retcode != 2 {
		t.Error("check process test with no parameters should have returned CRITICAL")
	}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{pgmName, "-name", "dummyprocess", "-type", "badtype"}
	_, retcode = CheckProcessFlags(testCheckProcess, new(testProcessHandler))

	if retcode != 2 {
		t.Error("check process test with invalid type should have returned CRITICAL")
	}

	os.Args = savedArgs
	flag.CommandLine = savedFlagCommandLine

	showHelp()
}
