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

package api

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

// eventsURL is the url for a Client.Events API call.
var eventsURL, _ = url.Parse("events")

// EventsReq has the options for a Client.Events API call.
type EventsReq struct {
	Start *time.Time
	End   *time.Time
}

// Query serializes the options in EventsReq to a new url.Values.
func (e EventsReq) Query() (v url.Values) {
	v = url.Values{}
	v.Add("format", "json")
	if e.Start != nil {
		v.Add("start", e.Start.Format(time.RFC3339))
	}
	if e.End != nil {
		v.Add("end", e.End.Format(time.RFC3339))
	}
	return v
}

// EventsResp is the expected response from an Client.Events API call.
type EventsResp = []Event

// Event represents a single event returned by the API.
type Event struct {
	Id     int       `json:"id"`
	Org    string    `json:"organization"`
	Tags   []Tag     `json:"tags"`
	Name   string    `json:"name"`
	Desc   string    `json:"description"`
	Start  time.Time `json:"start_date"`
	End    time.Time `json:"end_date"`
	Public bool      `json:"is_public"`
	Term   int       `json:"term"`
}

// AllDay returns whether the Event is an all-day event.
// Note: "all-day event" means that the event is at least 23:59:00 long and is shorter than 24:00:00.
func (e Event) AllDay() bool {
	diff := e.End.Sub(e.Start).Truncate(1 * time.Second)
	return diff == 24*time.Hour-1*time.Second
}

func (e Event) String() string {
	return fmt.Sprintf(
		"%s (id %d)\norg: %s\nall-day: %t\nframe: %s to %s (%s)\npublic: %t\nterm: %d\ntags: %s\n%s",
		e.Name,
		e.Id,
		e.Org,
		e.AllDay(),
		e.Start,
		e.End,
		e.End.Sub(e.Start),
		e.Public,
		e.Term,
		strings.Join(e.TagsString(), ", "),
		e.Desc,
	)
}

// TagsString generates a string slice with the Tag.Name of each Tag in Tags.
func (e Event) TagsString() []string {
	tags := make([]string, len(e.Tags))
	for i, tag := range e.Tags {
		tags[i] = tag.String()
	}
	return tags
}

// Tag represents a single tag from the API.
type Tag struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

func (t Tag) String() string {
	return fmt.Sprintf("%s (%s)", t.Name, t.Color)
}
