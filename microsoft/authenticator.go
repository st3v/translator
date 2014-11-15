package microsoft

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	SCOPE    = "http://api.microsofttranslator.com"
	AUTH_URL = "https://datamarket.accesscontrol.windows.net/v2/OAuth2-13"
)

type Authenticator interface {
	Authenticate(request *http.Request) error
}

type authenticator struct {
	clientId     string
	clientSecret string
	accessToken  *accessToken
}

type accessToken struct {
	Token     string `json:"access_token"`
	Type      string `json:"token_type"`
	Scope     string `json:"scope"`
	ExpiresIn string `json:"expires_in"`
	ExpiresAt time.Time
}

func NewAuthenticator(clientId, clientSecret string) Authenticator {
	return &authenticator{
		clientId:     clientId,
		clientSecret: clientSecret,
	}
}

func (a *authenticator) Authenticate(request *http.Request) error {
	authToken, err := a.authToken()
	if err != nil {
		return err
	}

	request.Header.Add("Authorization", authToken)
	return nil
}

func (a *authenticator) authToken() (string, error) {
	if a.accessToken == nil || a.accessToken.expired() {
		if err := a.requestAccessToken(AUTH_URL); err != nil {
			return "", err
		}
	}
	return "Bearer " + a.accessToken.Token, nil
}

func (a *authenticator) requestAccessToken(authUrl string) error {
	values := make(url.Values)
	values.Set("client_id", a.clientId)
	values.Set("client_secret", a.clientSecret)
	values.Set("scope", SCOPE)
	values.Set("grant_type", "client_credentials")

	response, err := http.PostForm(authUrl, values)
	if err != nil {
		log.Println(err)
		return err
	}

	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		log.Println(err)
		return err
	}

	token := &accessToken{}
	if err := json.Unmarshal(body, &token); err != nil {
		log.Println(err)
		return err
	}

	expiresInSeconds, err := strconv.Atoi(token.ExpiresIn)
	if err != nil {
		log.Println(err)
		return err
	}

	token.ExpiresAt = time.Now().Add(time.Duration(expiresInSeconds) * time.Second)
	a.accessToken = token

	return nil
}

func (t *accessToken) expired() bool {
	// be conservative and expire 10 seconds early
	return t.ExpiresAt.Before(time.Now().Add(time.Second * 10))
}
