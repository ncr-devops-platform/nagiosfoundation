package nagiosfoundation

import (
	"errors"
	"flag"
	"os"
	"testing"
)

func TestCheckAvailableMemory(t *testing.T) {
	pgmName := "TestCheckAvailableMemory"
	testReturnValid := func() (float64, error) { return 1000, nil }
	testReturnError := func() (float64, error) { return 1000, errors.New("GetFreeMemory() failure") }

	// Save args and flagset for restoration
	savedArgs := os.Args
	savedFlagCommandLine := flag.CommandLine

	// No "get memory" service passed
	os.Args = []string{pgmName}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	msg, retcode := CheckAvailableMemoryWithHandler(nil)

	if retcode != 2 || msg == "" {
		t.Error("CheckAvailableMemoryWithHandler() failed to handle nil service")
	}

	// Valid memory service with flag defaults
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	msg, retcode = CheckAvailableMemoryWithHandler(testReturnValid)
	if retcode != 0 || msg == "" {
		t.Error("CheckAvailableMemoryWithHandler() failed with valid GetFreeMemory() call")
	}

	// Valid memory service but service returns error
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	msg, retcode = CheckAvailableMemoryWithHandler(testReturnError)
	if retcode != 2 || msg == "" {
		t.Error("CheckAvailableMemoryWithHandler() failed with valid GetFreeMemory() call")
	}

	os.Args = savedArgs
	flag.CommandLine = savedFlagCommandLine
}
