package microsoft

import "github.com/st3v/translator"

type api struct {
	router              Router
	languageProvider    LanguageProvider
	translationProvider TranslationProvider
}

func NewTranslator(clientId, clientSecret string) translator.Translator {
	authenticator := newAuthenticator(clientId, clientSecret)
	return &api{
		languageProvider:    newLanguageProvider(authenticator),
		translationProvider: newTranslationProvider(authenticator),
	}
}

func (a *api) Translate(text, from, to string) (string, error) {
	return a.translationProvider.Translate(text, from, to)
}

func (a *api) Languages() ([]translator.Language, error) {
	return a.languageProvider.Languages()
}
