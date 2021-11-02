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
	flag.StringVar(&apiURL, "api-url", "", "api: base URL")
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
	cl := api.NewClient(&http.Client{Timeout: apiTimeout}, baseURL)
	resp := api.VersionResp{}
	err = cl.Do(api.VersionReq{}, &resp)
	if err != nil {
		return nil, fmt.Errorf("api version: %w", err)
	}
	log.Printf("api: conn'd (version: %s)", resp.Version)
	return cl, nil
}
