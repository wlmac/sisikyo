// Package db provides types for storing data in SQL databases using https://github.com/jmoiron/sqlx.
package db

// Schema is the SQL schema of everything.
const Schema = UserSchema

// UserSchema is the SQL schema of User.
const UserSchema = `
CREATE TABLE user (
	username text,
	code text
);
`

// User stores a user's identity and access to the API.
type User struct {
	Username string `db:"username"`
	Code     string `db:"code"`
}
