package feed

import (
	"fmt"
	"gitlab.com/mirukakoro/sisikyo/events/api"
	"image"
	"net/url"
	"time"
)

type Sink interface {
	Post(api.Ann) (url *url.URL, err error)
}

type FormatImager interface {
	FormatImage(api.Ann) (img image.Image, caption string, err error)
}

// ConsideredTime gets the considered time of an api.Ann.
// For example, to get post api.Ann created after t, use Created.
type ConsideredTime func(ann api.Ann) (considered time.Time)

func Created(ann api.Ann) (considered time.Time) {
	return ann.Created
}

func LastModified(ann api.Ann) (considered time.Time) {
	return ann.LastModified
}

type Pipe struct {
	ee             Sink
	c              *api.Client
	consideredTime ConsideredTime
}

func NewPipe(ee Sink, c *api.Client, consideredTime ConsideredTime) *Pipe {
	return &Pipe{ee: ee, c: c, consideredTime: consideredTime}
}

func (p *Pipe) PostAfter(threshold time.Time) (urls []*url.URL, err error) {
	toPost, err := p.After(threshold)
	if err != nil {
		return nil, fmt.Errorf("after: %w", err)
	}
	urls, err = p.Post(toPost)
	if err != nil {
		return nil, fmt.Errorf("post: %w", err)
	}
	return
}

func (p *Pipe) After(threshold time.Time) (toPost []api.Ann, err error) {
	resp := api.AnnResp{}
	err = p.c.Do(api.AnnReq{}, &resp)
	if err != nil {
		err = fmt.Errorf("api: %w", err)
		return
	}
	// make sure that latest ann is the last
	for i := len(resp) - 1; i >= 0; i-- {
		ann := resp[i]
		if p.consideredTime(ann).After(threshold) {
			toPost = append(toPost, ann)
		}
	}
	return
}

func (p *Pipe) Post(anns []api.Ann) (urls []*url.URL, err error) {
	urls = make([]*url.URL, len(anns))
	for i, ann := range anns {
		_, ann.XImageURL, _ = api.GetImageFromMd(ann.Body)
		urls[i], err = p.ee.Post(ann)
		if err != nil {
			return nil, fmt.Errorf("sink: %w", err)
		}
	}
	return
}
