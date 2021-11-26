package oauth

// RedirectParams stores all URL params relevant to OAuth redirects to OAuth clients.
type RedirectParams struct {
	State string `form:"state"`
	Code  string `form:"code"`

	Error     string `form:"error"`
	ErrorDesc string `form:"error_description"`
	ErrorURI  string `form:"error_uri"`
}

// IsValid returns whether the URL params signify that the redirect is valid.
func (r RedirectParams) IsValid() bool {
	return r.Code != "" || r.Error != ""
}

// IsError returns whether the URL params signify that the redirect describe an error.
//
// Deprecated: use HasError and check if the error is nil.
func (r RedirectParams) IsError() bool {
	return r.Error != ""
}

// HasError returns an error from the URL params if applicable.
func (r RedirectParams) HasError() error {
	if r.Error == "" {
		return nil
	}
	return RedirectError{
		Name: r.Error,
		Desc: r.ErrorDesc,
		URI:  r.ErrorURI,
	}
}
