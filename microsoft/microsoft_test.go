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

func newMockLanguageProvider() *mockLanguageProvider {
	return &mockLanguageProvider{
		callCounter: make(map[string]int),
		codes:       make([]string, 0),
		names:       make([]string, 0),
	}
}

type mockLanguageProvider struct {
	callCounter map[string]int
	codes       []string
	names       []string
}

func (p *mockLanguageProvider) Codes() ([]string, error) {
	p.callCounter["Codes"]++
	return p.codes, nil
}

func (p *mockLanguageProvider) Names(codes []string) ([]string, error) {
	p.callCounter["Names"]++
	return p.names, nil
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
