package nagiosfoundation

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

// The command name and version are injected into
// these variables at build time.
// See godel/config/dist-plugin.yml
var cmdName string
var cmdVersion string

// SetFlagIfNotProvided sets a command line flag if it wasn't
// provided. This overcomes a command line flag in library
// being set to a different default value than desired.
//
// Returns true if the command line flag was not provided
// and therefore was set to the value provided.
func SetFlagIfNotProvided(flagName string, flagValue string) bool {
	flagSet := false
	flag.Parse()

	flag.Visit(func(f *flag.Flag) {
		if (*f).Name == flagName {
			flagSet = true
		}
	})

	if !flagSet {
		flag.Set(flagName, flagValue)
	}

	return !flagSet
}

// SetDefaultGlogStderr will prevent the glog package from
// defaulting to creating new log files on every execution,
// set the logtostderr option for glog to true if it wasn't
// specified on the command line.
func SetDefaultGlogStderr() {
	SetFlagIfNotProvided("logtostderr", "true")
}

// GetVersion returns the executable version as a string
//
func GetVersion() string {
	const unknown = "<unknown>"

	if cmdName == "" {
		cmdName = unknown
	}

	if cmdVersion == "" {
		cmdVersion = unknown
	}

	version := cmdName + " version " + cmdVersion + " " + runtime.GOOS + "/" + runtime.GOARCH

	return version
}

// ShowVersion checks for "version" to be the only argument
// and if true, writes the version to the io.Writer passed
// in.
//
// Returns true if the version was output, else false
func ShowVersion(w io.Writer) bool {
	retval := false

	if len(os.Args) == 2 && strings.EqualFold(os.Args[1], "version") {
		if w == nil {
			w = os.Stdout
		}

		fmt.Fprintln(w, GetVersion())
		retval = true
	}

	return retval
}

// CheckExecutableVersion will attempt to show the version of
// the executable and if shown, will immediately exit.
func CheckExecutableVersion() {
	if ShowVersion(os.Stdout) {
		os.Exit(0)
	}
}
