package mysql

import (
	"database/sql"
	"dvhthomas/snippetbox/pkg/models"
	"errors"
)

// SnippetModel is a wrapper around sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// Insert will insert a new snippet in the database
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
		VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// Result has some handy information on what was executed. You can, for example,
	// get the ID of the record.
	// Postgres does NOT support LastInsertId so be warned.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// ID has type int64 so convert to int before returning
	return int(id), nil
}

// Get returns a single snippet based on it's ID
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires from snippets 
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)
	s := &models.Snippet{}
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err != nil {
		// If the query returns no rows then row.Scan() will return
		// a sql.ErrNoRows error. We use the errors.Is() function to check for that.
		// If that's the case we return our own error.

		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}

		return nil, err
	}

	return s, nil
}

// Latest returns the 10 most recently created snippets
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires from snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// We want to make sure the that rows is closed eventually. But make sure that
	// this comes _after_ the error check on Query. If we don't do it after but query
	// did return an error, we'll get a panic when trying to close a nil resultset.
	defer rows.Close()

	snippets := []*models.Snippet{}

	for rows.Next() {
		s := &models.Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	// Just because the loop finished doesn't mean we made it through the whole
	// list of results. For example, we could have lost the DB connection half way
	// through.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went OK then return the slice of Snippets
	return snippets, nil
}
