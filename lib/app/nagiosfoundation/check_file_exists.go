package nagiosfoundation

import (
	"fmt"
	"path/filepath"
	"strconv"
)

// CheckFileExists tests the assertion that one or more files matching specified pattern should or should not exist.
func CheckFileExists(pattern string, negate bool) (string, int) {
	var msg string
	var retCode int
	var checkStateText, msgString string

	matches, err := filepath.Glob(pattern)

	switch {
	case err != nil:
		checkStateText = statusTextUnknown
		msgString = fmt.Sprintf("Error matching pattern %s: %s", pattern, err)
		retCode = 3
	default:
		matchCount := len(matches)

		switch {
		case (matchCount == 0 && negate == false) ||
			(matchCount > 0 && negate == true):
			checkStateText = statusTextCritical
			retCode = 2
		case (matchCount == 0 && negate == true) ||
			(matchCount > 0 && negate == false):
			checkStateText = statusTextOK
			retCode = 0
		}

		msgString = fmt.Sprintf("%s files matched pattern %s", strconv.Itoa(len(matches)), pattern)
	}

	msg, _ = resultMessage("CheckFileExists", checkStateText, msgString)
	return msg, retCode
}
