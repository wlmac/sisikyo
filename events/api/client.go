package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)

type Client struct {
	client  HTTPClient
	baseURL *url.URL
	noCopy  [0]sync.Mutex
	auth    Auth
}

var DefaultBaseURL, _ = url.Parse("https://maclyonsden.com/api/")

func DefaultClient() *Client {
	return &Client{client: http.DefaultClient, baseURL: DefaultBaseURL}
}

func NewClient(client HTTPClient, baseURL *url.URL) *Client {
	return &Client{client: client, baseURL: baseURL}
}

func (c *Client) Rel(u *url.URL) *url.URL {
	if u.IsAbs() {
		return u
	}
	return c.baseURL.ResolveReference(u)
}

func (c *Client) HTTPClient() HTTPClient {
	return c.client
}

func (c *Client) BaseURL() *url.URL {
	return c.baseURL
}

func (c *Client) Do(req Req, v interface{}) (err error) {
	request, err := req.Req(c)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept-Encoding", "application/json")
	resp, err := c.client.Do(request)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err2 := Body.Close()
		if err2 != nil {
			err = err2
		}
	}(resp.Body)
	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		return errors.New(resp.Status)
	}
	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return fmt.Errorf("json: %s", err)
	}
	return
}
