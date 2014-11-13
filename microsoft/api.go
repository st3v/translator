package microsoft

import "github.com/st3v/translator"

type api struct {
	clientId     string
	clientSecret string
}

func NewTranslator(clientId, clientSecret string) translator.Translator {
	return &api{
		clientId:     clientId,
		clientSecret: clientSecret,
	}
}

func (a *api) Languages() ([]translator.Language, error) {
	return make([]translator.Language, 0), nil
}

func (a *api) Translate(text, from, to string) (string, error) {
	return "", nil
}
