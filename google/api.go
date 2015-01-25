package google

import "github.com/st3v/translator"

type api struct {
	lp languageProvider
	tp translationProvider
}

// NewTranslator instantiates a new Translator for Google's Translate API.
func NewTranslator(apiKey string) translator.Translator {
	authenticator := newAuthenticator(apiKey)
	router := newRouter()

	return &api{
		lp: newLanguageProvider(authenticator, router),
		tp: newTranslationProvider(authenticator, router),
	}
}

func (a *api) Languages() ([]translator.Language, error) {
	return a.lp.languages()
}

func (a *api) Detect(text string) (string, error) {
	return a.lp.detect(text)
}

func (a *api) Translate(text, from, to string) (string, error) {
	return a.tp.translate(text, from, to)
}
