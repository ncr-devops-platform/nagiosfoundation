// +build !windows

package nagiosfoundation

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func getPidName(pid int) (string, error) {
	procFile := fmt.Sprintf("/proc/%d/stat", pid)
	procDataBytes, err := ioutil.ReadFile(procFile)
	if err != nil {
		return "", err
	}

	procData := string(procDataBytes)

	procNameStart := strings.IndexRune(procData, '(') + 1
	procNameEnd := strings.IndexRune(procData, ')')

	if procNameStart >= procNameEnd {
		return "", errors.New("Could not parse process name")
	}

	procName := procData[procNameStart:procNameEnd]

	return procName, nil
}

func getProcessesByName(name string) ([]os.FileInfo, error) {
	var errorReturn error
	matchingEntries := make([]os.FileInfo, 0)

	dir, err := os.Open("/proc")
	if err != nil {
		matchingEntries = nil
		errorReturn = err
	}

	defer dir.Close()

	var procEntries []os.FileInfo
	if errorReturn == nil {
		procEntries, err = dir.Readdir(0)

		if err != nil {
			matchingEntries = nil
			errorReturn = err
		}
	}

	if errorReturn == nil {
		for _, procEntry := range procEntries {
			// Skip entries that aren't directories
			if !procEntry.IsDir() {
				continue
			}

			// Skip entries that aren't numbers
			pid, err := strconv.Atoi(procEntry.Name())
			if err != nil {
				continue
			}

			if procName, _ := getPidName(pid); procName == name {
				matchingEntries = append(matchingEntries, procEntry)
			}
		}
	}

	return matchingEntries, errorReturn
}

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
