package nagiosfoundation

import (
	"fmt"
	"net"
	"time"
)

// CheckPortProtocol is the type of protocol to use for checking the port
type CheckPortProtocol int

const (
	// CheckPortProtocolTCP specifies using TCP as the network protocol for
	// checking theport
	CheckPortProtocolTCP CheckPortProtocol = 0
)

func (cpp CheckPortProtocol) String() string {
	s := "unknown"

	switch cpp {
	case CheckPortProtocolTCP:
		s = "tcp"
	}

	return s
}

func checkPort(protocol CheckPortProtocol, address string, port, timeout int, invert bool) (string, int, string) {
	resultText := statusTextCritical
	resultCode := statusCodeCritical
	resultDesc := ""

	d := net.Dialer{Timeout: time.Duration(timeout) * time.Second}
	conn, err := d.Dial(protocol.String(), fmt.Sprintf("%s:%d", address, port))

	if err != nil {
		resultDesc = err.Error()

		if invert == true {
			resultText = statusTextOK
			resultCode = statusCodeOK
		}
	} else {
		conn.Close()

		if invert == false {
			resultText = statusTextOK
			resultCode = statusCodeOK
		}
	}

	return resultText, resultCode, resultDesc
}

// CheckPort checks for a listening port at an address
func CheckPort(protocol CheckPortProtocol, address string, port, timeout int, invert bool, metricName string) (string, int) {
	const checkName = "CheckPort"
	var retCode int
	var msg, desc, nagiosOutput, retText, retDesc string

	if protocol.String() == "unknown" {
		desc = fmt.Sprintf("Unknown protocol value of %d", protocol)
		retText = statusTextCritical
		retCode = statusCodeCritical
	} else if metricName == "" {
		desc = "Empty metric name"
		retText = statusTextCritical
		retCode = statusCodeCritical
	} else if address == "" {
		desc = "Empty address"
		retText = statusTextCritical
		retCode = statusCodeCritical
	} else {
		retText, retCode, retDesc = checkPort(protocol, address, port, timeout, invert)
		desc = fmt.Sprintf("port %d on %s using %s", port, address, protocol)
		if retDesc != "" {
			desc = desc + fmt.Sprintf(" (%s)", retDesc)
		}
	}

	nagiosOutput = fmt.Sprintf("%s=%d", metricName, retCode)
	msg, _ = resultMessage(checkName, retText, desc, nagiosOutput)

	return msg, retCode
}
