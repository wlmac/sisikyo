package api

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrMultipleLateStarts = errors.New("multiple late starts found")
	ErrSchoolNotOpen      = errors.New("school not open")
)

type StatusError struct {
	code int
}

var _ error = StatusError{}

func (err StatusError) Error() string {
	return fmt.Sprintf("status code %d %s", err.code, http.StatusText(err.code))
}
