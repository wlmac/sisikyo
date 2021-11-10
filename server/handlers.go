package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

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
	Url     string `form:"url"`
	Control string `form:"control"`
}

func (r *removeForm) getControl() (string, error) {
	if r.Control == "" && r.Url == "" {
		return "", errors.New("either control or url must be filled")
	}
	if r.Control != "" {
		return r.Control, nil
	}
	u, err := url.Parse(r.Url)
	if err != nil {
		return "", err
	}
	if !u.IsAbs() {
		return "", errors.New("url must be absolute")
	}
	if !strings.HasPrefix(u.Path, "/events/") || !strings.HasSuffix(strings.Split(u.Path, ".")[0], "/private") {
		return "", errors.New("url does not have expected format")
	}
	return strings.TrimSuffix(strings.Split(strings.TrimPrefix(u.Path, "/events/"), ".")[0], "/private"), nil
}

func (e *engineContext) Remove(c *gin.Context) {
	var err error
	var form removeForm
	if err = c.Bind(&form); err != nil {
		c.String(http.StatusBadRequest, "%s", err)
		return
	}
	var control string
	if control, err = form.getControl(); err != nil {
		c.String(http.StatusBadRequest, "%s", err)
		return
	}
	if err = e.userRemove(controlParam{Control: control}); err != nil {
		c.String(http.StatusInternalServerError, "%s", err)
		return
	}
	tmplRender(http.StatusOK, "removed", goview.M{
		"title":   "Removed If Present",
		"control": control,
	})(c)
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
		tmplRender(http.StatusBadRequest, "register-err", goview.M{
			"title": "Registration Failed",
			"desc":  "An invalid request to register was detected: failed to get the state from your cookies.",
			"err":   errors.New("state get: failed"),
		})(c)
		return
	}
	oauthClearState(c)
	var params oauth.RedirectParams
	if err := c.Bind(&params); err != nil {
		tmplRender(http.StatusBadRequest, "register-err", goview.M{
			"title": "Registration Failed",
			"desc":  "An invalid request to register was detected: failed to get the URL parameters.",
			"err":   errors.New("params get: failed"),
		})(c)
		return
	}
	err = e.O.CheckCode(storedState, params)
	if err != nil {
		tmplRender(http.StatusForbidden, "register-err", goview.M{
			"title": "Registration Failed",
			"desc":  "Stored state and state from URL params did not match.",
			"err":   err,
		})(c)
		return
	}

	tok, err := e.O.Auth(context.Background(), params.Code)
	if err != nil {
		tmplRender(http.StatusInternalServerError, "register-err", goview.M{
			"title": "Registration Failed",
			"desc":  "OAuth exchange failed for some reason (see below):",
			"err":   err,
		})(c)
		return
	}
	if !tok.Valid() {
		tmplRender(http.StatusInternalServerError, "register-err", goview.M{
			"title": "Registration Failed",
			"desc":  "OAuth token is invalid.",
		})(c)
		return
	}

	txx, err := e.Db.Beginx()
	if err != nil {
		tmplRender(http.StatusInternalServerError, "register-err", goview.M{
			"title": "Registration Failed",
			"desc":  "Connection attempt to database failed.",
			"err":   err,
		})(c)
		return
	}

	control, err := util.GenRandom(128)
	if err != nil {
		tmplRender(http.StatusInternalServerError, "register-err", goview.M{
			"title": "Registration Failed",
			"desc":  "Generation of random values for the control key failed.",
			"err":   err,
		})(c)
		return
	}
	user := db.User{Control: control}
	user.ApplyToken(tok)
	_, err = txx.NamedExec(db.UserRegister, &user)
	if err != nil {
		tmplRender(http.StatusInternalServerError, "register-err", goview.M{
			"title": "Registration Failed",
			"desc":  "Insertion into the database failed.",
			"err":   err,
		})(c)
		return
	}
	err = txx.Commit()
	if err != nil {
		tmplRender(http.StatusOK, "register-err", goview.M{
			"title": "Registration Failed",
			"desc":  "Committing changes to the database failed.",
			"err":   err,
		})(c)
		return
	}
	tmplRender(http.StatusOK, "registered", goview.M{
		"title":   "Registered",
		"control": user.Control,
	})(c)
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
