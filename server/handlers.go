package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/foolin/goview"
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

func (e *engineContext) userRemove(params controlParam) error {
	txx, err := e.Db.Beginx()
	if err != nil {
		return err
	}
	_, err = txx.NamedExec(db.UserRemove, params)
	if err != nil {
		return err
	}
	err = txx.Commit()
	if err != nil {
		return err
	}
	return nil
}

type removeForm struct {
	Control string `form:"control" binding:"required"`
}

func (e *engineContext) Remove(c *gin.Context) {
	var form removeForm
	if err := c.Bind(&form); err != nil {
		c.String(http.StatusBadRequest, "%s", err)
		return
	}
	if err := e.userRemove(controlParam{Control: form.Control}); err != nil {
		c.String(http.StatusInternalServerError, "%s", err)
		return
	}
	tmplRender(http.StatusOK, "removed", goview.M{})(c)
	return
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
	if err := e.userRemove(params); err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprint(err))
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
	oauthClearState(c)
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
	tmplRender(http.StatusOK, "registered", goview.M{"control": user.Control})(c)
}

func (e *engineContext) oauthAuthorize(c *gin.Context) (url string) {
	state, url, err := e.O.AuthorizeURL()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	oauthSetState(c, state)
	return
}

// OauthAuthorize redirects the user to the OAuth server to authorize this OAuth client.
func (e *engineContext) OauthAuthorize(c *gin.Context) {
	c.String(http.StatusOK, e.oauthAuthorize(c))
}

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
