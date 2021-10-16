package oauth

// TODO: use golang.org/x/oauth2

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"gitlab.com/mirukakoro/sisikyo/events/api"
)

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

func genRandom(l int) ([]byte, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	result := make([]byte, l)
	for i := 0; i < l; i++ {
		j, err := rand.Int(rand.Reader, new(big.Int).SetInt64(int64(len(letters))))
		if err != nil {
			return nil, err
		}
		result[i] = letters[j.Int64()]
	}
	return result, nil
}
