package server

import (
	"fmt"
	"regexp"
	"time"

	"gitlab.com/mirukakoro/sisikyo/events/api"
)

type ICSOptions struct {
	OAuthSecret string `form:"secret" binding:""`

	// in order to be applied

	Public *bool         `form:"public" binding:""`
	Start  *time.Time    `form:"start" binding:""`
	End    *time.Time    `form:"end" binding:""`
	Tags   []string      `form:"tag" binding:""`
	Term   []api.TermID  `form:"term" binding:""`
	Orgs   []api.OrgSlug `form:"org" binding:""`
	Name   string        `form:"name" binding:""`
	Desc   string        `form:"desc" binding:""`
}

func (o ICSOptions) List(c *api.Client) ([]api.Event, error) {
	events := api.EventsResp{}
	err := c.Do(api.EventsReq{
		Start: o.Start,
		End:   o.End,
	}, &events)
	if err != nil {
		return nil, err
	}

	var name, desc *regexp.Regexp
	if o.Name != "" {
		name, err = regexp.Compile(o.Name)
		if err != nil {
			return nil, fmt.Errorf("name: %w", err)
		}
	}
	if o.Desc != "" {
		name, err = regexp.Compile(o.Desc)
		if err != nil {
			return nil, fmt.Errorf("desc: %w", err)
		}
	}

	final := make([]api.Event, 0, len(events))
	for _, event := range events {
		if o.Public != nil && *o.Public != event.Public {
			continue
		}
		if o.Start != nil && !o.Start.Before(event.Start) {
			continue
		}
		if o.End != nil && !o.End.After(event.End) {
			continue
		}

		if o.Tags != nil {
			for _, tag := range o.Tags {
				for _, eTag := range event.Tags {
					if tag == eTag.Name {
						goto tagsOk
					}
				}
			}
			continue
		tagsOk:
		}

		if o.Term != nil {
			for _, term := range o.Term {
				if term == event.Term {
					goto termOk
				}
			}
			continue
		termOk:
		}

		if o.Orgs != nil {
			for _, org := range o.Orgs {
				if org == event.Org.Slug {
					goto orgsOk
				}
			}
			continue
		orgsOk:
		}

		if name != nil && !name.MatchString(event.Name) {
			continue
		}
		if desc != nil && !desc.MatchString(event.Desc) {
			continue
		}
		final = append(final, event)
	}
	return final, nil
}
