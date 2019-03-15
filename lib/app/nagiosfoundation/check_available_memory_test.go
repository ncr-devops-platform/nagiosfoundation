package nagiosfoundation

import (
	"flag"
	"os"
	"testing"
)

func TestCheckAvailableMemory(t *testing.T) {
	pgmName := "TestCheckAvailableMemory"
	testReturnValid := func() uint64 { return uint64(50) }
	testReturnZero := func() uint64 { return uint64(0) }

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
	msg, retcode = CheckAvailableMemoryWithHandler(testReturnZero)

	if retcode != 2 || msg == "" {
		t.Error("CheckAvailableMemoryWithHandler() failed with valid GetFreeMemory() call")
	}

	os.Args = []string{pgmName, "-warning", "40"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	msg, retcode = CheckAvailableMemoryWithHandler(testReturnValid)

	if retcode != 1 || msg == "" {
		t.Error("CheckAvailableMemoryWithHandler() should have emitted WARNING")
	}

	os.Args = []string{pgmName, "-critical", "45"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	msg, retcode = CheckAvailableMemoryWithHandler(testReturnValid)

	if retcode != 2 || msg == "" {
		t.Error("CheckAvailableMemoryWithHandler() should have emitted CRITICAL")
	}

	os.Args = savedArgs
	flag.CommandLine = savedFlagCommandLine
}
