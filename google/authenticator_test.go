package google

import (
	"fmt"
	"net/http"
	"testing"
)

func TestAuthenticate(t *testing.T) {
	apiKey := "my-secret-api-key"
	authenticator := newAuthenticator(apiKey)

	host := "localhost:6666"
	params := "one=1&two=2"
	originalURL := fmt.Sprintf("%s?%s", host, params)

	request, err := http.NewRequest("GET", originalURL, nil)
	if err != nil {
		t.Errorf("Unexpected error creating new request: %s", err.Error())
	}

	if request.URL.String() != originalURL {
		t.Errorf(
			"Unexpected request URL prior to authentication. Got: %s. Want: %s.",
			request.URL.String(),
			originalURL,
		)
	}

	err = authenticator.Authenticate(request)
	if err != nil {
		t.Errorf("Unexpected error authenticating request: %s", err.Error())
	}

	newURL := fmt.Sprintf("%s?key=%s&%s", host, apiKey, params)
	if request.URL.String() != newURL {
		t.Errorf(
			"Unexpected request URL prior to authentication. Got: %s. Want: %s.",
			request.URL.String(),
			newURL,
		)
	}

}
