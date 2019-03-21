package nagiosfoundation

import (
	"flag"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/spf13/cobra"
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

func TestVersion(t *testing.T) {
	pgmName := "TestCmd"
	savedArgs := os.Args

	if ShowVersion(nil) != false {
		t.Error("Version should not have been shown since \"version\" was not an argument.")
	}

	os.Args = []string{pgmName, "version"}
	cmdName = pgmName
	cmdVersion = "TestVersion"

	if ShowVersion(nil) == false {
		t.Error("Version should have been shown since \"version\" was an argument.")
	}

	var s strings.Builder
	if ShowVersion(&s) == false {
		t.Error("Version should have been shown since \"version\" was an argument.")
	}

	expectedResult := cmdName + " version " + cmdVersion + " " + runtime.GOOS + "/" + runtime.GOARCH + "\n"
	if s.String() != expectedResult {
		t.Errorf("Version string returned is not correct. Expected result: %s Actual Result: %s",
			expectedResult,
			s.String())
	}

	cmdName = ""
	cmdVersion = ""
	s.Reset()
	ShowVersion(&s)
	expectedResult = "<unknown> version <unknown> " + runtime.GOOS + "/" + runtime.GOARCH + "\n"
	if s.String() != expectedResult {
		t.Errorf("Version string returned is not correct. Expected result: %s Actual Result: %s",
			expectedResult,
			s.String())
	}

	testCmd := &cobra.Command{}
	AddVersionCommand(testCmd)
	cmdList := testCmd.Commands()
	cmdList[0].Execute()
	if cmdList[0].Use != "version" {
		t.Error("version command did not load into Cobra")
	}

	os.Args = savedArgs
}
