package google

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTranslate(t *testing.T) {
	expectedOriginal := "Rindfleischetikettierungsüberwachungsaufgabenübertragungsgesetz"
	expectedTranslation := "WTF!?!"

	expectedSource := "de"
	expectedTarget := "en"

	expectedAPIKey := "my-secret-key"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("Unexpected request method: %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "text/plain" {
			t.Fatalf("Unexpected content type in request header: %s", r.Header.Get("Content-Type"))
		}

		if r.FormValue("key") != expectedAPIKey {
			t.Fatalf("Unexpected `key` param in request. Got: %s. Want: %s", r.FormValue("key"), expectedAPIKey)
		}

		if r.FormValue("source") != expectedSource {
			t.Fatalf("Unexpected `source` param in request. Got: %s. Want: %s", r.FormValue("source"), expectedSource)
		}

		if r.FormValue("target") != expectedTarget {
			t.Fatalf("Unexpected `target` param in request. Got: %s. Want: %s", r.FormValue("target"), expectedTarget)
		}

		if r.FormValue("q") != expectedOriginal {
			t.Fatalf("Unexpected `q` param in request. Got: %s. Want: %s", r.FormValue("q"), expectedOriginal)
		}

		w.Header().Set("Content-Type", "application/json")

		jsonResponse := fmt.Sprintf(
			`{ "data": { "translations": [ { "translatedText": "%s" } ] } }`,
			expectedTranslation,
		)

		fmt.Fprint(w, jsonResponse)
		return
	}))
	defer server.Close()

	authenticator := newAuthenticator(expectedAPIKey)
	router := &router{translateEndpoint: server.URL}
	provider := newTranslationProvider(authenticator, router)

	actualTranslation, err := provider.translate(expectedOriginal, expectedSource, expectedTarget)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	if actualTranslation != expectedTranslation {
		t.Errorf(
			"Unexpected translation result. Got: '%s'. Want: '%s'.",
			actualTranslation,
			expectedTranslation,
		)
	}
}
