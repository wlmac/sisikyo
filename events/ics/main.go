package ics

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"time"

	ics "github.com/arran4/golang-ical"
	"gitlab.com/mirukakoro/sisikyo/events/api"
)

//go:embed tmpls/*.html
var tmpls embed.FS

func onlyDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func GetCal(tmpl *template.Template, resp api.EventsResp) (rendered string, err error) {
	cal := ics.NewCalendar()
	for i, event := range resp {
		ev := cal.AddEvent(fmt.Sprintf("%d-%d", i, event.Id))
		{
			descBuf := new(bytes.Buffer)
			err = tmpl.ExecuteTemplate(descBuf, "desc", event)
			if err != nil {
				return
			}
			ev.SetDescription(descBuf.String())
		}
		ev.SetSummary(event.Name)
		ev.SetOrganizer(event.Org.Name)
		if event.PerDay() {
			ev.SetAllDayStartAt(onlyDay(event.Start))
			ev.SetAllDayEndAt(onlyDay(event.End))
		} else {
			ev.SetStartAt(event.Start)
			ev.SetEndAt(event.End)
		}
		log.Println("event", event)
	}
	rendered = cal.Serialize()
	return
}

var Tmpl *template.Template

func init() {
	var err error
	tmpl := template.New("root")
	_, err = tmpl.Parse("")
	if err != nil {
		panic(fmt.Sprintf("parsing of blank template failed: %s", err))
	}
	tmpl.Funcs(template.FuncMap{
		"now": func() string { return time.Now().UTC().Format(time.RFC3339) },
	})
	Tmpl, err = tmpl.ParseFS(tmpls, "tmpls/*.html")
	if err != nil {
		err = fmt.Errorf("parse tmpls: %w", err)
		panic(err)
	}
}
