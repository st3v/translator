package auth

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/st3v/tracerr"
)

// The AccessTokenProvider handles access tokens for Microsoft's API endpoints.
type AccessTokenProvider interface {
	RefreshToken(*accessToken) error
}

type accessTokenProvider struct {
	clientID     string
	clientSecret string
	authURL      string
}

func newAccessTokenProvider(clientID, clientSecret, authURL string) AccessTokenProvider {
	return &accessTokenProvider{
		clientID:     clientID,
		clientSecret: clientSecret,
		authURL:      authURL,
	}
}

func (p *accessTokenProvider) RefreshToken(token *accessToken) error {
	req, err := http.NewRequest("POST", p.authURL, nil)
	if err != nil {
		return tracerr.Wrap(err)
	}
	req.Header.Add("Ocp-Apim-Subscription-Key", p.clientSecret)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return tracerr.Wrap(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return tracerr.Errorf("Unexpected Status: %s", response.Status)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return tracerr.Wrap(err)
	}

	token.Token = string(body)

	token.ExpiresAt = time.Now().Add(8 * time.Minute)

	return nil
}
