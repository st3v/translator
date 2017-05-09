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
	key     string
	authURL string
}

func newAccessTokenProvider(subscribtionKey, authURL string) AccessTokenProvider {
	return &accessTokenProvider{
		key:     subscribtionKey,
		authURL: authURL,
	}
}

func (p *accessTokenProvider) RefreshToken(token *accessToken) error {
	req, err := http.NewRequest("POST", p.authURL, nil)
	if err != nil {
		return tracerr.Wrap(err)
	}

	req.Header.Add("Ocp-Apim-Subscription-Key", p.key)

	client := new(http.Client)

	response, err := client.Do(req)
	if err != nil {
		return tracerr.Wrap(err)
	}

	if response.StatusCode != http.StatusOK {
		return tracerr.Errorf("Unexpected Status: %s", response.Status)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return tracerr.Wrap(err)
	}

	token.Token = string(body)
	token.ExpiresAt = time.Now().Add(10 * time.Minute)

	return nil
}
