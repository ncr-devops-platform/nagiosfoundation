package nagiosfoundation

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/thedevsaddam/gojsonq"
)

// CheckHTTP attempts an HTTP request against the provided url, reporting the HTTP response code and overall request state.
func CheckHTTP(url string, redirect bool, timeout int, format string, path string, expectedValue string) (string, int) {
	var accept string
	var msg string
	var retCode int
	var responseStateText string
	var responseCode string

	switch {
	case format == "json":
		accept = "application/json"
	default:
		accept = ""
	}

	status, body, _ := statusCode(url, timeout, accept)

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

	var checkMsg = ""
	if retCode == 0 &&
		len(expectedValue) > 0 &&
		len(format) > 0 &&
		len(path) > 0 {

		switch {
		case format == "json":
			queryValue := fmt.Sprintf("%v", gojsonq.New().JSONString(body).Find(path))
			if queryValue == expectedValue {
				checkMsg = fmt.Sprintf(". The value found at %s has expected value %s", path, expectedValue)
			} else {
				retCode = 2
				responseStateText = "CRITICAL"
				checkMsg = fmt.Sprintf(". The value found at %s has unexpected value %s", path, queryValue)
			}
		}
	}

	msg = fmt.Sprintf("CheckHttp %s - Url %s responded with %s%s", responseStateText, url, responseCode, checkMsg)
	strconv.Itoa(status)

	return msg, retCode
}

func statusCode(url string, timeout int, accept string) (int, string, error) {
	if !isURL(url) {
		return -1, "", fmt.Errorf("%s is not a valid url", url)
	}

	http.DefaultClient.Timeout = time.Duration(timeout) * time.Second

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, "", err
	}

	request.Header.Set("accept", accept)

	response, err := http.DefaultTransport.RoundTrip(request)
	if err != nil {
		return -1, "", err
	}

	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return -1, "", readErr
	}
	return response.StatusCode, string(body), nil
}

func isURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
