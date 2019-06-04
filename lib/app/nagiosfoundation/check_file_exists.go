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
	var checkStateText string

	matches, err := filepath.Glob(pattern)
	if err != nil {
		msg = fmt.Sprintf("UNKNOWN - Error matching pattern %s: %s", pattern, err)
		retCode = 3
		return msg, retCode
	}

	switch {
	case len(matches) == 0 && negate == false:
		checkStateText = "CRITICAL"
		retCode = 2
	case len(matches) == 0 && negate == true:
		checkStateText = "OK"
		retCode = 0
	case len(matches) > 0 && negate == false:
		checkStateText = "OK"
		retCode = 0
	case len(matches) > 0 && negate == true:
		checkStateText = "CRITICAL"
		retCode = 2
	}

	msg = fmt.Sprintf("%s: %s files matched pattern %s", checkStateText, strconv.Itoa(len(matches)), pattern)

	return msg, retCode
}
