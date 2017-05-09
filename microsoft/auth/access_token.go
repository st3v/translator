package auth

import "time"

type accessToken struct {
	Token     string `json:"access_token"`
	ExpiresAt time.Time
}

func (t *accessToken) expired() bool {
	// be conservative and expire 10 seconds early
	return t.ExpiresAt.Before(time.Now().Add(time.Second * 10))
}
