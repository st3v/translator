package microsoft

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHTTPClientSendRequest(t *testing.T) {
	expectedRequestMethod := "GET"
	expectedRequestContentType := "fake-content-type"
	expectedRequestAuthToken := "fake-auth-token"
	expectedRequestBody := "fake-request-body"

	expectedResponseBody := "fake-response-body"
	expectedResponseHeader := map[string]string{
		"fake-header-key-1": "fake-header-value-1",
		"fake-header-key-2": "fake-header-value-2",
		"fake-header-key-3": "fake-header-value-3",
	}

	authenticator := newMockAuthenticator(func(request *http.Request) error {
		request.Header.Set("Authorization", expectedRequestAuthToken)
		return nil
	})

	httpClient := newHTTPClient(authenticator)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != expectedRequestMethod {
			t.Errorf("Unexpected request method. Want: %s. Got: %s", expectedRequestMethod, r.Method)
		}

		if r.Header.Get("Content-Type") != expectedRequestContentType {
			t.Errorf(
				"Unexpected content type in request header. Want: '%s'. Got: '%s'",
				expectedRequestContentType,
				r.Header.Get("Content-Type"),
			)
		}

		if r.Header.Get("Authorization") != expectedRequestAuthToken {
			t.Errorf(
				"Unexpected auth token in request header. Want: '%s'. Got: '%s'",
				expectedRequestAuthToken,
				r.Header.Get("Authorization"),
			)
		}

		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			t.Errorf("Unexpected error reading request body: %s", err.Error())
		}

		if string(body) != expectedRequestBody {
			t.Errorf("Unexpected request body. Want: '%s'. Got: '%s'", expectedRequestBody, string(body))
		}

		for key, value := range expectedResponseHeader {
			w.Header().Set(key, value)
		}

		fmt.Fprint(w, expectedResponseBody)
		return
	}))
	defer server.Close()

	response, err := httpClient.SendRequest(
		expectedRequestMethod,
		server.URL,
		strings.NewReader(expectedRequestBody),
		expectedRequestContentType,
	)

	if err != nil {
		t.Fatalf("Unexpected error when sending request: %s", err.Error())
	}

	for key, value := range expectedResponseHeader {
		if response.Header.Get(key) != value {
			t.Fatalf(
				"Unexpected response header for key '%s'. Want: '%s'. Got: '%s'",
				key,
				value,
				response.Header.Get(key),
			)
		}
	}

	actualBody, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		t.Fatalf("Unexpected error reading response body: %s", err.Error())
	}

	if string(actualBody) != expectedResponseBody {
		t.Fatalf("Unexpected response body. Want: '%s'. Got: '%s'", expectedResponseBody, string(actualBody))
	}
}

func TestHTTPClientSendRequestAuthenticateError(t *testing.T) {
	authenticator := newMockAuthenticator(func(request *http.Request) error {
		return errors.New("fake-authentication-error")
	})

	httpClient := newHTTPClient(authenticator)

	_, err := httpClient.SendRequest(
		"POST",
		"fake-url",
		strings.NewReader("fake-body"),
		"fake-content-type",
	)

	if err == nil {
		t.Fatal("Expected error but got none.")
	}

	if !strings.HasPrefix(err.Error(), "fake-authentication-error") {
		t.Fatalf("Expected fake-authentication-error. Got: %s", err.Error())
	}
}

func newAuthenticatedHTTPClient() HTTPClient {
	authenticator := newMockAuthenticator(func(request *http.Request) error {
		request.Header.Set("Authorization", "fake-authorization")
		return nil
	})

	return newHTTPClient(authenticator)
}

func newMockAuthenticator(authenticate func(request *http.Request) error) *mockAuthenticator {
	return &mockAuthenticator{
		authenticate: authenticate,
	}
}

type mockAuthenticator struct {
	authenticate func(request *http.Request) error
}

func (a *mockAuthenticator) Authenticate(request *http.Request) error {
	if a.authenticate != nil {
		return a.authenticate(request)
	}
	return nil
}
