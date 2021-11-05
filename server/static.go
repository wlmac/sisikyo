/*
This program, Sisikyo, is a program that provides utilities for an API.
Copyright 2021 Ken Shibata
This program comes with ABSOLUTELY NO WARRANTY and this program is free software, and you are welcome to redistribute it under certain conditions; for details view the `license.md` file.
This program uses open source software; for details view the `license.md` file.
*/
package server

import (
	_ "embed"
	"fmt"
	"os"
	"runtime/debug"
)

func index() string {
	return licenseBrief + "\n\n" + debugInfo()
}

const licenseBrief = `This program, Sisikyo, is a program that provides utilities for an API.

Copyright 2021 Ken Shibata

This program comes with ABSOLUTELY NO WARRANTY and this program is free software, and you are welcome to redistribute it under certain conditions; for details view the [/license](/license) page or view the ` + "`license.md`" + ` file.
This program uses open source software; for details view the [/license](/license) page or view the ` + "`license.md`" + ` file.
The source code is available from the [/src](/src) page.

TODO: remove the registered page, and replace with a fat client.
`

//go:embed license.md
var licenseFull string

func printLicenseInfo() {
	_, _ = fmt.Fprint(os.Stderr, licenseBrief)
}

func debugInfo() string {
	res := "# デバッグ情報\n"
	res += fmt.Sprintf("API用のベースURL: [%s](%s)\n\n", apiURL, apiURL)
	res += fmt.Sprintf("OAuth用のベースURL: [%s](%s)\n\n", oauthBaseURL, oauthBaseURL)
	res += buildInfo()
	return res
}

func buildInfo() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "## ビルド情報\n\n情報を取得出来ませんでした。"
	}
	res := "## ビルド情報\n\n"
	res += fmt.Sprintf("パス: `%s`\n\n", info.Path)
	res += fmt.Sprintf("メイン: %s\n\n", modString(&info.Main))
	res += "### 構築依存\n\n"
	for i, mod := range info.Deps {
		res += fmt.Sprintf("%d. %s\n\n", i+1, modString(mod))
	}
	return res
}

func modString(mod *debug.Module) string {
	if mod == nil {
		return ""
	}
	res := fmt.Sprintf("`%s@%s` (チェックサム: `%s`)", mod.Path, mod.Version, mod.Sum)
	if mod.Replace != nil {
		res += fmt.Sprintf(" (取り替: %s)", modString(mod.Replace))
	}
	return res
}
