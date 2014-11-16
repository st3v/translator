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
		t.Errorf("Unexpected translation: %s", translation)
	}
}
