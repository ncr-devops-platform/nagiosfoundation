package nagiosfoundation

import (
	"flag"
)

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