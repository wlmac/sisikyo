package server

import (
	_ "embed"
	"fmt"
	"os"
)

const licenseBrief = "This program, Sisikyo, is a program that provides utilities for an API.\n" +
	"Copyright (C) 2021 Ken Shibata\n" +
	"This program comes with ABSOLUTELY NO WARRANTY and this program is free software, and you are welcome to " +
	"redistribute it under certain conditions; for details view the '/license' page or view the 'license.md' file. " +
	"This program uses open source software; for details view the '/license' page or view the 'license.md' file.\n\n" +
	"The source code is available from the '/src' page.\n"

//go:embed license.md
var licenseFull string

func printLicenseInfo() {
	_, _ = fmt.Fprint(os.Stderr, licenseBrief)
}
