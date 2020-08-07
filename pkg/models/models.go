package models

import (
	"errors"
	"time"
)

// ErrNoRecord reports if no records are found
var ErrNoRecord = errors.New("models: no matching records found")

// Snippet represents a single snippet in the app
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
