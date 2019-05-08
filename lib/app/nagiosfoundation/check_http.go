package nagiosfoundation

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// CheckHTTP attempts an HTTP request against the provided url, reporting the HTTP response code and overall request state.
func CheckHTTP(url string, redirect bool, timeout int) (string, int) {

	var msg string
	var retCode int
	var responseStateText string
	var responseCode string

	status, err := statusCode(url, timeout)
	if err != nil {
		retCode = 2
		responseStateText = "CRITICAL"
		responseCode = "UNHANDLED ERROR"
	}

	switch {
	case status >= 400:
		retCode = 2
		responseStateText = "CRITICAL"
		responseCode = strconv.Itoa(status)
	case status >= 300 && redirect:
		retCode = 0
		responseStateText = "OK"
		responseCode = strconv.Itoa(status)
	case status >= 300:
		retCode = 1
		responseStateText = "WARNING"
		responseCode = strconv.Itoa(status)
	case status == -1:
		retCode = 2
		responseStateText = "UNKNOWN ERROR"
		responseCode = strconv.Itoa(status)
	default:
		retCode = 0
		responseStateText = "OK"
		responseCode = strconv.Itoa(status)
	}

	msg = fmt.Sprintf("CheckHttp %s - Url %s responded with %s", responseStateText, url, responseCode)

	return msg, retCode
}

func statusCode(url string, timeout int) (int, error) {
	http.DefaultClient.Timeout = time.Duration(timeout) * time.Second

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1, err
	}

	response, err := http.DefaultTransport.RoundTrip(request)
	if err != nil {
		return -1, err
	}

	return response.StatusCode, nil
}
