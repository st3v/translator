package microsoft

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/st3v/tracerr"
	"github.com/st3v/translator/http"
	"io/ioutil"
	"net/url"
)

// The TranslationProvider communicates with Microsoft's
// API to provide a translation for a given text.
type TranslationProvider interface {
	Translate(text, from, to, version string) (string, error)
	Detect(text, version string) (string, error)
}

type translationProvider struct {
	router     Router
	httpClient http.Client
}

type Request []struct {
	Text string `json:"Text"`
}

type Translate []struct {
	Translations []Translations `json:"translations"`
}
type Translations struct {
	Text string `json:"text"`
	To   string `json:"to"`
}

type Detect []struct {
	Language                   string         `json:"language"`
	Score                      float64        `json:"score"`
	IsTranslationSupported     bool           `json:"isTranslationSupported"`
	IsTransliterationSupported bool           `json:"isTransliterationSupported"`
	Alternatives               []Alternatives `json:"alternatives"`
}

type Alternatives struct {
	Language                   string  `json:"language"`
	Score                      float64 `json:"score"`
	IsTranslationSupported     bool    `json:"isTranslationSupported"`
	IsTransliterationSupported bool    `json:"isTransliterationSupported"`
}

type Errors struct {
	Error Error `json:"error"`
}
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func newTranslationProvider(authenticator http.Authenticator, router Router) TranslationProvider {
	return &translationProvider{
		router:     router,
		httpClient: http.NewClient(authenticator),
	}
}

func (p *translationProvider) Translate(text, from, to, version string) (string, error) {
	apiVer := version

	if apiVer == "" {
		apiVer = p.router.ApiVersion()
	}
	request := make(Request, 1)
	request[0].Text = text
	b, _ := json.Marshal(&request)
	uri := fmt.Sprintf(
		"%s?api-version=%s&from=%s&to=%s&textType=html",
		p.router.TranslationURL(),
		apiVer,
		url.QueryEscape(from),
		url.QueryEscape(to))

	response, err := p.httpClient.SendRequest("POST", uri, bytes.NewBuffer(b), "application/json")
	if err != nil {
		return "", tracerr.Wrap(err)
	}

	body, err := ioutil.ReadAll(response.Body)

	defer response.Body.Close()
	if err != nil {
		return "", tracerr.Wrap(err)
	}

	errMsg := Errors{}
	err = json.Unmarshal(body, &errMsg)
	translation := Translate{}

	if err != nil {
		err = json.Unmarshal(body, &translation)
		if err != nil {
			return "", tracerr.Wrap(err)
		}
	} else {
		return "", errors.New(errMsg.Error.Message)
	}

	return translation[0].Translations[0].Text, nil
}

func (p *translationProvider) Detect(text, version string) (string, error) {
	apiVer := version

	if apiVer == "" {
		apiVer = p.router.ApiVersion()
	}
	request := make(Request, 1)
	request[0].Text = text
	b, _ := json.Marshal(&request)
	uri := fmt.Sprintf(
		"%s?api-version=%s",
		p.router.DetectURL(),
		apiVer)

	response, err := p.httpClient.SendRequest("POST", uri, bytes.NewBuffer(b), "application/json")
	if err != nil {
		return "", tracerr.Wrap(err)
	}

	body, err := ioutil.ReadAll(response.Body)

	defer response.Body.Close()
	if err != nil {
		return "", tracerr.Wrap(err)
	}

	detect := Detect{}
	err = json.Unmarshal(body, &detect)
	if err != nil {
		return "", tracerr.Wrap(err)
	}

	return detect[0].Language, nil
}
