package oauth

// TOOD: properly revoke token on UserRemove

import (
	"context"
	"errors"

	"gitlab.com/mirukakoro/sisikyo/util"
	"golang.org/x/oauth2"
)

// Client is a wrapper that provides convenience methods around oauth2.Config.
type Client struct {
	oauth2.Config
}

// NewClient makes a new Client with config.
func NewClient(config oauth2.Config) *Client {
	return &Client{
		Config: config,
	}
}

var (
	ErrStateMismatch = errors.New("state: unexpected mismatch") // given states do not match
	ErrParamsInvalid = errors.New("params: invalid")            // invalid params
)

// Auth returns a oauth2.Token from a code.
func (o *Client) Auth(ctx context.Context, code string) (*oauth2.Token, error) {
	return o.Config.Exchange(ctx, code)
}

// CheckCode checks if the states are correct (are equal) and returns an appropriate error if necessary
func (o *Client) CheckCode(expectedState string, params RedirectParams) error {
	if !params.IsValid() {
		return ErrParamsInvalid
	}
	if err := params.HasError(); err != nil {
		return err
	}
	if params.State != expectedState {
		return ErrStateMismatch
	}
	return nil
}

// AuthorizeURL returns an authorization URL that the user can access.
func (o *Client) AuthorizeURL() (state string, url string, err error) {
	state, err = util.GenRandom(64)
	if err != nil {
		return
	}
	url = o.Config.AuthCodeURL(string(state))
	return
}
