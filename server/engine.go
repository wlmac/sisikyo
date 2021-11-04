package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"gitlab.com/mirukakoro/sisikyo/events/api"
	"gitlab.com/mirukakoro/sisikyo/events/ics"
	"gitlab.com/mirukakoro/sisikyo/oauth"
)

type renderFunc func(c *gin.Context, evs []api.Event)

func makeEventsResp(render renderFunc, f func(o ICSOptions) ([]api.Event, error)) func(c *gin.Context) {
	return func(c *gin.Context) {
		o := ICSOptions{}
		err := c.ShouldBindQuery(&o)
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		evs, err := f(o)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("list: %w", err))
			return
		}

		render(c, evs)
	}
}

func renderICS(c *gin.Context, evs []api.Event) {
	cal, err := ics.GetCal(ics.Tmpl, evs)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cal: %w", err))
		return
	}
	c.Data(http.StatusOK, "text/calendar", []byte(cal))
}

func renderJSON(c *gin.Context, evs []api.Event) {
	c.JSON(http.StatusOK, evs)
}

func setupEngine(e *gin.Engine, cl *api.Client, o *oauth.Client, conn *sqlx.DB) {
	indexOutput := index()
	e.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/markdown")
		c.String(http.StatusOK, "%s", indexOutput)
	})
	e.GET("/license", func(c *gin.Context) { c.String(http.StatusOK, "%s", licenseFull) })
	e.GET("/src", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "https://gitlab.com/mirukakoro/sisikyo")
	})

	ec := engineContext{
		e:   e,
		api: cl,
		o:   o,
		db:  conn,
	}
	e.GET("/events/public.json", ec.PublicQuery(renderJSON))
	e.GET("/events/public.ics", ec.PublicQuery(renderICS))
	e.GET("/events/:control/private.json", ec.UserQuery(renderJSON))
	e.GET("/events/:control/private.ics", ec.UserQuery(renderICS))
	e.POST("/remove/:control", ec.UserRemove)
	e.GET("/o/redirect", ec.OauthRedirect)
	e.GET("/o/authorize", ec.OauthAuthorize)
}
