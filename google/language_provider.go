package google

import (
	"fmt"
	"net/url"

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

type detectionPayload struct {
	Data struct {
		Detections [][]struct {
			Language string
		}
	}
}

type languageProvider interface {
	languages() ([]translator.Language, error)
	detect(text string) (string, error)
}

type concreteLanguageProvider struct {
	router        *router
	authenticator http.Authenticator
	catalog       []translator.Language
}

func newLanguageProvider(a http.Authenticator, r *router) *concreteLanguageProvider {
	return &concreteLanguageProvider{
		router:        r,
		authenticator: a,
		catalog:       nil,
	}
}

func (p *concreteLanguageProvider) languages() ([]translator.Language, error) {
	if p.catalog == nil {
		httpClient := http.NewClient(p.authenticator)

		resp, err := httpClient.SendRequest(
			"GET",
			fmt.Sprintf("%s?target=en", p.router.languagesURL()),
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

		payload, ok := result.(*languagesPayload)
		if !ok {
			return nil, tracerr.Error("Invalid response.")
		}

		p.catalog = make([]translator.Language, len(payload.Data.Languages))
		for i, l := range payload.Data.Languages {
			p.catalog[i] = translator.Language{
				Code: l.Language,
				Name: l.Name,
			}
		}
	}

	return p.catalog, nil
}

func (p *concreteLanguageProvider) detect(text string) (string, error) {
	httpClient := http.NewClient(p.authenticator)

	resp, err := httpClient.SendRequest(
		"GET",
		fmt.Sprintf("%s?q=%s", p.router.detectURL(), url.QueryEscape(text)),
		nil,
		"text/plain",
	)

	if err != nil {
		return "", tracerr.Wrap(err)
	}

	result, err := parseResponse(resp, &detectionPayload{})
	if err != nil {
		return "", tracerr.Wrap(err)
	}

	payload, ok := result.(*detectionPayload)
	if !ok {
		return "", tracerr.Error("Invalid response.")
	}

	return payload.Data.Detections[0][0].Language, nil
}
