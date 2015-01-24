package microsoft

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	_http "github.com/st3v/translator/http"
)

func TestTranslationProviderTranslate(t *testing.T) {
	expectedOriginal := "Ich verstehe nur Bahnhof."
	expectedTranslation := "I only understand train station."
	expectedFrom := "de"
	expectedTo := "en"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("Unexpected request method: %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "text/plain" {
			t.Fatalf("Unexpected content type in request header: %s", r.Header.Get("Content-Type"))
		}

		if r.FormValue("text") != expectedOriginal {
			t.Fatalf("Unexpected `text` param in request: %s", r.FormValue("text"))
		}

		if r.FormValue("to") != expectedTo {
			t.Fatalf("Unexpected `to` param in request: %s", r.FormValue("to"))
		}

		if r.FormValue("from") != expectedFrom {
			t.Fatalf("Unexpected `from` param in request: %s", r.FormValue("from"))
		}

		response, err := xml.Marshal(newXMLString(expectedTranslation))
		if err != nil {
			t.Fatalf("Unexpected error marshalling xml repsonse: %s", err.Error())
		}

		w.Header().Set("Content-Type", "text/xml")

		fmt.Fprint(w, string(response))
		return
	}))
	defer server.Close()

	router := newMockRouter()
	router.translationURL = server.URL

	translationProvider := &translationProvider{
		router:     router,
		httpClient: _http.NewAuthenticatedClient(),
	}

	actualTranslation, err := translationProvider.Translate(expectedOriginal, expectedFrom, expectedTo)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}

	if actualTranslation != expectedTranslation {
		t.Fatalf("Unexpected translation: %s. Expected: %s.", actualTranslation, expectedTranslation)
	}
}

func TestTranslationProviderDetect(t *testing.T) {
	text := "Ich verstehe nur Bahnhof."
	expectedLanguage := "de"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("Unexpected request method: %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "text/plain" {
			t.Fatalf("Unexpected content type in request header: %s", r.Header.Get("Content-Type"))
		}

		if r.FormValue("text") != text {
			t.Fatalf("Unexpected `text` param in request: %s", r.FormValue("text"))
		}

		response, err := xml.Marshal(newXMLString(expectedLanguage))
		if err != nil {
			t.Fatalf("Unexpected error marshalling xml repsonse: %s", err.Error())
		}

		w.Header().Set("Content-Type", "text/xml")

		fmt.Fprint(w, string(response))
		return
	}))
	defer server.Close()

	router := newMockRouter()
	router.detectURL = server.URL

	translationProvider := &translationProvider{
		router:     router,
		httpClient: _http.NewAuthenticatedClient(),
	}

	actualLanguage, err := translationProvider.Detect(text)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}

	if actualLanguage != expectedLanguage {
		t.Fatalf("Unexpected language detected: %s. Expected: %s.", actualLanguage, expectedLanguage)
	}
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

func (p *mockTranslationProvider) Detect(text string) (string, error) {
	return p.from, nil
}
