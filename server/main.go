// This program, Sisikyo, is a program that provides utilities for an API.
// Copyright (C) 2021 Ken Shibata

package server

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"gitlab.com/mirukakoro/sisikyo/events/api"
	"gitlab.com/mirukakoro/sisikyo/events/ics"
	"gitlab.com/mirukakoro/sisikyo/oauth"
)

const licenseBrief = "This program, Sisikyo, is a program that provides utilities for an API.\n" +
	"Copyright (C) 2021 Ken Shibata\n" +
	"This program comes with ABSOLUTELY NO WARRANTY and this program is free software, and you are welcome to " +
	"redistribute it under certain conditions; for details view the '/license' page or view the 'license.md' file. " +
	"This program uses open source software; for details view the '/license' page or view the 'license.md' file.\n\n" +
	"The source code is available from the '/src' page."

//go:embed license.md
var licenseFull string

const username = ""
const password = ""

func printLicenseInfo() {
	_, _ = fmt.Fprint(os.Stderr, licenseBrief)
}

func setupFlag() (string, error) {
	var port int
	var host string
	var apiURL string
	flag.IntVar(&port, "port", 8080, "port to bind to")
	flag.StringVar(&host, "host", "", "host to bind to")
	flag.StringVar(&apiURL, "api-url", "", "URL of API to use")
	flag.Parse()

	if apiURL == "" {
		return "", errors.New("api-url: cannot be blank")
	}
	var err error
	api.DefaultBaseURL, err = url.Parse(apiURL)
	if err != nil {
		return "", fmt.Errorf("api-url: %w", err)
	}
	return fmt.Sprintf("%s:%d", host, port), nil
}

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

func setupEngine(e *gin.Engine) {
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
	e.GET("/o/redirect", oauth.RedirectJSON)
	e.GET("/o/authorize", oauth.Authorize)
}

func Main() error {
	printLicenseInfo()
	addr, err := setupFlag()
	if err != nil {
		return err
	}
	e := gin.Default()
	setupEngine(e)
	return e.Run(addr)
}
