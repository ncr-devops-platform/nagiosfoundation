// +build !windows

package nagiosfoundation

import (
	"fmt"
)

func showHelpOsConstrained() {
	fmt.Println("\nNote: Process names in POSIX systems are case sensitive.")
}

func isProcessRunningOsConstrained(name string) bool {
	retVal := false

	if processEntries, _ := getProcessesByName(name); len(processEntries) > 0 {
		retVal = true
	}

	return retVal
}
