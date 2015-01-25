package google

import (
	"fmt"
	"net/url"

	"github.com/st3v/tracerr"
	"github.com/st3v/translator/http"
)

type translationPayload struct {
	Data struct {
		Translations []struct {
			TranslatedText string
		}
	}
}

type translationProvider interface {
	translate(text, from, to string) (string, error)
}

type concreteTranslationProvider struct {
	authenticator http.Authenticator
	router        *router
}

func newTranslationProvider(a http.Authenticator, r *router) *concreteTranslationProvider {
	return &concreteTranslationProvider{
		authenticator: a,
		router:        r,
	}
}

func (t *concreteTranslationProvider) translate(text, from, to string) (string, error) {
	httpClient := http.NewClient(t.authenticator)

	uri := fmt.Sprintf(
		"%s?q=%s&source=%s&target=%s",
		t.router.translateURL(),
		url.QueryEscape(text),
		url.QueryEscape(from),
		url.QueryEscape(to))

	resp, err := httpClient.SendRequest("GET", uri, nil, "text/plain")
	if err != nil {
		return "", tracerr.Wrap(err)
	}

	result, err := parseResponse(resp, &translationPayload{})
	if err != nil {
		return "", tracerr.Wrap(err)
	}

	payload, ok := result.(*translationPayload)
	if !ok {
		return "", tracerr.Error("Invalid response.")
	}

	return payload.Data.Translations[0].TranslatedText, nil
}
