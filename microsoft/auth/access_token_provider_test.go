package auth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// and is able to generate a valid access token from the server's response.
func TestAccessTokenProviderRefreshToken(t *testing.T) {
	subscriptionKey := "private"

	expectedToken := "some-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("Unexpected request method: %s", r.Method)
		}

		if have, want := r.Header.Get("Ocp-Apim-Subscription-Key"), subscriptionKey; have != want {
			t.Fatalf("Unexpected 'Ocp-Apim-Subscription-Key' header: want %q, have %q.", want, have)
		}

		w.Header().Set("Content-Type", "application/jwt; charset=us-ascii")

		fmt.Fprint(w, expectedToken)
		return
	}))
	defer server.Close()

	accessTokenProvider := newAccessTokenProvider(subscriptionKey, server.URL)

	actualToken := new(accessToken)
	if err := accessTokenProvider.RefreshToken(actualToken); err != nil {
		t.Fatalf("Unexpected error returned by RefreshToken: %v", err.Error())
	}

	if have, want := actualToken.Token, expectedToken; have != want {
		t.Fatalf("Unexpected Token: want %q, have %q.", actualToken.Token)
	}

	if s := actualToken.ExpiresAt.Sub(time.Now()).Seconds(); s < 598 || s > 600 {
		t.Fatalf("Unexpected ExpiresAt %q for access token generated from http response.", actualToken.ExpiresAt)
	}
}
