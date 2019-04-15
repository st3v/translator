package microsoft

import (
	"encoding/json"
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
	expectedVersion := "3.0"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("Unexpected request method: %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("Unexpected content type in request header: %s", r.Header.Get("Content-Type"))
		}

		if r.FormValue("to") != expectedTo {
			t.Fatalf("Unexpected `to` param in request: %s", r.FormValue("to"))
		}

		if r.FormValue("from") != expectedFrom {
			t.Fatalf("Unexpected `from` param in request: %s", r.FormValue("from"))
		}
		var request interface{}
		tr := []byte(`[{"detectedLanguage":{"language": "en","score": 1.0},"translations":[{"text":"I only understand train station.","to": "en"},{"text": "Salve, mondo!","to": "it"}]}]`)
		err := json.Unmarshal(tr, &request)
		if err != nil {
			t.Fatalf("Unexpected error marshalling json response: %s", err.Error())
		}
		response, err := json.Marshal(&request)
		if err != nil {
			t.Fatalf("Unexpected error marshalling json response: %s", err.Error())
		}

		w.Header().Set("Content-Type", "application/json")

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

	actualTranslation, err := translationProvider.Translate(expectedOriginal, expectedFrom, expectedTo, expectedVersion)
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
	version := "3.0"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("Unexpected request method: %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("Unexpected content type in request header: %s", r.Header.Get("Content-Type"))
		}

		var request interface{}
		tr := []byte(`[{"language":"de","score":1.0,"isTranslationSupported":true,"isTransliterationSupported":false,"alternatives":[{"language":"en","score":0.75,"isTranslationSupported":false,"isTransliterationSupported":false},{"language":"pl","score":0.75,"isTranslationSupported":true,"isTransliterationSupported":false}]}]`)
		err := json.Unmarshal(tr, &request)
		if err != nil {
			t.Fatalf("Unexpected error marshalling json response: %s", err.Error())
		}

		response, err := json.Marshal(&request)
		if err != nil {
			t.Fatalf("Unexpected error marshalling json response: %s", err.Error())
		}

		w.Header().Set("Content-Type", "application/json")

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

	actualLanguage, err := translationProvider.Detect(text, version)
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

func (p *mockTranslationProvider) Translate(text, from, to, version string) (string, error) {
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

func (p *mockTranslationProvider) Detect(text, version string) (string, error) {
	return p.from, nil
}
