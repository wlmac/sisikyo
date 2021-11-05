package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"gitlab.com/mirukakoro/sisikyo/db"
	"gitlab.com/mirukakoro/sisikyo/events/api"
	"gitlab.com/mirukakoro/sisikyo/oauth"
	"gitlab.com/mirukakoro/sisikyo/util"
)

// engineContext contains the context required for the server methods
type engineContext struct { // yes, great naming
	E   *gin.Engine
	API *api.Client
	O   *oauth.Client
	Db  *sqlx.DB
}

func (e *engineContext) PublicQuery(render renderFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		o := ICSOptions{}
		err := c.ShouldBindQuery(&o)
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
			c.String(http.StatusBadRequest, fmt.Sprint(err))
			return
		}

		evs, err := o.List(e.API)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("list: %s", err))
			return
		}

		render(c, evs)
	}
}

// controlParam is a param binding for Gin.
type controlParam struct {
	Control string `uri:"control" binding:"required"`
}

func (e *engineContext) UserRemove(c *gin.Context) {
	var params controlParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, "params get: failed")
		return
	}

	txx, err := e.Db.Beginx()
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	_, err = txx.NamedExec(db.UserRemove, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	err = txx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	c.Status(http.StatusOK)
}

func (e *engineContext) UserQuery(render renderFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params controlParam
		if err := c.ShouldBindUri(&params); err != nil {
			c.JSON(http.StatusBadRequest, "params get: failed")
			return
		}

		// not using a transaction here because:
		// 	since this is a single operation, so atomicity shouldn't matter
		var user db.User
		err := e.Db.Get(&user, db.UserQuery, params.Control)
		if err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprint(err))
			return
		}

		// get calendar from API
		//scheduleWeek := api.MeScheduleWeekResp{}
		customCl := api.NewClient(e.O.Client(context.Background(), user.Token()), e.API.BaseURL())
		evs, err := customCl.CourseEvents()
		//	err = customCl.Do(api.MeScheduleWeekReq{}, &scheduleWeek)
		if err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprint(err))
			return
		}
		render(c, evs)
	}
}

// OauthRedirect verifies the OAuth response from the OAuth server.
func (e *engineContext) OauthRedirect(c *gin.Context) {
	storedState, err := oauthGetState(c.Request)
	if err != nil {
		c.String(http.StatusBadRequest, "state get: failed")
		return
	}
	var params oauth.RedirectParams
	if err := c.Bind(&params); err != nil {
		c.String(http.StatusBadRequest, "params get: failed")
		return
	}
	err = e.O.CheckCode(storedState, params)
	if err != nil {
		c.String(http.StatusForbidden, fmt.Sprint(err))
		return
	}

	tok, err := e.O.Auth(context.Background(), params.Code)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	if !tok.Valid() {
		c.String(http.StatusInternalServerError, "invalid token")
		return
	}

	txx, err := e.Db.Beginx()
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}

	control, err := util.GenRandom(128)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	user := db.User{Control: control}
	user.ApplyToken(tok)
	_, err = txx.NamedExec(db.UserRegister, &user)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	err = txx.Commit()
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	c.String(http.StatusOK, user.Control)
}

// OauthAuthorize redirects the user to the OAuth server to authorize this OAuth client.
func (e *engineContext) OauthAuthorize(c *gin.Context) {
	state, url, err := e.O.AuthorizeURL()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	oauthSetState(c, state)
	c.String(http.StatusOK, url)
}

// oauthStateName is the name of the cookie to store the OAuth state.
const oauthStateName = "oauth-state" // set __Host- prefix only if secure

// oauthSetState sets a Set-Cookie header to set a state cookie.
func oauthSetState(c *gin.Context, state string) {
	if gin.Mode() == "debug" {
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     oauthStateName,
			Value:    state,
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
