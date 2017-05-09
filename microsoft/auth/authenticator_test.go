package auth

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"
)

// Make sure a valid authToken is being generated from a given access token.
func TestAuthenticatorAuthToken(t *testing.T) {
	authenticator := newMockAuthenticator(newMockAccessToken(100))

	authToken, err := authenticator.authToken()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}

	expectedToken := <-authenticator.accessTokenChan
	if authToken != fmt.Sprintf("Bearer %s", expectedToken.Token) {
		t.Fatalf("Invalid authToken '%s'.", authToken)
	}
}

// Make sure Authenticate() correctly sets the authrorization header of a given request.
func TestAuthenticatorAuthenticate(t *testing.T) {
	authenticator := newMockAuthenticator(newMockAccessToken(10 * time.Minute))

	authToken, err := authenticator.authToken()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}

	r, err := http.NewRequest("GET", "http://foo.bar", nil)
	if err != nil {
		t.Fatalf("Unexpected error when getting new request: %s", err.Error())
	}

	if r.Header.Get("Authorization") != "" {
		t.Fatalf("Authorization header should not haven been set. Header: %s", r.Header.Get("Authorization"))
	}

	authenticator.Authenticate(r)

	if r.Header.Get("Authorization") != authToken {
		t.Fatalf("Unexpected authorization header. Header: %s", r.Header.Get("Authorization"))
	}
}

// Concurrency test for Authenticator.Authenticate().
// Run tests with race detector enabled
func TestAuthenticatorConcurrentAuthenticate(t *testing.T) {
	callCount := 0

	// mock out authentication provider, keep track of how many time the access token is being refreshed
	provider := newMockAccessTokenProvider()
	provider.refreshToken = func(token *accessToken) error {
		callCount++
		*token = *newMockAccessToken(10 * time.Minute)
		return nil
	}

	// create an authenticator that uses the mock provider and starts with an expired access token that needs to be refreshed
	authenticator := newMockAuthenticator(&accessToken{})
	authenticator.accessTokenProvider = provider

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
				err = fmt.Errorf("Unexpected authToken `%s`. Expected `%s`.", authToken, authenticator.expectedAuthToken(t))
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
			t.Error(err.Error())
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

func newMockAuthenticator(token *accessToken) *authenticator {
	// make buffered accessToken channel an pre-fill it with nil
	tokenChan := make(chan *accessToken, 1)
	tokenChan <- token

	// return new authenticator that uses the above accessToken channel
	return &authenticator{
		accessTokenChan:     tokenChan,
		accessTokenProvider: newMockAccessTokenProvider(),
	}
}
