package nagiosformatters

import (
	"strings"
	"testing"
)

func TestNagiosFormatters(t *testing.T) {
	const testCheck = "TestCheck"
	const testMetric = "Test Metric"
	const testOk = " OK"
	const testWarning = " WARNING"
	const testCritical = " CRITICAL"

	var msg string
	var retval int

	// Greater Check OK
	msg, retval = GreaterFormatNagiosCheck(testCheck, 50, 80, 90, testMetric)

	if retval != 0 {
		t.Error("Greater check OK failed with retval", retval)
	}

	if !strings.HasPrefix(msg, testCheck+testOk) {
		t.Error("Greater check OK failed with msg", msg)
	}

	// Greater Check WARNING
	msg, retval = GreaterFormatNagiosCheck(testCheck, 85, 80, 90, testMetric)

	if retval != 1 {
		t.Error("Greater check WARNING failed")
	}

	if !strings.HasPrefix(msg, testCheck+testWarning) {
		t.Error("Greater check WARNING failed with msg", msg)
	}

	// Greater Check CRITICAL
	msg, retval = GreaterFormatNagiosCheck(testCheck, 95, 80, 90, testMetric)

	if retval != 2 {
		t.Error("Greater check CRITICAL failed")
	}

	if !strings.HasPrefix(msg, testCheck+testCritical) {
		t.Error("Greater check CRITICAL failed with msg", msg)
	}

	// Lesser Check OK
	msg, retval = LesserFormatNagiosCheck(testCheck, 95, 90, 80, testMetric)
	if retval != 0 {
		t.Error("Lesser check OK failed")
	}

	if !strings.HasPrefix(msg, testCheck+testOk) {
		t.Error("Lesser check OK failed with msg", msg)
	}

	// Lesser Check WARNING
	msg, retval = LesserFormatNagiosCheck(testCheck, 85, 90, 80, testMetric)
	if retval != 1 {
		t.Error("Lesser check WARNING failed")
	}

	if !strings.HasPrefix(msg, testCheck+testWarning) {
		t.Error("Lesser check WARNING failed with msg", msg)
	}

	// Lesser Check CRITICAL
	msg, retval = LesserFormatNagiosCheck(testCheck, 50, 90, 80, testMetric)
	if retval != 2 {
		t.Error("Lesser check CRITICAL failed")
	}

	if !strings.HasPrefix(msg, testCheck+testCritical) {
		t.Error("Lesser check CRITICAL failed with msg", msg)
	}
}
