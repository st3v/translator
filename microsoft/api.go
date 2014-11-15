package microsoft

import (
	"io"
	"net/http"

	"github.com/st3v/translator"
)

type api struct {
	authenticator Authenticator
	languages     []translator.Language
}

const (
	serviceUri          = "http://api.microsofttranslator.com/v2/Http.svc/"
	translationEndpoint = serviceUri + "Translate"
)

func NewTranslator(clientId, clientSecret string) translator.Translator {
	return &api{
		authenticator: NewAuthenticator(clientId, clientSecret),
	}
}

func (a *api) Translate(text, from, to string) (string, error) {
	return "", nil
}

func (a *api) sendRequest(method, uri string, body io.Reader, contentType string) (*http.Response, error) {
	request, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", contentType)

	err = a.authenticator.Authenticate(request)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
