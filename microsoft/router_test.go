package microsoft

import "testing"

func TestRouterAuthURL(t *testing.T) {
	router := newRouter()

	expectedURL := "https://datamarket.accesscontrol.windows.net/v2/OAuth2-13"

	actualURL := router.AuthURL()

	if actualURL != expectedURL {
		t.Fatalf("Unexpected AuthURL. Want: '%s'. Got: '%s'.", expectedURL, actualURL)
	}
}

func TestRouterTranslationURL(t *testing.T) {
	router := newRouter()

	expectedURL := "http://api.microsofttranslator.com/v2/Http.svc/Translate"

	actualURL := router.TranslationURL()

	if actualURL != expectedURL {
		t.Fatalf("Unexpected AuthURL. Want: '%s'. Got: '%s'.", expectedURL, actualURL)
	}
}

func TestRouterLanguageNamesURL(t *testing.T) {
	router := newRouter()

	expectedURL := "http://api.microsofttranslator.com/v2/Http.svc/GetLanguageNames"

	actualURL := router.LanguageNamesURL()

	if actualURL != expectedURL {
		t.Fatalf("Unexpected AuthURL. Want: '%s'. Got: '%s'.", expectedURL, actualURL)
	}
}

func TestRouterLanguageCodesURL(t *testing.T) {
	router := newRouter()

	expectedURL := "http://api.microsofttranslator.com/v2/Http.svc/GetLanguagesForTranslate"

	actualURL := router.LanguageCodesURL()

	if actualURL != expectedURL {
		t.Fatalf("Unexpected AuthURL. Want: '%s'. Got: '%s'.", expectedURL, actualURL)
	}
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
