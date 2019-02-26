package nagiosformatters

import (
	"strings"
	"testing"
)

func TestNagiosFormatters(t *testing.T) {
	const testCheck = "TestCheck"
	const testMetric = "Test Metric"
	const testOk = "OK"
	const testWarning = "WARNING"
	const testCritical = "CRITICAL"

	type formatTests struct {
		name           string
		method         func(string, float64, float64, float64, string) (string, int)
		expectedValue  int
		expectedString string
		value          float64
		warning        float64
		critical       float64
	}

	testList := []formatTests{
		{
			name:           "Greater",
			method:         GreaterFormatNagiosCheck,
			expectedValue:  0,
			expectedString: testOk,
			value:          50,
			warning:        80,
			critical:       90,
		},
		{
			name:           "Greater",
			method:         GreaterFormatNagiosCheck,
			expectedValue:  1,
			expectedString: testWarning,
			value:          85,
			warning:        80,
			critical:       90,
		},
		{
			name:           "Greater",
			method:         GreaterFormatNagiosCheck,
			expectedValue:  2,
			expectedString: testCritical,
			value:          95,
			warning:        80,
			critical:       90,
		},
		{
			name:           "Lesser",
			method:         LesserFormatNagiosCheck,
			expectedValue:  0,
			expectedString: testOk,
			value:          95,
			warning:        90,
			critical:       80,
		},
		{
			name:           "Lesser",
			method:         LesserFormatNagiosCheck,
			expectedValue:  1,
			expectedString: testWarning,
			value:          85,
			warning:        90,
			critical:       80,
		},
		{
			name:           "Lesser",
			method:         LesserFormatNagiosCheck,
			expectedValue:  2,
			expectedString: testCritical,
			value:          50,
			warning:        90,
			critical:       80,
		},
	}

	var msg string
	var retval int

	for _, test := range testList {
		msg, retval = test.method(testCheck, test.value, test.warning, test.critical, testMetric)

		if retval != test.expectedValue {
			t.Errorf("%s check %s failed with retval %d", test.name, test.expectedString, retval)
		}

		if !strings.HasPrefix(msg, testCheck+" "+test.expectedString) {
			t.Errorf("%s check %s failed with msg %s", test.name, test.expectedString, msg)
		}
	}
}
