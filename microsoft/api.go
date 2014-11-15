package microsoft

import "github.com/st3v/translator"

type api struct {
	authenticator Authenticator
}

func NewTranslator(clientId, clientSecret string) translator.Translator {
	return &api{
		authenticator: NewAuthenticator(clientId, clientSecret),
	}
}

func (a *api) Languages() ([]translator.Language, error) {
	return make([]translator.Language, 0), nil
}

func (a *api) Translate(text, from, to string) (string, error) {
	return "", nil
}
