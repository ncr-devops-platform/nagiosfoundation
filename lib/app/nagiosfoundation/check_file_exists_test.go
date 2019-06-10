package nagiosfoundation

import (
	"os"
	"strings"
	"testing"
)

func validateTestCheckFileExistsResponse(t *testing.T, description string, expectedCode int, expectedMsg string, actualCode int, actualMsg string) {
	t.Helper()

	if expectedCode != actualCode {
		t.Errorf("%s: Expected Code: %d, Actual Code: %d", description, expectedCode, actualCode)
	}

	if strings.Contains(actualMsg, expectedMsg) == false {
		t.Errorf("%s: Expected Message: %s, Actual Message: %s", description, expectedMsg, actualMsg)
	}
}

func TestCheckFileExists(t *testing.T) {
	var msg string
	var code int
	const validTestFile = "validtestfile"
	const invalidTestFile = "invalidtestfile"
	const ok = "OK:"
	const critical = "CRITICAL:"
	const unknown = "UNKNOWN:"

	type testItem struct {
		description  string
		file         string
		inverted     bool
		expectedCode int
		expectedMsg  string
	}

	testList := []testItem{
		{
			description:  "Invalid Glob Pattern",
			file:         "[]a",
			inverted:     false,
			expectedCode: 3,
			expectedMsg:  unknown,
		},
		{
			description:  "File does not exist, not inverted",
			file:         invalidTestFile,
			inverted:     false,
			expectedCode: 2,
			expectedMsg:  critical,
		},
		{
			description:  "File does not exist, inverted",
			file:         invalidTestFile,
			inverted:     true,
			expectedCode: 0,
			expectedMsg:  ok,
		},
		{
			description:  "File exists, not inverted",
			file:         validTestFile,
			inverted:     false,
			expectedCode: 0,
			expectedMsg:  ok,
		},
		{
			description:  "File exists, inverted",
			file:         validTestFile,
			inverted:     true,
			expectedCode: 2,
			expectedMsg:  critical,
		},
	}

	// Create a valid file
	if fp, err := os.OpenFile(validTestFile, os.O_RDONLY|os.O_CREATE, 0666); err != nil {
		t.Errorf("Error creating test file: %s. Error: %s", validTestFile, err)
	} else {
		defer os.Remove(validTestFile)
		fp.Close()
	}

	for _, i := range testList {
		msg, code = CheckFileExists(i.file, i.inverted)
		validateTestCheckFileExistsResponse(t, i.description, i.expectedCode, i.expectedMsg, code, msg)
	}
}
