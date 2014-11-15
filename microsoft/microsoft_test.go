package microsoft

import (
	"fmt"
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

type mockRouter struct {
	authUrl          string
	translationUrl   string
	languageNamesUrl string
	languageCodesUrl string
}

func NewMockRouter() *mockRouter {
	return &mockRouter{
		authUrl:          "auth",
		translationUrl:   "translation",
		languageNamesUrl: "languages_names",
		languageCodesUrl: "languages_codes",
	}
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
