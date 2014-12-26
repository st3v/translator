package microsoft

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLanguageProviderCodes(t *testing.T) {
	expectedCodes := []string{"en", "de", "es", "ru", "jp"}

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

		response, err := xml.Marshal(newXMLArrayOfStrings(expectedCodes))
		if err != nil {
			t.Fatalf("Unexpected error marshalling xml repsonse: %s", err.Error())
		}

		w.Header().Set("Content-Type", "text/xml")

		fmt.Fprint(w, string(response))
		return
	}))
	defer server.Close()

	router := newMockRouter()
	router.languageCodesURL = server.URL

	languageProvider := &languageProvider{
		router:     router,
		httpClient: newHTTPClient(authenticator),
	}

	actualCodes, err := languageProvider.Codes()
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

func TestLanguageProviderNames(t *testing.T) {
	expectedCodes := []string{"en", "de", "es", "ru", "jp"}
	expectedNames := []string{"English", "German", "Spanish", "Russian", "Japanese"}

	authenticator := newMockAuthenticator(newMockAccessToken(100))

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("Unexpected request method: %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "text/xml" {
			t.Fatalf("Unexpected content type in request header: %s", r.Header.Get("Content-Type"))
		}

		if r.Header.Get("Authorization") != authenticator.expectedAuthToken(t) {
			t.Fatalf("Unexpected authorization header for request: %s", r.Header.Get("Authorization"))
		}

		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			t.Fatalf("Unexpected error reading request body: %s", err.Error())
		}

		actualCodes := &xmlArrayOfStrings{}
		if err := xml.Unmarshal(body, &actualCodes); err != nil {
			t.Fatalf("Unexpected error unmarshalling xml request body: %s", err.Error())
		}

		if len(actualCodes.Strings) != len(expectedCodes) {
			t.Fatalf("Unexpected number of languages codes in request: %q", actualCodes.Strings)
		}

		for i := range expectedCodes {
			if actualCodes.Strings[i] != expectedCodes[i] {
				t.Fatalf("Unexpected language code '%s' in request body. Expected '%s'", actualCodes.Strings[i], expectedCodes[i])
			}
		}

		response, err := xml.Marshal(newXMLArrayOfStrings(expectedNames))
		if err != nil {
			t.Fatalf("Unexpected error marshalling xml repsonse: %s", err.Error())
		}

		w.Header().Set("Content-Type", "text/xml")

		fmt.Fprint(w, string(response))
		return
	}))
	defer server.Close()

	router := newMockRouter()
	router.languageNamesURL = server.URL

	languageProvider := &languageProvider{
		router:     router,
		httpClient: newHTTPClient(authenticator),
	}

	actualNames, err := languageProvider.Names(expectedCodes)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}

	if len(actualNames) != len(expectedNames) {
		t.Fatalf("Unexpected number of languages codes: %q", actualNames)
	}

	for i := range expectedNames {
		if actualNames[i] != expectedNames[i] {
			t.Fatalf("Unexpected language code '%s'. Expected '%s'", actualNames[i], expectedNames[i])
		}
	}
}

func TestLanguageCatalogLanguages(t *testing.T) {
	expectedCodes := []string{"en", "de", "es", "ru", "jp"}
	expectedNames := []string{"English", "German", "Spanish", "Russian", "Japanese"}

	// intantiate language catalog and inject mocked out language provider
	languageProvider := newMockLanguageProvider()
	languageProvider.codes = expectedCodes
	languageProvider.names = expectedNames
	languageCatalog := newLanguageCatalog(languageProvider)

	// retrieve languages from catalog 3 times
	// make sure the catalog caches languages, i.e. it sends exactly one request to the language provider methods
	for _ = range make([]int, 3) {
		languages, err := languageCatalog.Languages()
		if err != nil {
			t.Fatalf("Unexpected error: %s", err.Error())
		}

		if languageProvider.callCounter["Codes"] != 1 {
			t.Fatalf("LanguagesProvider.Codes should have been called exactly once not %d times.", languageProvider.callCounter["Codes"])
		}

		if languageProvider.callCounter["Names"] != 1 {
			t.Fatalf("LanguagesProvider.Names should have been called exactly once not %d times.", languageProvider.callCounter["Names"])
		}

		if len(languages) != len(expectedCodes) {
			t.Fatalf("Unexpected number of languages: %q", languages)
		}

		for i := range expectedCodes {
			if languages[i].Code != expectedCodes[i] {
				t.Fatalf("Unexpected language code '%s'. Expected '%s'", languages[i].Code, expectedCodes[i])
			}

			if languages[i].Name != expectedNames[i] {
				t.Fatalf("Unexpected language code '%s'. Expected '%s'", languages[i].Name, expectedNames[i])
			}
		}
	}
}

func (a *authenticator) expectedAuthToken(t *testing.T) string {
	token, err := a.authToken()
	if err != nil {
		t.Fatalf("Unexpected error getting authToken from authenticator: %s", err.Error())
	}
	return token
}
