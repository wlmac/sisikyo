package server

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"gitlab.com/mirukakoro/sisikyo/events/api"
	"gitlab.com/mirukakoro/sisikyo/events/ics"
	"gitlab.com/mirukakoro/sisikyo/oauth"
)

type renderFunc func(c *gin.Context, evs []api.Event)

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

func (ec *engineContext) renderHTML(name string) func(c *gin.Context, evs []api.Event) {
	return func(c *gin.Context, evs []api.Event) {
		c.HTML(http.StatusOK, "evs.html", gin.H{
			"name": name,
			"ec":   ec,
			"evs":  evs,
		})
	}
}

func setupEngine(e *gin.Engine, cl *api.Client, o *oauth.Client, conn *sqlx.DB) {
	e.SetFuncMap(template.FuncMap{
		"licenseBrief": func() string { return licenseBrief },
	})

	ec := engineContext{E: e, API: cl, O: o, Db: conn}
	indexOutput := index()
	e.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "%s", licenseBrief+"\n\n"+indexOutput)
	})
	e.GET("/license", func(c *gin.Context) {
		c.String(http.StatusOK, "%s", licenseFull)
	})
	e.GET("/src", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "https://gitlab.com/mirukakoro/sisikyo")
	})
	e.LoadHTMLGlob("../tmpls/*.html")
	e.GET("/events/public.json", ec.PublicQuery(renderJSON))
	e.GET("/events/public.ics", ec.PublicQuery(renderICS))
	e.GET("/events/public.html", ec.PublicQuery(ec.renderHTML("public")))
	e.GET("/events/:control/private.json", ec.UserQuery(renderJSON))
	e.GET("/events/:control/private.ics", ec.UserQuery(renderICS))
	e.GET("/events/:control/private.html", ec.UserQuery(ec.renderHTML("private")))
	e.GET("/remove/:control", ec.UserRemove) // TODO: add a HTML form and replace this with a POST request
	e.GET("/o/redirect", ec.OauthRedirect)
	e.GET("/o/authorize", ec.OauthAuthorize)
}
