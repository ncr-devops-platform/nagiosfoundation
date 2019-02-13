package nagiosfoundation

import (
	"flag"
	"os"
	"testing"
)

func TestSetFlagIfNotProvided(t *testing.T) {
	// Save args and flagset for restoration
	savedArgs := os.Args
	savedFlagCommandLine := flag.CommandLine

	// Reset the default flag set
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	pgmName := "TestSetFlagIfNotProvided"
	flagName := "testflagname"
	flagDefault := "defaultValue"
	flagOverride := "overrideValue"
	flagNotProvided := "notProvidedValue"
	flagPtr := flag.String(flagName, flagDefault, "flag for testing")

	// Flag not provided on the command line therefore it is
	// set internally.
	os.Args = []string{pgmName}
	SetFlagIfNotProvided(flagName, flagNotProvided)
	flag.Parse()
	if *flagPtr != flagNotProvided {
		t.Error("Flag not provided and was not set")
	}

	// Flag provided on the command line therefore do not
	// set it internally.
	os.Args = []string{pgmName, "-" + flagName, flagOverride}
	SetFlagIfNotProvided(flagName, flagNotProvided)
	flag.Parse()
	if *flagPtr != flagOverride {
		t.Error("Flag was provided but was also set")
	}

	os.Args = savedArgs
	flag.CommandLine = savedFlagCommandLine
}
