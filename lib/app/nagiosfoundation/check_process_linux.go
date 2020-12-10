// +build !windows

package nagiosfoundation

import "github.com/ncr-devops-platform/nagiosfoundation/lib/pkg/process"

func isProcessRunningOsConstrained(name string) bool {
	retVal := false

	if processEntries, _ := process.GetProcessesByName(name); len(processEntries) > 0 {
		retVal = true
	}

	return retVal
}
