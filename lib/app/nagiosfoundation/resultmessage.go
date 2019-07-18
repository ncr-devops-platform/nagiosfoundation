package nagiosfoundation

import (
	"errors"
	"fmt"
)

const (
	statusTextOK       = "OK"
	statusTextWarning  = "WARNING"
	statusTextCritical = "CRITICAL"
	statusTextUnknown  = "UNKNOWN"
)

var errResultMsgNotEnoughArgs = errors.New("Not enough arguments")
var errResultMsgTooManyArgs = errors.New("Too many arguments")
var errResultMsgInvalidStatus = errors.New("Invalid status text")

func resultMessage(s ...string) (string, error) {
	// s[0] - check name
	// s[1] - status text
	// s[2] - result description
	// s[3] - nagios output

	const (
		checkNameOffset = iota
		statusTextOffset
		resultDescOffset
		nagiosOutputOffset
	)

	var msg string
	var err error

	argCount := len(s)

	if argCount < 2 {
		err = errResultMsgNotEnoughArgs
	} else if argCount > 4 {
		err = errResultMsgTooManyArgs
	} else if s[statusTextOffset] != statusTextOK && s[statusTextOffset] != statusTextWarning &&
		s[statusTextOffset] != statusTextCritical && s[statusTextOffset] != statusTextUnknown {
		err = errResultMsgInvalidStatus
	} else {
		// "CheckName OK"
		msg = fmt.Sprintf("%s %s", s[checkNameOffset], s[statusTextOffset])

		// "CheckName OK - Description of result"
		if argCount > resultDescOffset && len(s[resultDescOffset]) > 0 {
			msg += fmt.Sprintf(" - %s", s[resultDescOffset])
		}

		// "CheckName OK - Description of result | nagios output"
		if argCount > nagiosOutputOffset && len(s[nagiosOutputOffset]) > 0 {
			msg += fmt.Sprintf(" | %s", s[nagiosOutputOffset])
		}
	}

	return msg, err
}
