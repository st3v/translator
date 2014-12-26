package microsoft

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/st3v/tracerr"
)

type TranslationProvider interface {
	Translate(text, from, to string) (string, error)
}

type translationProvider struct {
	router     Router
	httpClient HttpClient
}

func newTranslationProvider(authenticator Authenticator) TranslationProvider {
	return &translationProvider{
		router:     newRouter(),
		httpClient: newHttpClient(authenticator),
	}
}

func (p *translationProvider) Translate(text, from, to string) (string, error) {
	uri := fmt.Sprintf(
		"%s?text=%s&from=%s&to=%s",
		p.router.TranslationUrl(),
		url.QueryEscape(text),
		url.QueryEscape(from),
		url.QueryEscape(to))

	response, err := p.httpClient.SendRequest("GET", uri, nil, "text/plain")
	if err != nil {
		return "", tracerr.Wrap(err)
	}

	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return "", tracerr.Wrap(err)
	}

	translation := &xmlString{}
	err = xml.Unmarshal(body, &translation)
	if err != nil {
		return "", tracerr.Wrap(err)
	}

	return translation.Value, nil
}
