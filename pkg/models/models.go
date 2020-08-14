package models

import (
	"errors"
	"time"
)

// ErrNoRecord reports if no records are found
var ErrNoRecord = errors.New("models: no matching records found")

// ErrDuplicateEmail reports if a user already exists based on email address
var ErrDuplicateEmail = errors.New("models: existing user with that email")

// ErrInvalidCredentials when the user does not exist in a login or the password is invalid
var ErrInvalidCredentials = errors.New("models: invalid user credentials")

// Snippet represents a single snippet in the app
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// User that owns snippets and can log in
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Active         bool
}
