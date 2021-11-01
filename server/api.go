package server

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"

	"gitlab.com/mirukakoro/sisikyo/events/api"
)

var apiURL string

func init() {
	flag.StringVar(&apiURL, "api-url", "", "URL of API to use")
}

func setupAPI() error {
	if apiURL == "" {
		return errors.New("api-url: cannot be blank")
	}
	var err error
	api.DefaultBaseURL, err = url.Parse(apiURL)
	if err != nil {
		return fmt.Errorf("api-url: %w", err)
	}

	cl := api.DefaultClient()
	resp := api.VersionResp{}
	cl.Do(api.VersionReq{}, &resp)
	log.Printf("api: conn'd (version: %s)", resp.Version)
	return nil
}
