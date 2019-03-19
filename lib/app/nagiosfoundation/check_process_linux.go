// +build !windows

package nagiosfoundation

func getHelpOsConstrained() string {
	return "\nNote: Process names in POSIX systems are case sensitive."
}

func isProcessRunningOsConstrained(name string) bool {
	retVal := false

	if processEntries, _ := getProcessesByName(name); len(processEntries) > 0 {
		retVal = true
	}

	return retVal
}
