package oauth

type RedirectParams struct {
	State string `form:"state"`
	Code  string `form:"code"`

	Error     string `form:"error"`
	ErrorDesc string `form:"error_description"`
	ErrorURI  string `form:"error_uri"`
}

func (r RedirectParams) IsValid() bool {
	return r.Code != "" || r.Error != ""
}

func (r RedirectParams) IsError() bool {
	return r.Error != ""
}

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
