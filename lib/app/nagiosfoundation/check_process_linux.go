// +build !windows

package nagiosfoundation

func isProcessRunningOsConstrained(name string) bool {
	retVal := false

	if processEntries, _ := getProcessesByName(name); len(processEntries) > 0 {
		retVal = true
	}

	return retVal
}
