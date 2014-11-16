package microsoft

import (
	"encoding/xml"
	"io/ioutil"
	"strings"

	"github.com/st3v/translator"
)

type LanguageProvider interface {
	Codes() ([]string, error)
	Names(codes []string) ([]string, error)
	Languages() ([]translator.Language, error)
}

type languageProvider struct {
	router     Router
	httpClient HttpClient
	languages  []translator.Language
}

func newLanguageProvider(authenticator Authenticator) LanguageProvider {
	return &languageProvider{
		router:     newRouter(),
		httpClient: newHttpClient(authenticator),
	}
}

func (p *languageProvider) Languages() ([]translator.Language, error) {
	if p.languages == nil {
		codes, err := p.Codes()
		if err != nil {
			return nil, err
		}

		names, err := p.Names(codes)
		if err != nil {
			return nil, err
		}

		for i := range codes {
			p.languages = append(
				p.languages,
				translator.Language{
					Code: codes[i],
					Name: names[i],
				})
		}
	}
	return p.languages, nil
}

func (p *languageProvider) Names(codes []string) ([]string, error) {
	payload, _ := xml.Marshal(newArrayOfStrings(codes))
	uri := p.router.LanguageNamesUrl() + "?locale=en"

	response, err := p.httpClient.SendRequest("POST", uri, strings.NewReader(string(payload)), "text/xml")
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}

	result := &arrayOfStrings{}
	if err := xml.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Strings, nil
}

func (p *languageProvider) Codes() ([]string, error) {
	response, err := p.httpClient.SendRequest("GET", p.router.LanguageCodesUrl(), nil, "text/plain")
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}

	result := &arrayOfStrings{}
	if err = xml.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Strings, nil
}
