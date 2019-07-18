package nagiosfoundation

import (
	"fmt"
	"testing"
)

func TestResultMessage(t *testing.T) {
	checkName := "TestName"
	statusText := "OK"
	description := "Description"
	nagiosOutput := "nagios output"

	// All parameters populated, valid result
	expectedValue := fmt.Sprintf("%s %s - %s | %s", checkName, statusText, description, nagiosOutput)
	actualValue, err := resultMessage(checkName, statusText, description, nagiosOutput)
	if err != nil {
		t.Error("resultMessage() with all parameters should have returned valid result. Err:", err)
	}
	if actualValue != expectedValue {
		t.Errorf("resultMessage() with all parameters did not return expected message. Expected: %s, Actual: %s", expectedValue, actualValue)
	}

	// Nagios output not populated, valid result
	expectedValue = fmt.Sprintf("%s %s - %s", checkName, statusText, description)
	actualValue, err = resultMessage(checkName, statusText, description)
	if err != nil {
		t.Error("resultMessage() with 3 parameters should have returned valid result. Err:", err)
	}
	if actualValue != expectedValue {
		t.Errorf("resultMessage() with 3 parameters did not return expected message. Expected: %s, Actual: %s", expectedValue, actualValue)
	}

	// Only required parameters populated, valid result
	expectedValue = fmt.Sprintf("%s %s", checkName, statusText)
	actualValue, err = resultMessage(checkName, statusText)
	if err != nil {
		t.Error("resultMessage() with all parameters should have returned valid result. Err:", err)
	}
	if actualValue != expectedValue {
		t.Errorf("resultMessage() with all parameters did not return expected message. Expected: %s, Actual: %s", expectedValue, actualValue)
	}

	// Invalid status text
	_, err = resultMessage("TestName", "INVALIDTEXTRESULT")
	if err != errResultMsgInvalidStatus {
		t.Error("resultMessage() with invalid status text should have returned errResultMsgInvalidStatus")
	}

	// Only one parameter
	_, err = resultMessage("TestName")
	if err != errResultMsgNotEnoughArgs {
		t.Error("resultMessage() with one parameter should have returned errResultMsgNotEnoughArgs")
	}

	// No parameters
	_, err = resultMessage()
	if err != errResultMsgNotEnoughArgs {
		t.Error("resultMessage() with no parameters should have returned errResultMsgNotEnoughArgs")
	}

	// Too many parameters
	_, err = resultMessage("", "", "", "", "")
	if err != errResultMsgTooManyArgs {
		t.Error("resultMessage() with five parameters should have returned errResultMsgTooManyArgs")
	}

	// Test remaining status text values
	for _, statusText := range []string{statusTextCritical, statusTextWarning, statusTextUnknown} {
		_, err = resultMessage("TestName", statusText)
		if err != nil {
			t.Errorf("resultMessage() with valid status text of %s failed", statusText)
		}
	}
}
