package server

import (
	"errors"
	"flag"
	"fmt"
	"net/url"

	"golang.org/x/oauth2"
)

var oauthBaseURL string
var oauthAuthURL string
var oauthTokenURL string
var oauthClientID string
var oauthClientSecret string

func init() {
	flag.StringVar(&oauthBaseURL, "o-url", "", "oauth: absolute base URL")
	flag.StringVar(&oauthAuthURL, "o-auth", "authorize", "oauth: relative auth URL")
	flag.StringVar(&oauthTokenURL, "o-token", "token", "oauth: relative redirect URL")
	flag.StringVar(&oauthClientID, "o-id", "changeme", "oauth: client ID")
	flag.StringVar(&oauthClientSecret, "o-secret", "changeme", "oauth: client secret")
}

func setupOauth() (*oauth2.Config, error) {
	if oauthBaseURL == "" {
		return nil, nil
	}
	baseURL, err := url.Parse(oauthBaseURL)
	if err != nil {
		return nil, fmt.Errorf("oauth url: %w", err)
	}
	authURL, err := url.Parse(oauthAuthURL)
	if err != nil {
		return nil, fmt.Errorf("oauth auth: %w", err)
	}
	tokenURL, err := url.Parse(oauthTokenURL)
	if err != nil {
		return nil, fmt.Errorf("oauth token: %w", err)
	}

	if oauthClientID == "" {
		return nil, errors.New("oauth id: must not be blank")
	}
	if oauthClientSecret == "" {
		return nil, errors.New("oauth secret: must not be blank")
	}
	return &oauth2.Config{
		ClientID:     oauthClientID,
		ClientSecret: oauthClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   baseURL.ResolveReference(authURL).String() + "/",
			TokenURL:  baseURL.ResolveReference(tokenURL).String() + "/",
			AuthStyle: oauth2.AuthStyleInHeader,
		},
		Scopes: []string{"me_schedule", "me_timetable"},
	}, nil
}
