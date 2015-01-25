package google

import (
	"net/http"

	_http "github.com/st3v/translator/http"
)

type authenticator struct {
	apiKey string
}

// NewAuthenticator instantiates a new Authenticator for Google's Translate API.
func newAuthenticator(apiKey string) _http.Authenticator {
	return &authenticator{
		apiKey: apiKey,
	}
}

func (a *authenticator) Authenticate(request *http.Request) error {
	params := request.URL.Query()
	params.Add("key", a.apiKey)
	request.URL.RawQuery = params.Encode()
	return nil
}
