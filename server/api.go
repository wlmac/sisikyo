package server

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"gitlab.com/mirukakoro/sisikyo/events/api"
)

var apiURL string
var apiTimeout time.Duration

func init() {
	flag.StringVar(&apiURL, "api-url", "https://maclyonsden.com/api", "api: base URL")
	flag.DurationVar(&apiTimeout, "api-timeout", 1*time.Second, "api: response timeout")
}

func setupAPI() (*api.Client, error) {
	if apiURL == "" {
		return nil, errors.New("api-url: cannot be blank")
	}
	baseURL, err := url.Parse(apiURL)
	if err != nil {
		return nil, fmt.Errorf("api-url: %w", err)
	}

	// version check
	log.Printf("api version: compatible with %s", api.APIVersion)
	cl := api.NewClient(&http.Client{Timeout: apiTimeout}, baseURL)
	ver, ok, err := cl.CheckAPIVersion()
	if err != nil {
		return nil, fmt.Errorf("api version check: %w", err)
	}
	if ok {
		log.Printf("api version: compatible with server (version: %s)", ver)
	} else {
		log.Printf("api version: incompatible with server (version: %s)", ver)
	}
	return cl, nil
}
