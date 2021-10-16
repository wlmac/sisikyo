package api

import (
	"net/http"
	"time"
)

type Auth struct {
	Access    string    `json:"access"`
	Refresh   string    `json:"refresh"`
	Generated time.Time `json:"-"`
}

func (a Auth) aHeader() string               { return "Bearer " + a.Access }
func (a Auth) rHeader() string               { return "Bearer " + a.Refresh }
func (a Auth) setRHeader(header http.Header) { header.Set("Authorization", a.rHeader()) }
func (a Auth) setAHeader(header http.Header) { header.Set("Authorization", a.aHeader()) }
