package mock

import (
	"dvhthomas/snippetbox/pkg/models"
	"time"
)

var mockSnippet = &models.Snippet{
	ID:      1,
	Title:   "An old silent pond",
	Content: "And old silent pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

// SnippetModel for non-existent database
type SnippetModel struct{}

// Insert a fake record
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	return 2, nil
}

// Get a predictable value
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}

// Latest containing known records
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}
