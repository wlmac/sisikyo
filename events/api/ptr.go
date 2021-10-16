package api

import (
	"errors"
)

type AnnID uint

func (id AnnID) Deref(c *Client) (Ann, error) {
	resp := AnnResp{}
	err := c.Do(AnnReq{}, &resp)
	if err != nil {
		return Ann{}, err
	}
	for _, ann := range resp {
		if ann.Id == id {
			return ann, nil
		}
	}
	return Ann{}, errors.New("not found")
}

type Username string

func (username Username) Deref(c *Client) (resp UserResp, err error) {
	err = c.Do(UserReq{
		Username: string(username),
	}, &resp)
	return
}

type TagID = uint

type OrgName = string

type OrgSlug string

func (slug OrgSlug) Deref(c *Client) (org Org, err error) {
	resp := OrgsResp{}
	err = c.Do(OrgsReq{}, &resp)
	if err != nil {
		return Org{}, err
	}
	for _, org := range resp {
		if org.Slug == slug {
			return org, nil
		}
	}
	return Org{}, errors.New("not found")
}

type EventID = int // negative is evs not on server (generated)

type TermID = uint

type CourseID = uint

type CourseCode string

func (code CourseCode) Deref(c *Client) (Course, error) {
	resp := MeTimetableResp{}
	err := c.Do(MeTimetableReq{c.auth}, &resp)
	if err != nil {
		return Course{}, err
	}
	for _, course := range resp.Courses {
		if course.Code == code {
			return course, nil
		}
	}
	return Course{}, errors.New("not found")
}
