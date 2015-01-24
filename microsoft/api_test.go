package microsoft

import (
	"testing"

	"github.com/st3v/translator"
)

func TestAPITranslate(t *testing.T) {
	original := "Mein Englisch ist unter aller Sau."
	expectedTranslation := "My English is under all pig."
	from := "de"
	to := "en"

	api := &api{
		translationProvider: newMockTranslationProvider(original, from, to, expectedTranslation, t),
	}

	actualTranslation, err := api.Translate(original, from, to)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	if actualTranslation != expectedTranslation {
		t.Errorf("Unexpected translation: %s", actualTranslation)
	}
}

func TestAPILanguages(t *testing.T) {
	expectedLanguages := []translator.Language{
		translator.Language{
			Code: "en",
			Name: "English",
		},
		translator.Language{
			Code: "de",
			Name: "German",
		},
		translator.Language{
			Code: "es",
			Name: "Spanish",
		},
	}

	api := &api{
		languageCatalog: &languageCatalog{
			languages: expectedLanguages,
		},
	}

	actualLanguages, err := api.Languages()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}

	if len(actualLanguages) != len(expectedLanguages) {
		t.Fatalf("Unexpected number of languages: %q", actualLanguages)
	}

	for i := range expectedLanguages {
		if actualLanguages[i].Code != expectedLanguages[i].Code {
			t.Fatalf("Unexpected language code '%s'. Expected '%s'", actualLanguages[i].Code, expectedLanguages[i].Code)
		}

		if actualLanguages[i].Name != expectedLanguages[i].Name {
			t.Fatalf("Unexpected language code '%s'. Expected '%s'", actualLanguages[i].Name, expectedLanguages[i].Name)
		}
	}
}

func TestAPIDetect(t *testing.T) {
	text := "Mein Englisch ist unter aller Sau."
	expectedLanguage := "de"
	from := "de"

	api := &api{
		translationProvider: newMockTranslationProvider(text, from, "", "", t),
	}

	actualLanguage, err := api.Detect(text)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	if actualLanguage != expectedLanguage {
		t.Errorf("api.Detect(`%v`) failed:\n\twant: %v\n\t got: %v",
			text, expectedLanguage, actualLanguage)
	}
}
