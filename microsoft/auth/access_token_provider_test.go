package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// and is able to generate a valid access token from the server's response.
func TestAccessTokenProviderRefreshToken(t *testing.T) {
	clientID := "foobar"
	clientSecret := "private"

	expectedToken := newMockAccessToken(100)
	expectedScope := "mock-scope"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("Unexpected request method: %s", r.Method)
		}

		if r.PostFormValue("client_id") != clientID {
			t.Fatalf("Unexpected client_id '%s' in post request.", r.PostFormValue("client_id"))
		}

		if r.PostFormValue("client_secret") != clientSecret {
			t.Fatalf("Unexpected client_secret '%s' in post request.", r.PostFormValue("client_secret"))
		}

		if r.PostFormValue("scope") != expectedScope {
			t.Fatalf("Unexpected scope '%s' in post request.", r.PostFormValue("scope"))
		}

		if r.PostFormValue("grant_type") != "client_credentials" {
			t.Fatalf("Unexpected grant_type '%s' in post request.", r.PostFormValue("grant_type"))
		}

		response, err := json.Marshal(expectedToken)
		if err != nil {
			t.Fatalf("Unexpected error marshalling json repsonse: %s", err.Error())
		}

		w.Header().Set("Content-Type", "application/json")

		fmt.Fprint(w, string(response))
		return
	}))
	defer server.Close()

	accessTokenProvider := newAccessTokenProvider(clientID, clientSecret, server.URL)

	actualToken := newAccessToken(expectedScope)

	if err := accessTokenProvider.RefreshToken(actualToken); err != nil {
		t.Fatalf("Unexpected error returned by requestAccessToken: %s", err.Error())
	}

	if actualToken.Token != expectedToken.Token {
		t.Fatalf("Unexpected Token '%s' in access token generated from http response.", actualToken.Token)
	}

	if actualToken.Type != expectedToken.Type {
		t.Fatalf("Unexpected Type '%s' in access token generated from http response.", actualToken.Type)
	}

	if actualToken.ExpiresIn != expectedToken.ExpiresIn {
		t.Fatalf("Unexpected ExpiresIn '%s' in access token generated from http response.", actualToken.ExpiresIn)
	}

	if actualToken.Scope != expectedToken.Scope {
		t.Fatalf("Unexpected Scope '%s' in access token generated from http response.", actualToken.Scope)
	}

	if actualToken.ExpiresAt.After(time.Now().Add(100*time.Second)) ||
		actualToken.ExpiresAt.Before(time.Now().Add(97*time.Second)) {
		t.Fatalf("Unexpected ExpiresAt '%s' in access token generated from http response.", actualToken.ExpiresAt)
	}
}
