package tmpls

import (
	"embed"
	"fmt"
	"html/template"
)

//go:embed *.tmpl
var tmpls embed.FS

var Tmpl *template.Template

func init() {
	var err error
	tmpl := template.New("root")
	_, err = tmpl.Parse("")
	if err != nil {
		err = fmt.Errorf("parsing of blank template failed: %s", err)
		panic(err)
	}
	Tmpl, err = tmpl.ParseFS(tmpls, "*.tmpl")
	if err != nil {
		err = fmt.Errorf("parse: %w", err)
		panic(err)
	}
}
