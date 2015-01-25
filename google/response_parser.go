package google

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/st3v/tracerr"
)

type errorResult struct {
	Error struct {
		Errors []struct {
			Domain  string
			Reason  string
			Message string
		}
		Code    int
		Message string
	}
}

var parseResponse = func(resp *http.Response, target interface{}) (interface{}, error) {
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	errorResult := &errorResult{}
	err = json.Unmarshal(body, errorResult)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	if errorResult.Error.Code != 0 {
		return nil, tracerr.Errorf(
			"API Error. Code: %d, Message: %s, Domain: %s, Reason: %s",
			errorResult.Error.Code,
			errorResult.Error.Message,
			errorResult.Error.Errors[0].Domain,
			errorResult.Error.Errors[0].Reason,
		)
	}

	err = json.Unmarshal(body, target)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	return target, nil
}
