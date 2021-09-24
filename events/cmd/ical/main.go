// This file is part of a program named Sisikyō or Sisikyo.
//
// Copyright (C) 2019 Ken Shibata <kenxshibata@gmail.com>
//
// License as published by the Free Software Foundation, either version 1 of the License, or (at your option) any later
// version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied
// warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see
// <https://www.gnu.org/licenses/>.

// Package main is a command that generates an iCalendar file from the API.
package main

// TODO: auto-detect all-day events

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	ics "github.com/arran4/golang-ical"
	"gitlab.com/mirukakoro/sisikyo/events/api"
	"gitlab.com/mirukakoro/sisikyo/events/cmd"
)

//go:embed tmpls/*.html
var tmpls embed.FS

func setupFlag() {
	_, _ = fmt.Fprint(os.Stderr, cmd.StartupInfo)
	oldUsage := flag.Usage
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, cmd.LicenseInfo)
		oldUsage()
	}
}

func readFlags() (baseURL *url.URL, timeout time.Duration, start *time.Time, duration time.Duration, err error) {
	var rawBaseURL string
	var rawStart string
	flag.StringVar(&rawBaseURL, "base-url", api.DefaultBaseURL.String(), "base URL of the API to use")
	flag.DurationVar(&timeout, "timeout", 1*time.Second, "timeout for API")
	flag.StringVar(&rawStart, "start", "", "start of time frame in RFC 3339 foramt")
	flag.DurationVar(&duration, "duration", 0, "duration of time frame")
	flag.Parse()
	baseURL, err = url.Parse(rawBaseURL)
	if err != nil {
		err = fmt.Errorf("timeout format: %w", err)
		return
	}
	if rawStart != "" {
		var start2 time.Time
		start2, err = time.Parse(time.RFC3339, rawStart)
		if err != nil {
			err = fmt.Errorf("start format: %w", err)
			return
		}
		start = &start2
	}
	return
}

func getEvents(baseURL *url.URL, timeout time.Duration, start *time.Time, duration time.Duration) (resp api.EventsResp, err error) {
	client := api.NewClient(&http.Client{Timeout: timeout}, baseURL)
	var end *time.Time
	if start != nil {
		if duration != 0 {
			end2 := start.Add(duration)
			end = &end2
		}
	}
	resp, err = client.Events(api.EventsReq{
		Start: start,
		End:   end,
	})
	if err != nil {
		err = fmt.Errorf("api: %w", err)
		return
	}
	return
}

func getCal(tmpl *template.Template, resp api.EventsResp) (rendered string, err error) {
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
		ev.SetOrganizer(event.Org)
		if event.PerDay() {
			ev.SetAllDayStartAt(event.Start.Truncate(24 * time.Hour))
			ev.SetAllDayEndAt(event.End.Truncate(24 * time.Hour))
		} else {
			ev.SetStartAt(event.Start)
			ev.SetEndAt(event.End)
		}
	}
	rendered = cal.Serialize()
	return
}

func main_() (err error) {
	setupFlag()
	baseURL, timeout, start, duration, err := readFlags()
	if err != nil {
		err = fmt.Errorf("flags: %w", err)
		return
	}
	resp, err := getEvents(baseURL, timeout, start, duration)
	if err != nil {
		err = fmt.Errorf("get events: %w", err)
		return
	}
	tmpl := template.New("root")
	_, err = tmpl.Parse("")
	if err != nil {
		panic(fmt.Sprintf("parsing of blank template failed: %s", err))
	}
	tmpl.Funcs(template.FuncMap{
		"now": func() string { return time.Now().UTC().Format(time.RFC3339) },
	})
	tmpl, err = tmpl.ParseFS(tmpls, "tmpls/*.html")
	if err != nil {
		err = fmt.Errorf("parse tmpls: %w", err)
		return
	}

	cal, err := getCal(tmpl, resp)
	if err != nil {
		err = fmt.Errorf("get cal: %w", err)
		return
	}
	_, _ = fmt.Fprint(os.Stdout, cal)
	return nil
}

func main() {
	err := main_()
	if err != nil {
		log.Fatal(err)
	}
}
