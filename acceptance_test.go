package translator_test

import (
	"os"
	"testing"

	"github.com/st3v/translator"
	"github.com/st3v/translator/google"
	"github.com/st3v/translator/microsoft"
)

func TestAcceptanceGoogleTranslate(t *testing.T) {
	testTranslate(t, googleTranslator(t))
}

func TestAcceptanceGoogleDetect(t *testing.T) {
	testDetect(t, googleTranslator(t))
}

func TestAcceptanceGoogleLanguages(t *testing.T) {
	testLanguages(t, googleTranslator(t))
}

func TestAcceptanceMicrosoftTranslate(t *testing.T) {
	testTranslate(t, microsoftTranslator(t))
}

func TestAcceptanceMicrosoftDetect(t *testing.T) {
	testDetect(t, microsoftTranslator(t))
}

func TestAcceptanceMicrosoftLanguages(t *testing.T) {
	testLanguages(t, microsoftTranslator(t))
}

func googleTranslator(t *testing.T) translator.Translator {
	key := os.Getenv("GOOGLE_API_KEY")

	if key == "" {
		t.Skip("Skipping acceptance tests for Google. Set environment variable GOOGLE_API_KEY.")
	}

	return google.NewTranslator(key)
}

func microsoftTranslator(t *testing.T) translator.Translator {
	key := os.Getenv("MS_SUBSCRIPTION_KEY")

	if key == "" {
		t.Skip("Skipping acceptance tests for Microsoft. Set environment variables MS_SUBSCRIPTION_KEY.")
	}

	return microsoft.NewTranslator(key)
}

func testTranslate(t *testing.T, translator translator.Translator) {
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

func testDetect(t *testing.T, translator translator.Translator) {
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

func testLanguages(t *testing.T, translator translator.Translator) {
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
