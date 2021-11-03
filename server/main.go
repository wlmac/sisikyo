// This program, Sisikyo, is a program that provides utilities for an API.
// Copyright (C) 2021 Ken Shibata

package server

import (
	_ "embed"
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
	"gitlab.com/mirukakoro/sisikyo/oauth"
)

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
	oCl := oauth.NewClient(*oCfg)
	e := gin.Default()
	setupEngine(e, cl, oCl, conn)
	return e.Run(addr)
}
