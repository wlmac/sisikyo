package api

import (
	"net/http"
	"net/url"
)

var userURL, _ = url.Parse("user/")

type UserReq struct {
	Username string
}

func (req UserReq) Req(c *Client) (*http.Request, error) {
	username, err := url.Parse(url.PathEscape(req.Username))
	if err != nil {
		return nil, err
	}
	return http.NewRequest(http.MethodGet, c.baseURL.ResolveReference(userURL).ResolveReference(username).String(), nil)
}

type UserResp struct {
	Username       string        `json:"username"`
	FirstName      string        `json:"first_name"`
	LastName       string        `json:"last_name"`
	Bio            string        `json:"bio"`
	Timezone       string        `json:"timezone"`
	GraduatingYear int           `json:"graduating_year"`
	Organizations  []string      `json:"organizations"`
	TagsFollowing  []interface{} `json:"tags_following"`
}

func (u UserResp) Name() string {
	return u.FirstName + " " + u.LastName
}
