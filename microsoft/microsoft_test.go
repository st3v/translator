package microsoft

import (
	"fmt"
	"time"
)

func newMockAccessToken(expiresIn int) *accessToken {
	return &accessToken{
		Token:     "token",
		Type:      "token_type",
		Scope:     "token_scope",
		ExpiresIn: fmt.Sprintf("%d", expiresIn),
		ExpiresAt: time.Now().Add(time.Duration(expiresIn) * time.Second),
	}
}

func newMockAuthenticator(token *accessToken) *authenticator {
	return &authenticator{
		accessToken: token,
	}
}

func newMockRouter() *mockRouter {
	return &mockRouter{
		authUrl:          "auth",
		translationUrl:   "translation",
		languageNamesUrl: "languages_names",
		languageCodesUrl: "languages_codes",
	}
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
