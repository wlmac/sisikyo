package oauth

// RedirectError is an error from a redirect.
type RedirectError struct {
	Name string
	Desc string
	URI  string
}

func (r RedirectError) Error() string {
	res := r.Name
	if r.Desc != "" {
		res += "\n\n" + r.Desc
	}
	if r.URI != "" {
		res += "\n\n" + r.URI
	}
	return res
}
