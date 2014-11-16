package microsoft

import (
	"encoding/xml"
	"io/ioutil"
	"strings"

	"github.com/st3v/translator"
)

type LanguageCatalog interface {
	Languages() ([]translator.Language, error)
}

type LanguageProvider interface {
	Codes() ([]string, error)
	Names(codes []string) ([]string, error)
}

type languageCatalog struct {
	provider  LanguageProvider
	languages []translator.Language
}

type languageProvider struct {
	router     Router
	httpClient HttpClient
}

func newLanguageCatalog(provider LanguageProvider) LanguageCatalog {
	return &languageCatalog{
		provider: provider,
	}
}

func newLanguageProvider(authenticator Authenticator) LanguageProvider {
	return &languageProvider{
		router:     newRouter(),
		httpClient: newHttpClient(authenticator),
	}
}

func (c *languageCatalog) Languages() ([]translator.Language, error) {
	if c.languages == nil {
		codes, err := c.provider.Codes()
		if err != nil {
			return nil, err
		}

		names, err := c.provider.Names(codes)
		if err != nil {
			return nil, err
		}

		for i := range codes {
			c.languages = append(
				c.languages,
				translator.Language{
					Code: codes[i],
					Name: names[i],
				})
		}
	}
	return c.languages, nil
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
