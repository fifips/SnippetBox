package models

import (
	"database/sql"
	"errors"
	"time"
)

type SnippetModelInterface interface {
	Insert(title, content string, expires int) (int, error)
	Get(id int) (*Snippet, error)
	Latest() ([]*Snippet, error)
}

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

// Insert into database snippet with given title, content and
// expiration date set x (specified by expires parameter) days form current date
func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
			VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	res, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Get returns snippet with given id
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets 
				WHERE id = ? AND EXPIRES > UTC_TIMESTAMP`
	row := m.DB.QueryRow(stmt, id)

	var res Snippet
	err := row.Scan(&res.ID, &res.Title, &res.Content, &res.Created, &res.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return &res, nil
}

// Latest returns max 10 latest snippets ordered by creation order from latest to oldest
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	var res []*Snippet
	stmt := `SELECT id, title, content, created, expires FROM snippets 
				WHERE EXPIRES > UTC_TIMESTAMP ORDER BY id DESC LIMIT 10`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		snippet := &Snippet{}
		err := rows.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
		if err != nil {
			return res, err
		}
		res = append(res, snippet)
	}

	if err = rows.Err(); err != nil {
		return res, err
	}

	return res, nil
}
