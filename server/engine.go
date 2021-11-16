package server

import (
	"fmt"
	"html/template"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"

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

func setupCache(timeout time.Duration) *persistence.InMemoryStore {
	return persistence.NewInMemoryStore(time.Second)
}

func setupViews() {
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
}

func (ec *engineContext) setupIndex() {
	indexOutput := index()
	buildInfo, _ := debug.ReadBuildInfo()
	ec.E.GET("/", cache.CachePage(ec.Store, time.Minute, tmplRender(http.StatusOK, "index", goview.M{
		"title":        "Sisikyō",
		"index":        indexOutput,
		"apiURL":       apiURL,
		"oauthBaseURL": oauthBaseURL,
		"buildInfo":    buildInfo,
	})))
}

func setupEngine(e *gin.Engine, cl *api.Client, o *oauth.Client, conn *sqlx.DB) {
	setupViews()
	e.Static("/static", "server/static") // TODO: cache for static?
	ec := engineContext{E: e, API: cl, O: o, Db: conn, Store: setupCache(cacheTimeout)}
	ec.setupIndex()
	e.GET("/register", cache.CachePage(ec.Store, cacheTimeout, func(c *gin.Context) {
		tmplRender(http.StatusOK, "register", goview.M{
			"title": "Register",
			"url":   ec.oauthAuthorize(c),
		})(c)
	}))
	e.GET("/remove", cache.CachePage(ec.Store, cacheTimeout, tmplRender(http.StatusOK, "remove", goview.M{
		"title": "Remove",
	})))
	e.POST("/remove", ec.Remove)
	e.GET("/about", cache.CachePage(ec.Store, cacheTimeout, tmplRender(http.StatusOK, "about", goview.M{
		"title": "About",
	})))
	e.GET("/license", cache.CachePage(ec.Store, cacheTimeout, func(c *gin.Context) {
		c.String(http.StatusOK, "%s", licenseFull)
	}))
	e.GET("/src", cache.CachePage(ec.Store, cacheTimeout, func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "https://gitlab.com/mirukakoro/sisikyo")
	}))
	e.GET("/events/public.json", cache.CachePage(ec.Store, cacheTimeout, ec.PublicQuery(renderJSON)))
	e.GET("/events/public.ics", cache.CachePage(ec.Store, cacheTimeout, ec.PublicQuery(renderICS)))
	e.GET("/events/public.html", cache.CachePage(ec.Store, cacheTimeout, ec.PublicQuery(ec.renderHTML("Public"))))
	e.GET("/events/:control/private.json", cache.CachePage(ec.Store, cacheTimeout, ec.UserQuery(renderJSON)))
	e.GET("/events/:control/private.ics", cache.CachePage(ec.Store, cacheTimeout, ec.UserQuery(renderICS)))
	e.GET("/events/:control/private.html", cache.CachePage(ec.Store, cacheTimeout, ec.UserQuery(ec.renderHTML("Private"))))
	e.GET("/o/redirect", ec.OauthRedirect)
	e.GET("/o/authorize", ec.OauthAuthorize)
}
