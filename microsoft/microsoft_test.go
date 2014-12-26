package microsoft

import (
	"fmt"
	"testing"
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
	// make buffered accessToken channel an pre-fill it with nil
	tokenChan := make(chan *accessToken, 1)
	tokenChan <- token

	// return new authenticator that uses the above accessToken channel
	return &authenticator{
		accessTokenChan: tokenChan,
	}
}

func newMockAuthenticationProvider() *mockAuthenticationProvider {
	return &mockAuthenticationProvider{
		refreshAccessToken: func(token *accessToken) error {
			return nil
		},
	}
}

type mockAuthenticationProvider struct {
	refreshAccessToken func(token *accessToken) error
}

func (p *mockAuthenticationProvider) RefreshAccessToken(token *accessToken) error {
	return p.refreshAccessToken(token)
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

func newMockTranslationProvider(text, from, to, translation string, t *testing.T) *mockTranslationProvider {
	return &mockTranslationProvider{
		text:        text,
		from:        from,
		to:          to,
		translation: translation,
		t:           t,
	}
}

type mockTranslationProvider struct {
	text        string
	from        string
	to          string
	translation string
	t           *testing.T
}

func (p *mockTranslationProvider) Translate(text, from, to string) (string, error) {
	if p.text != text {
		p.t.Fatalf("Unexpected text value: `%s`", text)
	}

	if p.from != from {
		p.t.Fatalf("Unexpected from value: `%s`", from)
	}

	if p.to != to {
		p.t.Fatalf("Unexpected to value: `%s`", to)
	}
	return p.translation, nil
}

func newMockRouter() *mockRouter {
	return &mockRouter{
		authURL:          "auth",
		translationURL:   "translation",
		languageNamesURL: "languages_names",
		languageCodesURL: "languages_codes",
	}
}

type mockRouter struct {
	authURL          string
	translationURL   string
	languageNamesURL string
	languageCodesURL string
}

func (m *mockRouter) AuthURL() string {
	return m.authURL
}

func (m *mockRouter) TranslationURL() string {
	return m.translationURL
}

func (m *mockRouter) LanguageNamesURL() string {
	return m.languageNamesURL
}

func (m *mockRouter) LanguageCodesURL() string {
	return m.languageCodesURL
}
