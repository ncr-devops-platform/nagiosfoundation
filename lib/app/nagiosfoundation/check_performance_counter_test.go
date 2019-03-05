package nagiosfoundation

import (
	"errors"
	"flag"
	"os"
	"testing"

	"github.com/ncr-devops-platform/nagiosfoundation/lib/pkg/perfcounters"
)

func TestCheckPerformanceCounter(t *testing.T) {
	pgmName := "TestCheckPerformanceCounter"
	testReturnValid := func(string, int, int) (perfcounters.PerformanceCounter, error) {
		return perfcounters.PerformanceCounter{Value: 5.0}, nil
	}
	testReturnError := func(string, int, int) (perfcounters.PerformanceCounter, error) {
		return perfcounters.PerformanceCounter{Value: 5.0}, errors.New("GetPerformanceCounter() failure")
	}
	// Save args and flagset for restoration
	savedArgs := os.Args
	savedFlagCommandLine := flag.CommandLine

	// No "ReadPerformanceCounter" service passed
	os.Args = []string{pgmName}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	msg, retcode := CheckPerformanceCounterWithHandler(nil)

	if retcode != 2 || msg == "" {
		t.Error("CheckPerformanceCounterWithHandler() failed to handle nil service")
	}

	// Valid ReadPerformanceCounter service with valid returns from service
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	msg, retcode = CheckPerformanceCounterWithHandler(testReturnValid)

	if retcode != 0 || msg == "" {
		t.Error("CheckPerformanceCounterWithHandler() failed with valid returns from service")
	}

	// Valid memory service with flag defaults
	os.Args = []string{pgmName, "-greater_than", "-critical", "10", "-warning", "10"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	msg, retcode = CheckPerformanceCounterWithHandler(testReturnValid)

	if retcode != 0 || msg == "" {
		t.Error("CheckPerformanceCounterWithHandler() failed with with -greater_than set")
	}

	// Valid ReadPerformanceCounter service with error returned from service
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	msg, retcode = CheckPerformanceCounterWithHandler(testReturnError)

	if retcode != 2 || msg == "" {
		t.Error("CheckPerformanceCounterWithHandler() failed with error return from service")
	}

	os.Args = savedArgs
	flag.CommandLine = savedFlagCommandLine
}
