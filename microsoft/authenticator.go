package microsoft

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/st3v/tracerr"
)

const scope = "http://api.microsofttranslator.com"

// Authenticator is used to authenicate HTTP requests to Microsoft's
// Translation API.
type Authenticator interface {
	// Authenticate a given HTTP request.
	Authenticate(request *http.Request) error
}

type authenticator struct {
	provider        AuthenticationProvider
	accessTokenChan chan *accessToken
}

func newAuthenticator(clientID, clientSecret string) Authenticator {
	// make buffered accessToken channel and pre-fill it with an expired token
	tokenChan := make(chan *accessToken, 1)
	tokenChan <- &accessToken{}

	// return new authenticator that uses the above accessToken channel
	return &authenticator{
		provider:        newAuthenticationProvider(clientID, clientSecret),
		accessTokenChan: tokenChan,
	}
}

func (a *authenticator) Authenticate(request *http.Request) error {
	authToken, err := a.authToken()
	if err != nil {
		return tracerr.Wrap(err)
	}

	request.Header.Add("Authorization", authToken)
	return nil
}

func (a *authenticator) authToken() (string, error) {
	// grab the token
	accessToken := <-a.accessTokenChan

	// make sure it's valid, otherwise request a new one
	if accessToken == nil || accessToken.expired() {
		err := a.provider.RefreshAccessToken(accessToken)
		if err != nil || accessToken == nil {
			a.accessTokenChan <- nil
			return "", tracerr.Wrap(err)
		}
	}

	// put the token back on the channel
	a.accessTokenChan <- accessToken

	// return authToken
	return "Bearer " + accessToken.Token, nil
}

type accessToken struct {
	Token     string `json:"access_token"`
	Type      string `json:"token_type"`
	Scope     string `json:"scope"`
	ExpiresIn string `json:"expires_in"`
	ExpiresAt time.Time
}

func (t *accessToken) expired() bool {
	// be conservative and expire 10 seconds early
	return t.ExpiresAt.Before(time.Now().Add(time.Second * 10))
}

// The AuthenticationProvider is used to refresh access tokens for
// Microsoft's API endpoints.
type AuthenticationProvider interface {
	RefreshAccessToken(*accessToken) error
}

type authenticationProvider struct {
	clientID     string
	clientSecret string
	router       Router
}

func newAuthenticationProvider(clientID, clientSecret string) AuthenticationProvider {
	return &authenticationProvider{
		clientID:     clientID,
		clientSecret: clientSecret,
		router:       newRouter(),
	}
}

func (p *authenticationProvider) RefreshAccessToken(token *accessToken) error {
	values := make(url.Values)
	values.Set("client_id", p.clientID)
	values.Set("client_secret", p.clientSecret)
	values.Set("scope", scope)
	values.Set("grant_type", "client_credentials")

	response, err := http.PostForm(p.router.AuthURL(), values)
	if err != nil {
		return tracerr.Wrap(err)
	}

	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return tracerr.Wrap(err)
	}

	if err := json.Unmarshal(body, token); err != nil {
		return tracerr.Wrap(err)
	}

	expiresInSeconds, err := strconv.Atoi(token.ExpiresIn)
	if err != nil {
		return tracerr.Wrap(err)
	}

	token.ExpiresAt = time.Now().Add(time.Duration(expiresInSeconds) * time.Second)

	return nil
}
