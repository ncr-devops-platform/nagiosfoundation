package nagiosfoundation

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func validateCheckPortMsg(t *testing.T, msg string) {
	t.Helper()

	if msg == "" {
		t.Error("CheckPort() returned empty message")
	}
}

func TestCheckPort(t *testing.T) {
	var msg string
	var ret int

	httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer httpServer.Close()
	url := strings.Split(strings.TrimPrefix(httpServer.URL, "http://"), ":")
	address := url[0]
	port, _ := strconv.Atoi(url[1])

	// happy path
	msg, ret = CheckPort(CheckPortProtocolTCP, address, port, 30, false, "metric")
	if ret != statusCodeOK {
		t.Error("CheckPort() should return OK code when given valid machine and port")
	}

	validateCheckPortMsg(t, msg)

	// inverted happy path
	msg, ret = CheckPort(CheckPortProtocolTCP, address, port, 30, true, "metric")
	if ret != statusCodeCritical {
		t.Error("CheckPort() should return critical code when given valid machine and port and is inverted")
	}

	validateCheckPortMsg(t, msg)

	// invalid protocol
	msg, ret = CheckPort(42, address, port, 30, false, "metric")
	if ret != statusCodeCritical {
		t.Error("CheckPort() should return critical code when given invalid protocol")
	}

	validateCheckPortMsg(t, msg)

	// invalid metric name
	msg, ret = CheckPort(CheckPortProtocolTCP, address, port, 30, false, "")
	if ret != statusCodeCritical {
		t.Error("CheckPort() should return critical code when given invalid metric name")
	}

	validateCheckPortMsg(t, msg)

	// invalid machine name
	msg, ret = CheckPort(CheckPortProtocolTCP, "", port, 30, false, "metric")
	if ret != statusCodeCritical {
		t.Error("CheckPort() should return critical code when given invalid machine name")
	}

	validateCheckPortMsg(t, msg)

	// no listener
	httpServer.Close()
	msg, ret = CheckPort(CheckPortProtocolTCP, address, port, 30, false, "metric")
	if ret != statusCodeCritical {
		t.Error("CheckPort() should return critical code when there is no listener")
	}

	validateCheckPortMsg(t, msg)

	// no listener with invert on
	httpServer.Close()
	msg, ret = CheckPort(CheckPortProtocolTCP, address, port, 30, true, "metric")
	if ret != statusCodeOK {
		t.Error("CheckPort() should return OK code when there is no listener and is inverted")
	}

	validateCheckPortMsg(t, msg)
}
