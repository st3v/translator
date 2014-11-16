package microsoft

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLanguageCodes(t *testing.T) {
	expectedCodes := []string{"en", "de", "es", "ru", "jp"}

	authenticator := NewMockAuthenticator()
	authenticator.accessToken = NewMockAccessToken(100)

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

		response, err := xml.Marshal(newArrayOfStrings(expectedCodes))
		if err != nil {
			t.Fatalf("Unexpected error marshalling xml repsonse: %s", err)
		}

		w.Header().Set("Content-Type", "text/xml")

		fmt.Fprint(w, string(response))
		return
	}))
	defer server.Close()

	router := NewMockRouter()
	router.languageCodesUrl = server.URL

	api := NewMockApi()
	api.router = router
	api.authenticator = authenticator

	actualCodes, err := api.languageCodes()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if len(actualCodes) != len(expectedCodes) {
		t.Fatalf("Unexpected number of languages codes: %q", actualCodes)
	}

	for i := range expectedCodes {
		if actualCodes[i] != expectedCodes[i] {
			t.Fatalf("Unexpected language code '%s'. Expected '%s'", actualCodes[i], expectedCodes[i])
		}
	}
}
