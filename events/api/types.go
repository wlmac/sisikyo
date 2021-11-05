package api

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type Data interface {
	URL(c *Client) *url.URL
}

var _ = [...]Data{
	User{},
	Timeframe{},
	Ann{},
	Course{},
	Term{},
	Event{},
	Org{},
	Tag{},
	Schedule{},
}

type User struct {
	Username       Username      `json:"username"`
	FirstName      string        `json:"first_name"`
	LastName       string        `json:"last_name"`
	Bio            string        `json:"bio"`
	Timezone       string        `json:"timezone"`
	GraduatingYear int           `json:"graduating_year"`
	Organizations  []string      `json:"organizations"`
	TagsFollowing  []interface{} `json:"tags_following"`
}

func (u User) Name() string {
	return u.FirstName + " " + u.LastName
}

func (u User) URL(c *Client) *url.URL {
	sub, _ := url.Parse(string(u.Username))
	return c.BaseURL().ResolveReference(sub)
}

type Timeframe struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func (t Timeframe) URL(_ *Client) *url.URL { return nil }

// Ann represents an announcement.
// Model: core.models.post.Announcement
// Serializer: core.api.serializers.announcement.AnnouncementSerializer
type Ann struct {
	Id           AnnID     `json:"id"`
	Author       Username  `json:"author"`
	Org          OrgName   `json:"organization"`
	Tags         []Tag     `json:"tags"`
	Created      time.Time `json:"created_date"`
	LastModified time.Time `json:"last_modified_date"`
	Title        string    `json:"title"`
	Body         string    `json:"body"`
	Public       bool      `json:"is_public"`

	// fields below are not in the API
	XAuthor    UserResp
	XOrg       Org
	XImageURL  string
	XImagePath string
	XImageAlt  string
	XURL       string
}

var annBaseURL, _ = url.Parse("/announcement/")

func (a Ann) URL(c *Client) *url.URL {
	id, _ := url.Parse(strconv.FormatInt(int64(a.Id), 10))
	return c.baseURL.ResolveReference(annBaseURL).ResolveReference(id)
}

func (a Ann) String() string {
	return fmt.Sprintf(
		"(%d) %s\nauthor: %s\norg: %s\ntags: %v\ncreated: %s\nlast modified: %s\npublic: %t\n\n%s",
		a.Id,
		a.Title,
		a.Author,
		a.Org,
		a.Tags,
		a.Created,
		a.LastModified,
		a.Public,
		a.Body,
	)
}

func (a Ann) ReqOrg(c *Client) (Org, error) {
	resp := OrgsResp{}
	err := c.Do(OrgsReq{}, &resp)
	if err != nil {
		return Org{}, err
	}
	for _, org := range resp {
		if org.Name == a.Org {
			return org, nil
		}
	}
	return Org{}, errors.New("not found")
}

// Course represents a course.
// TODO: fill model
// Model: core.models.course.Course
// Serializer: core.api.serializers.course.CourseSerializer
type Course struct {
	Id        CourseID    `json:"id"`
	Code      CourseCode  `json:"code"`
	Term      TermID      `json:"term"`
	Desc      string      `json:"description"`
	Position  int         `json:"position"`
	Submitter interface{} `json:"submitter"`
}

func (c Course) URL(c2 *Client) *url.URL { return c2.baseURL }

func (c Course) String() string {
	return fmt.Sprintf("(%d) %s (term %d) (position %d) (submitter %v): %s", c.Id, c.Code, c.Term, c.Position, c.Submitter, c.Desc)
}

// Term represents a term.
// TODO: fill model
// Model: core.models.course.Term
// Serializer: core.api.serializers.course.TermSerializer
type Term struct {
	Id     TermID    `json:"id"`
	Name   string    `json:"name"`
	Desc   string    `json:"description"`
	Fmt    string    `json:"timetable_format"`
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
	Frozen bool      `json:"is_frozen"`
}

func (t Term) URL(c *Client) *url.URL { return c.baseURL }

// Event represents an event,
// Model: core.models.course.Event
// Serializer: core.api.serializers.course.EventSerializer
type Event struct {
	Id     EventID   `json:"id"`
	Name   string    `json:"name"`
	Desc   string    `json:"description"`
	Tags   []Tag     `json:"tags"`
	Term   TermID    `json:"term"`
	Org    Org       `json:"organization"`
	Start  time.Time `json:"start_date"`
	End    time.Time `json:"end_date"`
	Public bool      `json:"is_public"`
}

func (e Event) URL(c *Client) *url.URL { return c.baseURL }

func (e Event) String() string {
	return fmt.Sprintf("(%d) %s (%s~%s,%s): %s", e.Id, e.Name, e.Start, e.End, e.End.Sub(e.Start), e.Desc)
}

// PerDay returns whether the Event is an all-day event.
// Note: "all-day event" means that the event is at least 23:59:00 long and is shorter than 24:00:00.
func (e Event) PerDay() bool {
	diff := e.End.Sub(e.Start).Truncate(24 * time.Hour).Truncate(1 * time.Second)
	return diff == 24*time.Hour-1*time.Second
}

// Org represents an organization.
// Model:core.models.organization.Org
// Serializer: core.api.serializers.organization.OrganizationSerializer
type Org struct {
	Id          int       `json:"id"`
	Owner       User      `json:"owner"`
	Supervisors []User    `json:"supervisors"`
	Execs       []User    `json:"execs"`
	Tags        []Tag     `json:"tags"`
	Name        OrgName   `json:"name"`
	Bio         string    `json:"bio"`
	Extra       string    `json:"extra_content"`
	Slug        OrgSlug   `json:"slug"`
	Registered  time.Time `json:"registered_date"`
	Open        bool      `json:"is_open"`
	AppsOpen    bool      `json:"applications_open"`
	Banner      string    `json:"banner"`
	Icon        string    `json:"icon"`
}

var orgBaseURL, _ = url.Parse("/club/")

func (o Org) URL(c *Client) *url.URL {
	slug, _ := url.Parse(url.PathEscape(string(o.Slug)))
	return c.baseURL.ResolveReference(orgBaseURL).ResolveReference(slug)
}

func (o Org) IconURL(c *Client) *url.URL {
	iconURL, _ := url.Parse(o.Icon)
	return c.baseURL.ResolveReference(iconURL)
}

// Tag represents a single tag from the API.
type Tag struct {
	ID    TagID  `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

func (t Tag) URL(c *Client) *url.URL { return c.baseURL }

type Schedule struct {
	Description struct {
		Time   string `json:"time"`
		Course string `json:"course"`
	} `json:"description"`
	Time   Timeframe   `json:"time"`
	Pos    []int       `json:"position"`
	Cycle  string      `json:"cycle"`
	Course *CourseCode `json:"course"`
}

func (s Schedule) URL(_ *Client) *url.URL { return nil }

func (s Schedule) Event(c *Client) (Event, error) {
	course, err := s.Course.Deref(c)
	if err != nil {
		return Event{}, err
	}
	return Event{
		Id:    -1,
		Name:  string(course.Code),
		Desc:  fmt.Sprintf("ID: %d\nTerm: %d\nPosition: %d, %d\nCycle: %s\n\n%s", course.Id, course.Term, course.Position, s.Pos, s.Cycle, course.Desc),
		Start: s.Time.Start,
		End:   s.Time.End,
	}, nil
}
