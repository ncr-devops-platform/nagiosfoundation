package nagiosfoundation

import (
	"errors"
	"flag"
	"os"
	"testing"
)

func TestCheckCpu(t *testing.T) {
	pgmName := "TestCheckCpu"
	testReturnValid := func() (float64, error) { return 0.5, nil }
	testReturnError := func() (float64, error) { return 0.5, errors.New("GetCPULoad() failure") }

	// Save args and flagset for restoration
	savedArgs := os.Args
	savedFlagCommandLine := flag.CommandLine

	// No "get memory" service passed
	os.Args = []string{pgmName}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	msg, retcode := CheckCPUWithHandler(nil)

	if retcode != 2 || msg == "" {
		t.Error("CheckCPUWithHandler() failed to handle nil service")
	}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	msg, retcode = CheckCPUWithHandler(testReturnValid)

	if retcode != 0 || msg == "" {
		t.Error("CheckCPUWithHandler() failed with valid returns from service")
	}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	msg, retcode = CheckCPUWithHandler(testReturnError)

	if retcode != 2 || msg == "" {
		t.Error("CheckCPUWithHandler() failed with error returned from service")
	}

	os.Args = savedArgs
	flag.CommandLine = savedFlagCommandLine
}
