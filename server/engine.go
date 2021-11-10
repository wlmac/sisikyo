package server

import (
	"fmt"
	"html/template"
	"net/http"
	"runtime/debug"

	"github.com/Masterminds/sprig"
	"github.com/foolin/goview"
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
		tmplRender(http.StatusOK, "evs", goview.M{
			"title": fmt.Sprintf("%s Events", name),
			"ec":    ec,
			"evs":   evs,
		})(c)
	}
}

func tmplRender(status int, name string, m goview.M) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := goview.Render(c.Writer, status, name, m)
		if err != nil {
			panic(err)
		}
	}
}

func setupEngine(e *gin.Engine, cl *api.Client, o *oauth.Client, conn *sqlx.DB) {
	e.SetFuncMap(template.FuncMap{
		"licenseBrief": func() string { return licenseBrief },
	})

	funcMap := sprig.FuncMap()
	funcMap["licenseBrief"] = func() string { return licenseBrief }
	funcMap["licenseBriefHTML"] = func() template.HTML { return template.HTML(licenseBriefHTML) }

	gv := goview.New(goview.Config{
		Root:         "server/views",
		Extension:    ".tpl",
		Master:       "/layouts/base",
		Partials:     []string{"/partials/links"},
		Funcs:        funcMap,
		DisableCache: gin.Mode() == "debug",
	})

	goview.Use(gv)

	//render index use `index` without `.tpl` extension, that will render with master layout.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := goview.Render(w, http.StatusOK, "index", goview.M{})
		if err != nil {
			fmt.Fprintf(w, "Render index error: %v!", err)
		}

	})

	e.Static("/static", "server/static")

	ec := engineContext{E: e, API: cl, O: o, Db: conn}
	indexOutput := index()
	{
		buildInfo, _ := debug.ReadBuildInfo()
		e.GET("/", tmplRender(http.StatusOK, "index", goview.M{
			"title":        "Sisikyō",
			"index":        indexOutput,
			"apiURL":       apiURL,
			"oauthBaseURL": oauthBaseURL,
			"buildInfo":    buildInfo,
		}))
	}
	e.GET("/register", func(c *gin.Context) {
		tmplRender(http.StatusOK, "register", goview.M{
			"title": "Register",
			"url":   ec.oauthAuthorize(c),
		})(c)
	})
	e.GET("/remove", tmplRender(http.StatusOK, "remove", goview.M{
		"title": "Remove",
	}))
	e.POST("/remove", ec.Remove)
	e.GET("/about", tmplRender(http.StatusOK, "about", goview.M{
		"title": "About",
	}))
	e.GET("/license", func(c *gin.Context) {
		c.String(http.StatusOK, "%s", licenseFull)
	})
	e.GET("/src", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "https://gitlab.com/mirukakoro/sisikyo")
	})
	e.GET("/events/public.json", ec.PublicQuery(renderJSON))
	e.GET("/events/public.ics", ec.PublicQuery(renderICS))
	e.GET("/events/public.html", ec.PublicQuery(ec.renderHTML("Public")))
	e.GET("/events/:control/private.json", ec.UserQuery(renderJSON))
	e.GET("/events/:control/private.ics", ec.UserQuery(renderICS))
	e.GET("/events/:control/private.html", ec.UserQuery(ec.renderHTML("Private")))
	e.GET("/o/redirect", ec.OauthRedirect)
	e.GET("/o/authorize", ec.OauthAuthorize)
}
