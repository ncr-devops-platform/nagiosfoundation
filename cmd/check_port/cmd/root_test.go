package cmd

import (
	"os"
	"testing"

	nf "github.com/ncr-devops-platform/nagiosfoundation/lib/app/nagiosfoundation"
)

func TestCheckPortCmd(t *testing.T) {
	var actualExitCode, expectedExitCode int
	checkName := "check_port"

	apiCheckPort := func(protocol nf.CheckPortProtocol, address string, port, timeout int, invert bool, metricName string) (string, int) {
		return "Test Message", 0
	}

	savedArgs := os.Args

	os.Args = []string{checkName, "--address", "localhost", "--port", "80"}
	expectedExitCode = 0
	actualExitCode = Execute(apiCheckPort)
	if actualExitCode != expectedExitCode {
		t.Errorf("%s: Expected Code: %d, Actual Code: %d", "check_port_1", expectedExitCode, actualExitCode)
	}

	os.Args = []string{checkName, "invalidcommand"}
	expectedExitCode = 1
	actualExitCode = Execute(apiCheckPort)
	if actualExitCode != expectedExitCode {
		t.Errorf("%s: Expected Code: %d, Actual Code: %d", "check_port_2", expectedExitCode, actualExitCode)
	}

	os.Args = savedArgs
}
