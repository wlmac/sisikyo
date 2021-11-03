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

func makeEventsResp(render func(c *gin.Context, evs []api.Event), f func(o ICSOptions) ([]api.Event, error)) func(c *gin.Context) {
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

func renderIcs(c *gin.Context, evs []api.Event) {
	cal, err := ics.GetCal(ics.Tmpl, evs)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cal: %w", err))
		return
	}
	c.Data(http.StatusOK, "text/calendar", []byte(cal))
}

func setupEngine(e *gin.Engine, cl *api.Client, o *oauth.Client, conn *sqlx.DB) {
	e.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "%s", licenseBrief)
	})
	e.GET("/license", func(c *gin.Context) {
		c.String(http.StatusOK, "%s", licenseFull)
	})
	e.GET("/src", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "https://gitlab.com/mirukakoro/sisikyo")
	})

	json := func(c *gin.Context, evs []api.Event) {
		c.JSON(http.StatusOK, evs)
	}

	e.GET("/events/public.json", makeEventsResp(
		json,
		func(o ICSOptions) ([]api.Event, error) {
			return o.List(cl)
		},
	))
	e.GET("/events/public.ics", makeEventsResp(
		renderIcs,
		func(o ICSOptions) ([]api.Event, error) {
			return o.List(cl)
		},
	))

	ec := engineContext{
		e:   e,
		api: cl,
		o:   o,
		db:  conn,
	}
	e.GET("/events/:control/private.json", ec.UserQuery)
	e.POST("/remove/:control", ec.UserRemove)
	e.GET("/o/redirect", ec.OauthRedirect)
	e.GET("/o/authorize", ec.OauthAuthorize)
	/*
		e.GET("/o/redirect.txt", oauth.RedirectJSON)
		e.GET("/o/redirect", oauth.Redirect)
		e.GET("/o/authorize.txt", oauth.Authorize)
		e.GET("/o/authorize", func(c *gin.Context) {
			url, err := oauth.Authorize2(cl)
			if err != nil {
				c.HTML(http.StatusInternalServerError, "error.html.tmpl", gin.H{
					"err": err,
				})
				return
			}
			c.HTML(http.StatusOK, "authorize.html.tmpl", gin.H{
				"url": url,
			})
		})
		e.StaticFS("/static", http.FS(static.Static))
		e.GET("/:path", func(c *gin.Context) {
			c.String(http.StatusOK, "%s", c.Request.URL.Path)
		})
	*/
}
