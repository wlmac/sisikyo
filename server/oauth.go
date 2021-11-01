package server

import (
	"errors"
	"flag"
)

var oauthClientID string
var oauthClientSecret string

const oauthEndpoint = "/o/authorize"
const redirectURL = "/o/redirect"

func init() {
	flag.StringVar(&oauthClientID, "o-id", "", "oauth: client ID")
	flag.StringVar(&oauthClientSecret, "o-secret", "", "oauth: client secret")
}

func setupOauth() error {
	if oauthClientID == "" {
		return errors.New("oauth id: must not be blank")
	}
	if oauthClientSecret == "" {
		return errors.New("oauth secret: must not be blank")
	}
}
