package microsoft

import "github.com/st3v/translator"

type api struct {
	router              Router
	languageCatalog     LanguageCatalog
	translationProvider TranslationProvider
}

// NewTranslator returns a struct that implements the Translator
// interface by exposing a Translate and a Languages function that
// are backed by Microsoft's translation API.
// The function takes the clientID and clientSecret for an existing
// app registered in Microsoft's Azure DataMarket.
func NewTranslator(clientID, clientSecret string) translator.Translator {
	authenticator := newAuthenticator(clientID, clientSecret)
	return &api{
		languageCatalog:     newLanguageCatalog(newLanguageProvider(authenticator)),
		translationProvider: newTranslationProvider(authenticator),
	}
}

func (a *api) Translate(text, from, to string) (string, error) {
	return a.translationProvider.Translate(text, from, to)
}

func (a *api) Languages() ([]translator.Language, error) {
	return a.languageCatalog.Languages()
}
