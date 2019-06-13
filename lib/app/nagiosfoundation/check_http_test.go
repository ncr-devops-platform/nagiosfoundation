package nagiosfoundation

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckHTTP(t *testing.T) {
	var httpStatus int

	httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(httpStatus)
	}))
	defer httpServer.Close()
	var format = ""
	var path = ""
	var expectedValue = ""
	// Code 200
	httpStatus = http.StatusOK
	msg, code := CheckHTTP(httpServer.URL, false, 1, format, path, expectedValue)
	fmt.Println(msg)
	fmt.Println(code)
	if code != 0 {
		t.Error("CheckHTTP() should return code of 0 when on OK (200) response")
	}

	// Code 400
	httpStatus = http.StatusBadRequest
	msg, code = CheckHTTP(httpServer.URL, false, 1, format, path, expectedValue)
	fmt.Println(msg)
	fmt.Println(code)
	if code != 2 {
		t.Error("CheckHTTP() should return code of 2 when on bad request (400) response")
	}

	// Code 300
	httpStatus = http.StatusMultipleChoices
	msg, code = CheckHTTP(httpServer.URL, false, 1, format, path, expectedValue)
	fmt.Println(msg)
	fmt.Println(code)
	if code != 1 {
		t.Error("CheckHTTP() should return code of 2 when on multiple choices (300) response")
	}

	// Code 300 with redirect on
	httpStatus = http.StatusMultipleChoices
	msg, code = CheckHTTP(httpServer.URL, true, 1, format, path, expectedValue)
	fmt.Println(msg)
	fmt.Println(code)

	// Shut down test server to generate errors
	httpServer.Close()

	// No server for connection
	httpStatus = http.StatusOK
	msg, code = CheckHTTP(httpServer.URL, false, 1, format, path, expectedValue)
	fmt.Println(msg)
	fmt.Println(code)
	if code != 2 {
		t.Error("CheckHTTP() should return code of 2 when no server is available")
	}

	// Invalid URL
	httpStatus = http.StatusOK
	msg, code = CheckHTTP("invalid%url", false, 1, format, path, expectedValue)
	fmt.Println(msg)
	fmt.Println(code)
	if code != 2 {
		t.Error("CheckHTTP() should return code of 2 when given an unparseable URL")
	}
}
