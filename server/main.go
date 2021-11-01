// This program, Sisikyo, is a program that provides utilities for an API.
// Copyright (C) 2021 Ken Shibata

package server

import (
	_ "embed"
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
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
	err = setupAPI()
	if err != nil {
		return fmt.Errorf("api: %w", err)
	}
	e := gin.Default()
	setupEngine(e, conn)
	return e.Run(addr)
}
