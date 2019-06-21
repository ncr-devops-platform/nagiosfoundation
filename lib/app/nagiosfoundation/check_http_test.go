package nagiosfoundation

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestCheckHTTP(t *testing.T) {
	var httpStatus int
	idValue := 1337
	idValueString := strconv.Itoa(idValue)

	responseBody := `{"id":"` + idValueString + `"}`
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
	_, code := CheckHTTP(httpServer.URL, false, 1, format, path, expectedValue, "")
	if code != 0 {
		t.Error("CheckHTTP() should return code of 0 when on OK (200) response")
	}

	// Code 400
	httpStatus = http.StatusBadRequest
	_, code = CheckHTTP(httpServer.URL, false, 1, format, path, expectedValue, "")
	if code != 2 {
		t.Error("CheckHTTP() should return code of 2 when on bad request (400) response")
	}

	// Code 300
	httpStatus = http.StatusMultipleChoices
	_, code = CheckHTTP(httpServer.URL, false, 1, format, path, expectedValue, "")
	if code != 1 {
		t.Error("CheckHTTP() should return code of 2 when on multiple choices (300) response")
	}

	// Code 300 with redirect on
	httpStatus = http.StatusMultipleChoices
	_, code = CheckHTTP(httpServer.URL, true, 1, format, path, expectedValue, "")
	if code != 0 {
		t.Error("CheckHTTP() should return code of 0 when on multiple choices (300) with redirect response")
	}

	// Code 200 with format json and a match on expected value
	httpStatus = http.StatusOK
	_, code = CheckHTTP(httpServer.URL, false, 1, "json", "id", idValueString, "")
	if code != 0 {
		t.Error("CheckHTTP() should return code of 0 when json path matches expected value")
	}

	// Code 200 with format json and failed match on expected value
	httpStatus = http.StatusOK
	_, code = CheckHTTP(httpServer.URL, false, 1, "json", "id", "failmatch", "")
	if code != 2 {
		t.Error("CheckHTTP() should return code of 2 when json path does not match expected value")
	}

	// Code 200 with format json and expression true
	httpStatus = http.StatusOK
	_, code = CheckHTTP(httpServer.URL, false, 1, "json", "id", "", "!= \""+idValueString+"\"")
	if code != 2 {
		t.Error("CheckHTTP() should return code of 2 when json path causes expression to return false")
	}

	// Code 200 with format json and no expected value or expression
	httpStatus = http.StatusOK
	_, code = CheckHTTP(httpServer.URL, false, 1, "json", "id", "", "")
	if code != 2 {
		t.Error("CheckHTTP() should return code of 2 with json path but no expected value or expression")
	}

	// Code 200 with format json and but both expected value and expression given
	httpStatus = http.StatusOK
	_, code = CheckHTTP(httpServer.URL, false, 1, "json", "id", "expectedvalue", "expression")
	if code != 2 {
		t.Error("CheckHTTP() should return code of 2 with json path but no expected value or expression")
	}

	// Code 200 with format json and expression true using integer
	responseBody = `{"id":` + idValueString + `}`
	httpStatus = http.StatusOK
	_, code = CheckHTTP(httpServer.URL, false, 1, "json", "id", "", "== "+idValueString)
	if code != 0 {
		t.Errorf("CheckHTTP() should return code of 0 with json path but and comparison to int %d", idValue)
	}

	// Invalid format
	responseBody = `{"id":` + idValueString + `}`
	httpStatus = http.StatusOK
	_, code = CheckHTTP(httpServer.URL, false, 1, "invalidformat", "id", "", "== "+idValueString)
	if code != 2 {
		t.Errorf("CheckHTTP() should return code of 2 when given an invalid format")
	}

	// Invalid path
	responseBody = `{"id":` + idValueString + `}`
	httpStatus = http.StatusOK
	_, code = CheckHTTP(httpServer.URL, false, 1, "json", "invalidpath", "", "== "+idValueString)
	if code != 2 {
		t.Errorf("CheckHTTP() should return code of 2 when given an invalid path")
	}

	// Shut down test server to generate errors
	httpServer.Close()

	// No server for connection
	httpStatus = http.StatusOK
	_, code = CheckHTTP(httpServer.URL, false, 1, format, path, expectedValue, "")
	if code != 2 {
		t.Error("CheckHTTP() should return code of 2 when no server is available")
	}

	// Invalid URL
	httpStatus = http.StatusOK
	_, code = CheckHTTP("invalid%url", false, 1, format, path, expectedValue, "")
	if code != 2 {
		t.Error("CheckHTTP() should return code of 2 when given an unparseable URL")
	}
}

func TestAcceptText(t *testing.T) {
	// Format is an empty string
	expectedValue := ""
	actualValue, err := getAcceptText("")
	if actualValue != expectedValue {
		t.Errorf("getAcceptText() should return empty string on empty format but returned \"%s\"", actualValue)
	}

	// Format is an invalid format
	actualValue, err = getAcceptText("invalidformat")
	if err == nil {
		t.Errorf("getAcceptText() should return error on invalid format")
	}

	if actualValue != expectedValue {
		t.Errorf("getAcceptText() should return empty string on invalid format but returned \"%s\"", actualValue)
	}

	// Format is json
	expectedValue = "application/json"
	actualValue, err = getAcceptText("json")
	if actualValue != expectedValue {
		t.Errorf("getAcceptText() should return \"%s\" on json format but returned \"%s\"", expectedValue, actualValue)
	}
}

