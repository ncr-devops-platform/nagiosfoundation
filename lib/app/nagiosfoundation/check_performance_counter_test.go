package nagiosfoundation

import (
	"errors"
	"testing"

	"github.com/ncr-devops-platform/nagiosfoundation/lib/pkg/perfcounters"
)

func TestCheckPerformanceCounter(t *testing.T) {
	testReturnValid := func(string, int, int) (perfcounters.PerformanceCounter, error) {
		return perfcounters.PerformanceCounter{Value: 5.0}, nil
	}
	testReturnError := func(string, int, int) (perfcounters.PerformanceCounter, error) {
		return perfcounters.PerformanceCounter{Value: 5.0}, errors.New("GetPerformanceCounter() failure")
	}

	// No "ReadPerformanceCounter" service passed
	msg, retcode := CheckPerformanceCounterWithHandler(0, 0, false, 2, 1, "test metric name", "test counter name", nil)

	if retcode != 2 || msg == "" {
		t.Error("CheckPerformanceCounterWithHandler() failed to handle nil service")
	}

	// Valid ReadPerformanceCounter service with valid returns from service
	msg, retcode = CheckPerformanceCounterWithHandler(0, 0, false, 2, 1, "test metric name", "test counter name", testReturnValid)

	if retcode != 0 || msg == "" {
		t.Error("CheckPerformanceCounterWithHandler() failed with valid returns from service")
	}

	// Valid memory service with flag defaults
	msg, retcode = CheckPerformanceCounterWithHandler(10, 10, true, 2, 1, "test metric name", "test counter name", testReturnValid)

	if retcode != 0 || msg == "" {
		t.Error("CheckPerformanceCounterWithHandler() failed with with -greater_than set")
	}

	// Valid ReadPerformanceCounter service with error returned from service
	msg, retcode = CheckPerformanceCounterWithHandler(10, 10, true, 2, 1, "test metric name", "test counter name", testReturnError)

	if retcode != 2 || msg == "" {
		t.Error("CheckPerformanceCounterWithHandler() failed with error return from service")
	}

	GetHelpPerformanceCounter()
}
