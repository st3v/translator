package google

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestResponseParserAPIError(t *testing.T) {
	errorPayload := `
    {
			"error": {
	      "errors": [
	        {
	          "domain": "error-domain",
	          "reason": "error-reason",
	          "message": "error-specific-message"
	        }
	      ],
	      "code": 666,
	      "message": "error-generic-message"
	    }
		}
  `
	response := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(strings.NewReader(errorPayload)),
	}

	result, err := parseResponse(response, &struct{}{})

	if result != nil {
		t.Errorf("Expected nil result but got: %#v", result)
	}

	expectedError := "API Error. Code: 666, Message: error-generic-message, Domain: error-domain, Reason: error-reason"

	if !strings.HasPrefix(err.Error(), expectedError) {
		t.Errorf("Unexpected Error. Got: '%s'. Want: '%s'.", err.Error(), expectedError)
	}
}
