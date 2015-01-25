package google

import (
	"fmt"

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

		result, err := parseResponse(resp, &languagesPayload{})
		if err != nil {
			return nil, tracerr.Wrap(err)
		}

		lr, ok := result.(*languagesPayload)
		if !ok {
			return nil, tracerr.Error("Invalid response.")
		}

		p.catalog = make([]translator.Language, len(lr.Data.Languages))
		for i, l := range lr.Data.Languages {
			p.catalog[i] = translator.Language{
				Code: l.Language,
				Name: l.Name,
			}
		}
	}

	return p.catalog, nil
}
