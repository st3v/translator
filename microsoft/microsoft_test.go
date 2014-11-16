package microsoft

import (
	"fmt"
	"testing"
	"time"
)

func NewMockAccessToken(expiresIn int) *accessToken {
	return &accessToken{
		Token:     "token",
		Type:      "token_type",
		Scope:     "token_scope",
		ExpiresIn: fmt.Sprintf("%d", expiresIn),
		ExpiresAt: time.Now().Add(time.Duration(expiresIn) * time.Second),
	}
}

func NewMockApi() *api {
	return &api{
		authenticator: NewMockAuthenticator(),
		router:        NewMockRouter(),
	}
}

func NewMockAuthenticator() *authenticator {
	return &authenticator{}
}

func NewMockRouter() *mockRouter {
	return &mockRouter{
		authUrl:          "auth",
		translationUrl:   "translation",
		languageNamesUrl: "languages_names",
		languageCodesUrl: "languages_codes",
	}
}

func (a *authenticator) expectedAuthToken(t *testing.T) string {
	token, err := a.authToken()
	if err != nil {
		t.Fatalf("Unexpected error getting authToken from authenticator: %s", err)
	}
	return token
}

type mockRouter struct {
	authUrl          string
	translationUrl   string
	languageNamesUrl string
	languageCodesUrl string
}

func (m *mockRouter) AuthUrl() string {
	return m.authUrl
}

func (m *mockRouter) TranslationUrl() string {
	return m.translationUrl
}

func (m *mockRouter) LanguageNamesUrl() string {
	return m.languageNamesUrl
}

func (m *mockRouter) LanguageCodesUrl() string {
	return m.languageCodesUrl
}

// func (m *mockRouter) SetAuthUrl(url string) {
// 	m.authUrl = url
// }

// func (m *mockRouter) SetTranslationUrl(url string) {
// 	m.translationUrl = url
// }

// func (m *mockRouter) SetLanguageNamesUrl(url string) {
// 	m.languageNamesUrl = url
// }

// func (m *mockRouter) SeLanguageCodesUrl(url string) {
// 	m.languageCodesUrl = url
// }
