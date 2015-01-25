package google

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/st3v/tracerr"
	"github.com/st3v/translator"
	"github.com/st3v/translator/http"
)

type languagesPayload struct {
	Data struct {
		Languages []struct {
			Language string
			Name     string
		}
	}
}

// LanguageProvider retrieves all languages supported by Google's Translate API.
type LanguageProvider interface {
	Languages() ([]translator.Language, error)
}

type languageProvider struct {
	router        Router
	authenticator http.Authenticator
	catalog       []translator.Language
}

func newLanguageProvider(a http.Authenticator, r Router) LanguageProvider {
	return &languageProvider{
		router:        r,
		authenticator: a,
		catalog:       nil,
	}
}

func (p *languageProvider) Languages() ([]translator.Language, error) {
	if p.catalog == nil {
		httpClient := http.NewClient(p.authenticator)

		resp, err := httpClient.SendRequest(
			"GET",
			fmt.Sprintf("%s?target=en", p.router.LanguagesURL()),
			nil,
			"text/plain",
		)

		if err != nil {
			return nil, tracerr.Wrap(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return nil, tracerr.Wrap(err)
		}

		result := &languagesPayload{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}

		p.catalog = make([]translator.Language, len(result.Data.Languages))
		for i, l := range result.Data.Languages {
			p.catalog[i] = translator.Language{
				Code: l.Language,
				Name: l.Name,
			}
		}
	}

	return p.catalog, nil
}
