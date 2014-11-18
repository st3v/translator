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
