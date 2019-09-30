package cmd

import (
	"os"
	"testing"
)

func TestCheckFileExistsCmd(t *testing.T) {
	type testItem struct {
		description string
		arguments   []string
		exitCode    int
	}

	var expectedExitCode int

	testList := []testItem{
		{
			description: "Invalid commmand",
			arguments:   []string{"check_file_exists", "invalidcommand"},
			exitCode:    1,
		},
		{
			description: "Valid command, exit code 0",
			arguments:   []string{"check_file_exists"},
			exitCode:    0,
		},
		{
			description: "Valid command, exit code 2",
			arguments:   []string{"check_file_exists"},
			exitCode:    2,
		},
	}

	apiCheckFileExists := func(pattern string, negate bool) (string, int) {
		return "Test Message", expectedExitCode
	}

	savedArgs := os.Args

	for _, i := range testList {
		os.Args = i.arguments
		expectedExitCode = i.exitCode
		actualExitCode := Execute(apiCheckFileExists)
		if actualExitCode != i.exitCode {
			t.Errorf("%s: Expected Code: %d, Actual Code: %d", i.description, i.exitCode, actualExitCode)
		}
	}

	os.Args = savedArgs
}
