package microsoft

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Make sure function requestAccessToken sends the expected request to the server
// and is able to generate a valid access token from the server's response.
func TestRequestAccessToken(t *testing.T) {
	clientId := "foobar"
	clientSecret := "private"

	accessToken := fakeAccessToken(100)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.PostFormValue("client_id") != clientId {
			t.Fatalf("Unexpected client_id '%s' in post request.", r.PostFormValue("client_id"))
		}

		if r.PostFormValue("client_secret") != clientSecret {
			t.Fatalf("Unexpected client_secret '%s' in post request.", r.PostFormValue("client_secret"))
		}

		if r.PostFormValue("scope") != SCOPE {
			t.Fatalf("Unexpected scope '%s' in post request.", r.PostFormValue("scope"))
		}

		if r.PostFormValue("grant_type") != "client_credentials" {
			t.Fatalf("Unexpected grant_type '%s' in post request.", r.PostFormValue("grant_type"))
		}

		response, err := json.Marshal(accessToken)
		if err != nil {
			t.Fatalf("Unexpected error marshalling json repsonse: %s", err)
		}

		fmt.Fprint(w, string(response))
		return
	}))
	defer server.Close()

	authenticator := &authenticator{
		clientId:     clientId,
		clientSecret: clientSecret,
	}

	err := authenticator.requestAccessToken(server.URL)

	if err != nil {
		t.Fatalf("Unexpected error returned by requestAccessToken: %s", err)
	}

	if authenticator.accessToken.Token != accessToken.Token {
		t.Fatalf("Unexpected Token '%s' in access token generated from http response.", authenticator.accessToken.Token)
	}

	if authenticator.accessToken.Type != accessToken.Type {
		t.Fatalf("Unexpected Type '%s' in access token generated from http response.", authenticator.accessToken.Type)
	}

	if authenticator.accessToken.ExpiresIn != accessToken.ExpiresIn {
		t.Fatalf("Unexpected ExpiresIn '%s' in access token generated from http response.", authenticator.accessToken.ExpiresIn)
	}

	if authenticator.accessToken.Scope != accessToken.Scope {
		t.Fatalf("Unexpected Scope '%s' in access token generated from http response.", authenticator.accessToken.Scope)
	}

	// verify that the expiration time is wihin 3 seconds of what is expected
	if authenticator.accessToken.ExpiresAt.After(time.Now().Add(100*time.Second)) ||
		authenticator.accessToken.ExpiresAt.Before(time.Now().Add(97*time.Second)) {
		t.Fatalf("Unexpected ExpiresAt '%s' in access token generated from http response.", authenticator.accessToken.ExpiresAt)
	}
}

// Make sure the access token expires as expected.
func TestExpired(t *testing.T) {
	accessToken := fakeAccessToken(12)
	if accessToken.expired() {
		t.Fatalf("Access token should not have expired. Now: %s. ExpiresAt: %s.", time.Now().String(), accessToken.ExpiresAt.String())
	}

	accessToken = fakeAccessToken(0)
	if !accessToken.expired() {
		t.Fatalf("Access token should have expired. Now: %s. ExpiresAt: %s.", time.Now().String(), accessToken.ExpiresAt.String())
	}
}

// Make sure a valid authToken is being generated from a given access token.
func TestAuthToken(t *testing.T) {
	authenticator := &authenticator{
		clientId:     "clientId",
		clientSecret: "clientSecret",
		accessToken:  fakeAccessToken(100),
	}

	authToken, err := authenticator.authToken()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if authToken != fmt.Sprintf("Bearer %s", authenticator.accessToken.Token) {
		t.Fatalf("Invalid authToken ''.", authToken)
	}
}

// Make sure Authenticate() correctly sets the authrorization header of a given request.
func TestAuthenticate(t *testing.T) {
	authenticator := &authenticator{
		clientId:     "clientId",
		clientSecret: "clientSecret",
		accessToken:  fakeAccessToken(100),
	}

	authToken, err := authenticator.authToken()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	r, err := http.NewRequest("GET", "http://foo.bar", nil)
	if err != nil {
		t.Fatalf("Unexpected error when getting new request: %s", err)
	}

	if r.Header.Get("Authorization") != "" {
		t.Fatalf("Authorization header should not haven been set. Header: ", r.Header.Get("Authorization"))
	}

	authenticator.Authenticate(r)

	if r.Header.Get("Authorization") != authToken {
		t.Fatalf("Unexpected authorization header. Header: ", r.Header.Get("Authorization"))
	}
}

func fakeAccessToken(expiresIn int) *accessToken {
	return &accessToken{
		Token:     "token",
		Type:      "token_type",
		Scope:     "token_scope",
		ExpiresIn: fmt.Sprintf("%d", expiresIn),
		ExpiresAt: time.Now().Add(time.Duration(expiresIn) * time.Second),
	}
}
