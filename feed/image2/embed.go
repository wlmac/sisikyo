package image2

import (
	"embed"
	"html/template"
)

//go:embed tmpls/*
var fs embed.FS

var Tmpl *template.Template

func init() {
	tmpl := template.New("root")
	tmpl, err := tmpl.Parse("")
	if err != nil {
		panic(err)
	}
	funcMap := template.FuncMap{
		"safeAttr": func(s string) template.HTMLAttr {
			return template.HTMLAttr(s)
		},
		"safeURL": func(s string) template.URL {
			return template.URL(s)
		},
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
	}
	tmpl = tmpl.Funcs(funcMap)
	tmpl, err = tmpl.ParseFS(fs, "tmpls/*.tmpl")
	if err != nil {
		panic(err)
	}
	Tmpl = tmpl
}
