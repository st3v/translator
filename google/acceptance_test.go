package google

import (
	"os"
	"testing"
)

func apiKey(t *testing.T) string {
	apiKey := os.Getenv("GOOGLE_API_KEY")

	if apiKey == "" {
		t.Skip("Skipping acceptance tests for Google. Set environment variable GOOGLE_API_KEY.")
	}

	return apiKey
}

func TestTranslateAcceptance(t *testing.T) {
	translator := NewTranslator(apiKey(t))
	translation, err := translator.Translate("Hello World!", "en", "de")

	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	expectedTranslation := "Hallo Welt!"

	if translation != expectedTranslation {
		t.Errorf(
			"Unexpected translation. Got: '%s'. Want: '%s'.",
			translation,
			expectedTranslation,
		)
	}
}

func TestDetectAcceptance(t *testing.T) {
	translator := NewTranslator(apiKey(t))
	languageCode, err := translator.Detect("¿cómo está?")

	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	expectedCode := "es"

	if languageCode != expectedCode {
		t.Errorf(
			"Unexpected language detected. Got: %s. Want: %s.",
			languageCode,
			expectedCode,
		)
	}

}

func TestLanguagesAcceptance(t *testing.T) {
	translator := NewTranslator(apiKey(t))
	languages, err := translator.Languages()

	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	if len(languages) == 0 {
		t.Error("Expected some languages but got none.")
	}

	expectedLanguages := []struct{ Code, Name string }{
		{"en", "English"},
		{"de", "German"},
		{"fr", "French"},
		{"es", "Spanish"},
		{"it", "Italian"},
		{"pt", "Portuguese"},
		{"ja", "Japanese"},
		{"ko", "Korean"},
		{"zh", "Chinese (Simplified)"},
		{"zh-TW", "Chinese (Traditional)"},
		{"ru", "Russian"},
	}

	for _, actual := range languages {
		for i, expected := range expectedLanguages {
			if actual.Code == expected.Code && actual.Name == expected.Name {
				expectedLanguages = append(expectedLanguages[:i], expectedLanguages[i+1:]...)
				break
			}
		}
	}

	if len(expectedLanguages) != 0 {
		t.Errorf("Languages not found: %v\nGot: %v", expectedLanguages, languages)
	}
}
