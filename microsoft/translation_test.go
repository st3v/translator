package microsoft

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTranslationProviderTranslate(t *testing.T) {
	expectedOriginal := "Ich verstehe nur Bahnhof."
	expectedTranslation := "I only understand train station."
	expectedFrom := "de"
	expectedTo := "en"

	authenticator := newMockAuthenticator(newMockAccessToken(100))

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("Unexpected request method: %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "text/plain" {
			t.Fatalf("Unexpected content type in request header: %s", r.Header.Get("Content-Type"))
		}

		if r.Header.Get("Authorization") != authenticator.expectedAuthToken(t) {
			t.Fatalf("Unexpected authorization header for request: %s", r.Header.Get("Authorization"))
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

		response, err := xml.Marshal(newXmlString(expectedTranslation))
		if err != nil {
			t.Fatalf("Unexpected error marshalling xml repsonse: %s", err)
		}

		w.Header().Set("Content-Type", "text/xml")

		fmt.Fprint(w, string(response))
		return
	}))
	defer server.Close()

	router := newMockRouter()
	router.translationUrl = server.URL

	translationProvider := &translationProvider{
		router:     router,
		httpClient: newHttpClient(authenticator),
	}

	actualTranslation, err := translationProvider.Translate(expectedOriginal, expectedFrom, expectedTo)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if actualTranslation != expectedTranslation {
		t.Fatalf("Unexpected translatipon: %s. Expected: %s.", actualTranslation, expectedTranslation)
	}
}
