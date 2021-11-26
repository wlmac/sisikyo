package text

import (
	"embed"
	"text/template"
)

//go:embed tmpls/*
var fs embed.FS

// Tmpl is used for rendering text.
var Tmpl *template.Template

func init() {
	tmpl := template.New("root")
	tmpl, err := tmpl.Parse("")
	if err != nil {
		panic(err)
	}
	tmpl, err = tmpl.ParseFS(fs, "tmpls/*.tmpl")
	if err != nil {
		panic(err)
	}
	Tmpl = tmpl
}
