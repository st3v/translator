package microsoft

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	_http "github.com/st3v/translator/http"
)

func TestLanguageProviderCodes(t *testing.T) {
	expectedCodes := []string{"en", "de", "es", "ru", "jp"}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("Unexpected request method: %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "text/plain" {
			t.Fatalf("Unexpected content type in request header: %s", r.Header.Get("Content-Type"))
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
		httpClient: _http.NewAuthenticatedClient(),
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

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("Unexpected request method: %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "text/xml" {
			t.Fatalf("Unexpected content type in request header: %s", r.Header.Get("Content-Type"))
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
		httpClient: _http.NewAuthenticatedClient(),
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

func newMockLanguageProvider() *mockLanguageProvider {
	return &mockLanguageProvider{
		callCounter: make(map[string]int),
		codes:       make([]string, 0),
		names:       make([]string, 0),
	}
}

type mockLanguageProvider struct {
	callCounter map[string]int
	codes       []string
	names       []string
}

func (p *mockLanguageProvider) Codes() ([]string, error) {
	p.callCounter["Codes"]++
	return p.codes, nil
}

func (p *mockLanguageProvider) Names(codes []string) ([]string, error) {
	p.callCounter["Names"]++
	return p.names, nil
}
