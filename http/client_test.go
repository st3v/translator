package http

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClientSendRequest(t *testing.T) {
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

	client := NewClient(authenticator)

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

	response, err := client.SendRequest(
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

func TestClientSendRequestAuthenticateError(t *testing.T) {
	authenticator := newMockAuthenticator(func(request *http.Request) error {
		return errors.New("fake-authentication-error")
	})

	client := NewClient(authenticator)

	_, err := client.SendRequest(
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
