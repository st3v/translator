package microsoft

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLanguageCodes(t *testing.T) {
	expectedCodes := []string{"ab", "cd", "ef", "gh", "ij"}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/xml")

		response, err := xml.Marshal(newArrayOfStrings(expectedCodes))
		if err != nil {
			t.Fatalf("Unexpected error marshalling xml repsonse: %s", err)
		}

		fmt.Fprint(w, string(response))
		return
	}))
	defer server.Close()

	authenticator := NewMockAuthenticator()
	authenticator.accessToken = NewMockAccessToken(100)

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
