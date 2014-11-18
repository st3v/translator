package microsoft

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// Make sure function requestAccessToken sends the expected request to the server
// and is able to generate a valid access token from the server's response.
func TestAuthenticatorRefreshAccessToken(t *testing.T) {
	clientId := "foobar"
	clientSecret := "private"

	expectedToken := newMockAccessToken(100)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("Unexpected request method: %s", r.Method)
		}

		if r.PostFormValue("client_id") != clientId {
			t.Fatalf("Unexpected client_id '%s' in post request.", r.PostFormValue("client_id"))
		}

		if r.PostFormValue("client_secret") != clientSecret {
			t.Fatalf("Unexpected client_secret '%s' in post request.", r.PostFormValue("client_secret"))
		}

		if r.PostFormValue("scope") != scope {
			t.Fatalf("Unexpected scope '%s' in post request.", r.PostFormValue("scope"))
		}

		if r.PostFormValue("grant_type") != "client_credentials" {
			t.Fatalf("Unexpected grant_type '%s' in post request.", r.PostFormValue("grant_type"))
		}

		response, err := json.Marshal(expectedToken)
		if err != nil {
			t.Fatalf("Unexpected error marshalling json repsonse: %s", err)
		}

		w.Header().Set("Content-Type", "application/json")

		fmt.Fprint(w, string(response))
		return
	}))
	defer server.Close()

	router := newMockRouter()
	router.authUrl = server.URL

	authenticationProvider := &authenticationProvider{
		clientId:     clientId,
		clientSecret: clientSecret,
		router:       router,
	}

	actualToken := &accessToken{}
	if err := authenticationProvider.RefreshAccessToken(actualToken); err != nil {
		t.Fatalf("Unexpected error returned by requestAccessToken: %s", err)
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

	// verify that the expiration time is wihin 3 seconds of what is expected
	if actualToken.ExpiresAt.After(time.Now().Add(100*time.Second)) ||
		actualToken.ExpiresAt.Before(time.Now().Add(97*time.Second)) {
		t.Fatalf("Unexpected ExpiresAt '%s' in access token generated from http response.", actualToken.ExpiresAt)
	}
}

// Make sure the access token expires as expected.
func TestAccessTokenExpired(t *testing.T) {
	accessToken := newMockAccessToken(12)
	if accessToken.expired() {
		t.Fatalf("Access token should not have expired. Now: %s. ExpiresAt: %s.", time.Now().String(), accessToken.ExpiresAt.String())
	}

	accessToken = newMockAccessToken(0)
	if !accessToken.expired() {
		t.Fatalf("Access token should have expired. Now: %s. ExpiresAt: %s.", time.Now().String(), accessToken.ExpiresAt.String())
	}
}

// Make sure a valid authToken is being generated from a given access token.
func TestAuthenticatorAuthToken(t *testing.T) {
	authenticator := newMockAuthenticator(newMockAccessToken(100))

	authToken, err := authenticator.authToken()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	expectedToken := <-authenticator.accessTokenChan
	if authToken != fmt.Sprintf("Bearer %s", expectedToken.Token) {
		t.Fatalf("Invalid authToken ''.", authToken)
	}
}

// Make sure Authenticate() correctly sets the authrorization header of a given request.
func TestAuthenticatorAuthenticate(t *testing.T) {
	authenticator := newMockAuthenticator(newMockAccessToken(100))

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

// Concurrency test for Authenticator.Authenticate().
// Run tests with race detector enabled
func TestAuthenticatorConcurrentAuthenticate(t *testing.T) {
	callCount := 0

	// mock out authentication provider, keep track of how many time the access token is being refreshed
	provider := newMockAuthenticationProvider()
	provider.refreshAccessToken = func(token *accessToken) error {
		callCount++
		*token = *newMockAccessToken(100)
		return nil
	}

	// create an authenticator that uses the mock provider and starts with an expired access token that needs to be refreshed
	authenticator := newMockAuthenticator(&accessToken{})
	authenticator.provider = provider

	// channel to make sure all concurrent go routines start at the same time
	readyGo := make(chan bool)

	// one error channel for each go routine
	errorChans := make([]chan error, 10)

	// spin up 10 concurrent go routines that each call Authenticate()
	// only one should trigger an access token refresh in the provider
	for i := 0; i < 10; i++ {
		errChan := make(chan error)
		go func() {
			<-readyGo
			authToken, err := authenticator.authToken()
			if err == nil && authToken != authenticator.expectedAuthToken(t) {
				err = errors.New(fmt.Sprintf("Unexpected authToken `%s`. Expected `%s`.", authToken, authenticator.expectedAuthToken(t)))
			}
			errChan <- err
			close(errChan)
		}()
		errorChans[i] = errChan
	}

	// ready set go!
	close(readyGo)

	// merge all error channels into a single channel
	// loop over errors in the channel and make sure they are all nil
	for err := range mergeErrorChans(errorChans) {
		if err != nil {
			t.Error(err)
		}
	}

	// verify call count for RefreshAccessToken
	if callCount != 1 {
		t.Fatalf("Expected RefreshAccessToken to be called exactly once. Looks like it was called %d times.", callCount)
	}
}

// Merges a slice of channels of errors into a single incoming channel of errors.
func mergeErrorChans(cs []chan error) <-chan error {
	var wg sync.WaitGroup
	out := make(chan error)

	output := func(c <-chan error) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