func TestEvaluateStatusCode(t *testing.T) {
	// http.StatusBadRequest
	expectedCode := 2
	expectedText := "CRITICAL"
	actualCode, actualText := evaluateStatusCode(http.StatusBadRequest, false)
	if actualCode != expectedCode {
		t.Errorf("evaluateStatusCode() with http.StatusBadRequest expected code of %d but returned %d", expectedCode, actualCode)
	}

	if actualText != expectedText {
		t.Errorf("evaluateStatusCode() with http.StatusBadRequest expected text of %s but returned %s", expectedText, actualText)
	}

	// http.StatusMultipleChoices with redirect false
	expectedCode = 1
	expectedText = "WARNING"
	actualCode, actualText = evaluateStatusCode(http.StatusMultipleChoices, false)
	if actualCode != expectedCode {
		t.Errorf("evaluateStatusCode() with http.StatusMultipleChoices and redirect false expected code of %d but returned %d", expectedCode, actualCode)
	}

	if actualText != expectedText {
		t.Errorf("evaluateStatusCode() with http.StatusMultipleChoices and redirect false expected text of %s but returned %s", expectedText, actualText)
	}

	// http.StatusMultipleChoices with redirect true
	expectedCode = 0
	expectedText = "OK"
	actualCode, actualText = evaluateStatusCode(http.StatusMultipleChoices, true)
	if actualCode != expectedCode {
		t.Errorf("evaluateStatusCode() with http.StatusMultipleChoices and redirect true expected code of %d but returned %d", expectedCode, actualCode)
	}

	if actualText != expectedText {
		t.Errorf("evaluateStatusCode() with http.StatusMultipleChoices and redirect true expected text of %s but returned %s", expectedText, actualText)
	}

	// http.StatusOK
	expectedCode = 0
	expectedText = "OK"
	actualCode, actualText = evaluateStatusCode(http.StatusOK, false)
	if actualCode != expectedCode {
		t.Errorf("evaluateStatusCode() with http.StatusOK expected code of %d but returned %d", expectedCode, actualCode)
	}

	if actualText != expectedText {
		t.Errorf("evaluateStatusCode() with http.StatusOK expected text of %s but returned %s", expectedText, actualText)
	}

	// Unknown (-1)
	expectedCode = 2
	expectedText = "UNKNOWN ERROR"
	actualCode, actualText = evaluateStatusCode(-1, false)
	if actualCode != expectedCode {
		t.Errorf("evaluateStatusCode() with -1 expected code of %d but returned %d", expectedCode, actualCode)
	}

	if actualText != expectedText {
		t.Errorf("evaluateStatusCode() with -1 expected text of %s but returned %s", expectedText, actualText)
	}
}

func TestEvaluateExpectedValue(t *testing.T) {
	actualValue := "testvalue"
	expectedGoodValue := actualValue
	expectedBadValue := "nomatch"
	testPath := "testpath"

	expectedCode := 0
	expectedText := "OK"
	actualCode, actualText, _ := evaluateExpectedValue(actualValue, expectedGoodValue, testPath)
	if actualCode != expectedCode {
		t.Errorf("evaluateExpectedValue() returned actual code of %d when expecting %d", actualCode, expectedCode)
	}

	if actualText != expectedText {
		t.Errorf("evaluateExpectedValue() returned actual text of %s when expecting %s", actualText, expectedText)
	}

	expectedCode = 2
	expectedText = "CRITICAL"
	actualCode, actualText, _ = evaluateExpectedValue(actualValue, expectedBadValue, testPath)
	if actualCode != expectedCode {
		t.Errorf("evaluateExpectedValue() returned actual code of %d when expecting %d", actualCode, expectedCode)
	}

	if actualText != expectedText {
		t.Errorf("evaluateExpectedValue() returned actual text of %s when expecting %s", actualText, expectedText)
	}
}

func TestEvaluateExpression(t *testing.T) {
	actualValue := "testvalue"
	trueExpression := "== \"" + actualValue + "\""
	falseExpression := "!= \"" + actualValue + "\""
	errorExpression := "+- \"" + actualValue + "\""
	testPath := "testpath"

	expectedCode := 0
	expectedText := "OK"
	actualCode, actualText, _ := evaluateExpression(actualValue, trueExpression, testPath)
	if actualCode != expectedCode {
		t.Errorf("evaluateExpression() returned actual code of %d when expecting %d", actualCode, expectedCode)
	}

	if actualText != expectedText {
		t.Errorf("evaluateExpression() returned actual text of %s when expecting %s", actualText, expectedText)
	}

	expectedCode = 2
	expectedText = "CRITICAL"
	actualCode, actualText, _ = evaluateExpression(actualValue, falseExpression, testPath)
	if actualCode != expectedCode {
		t.Errorf("evaluateExpression() returned actual code of %d when expecting %d", actualCode, expectedCode)
	}

	if actualText != expectedText {
		t.Errorf("evaluateExpression() returned actual text of %s when expecting %s", actualText, expectedText)
	}

	expectedCode = 2
	expectedText = "CRITICAL"
	actualCode, actualText, _ = evaluateExpression(actualValue, falseExpression, testPath)
	if actualCode != expectedCode {
		t.Errorf("evaluateExpression() returned actual code of %d when expecting %d", actualCode, expectedCode)
	}

	if actualText != expectedText {
		t.Errorf("evaluateExpression() returned actual text of %s when expecting %s", actualText, expectedText)
	}

	expectedCode = 2
	expectedText = "CRITICAL"
	actualCode, actualText, _ = evaluateExpression(actualValue, errorExpression, testPath)
	if actualCode != expectedCode {
		t.Errorf("evaluateExpression() returned actual code of %d when expecting %d", actualCode, expectedCode)
	}

	if actualText != expectedText {
		t.Errorf("evaluateExpression() returned actual text of %s when expecting %s", actualText, expectedText)
	}
}
