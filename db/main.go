// Package db provides types for storing data in SQL databases using https://github.com/jmoiron/sqlx.
package db

import (
	"time"

	"golang.org/x/oauth2"
)

// Schema is the SQL schema of everything.
const Schema = UserSchema

const UserRemove = `
DELETE FROM user_table WHERE control = :control;
`

const UserRegister = `
INSERT INTO user_table
(control, access_token, refresh_token, token_type, expiry)
VALUES
(:control, :access_token, :refresh_token, :token_type, :expiry);
`

const UserQuery = `
SELECT * FROM user_table WHERE control = :control;
`

// UserSchema is the SQL schema of User.
const UserSchema = `
CREATE TABLE IF NOT EXISTS user_table (
	control varchar(128) NOT NULL PRIMARY KEY,
	access_token text NOT NULL,
	refresh_token text NOT NULL,
	token_type text NOT NULL,
	expiry timestamp NOT NULL
);
`

// User stores a user's identity and access to the API.
type User struct {
	Control      string    `db:"control"` // random string of base64 characters
	AccessToken  string    `db:"access_token"`
	RefreshToken string    `db:"refresh_token"`
	TokenType    string    `db:"token_type"`
	Expiry       time.Time `db:"expiry"`
}

// Token returns a partial oauth2.Token (extra data is dropped) from User.
func (u User) Token() *oauth2.Token {
	return &oauth2.Token{
		AccessToken:  u.AccessToken,
		RefreshToken: u.RefreshToken,
		TokenType:    u.TokenType,
		Expiry:       u.Expiry,
	}
}

// ApplyToken sets some fields from an oauth2.Token.
func (u *User) ApplyToken(tok *oauth2.Token) {
	u.AccessToken = tok.AccessToken
	u.RefreshToken = tok.RefreshToken
	u.TokenType = tok.TokenType
	u.Expiry = tok.Expiry
}
