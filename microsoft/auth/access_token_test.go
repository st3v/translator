package auth

import (
	"testing"
	"time"
)

// Make sure the access token expires as expected.
func TestAccessTokenExpired(t *testing.T) {
	accessToken := newMockAccessToken(12 * time.Second)
	if accessToken.expired() {
		t.Fatalf("Access token should not have expired. Now: %s. ExpiresAt: %s.", time.Now().String(), accessToken.ExpiresAt.String())
	}

	accessToken = newMockAccessToken(0 * time.Second)
	if !accessToken.expired() {
		t.Fatalf("Access token should have expired. Now: %s. ExpiresAt: %s.", time.Now().String(), accessToken.ExpiresAt.String())
	}
}

func newMockAccessToken(expiresIn time.Duration) *accessToken {
	return &accessToken{
		Token:     "token",
		ExpiresAt: time.Now().Add(expiresIn),
	}
}

func newMockAccessTokenProvider() *mockAccessTokenProvider {
	return &mockAccessTokenProvider{
		refreshToken: func(token *accessToken) error {
			return nil
		},
	}
}

type mockAccessTokenProvider struct {
	refreshToken func(token *accessToken) error
}

func (p *mockAccessTokenProvider) RefreshToken(token *accessToken) error {
	return p.refreshToken(token)
}

func (a *authenticator) expectedAuthToken(t *testing.T) string {
	token, err := a.authToken()
	if err != nil {
		t.Fatalf("Unexpected error getting authToken from authenticator: %s", err.Error())
	}
	return token
}
