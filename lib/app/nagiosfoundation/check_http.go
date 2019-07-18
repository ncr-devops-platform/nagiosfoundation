package nagiosfoundation

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/PaesslerAG/gval"
	"github.com/thedevsaddam/gojsonq"
)

func getAcceptText(format string) (string, error) {
	var accept string
	var err error

	switch format {
	case "json":
		accept = "application/json"
	case "":
		accept = ""
	default:
		accept = ""
		err = errors.New("Invalid accept type")
	}

	return accept, err
}

func evaluateStatusCode(status int, redirect bool) (int, string) {
	var retCode int
	var responseStateText string

	switch {
	case status >= http.StatusBadRequest:
		retCode = 2
		responseStateText = "CRITICAL"
	case status >= http.StatusMultipleChoices && redirect:
		retCode = 0
		responseStateText = "OK"
	case status >= http.StatusMultipleChoices:
		retCode = 1
		responseStateText = "WARNING"
	case status == -1:
		retCode = 2
		responseStateText = "UNKNOWN ERROR"
	default:
		retCode = 0
		responseStateText = "OK"
	}

	return retCode, responseStateText
}

func evaluateExpectedValue(actualValue, expectedValue, path string) (int, string, string) {
	var retCode int
	var responseStateText, checkMsg string

	if actualValue == expectedValue {
		retCode = 0
		responseStateText = "OK"
		checkMsg = fmt.Sprintf(". The value found at %s has expected value %s", path, expectedValue)
	} else {
		retCode = 2
		responseStateText = "CRITICAL"
		checkMsg = fmt.Sprintf(". The value found at %s has unexpected value %s", path, actualValue)
	}

	return retCode, responseStateText, checkMsg
}

func evaluateExpression(actualValue interface{}, expression, path string) (int, string, string) {
	var retCode int
	var responseStateText, checkMsg string

	evalResult, err := gval.Evaluate("value "+expression, map[string]interface{}{"value": actualValue})
	if err == nil {
		if evalResult == true {
			retCode = 0
			responseStateText = "OK"
			checkMsg = fmt.Sprintf(". The value found at %s with value %v and expression \"%s\" yields true", path, actualValue, expression)
		} else {
			retCode = 2
			responseStateText = "CRITICAL"
			checkMsg = fmt.Sprintf(". The value found at %s with value %v does not match expression \"%s\"", path, actualValue, expression)
		}
	} else {
		retCode = 2
		responseStateText = "CRITICAL"
		checkMsg = fmt.Sprintf(". Error processing value found at %s with value %v using expression \"%s\": %s", path, actualValue, expression, err)
	}

	return retCode, responseStateText, checkMsg
}

// CheckHTTP attempts an HTTP request against the provided url, reporting the HTTP response code and overall request state.
func CheckHTTP(url string, redirect bool, timeout int, format, path, expectedValue, expression string) (string, int) {
	var retCode int
	var msg string

	acceptText, err := getAcceptText(format)
	if err != nil {
		return fmt.Sprintf("CheckHttp CRITICAL - The format (--format) \"%s\" is not valid. The only valid value is \"json\".", format),
			2
	}

	status, body, _ := statusCode(url, timeout, acceptText)

	retCode, responseStateText := evaluateStatusCode(status, redirect)
	responseCode := strconv.Itoa(status)

	var checkMsg = ""
	if retCode == 0 && len(format) > 0 && len(path) > 0 {
		var queryValue string

		switch {
		case format == "json":
			expectedValueLen := len(expectedValue)
			expressionLen := len(expression)

			value := gojsonq.New().JSONString(body).Find(path)

			if value == nil {
				retCode = 2
				responseStateText = "CRITICAL"
				checkMsg = fmt.Sprintf(". No entry at path %s", path)
			} else if expectedValueLen > 0 && expressionLen > 0 {
				retCode = 2
				responseStateText = "CRITICAL"
				checkMsg = fmt.Sprintf(". Both --expectedValue and --expression given but only one is used")
			} else if expectedValueLen > 0 {
				queryValue = fmt.Sprintf("%v", value)
				retCode, responseStateText, checkMsg = evaluateExpectedValue(queryValue, expectedValue, path)
			} else if expressionLen > 0 {
				retCode, responseStateText, checkMsg = evaluateExpression(value, expression, path)
			} else {
				retCode = 2
				responseStateText = "CRITICAL"
				checkMsg = fmt.Sprintf(". --expectedValue or --expression not given")
			}
		}
	}

	msg = fmt.Sprintf("CheckHttp %s - Url %s responded with %s%s", responseStateText, url, responseCode, checkMsg)

	return msg, retCode
}

func statusCode(url string, timeout int, accept string) (int, string, error) {
	http.DefaultClient.Timeout = time.Duration(timeout) * time.Second

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1, "", err
	}

	request.Header.Set("accept", accept)

	response, err := http.DefaultTransport.RoundTrip(request)
	if err != nil {
		return -1, "", err
	}
	defer response.Body.Close()

	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return -1, "", readErr
	}
	return response.StatusCode, string(body), nil
}
