package microsoft

import (
	"testing"

	"github.com/st3v/translator"
)

// func TestTranslate(t *testing.T) {
// 	api := NewTranslator("", "")

// 	original := "dog"
// 	translation, err := api.Translate(original, "en", "de")

// 	if err != nil {
// 		t.Errorf("Unexpected error: %s", err)
// 	}

// 	if translation != "" {
// 		t.Errorf("Unexpected translation: %s", translation)
// 	}
// }

func TestApiLanguages(t *testing.T) {
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
		t.Fatalf("Unexpected error: %s", err)
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
