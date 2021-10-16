package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type query interface {
	Query() url.Values
}

func queryReq(u2 *url.URL, q query, c *Client) (*http.Request, error) {
	var u *url.URL
	{
		u_ := *u2
		u = &u_
	}
	u.RawQuery = q.Query().Encode()
	return http.NewRequest(http.MethodGet, c.baseURL.ResolveReference(u).String(), nil)
}

func jsonEncode(v interface{}) (io.Reader, error) {
	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(v)
	if err != nil {
		return nil, err
	}
	return body, nil
}
