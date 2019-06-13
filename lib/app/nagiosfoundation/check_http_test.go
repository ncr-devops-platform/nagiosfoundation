package nagiosfoundation

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckHTTP(t *testing.T) {
	var httpStatus int
	var responseBody string

	httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(httpStatus)
		w.Write([]byte(responseBody))
	}))
	defer httpServer.Close()
	var format = ""
	var path = ""
	var expectedValue = ""

	// Code 200
	httpStatus = http.StatusOK
	_, code := CheckHTTP(httpServer.URL, false, 1, format, path, expectedValue)
	if code != 0 {
		t.Error("CheckHTTP() should return code of 0 when on OK (200) response")
	}

	// Code 400
	httpStatus = http.StatusBadRequest
	_, code = CheckHTTP(httpServer.URL, false, 1, format, path, expectedValue)
	if code != 2 {
		t.Error("CheckHTTP() should return code of 2 when on bad request (400) response")
	}

	// Code 300
	httpStatus = http.StatusMultipleChoices
	_, code = CheckHTTP(httpServer.URL, false, 1, format, path, expectedValue)
	if code != 1 {
		t.Error("CheckHTTP() should return code of 2 when on multiple choices (300) response")
	}

	// Code 300 with redirect on
	httpStatus = http.StatusMultipleChoices
	_, code = CheckHTTP(httpServer.URL, true, 1, format, path, expectedValue)
	if code != 0 {
		t.Error("CheckHTTP() should return code of 0 when on multiple choices (300) with redirect response")
	}

	// Code 200 with format json and a match
	httpStatus = http.StatusOK
	responseBody = `{"id":"HeaFdiyIJe","joke":"What kind of magic do cows believe in? MOODOO.","status":200}`
	_, code = CheckHTTP(httpServer.URL, false, 1, "json", "id", "HeaFdiyIJe")
	if code != 0 {
		t.Error("CheckHTTP() should return code of 0 when json path matches expected value")
	}

	// Code 200 with format json and failed match
	httpStatus = http.StatusOK
	responseBody = `{"id":"HeaFdiyIJe","joke":"What kind of magic do cows believe in? MOODOO.","status":200}`
	_, code = CheckHTTP(httpServer.URL, false, 1, "json", "id", "failmatch")
	if code != 2 {
		t.Error("CheckHTTP() should return code of 2 when json path matches expected value")
	}

	// Shut down test server to generate errors
	httpServer.Close()

	// No server for connection
	httpStatus = http.StatusOK
	_, code = CheckHTTP(httpServer.URL, false, 1, format, path, expectedValue)
	if code != 2 {
		t.Error("CheckHTTP() should return code of 2 when no server is available")
	}

	// Invalid URL
	httpStatus = http.StatusOK
	_, code = CheckHTTP("invalid%url", false, 1, format, path, expectedValue)
	if code != 2 {
		t.Error("CheckHTTP() should return code of 2 when given an unparseable URL")
	}
}
