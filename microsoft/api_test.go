package microsoft

import "testing"

func TestTranslate(t *testing.T) {
	api := NewTranslator("", "")

	original := "dog"
	translation, err := api.Translate(original, "en", "de")

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if translation != "" {
		t.Errorf("Unexpected transaltion: %s", translation)
	}
}

func TestLanguages(t *testing.T) {
	api := NewTranslator("", "")

	languages, err := api.Languages()

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if len(languages) != 0 {
		t.Errorf("Unexpected number of languages: %d", len(languages))
	}
}
