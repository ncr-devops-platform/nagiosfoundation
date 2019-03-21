package nagiosfoundation

import (
	"errors"
	"testing"
)

func TestCheckCpu(t *testing.T) {
	testReturnValid := func() (float64, error) { return 0.5, nil }
	testReturnError := func() (float64, error) { return 0.5, errors.New("GetCPULoad() failure") }

	// No "get memory" service passed
	msg, retcode := CheckCPUWithHandler(85, 95, "", nil)

	if retcode != 2 || msg == "" {
		t.Error("CheckCPUWithHandler() failed to handle nil service")
	}

	msg, retcode = CheckCPUWithHandler(85, 95, "pct_processor_time", testReturnValid)

	if retcode != 0 || msg == "" {
		t.Error("CheckCPUWithHandler() failed with valid returns from service")
	}

	msg, retcode = CheckCPUWithHandler(85, 95, "pct_processor_time", testReturnError)

	if retcode != 2 || msg == "" {
		t.Error("CheckCPUWithHandler() failed with error returned from service")
	}
}
