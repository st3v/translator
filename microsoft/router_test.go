package microsoft

import "testing"

func TestRouterAuthURL(t *testing.T) {
	router := newRouter()

	expectedURL := "https://api.cognitive.microsoft.com/sts/v1.0/issueToken"

	actualURL := router.AuthURL()

	if actualURL != expectedURL {
		t.Fatalf("Unexpected AuthURL. Want: %q. Got: %q.", expectedURL, actualURL)
	}
}

func TestRouterTranslationURL(t *testing.T) {
	router := newRouter()

	expectedURL := "https://api.microsofttranslator.com/v2/Http.svc/Translate"

	actualURL := router.TranslationURL()

	if actualURL != expectedURL {
		t.Fatalf("Unexpected TranslationURL. Want: %q. Got: %q.", expectedURL, actualURL)
	}
}

func TestRouterLanguageNamesURL(t *testing.T) {
	router := newRouter()

	expectedURL := "https://api.microsofttranslator.com/v2/Http.svc/GetLanguageNames"

	actualURL := router.LanguageNamesURL()

	if actualURL != expectedURL {
		t.Fatalf("Unexpected LanguageNamesURL. Want: %q. Got: %q.", expectedURL, actualURL)
	}
}

func TestRouterLanguageCodesURL(t *testing.T) {
	router := newRouter()

	expectedURL := "https://api.microsofttranslator.com/v2/Http.svc/GetLanguagesForTranslate"

	actualURL := router.LanguageCodesURL()

	if actualURL != expectedURL {
		t.Fatalf("Unexpected LanguageCodesURL. Want: %q. Got: %q.", expectedURL, actualURL)
	}
}

func newMockRouter() *mockRouter {
	return &mockRouter{
		authURL:          "auth",
		translationURL:   "translation",
		languageNamesURL: "languages_names",
		languageCodesURL: "languages_codes",
		detectURL:        "detect",
	}
}

type mockRouter struct {
	authURL          string
	translationURL   string
	languageNamesURL string
	languageCodesURL string
	detectURL        string
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

func (m *mockRouter) DetectURL() string {
	return m.detectURL
}
