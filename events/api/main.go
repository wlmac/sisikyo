// This file is part of a program named Sisikyō or Sisikyo.
//
// Copyright (C) 2019 Ken Shibata <kenxshibata@gmail.com>
//
// License as published by the Free Software Foundation, either version 1 of the License, or (at your option) any later
// version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied
// warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see
// <https://www.gnu.org/licenses/>.

// Package api provides abstractions for the REST API.
package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

// HTTPClient is a subset of the *http.Client interface.
type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
}

type Client struct {
	client  HTTPClient
	baseURL *url.URL
}

var DefaultBaseURL, _ = url.Parse("https://maclyonsden.com/api/")

func DefaultClient() *Client {
	return &Client{client: http.DefaultClient, baseURL: DefaultBaseURL}
}

func NewClient(client HTTPClient, baseURL *url.URL) *Client {
	return &Client{client: client, baseURL: baseURL}
}

func (c *Client) do(url string, v interface{}) (err error) {
	resp, err := c.client.Get(url)
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
		return err
	}
	return nil
}

func (c *Client) Events(req EventsReq) (resp EventsResp, err error) {
	u := c.baseURL.ResolveReference(eventsURL)
	u.RawQuery = req.Query().Encode()
	resp = EventsResp{}
	err = c.do(u.String(), &resp)
	if err != nil {
		return
	}
	return
}
