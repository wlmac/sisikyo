package oauth

// TODO: use golang.org/x/oauth2

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"gitlab.com/mirukakoro/sisikyo/events/api"
	"golang.org/x/oauth2"
)

type OAuth struct {
	config oauth2.Config
}

func NewOAuth(config oauth2.Config) *OAuth {
	return &OAuth{
		config: config,
	}
}

type RedirectParams struct {
	State string `form:"state"`
	Code  string `form:"code"`

	Error     string `form:"error"`
	ErrorDesc string `form:"error_description"`
	ErrorURI  string `form:"error_uri"`
}

func (r RedirectParams) IsValid() bool {
	return r.Code != "" || r.Error != ""
}

func (r RedirectParams) IsError() bool {
	return r.Error != ""
}

func Redirect(c *gin.Context) {
	params := RedirectParams{}
	err := c.Bind(&params)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html.tmpl", fmt.Sprintf("bad params: %s", err.Error()))
		return
	}
	if !params.IsValid() {
		c.HTML(http.StatusBadRequest, "error.html.tmpl", "invalid params")
		return
	}
	if params.IsError() {
		c.HTML(http.StatusOK, "error.html.tmpl", gin.H{
			"name": params.Error,
			"desc": params.ErrorDesc,
			"url":  params.ErrorURI,
		})
		return
	}
	//c.HTML(http.StatusOK, "ok.html.tmpl", params.Code)
}

func RedirectJSON(c *gin.Context) {
	params := RedirectParams{}
	err := c.Bind(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("bad params: %s", err.Error()))
		return
	}
	if !params.IsValid() {
		c.JSON(http.StatusBadRequest, "invalid params")
		return
	}
	if params.IsError() {
		c.JSON(http.StatusInternalServerError, gin.H{
			"name": params.Error,
			"desc": params.ErrorDesc,
		})
		return
	}
	c.JSON(http.StatusOK, nil)
}

var authorizeURL, _ = url.Parse("o/authorize")

func Authorize(c *gin.Context) {
	clientID := "aD87ahyBviMiz4kkswBIrTyp7sjbhN7paYxXF0kf"
	u := api.DefaultBaseURL.ResolveReference(authorizeURL)
	q := u.Query()
	q.Set("response_type", "code")
	q.Set("client_id", clientID)
	q.Set("scope", "me_schedule")
	state, err := genRandom(64)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	q.Set("state", base64.URLEncoding.EncodeToString(state))
	u.RawQuery = q.Encode()
	c.String(http.StatusOK, u.String())
}

func Authorize2(c *api.Client) (string, error) {
	authURL, _ := url.Parse("/o/authorize")
	tokenURL, _ := url.Parse("/o/token")
	cfg := oauth2.Config{
		ClientID:     "aD87ahyBviMiz4kkswBIrTyp7sjbhN7paYxXF0kf",
		ClientSecret: "MwJXfcdjcKELGTJA8Ak6Wg7Ty5cWv4TSMLiVfiIP4aBgy0JMzJTE5R4dKUxxuwC9H5WgVzfjxpF7UfsOe4qrXj3EsJnQZsfPmLQCqUdpa0UgTFbni2PAnAabBlKV8lyF",
		Endpoint: oauth2.Endpoint{
			AuthURL:   c.BaseURL().ResolveReference(authURL).String(),
			TokenURL:  c.BaseURL().ResolveReference(tokenURL).String(),
			AuthStyle: oauth2.AuthStyleInHeader,
		},
		Scopes: []string{"me_meta", "me_schedule", "me_timetable"},
	}
	state, err := genRandom(64)
	if err != nil {
		return "", err
	}
	return cfg.AuthCodeURL(string(state)), nil
}
