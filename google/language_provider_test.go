package google

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLanguages(t *testing.T) {
	expectedAPIKey := "my-secret-key"

	expectedLanguages := []struct{ Language, Name string }{
		{"en", "English"},
		{"de", "German"},
		{"zh-TW", "Chinese (Traditional)"},
	}

	authenticator := newAuthenticator(expectedAPIKey)

	requestCounter := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCounter++

		if r.Method != "GET" {
			t.Fatalf("Unexpected request method: %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "text/plain" {
			t.Fatalf("Unexpected content type in request header: %s", r.Header.Get("Content-Type"))
		}

		if r.FormValue("key") != expectedAPIKey {
			t.Fatalf("Unexpected `key` param in request. Got: %s. Want: %s", r.FormValue("key"), expectedAPIKey)
		}

		expectedTarget := "en"
		if r.FormValue("target") != expectedTarget {
			t.Fatalf("Unexpected `target` param in request. Got: %s. Want: %s", r.FormValue("target"), expectedTarget)
		}

		result := languagesPayload{}
		result.Data.Languages = expectedLanguages
		response, err := json.Marshal(languagesPayload(result))
		if err != nil {
			t.Fatalf("Unexpected error marshalling json repsonse: %s", err.Error())
		}

		w.Header().Set("Content-Type", "application/json")

		fmt.Fprint(w, string(response))
		return
	}))
	defer server.Close()

	router := &router{languagesEndpoint: server.URL}

	provider := newLanguageProvider(authenticator, router)
	languages, err := provider.languages()
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	for i := 0; i < 3; i++ {
		if len(languages) != len(expectedLanguages) {
			t.Errorf(
				"Unexpected number of languages. Got: %d. Want: %d.",
				len(languages),
				len(expectedLanguages),
			)
		}

		for i, l := range languages {
			if l.Code != expectedLanguages[i].Language {
				t.Errorf(
					"Unexpected language code. Got: %s. Want: %s.",
					l.Code,
					expectedLanguages[i].Language,
				)
			}

			if l.Name != expectedLanguages[i].Name {
				t.Errorf(
					"Unexpected language name. Got: %s. Want: %s.",
					l.Name,
					expectedLanguages[i].Name,
				)
			}
		}
	}

	if requestCounter != 1 {
		t.Errorf("Expected 1 http request but counted %d.", requestCounter)
	}
}

func TestDetect(t *testing.T) {
	expectedAPIKey := "my-secret-key"

	expectedText := "foo"
	expectedLanguageCode := "bar"

	authenticator := newAuthenticator(expectedAPIKey)

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

		if r.FormValue("q") != expectedText {
			t.Fatalf("Unexpected `q` param in request. Got: %s. Want: %s", r.FormValue("q"), expectedText)
		}

		w.Header().Set("Content-Type", "application/json")

		jsonPayload := fmt.Sprintf(`
			{
				"data": {
					"detections": [
						[
							{
								"language": "%s",
								"isReliable": false,
								"confidence": 0.66
							}
						]
					]
				}
			}
		`, expectedLanguageCode)

		fmt.Fprint(w, jsonPayload)
		return
	}))
	defer server.Close()

	router := &router{detectEndpoint: server.URL}

	provider := newLanguageProvider(authenticator, router)

	languageCode, err := provider.detect(expectedText)

	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	if languageCode != expectedLanguageCode {
		t.Errorf(
			"Unexpected language code. Got: %s. Want: %s.",
			languageCode,
			expectedLanguageCode,
		)
	}
}
