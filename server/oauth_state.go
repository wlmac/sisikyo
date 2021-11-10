package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// oauthStateName is the name of the cookie to store the OAuth state.
const oauthStateName = "oauth-state" // set __Host- prefix only if secure

// oauthClearState sets a Set-Cookie header to set a state cookie.
func oauthClearState(c *gin.Context) {
	// why not use `Max-Age: 0`?
	// 	see <https://www.rfc-editor.org/errata/eid3430> which is part of <https://stackoverflow.com/a/20320610>
	if gin.Mode() == "debug" {
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     oauthStateName,
			Value:    "",
			Path:     "/o",
			MaxAge:   1, // 1 second
			SameSite: http.SameSiteLaxMode,
		})
	} else {
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     oauthStateName,
			Value:    "",
			Path:     "/o",
			MaxAge:   1, // 1 second
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
	}
}

// oauthSetState sets a Set-Cookie header to set a state cookie.
func oauthSetState(c *gin.Context, state string) {
	if gin.Mode() == "debug" {
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     oauthStateName,
			Value:    state,
			Path:     "/o",
			SameSite: http.SameSiteLaxMode,
		})
	} else {
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     oauthStateName,
			Value:    state,
			Path:     "/o",
			MaxAge:   60, // 1 minute
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
	}
}

// oauthGetState gets the state stored in the cookie.
func oauthGetState(resp *http.Request) (state string, err error) {
	cookie, err := resp.Cookie(oauthStateName)
	if err != nil {
		return
	}
	state = cookie.Value
	return
}
