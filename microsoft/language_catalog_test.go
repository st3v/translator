package microsoft

import "testing"

func TestLanguageCatalogLanguages(t *testing.T) {
	expectedCodes := []string{"en", "de", "es", "ru", "jp"}
	expectedNames := []string{"English", "German", "Spanish", "Russian", "Japanese"}

	// instantiate language catalog and inject mocked out language provider
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
