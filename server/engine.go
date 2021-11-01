package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"gitlab.com/mirukakoro/sisikyo/events/api"
	"gitlab.com/mirukakoro/sisikyo/events/ics"
	"gitlab.com/mirukakoro/sisikyo/oauth"
	"gitlab.com/mirukakoro/sisikyo/server/static"
	"gitlab.com/mirukakoro/sisikyo/server/tmpls"
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

		//	auth := api.Auth{}
		//	{
		//		authResp := api.AuthResp{}
		//		err = cl.Do(api.AuthReq{
		//			Username: username,
		//			Password: password,
		//		}, &authResp)
		//		if err != nil {
		//			_ = c.AbortWithError(http.StatusForbidden, fmt.Errorf("auth: %w", err))
		//			return
		//		}
		//		authResp.UpdateAuth(&auth)
		//	}

		//	courseEvents, err := cl.CourseEvents(auth)
		//	if err != nil {
		//		_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("course events: %w", err))
		//		return
		//	}
		//	events = append(events, courseEvents...)
	}
}

func setupEngine(e *gin.Engine, conn *sqlx.DB) {
	e.SetHTMLTemplate(tmpls.Tmpl)

	cl := api.DefaultClient()
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

	ics := func(c *gin.Context, evs []api.Event) {
		cal, err := ics.GetCal(ics.Tmpl, evs)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cal: %w", err))
			return
		}
		c.Data(http.StatusOK, "text/calendar", []byte(cal))
	}

	e.GET("/events/public.json", makeEventsResp(
		json,
		func(o ICSOptions) ([]api.Event, error) {
			return o.List(cl)
		},
	))
	e.GET("/events/public.ics", makeEventsResp(
		ics,
		func(o ICSOptions) ([]api.Event, error) {
			return o.List(cl)
		},
	))
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
}
