/* snippets.go */

package models

import (
	"database/sql"
	"log"
)

type Snippet struct {
	Id      int
	Title   string
	Content string
	Created string
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string) (int, error) {
	stmt := "INSERT INTO snippets (title, content) VALUES(?, ?)"

	result, err := m.DB.Exec(stmt, title, content)
	if err != nil {
		return 0, err
	}

	// Use the LastInsertId() method to get the ID of our newly inserted record.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (Snippet, error) {
	stmt := "SELECT id, title, content, created  FROM snippets WHERE id = ?"
	row := m.DB.QueryRow(stmt, id)

	// Initialize a new Snippet struct.
	var snippet Snippet
	err := row.Scan(&snippet.Id, &snippet.Title, &snippet.Content, &snippet.Created)
	if err != nil {
		log.Println(err)
	}

	return snippet, err
}

func (m *SnippetModel) Latest() ([]Snippet, error) {
	stmt := "SELECT id, title, content, created FROM snippets ORDER BY id DESC LIMIT 10"
	rows, err := m.DB.Query(stmt)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	// Initialize an empty slice to hold the Snippet structs.
	var snippets []Snippet

	// Once iteration over all the rows completes, the resultset auto closes.
	for rows.Next() {
		// Create a new Snippet struct.
		var snippet Snippet

		err = rows.Scan(&snippet.Id, &snippet.Title, &snippet.Content, &snippet.Created)
		if err != nil {
			log.Println(err)
		}
		// Append it to the slice of snippets structs.
		snippets = append(snippets, snippet)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went OK then return the Snippets slice.
	return snippets, nil
}

func (m *SnippetModel) Update(id int, title string, content string) error {
	stmt := "UPDATE snippets set title=?, content=? where id=?"

	log.Println(id)
	log.Println(title)
	log.Println(content)
	
	_, err := m.DB.Exec(stmt, title, content, id)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (m *SnippetModel) Delete(id int) error {
	stmt := "DELETE FROM snippets WHERE id = ?"
	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		log.Println(err)
	}
	return err
}
