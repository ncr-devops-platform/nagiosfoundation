package perfcounters

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"testing"
)

const testPerfCounterValue = "4.2"

type testPowerShellService struct {
	testType int
}

func (pss testPowerShellService) Execute(args ...string) (string, string, error) {
	var stdOut string
	var stdErr string
	var err error

	switch pss.testType {
	case 1:
		stdOut = ""
		stdErr = ""
		err = nil
	case 2:
		stdOut = testPerfCounterValue
		stdErr = ""
		err = nil
	case 3:
		stdOut = testPerfCounterValue
		stdErr = ""
		err = errors.New("testPowerShellService Error")
	}

	return stdOut, stdErr, err
}

func TestReadPerformanceCounter(t *testing.T) {
	counterName := "TestCounter"
	type pss struct {
	}

	svc := new(testPowerShellService)

	var pcResult PerformanceCounter
	var err error
	var nul *os.File

	stdErr := os.Stderr

	if nul, err = os.OpenFile(os.DevNull, os.O_WRONLY, 0444); err != nil {
		nul = os.Stderr
	}

	// Passing no handler should yield an error
	_, err = ReadPerformanceCounterWithHandler(nil, counterName, 0, 0)

	if err == nil {
		t.Error("nil Powershell service failed to yield an error")
	}

	// Not returning a number should yield an error
	svc.testType = 1
	os.Stderr = nul
	_, err = ReadPerformanceCounterWithHandler(svc, counterName, 0, 0)
	os.Stderr = stdErr
	if err == nil {
		t.Error("Powershell service returning empty string did not yield an error")
	}

	// Returning valid data should return good Powershell struct
	svc.testType = 2
	pcResult, err = ReadPerformanceCounterWithHandler(svc, counterName, 0, 0)

	if err != nil {
		t.Error("Powershell service returned valid data but yielded an error")
	}

	value, _ := strconv.ParseFloat(testPerfCounterValue, 64)
	if pcResult.Value != value {
		t.Error("PerformanceCounter value not correct on valid data from Powershell service")
	}

	if pcResult.Name != counterName {
		t.Error("PerformanceCounter name was not properly populated")
	}

	// Returning an error should bubble back up
	svc.testType = 3
	os.Stderr = nul
	_, err = ReadPerformanceCounterWithHandler(svc, counterName, 0, 0)
	os.Stderr = stdErr
	if err == nil {
		t.Error("Powershell service returned an error but ReadPerformanceCounterWithHandler did not pass it up")
	}

	nul.Close()
}

// Lifted from exec_test.go
func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

const powerShellRunResult = "great success"

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	fmt.Fprintf(os.Stdout, powerShellRunResult)
	fmt.Fprintf(os.Stderr, powerShellRunResult)
	os.Exit(0)
}

func TestPowerShell(t *testing.T) {
	path := "/path/to/"

	ps := newPowerShell(func(file string) (string, error) { return path + file, nil },
		func(string, ...string) *exec.Cmd { return fakeExecCommand("") })

	if ps.powerShell != path+"powershell.exe" {
		t.Error("PowerShell did not properly populate lookup string")
	}

	if ps.command == nil {
		t.Error("PowerShell command was not initialized")
	}

	stdout, stderr, err := ps.Execute("")

	if stdout != powerShellRunResult {
		t.Error("Stdout wasn't populated with the proper run result")
	}

	if stderr != powerShellRunResult {
		t.Error("Stderr wasn't populated with the proper run result")
	}

	if err != nil {
		t.Error("PowerShell execute should have not returned an error")
	}
}
