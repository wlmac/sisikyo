// This program, Sisikyo, is a program that provides utilities for an API.
// Copyright (C) 2021 Ken Shibata

// Package server implements a server.
package server

import (
	_ "embed"
	"flag"
	"fmt"

	_ "github.com/gin-contrib/cache"

	"github.com/gin-gonic/gin"
	"gitlab.com/mirukakoro/sisikyo/oauth"
)

// Main runs the server. It should only be called once in a program from main().
func Main() error {
	flag.Parse()
	printLicenseInfo()
	addr, err := setupWeb()
	if err != nil {
		return fmt.Errorf("web: %w", err)
	}
	conn, err := setupDb()
	if err != nil {
		return fmt.Errorf("db: %w", err)
	}
	cl, err := setupAPI()
	if err != nil {
		return fmt.Errorf("api: %w", err)
	}
	oCfg, err := setupOauth()
	if err != nil {
		return fmt.Errorf("oauth: %w", err)
	}
	var oCl *oauth.Client
	if oCfg != nil {
		oCl = oauth.NewClient(*oCfg)
	}
	e := gin.Default()
	setupEngine(e, cl, oCl, conn)
	return e.Run(addr)
}
