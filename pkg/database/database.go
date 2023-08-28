package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Snippet struct {
	ID          int
	Title       string
	Description string
	Tags        string
	Code        string
}

type Database struct {
	conn *sql.DB
}

func NewDatabase(dbPath string) (*Database, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	return &Database{conn}, nil
}

func (db *Database) Close() error {
	return db.conn.Close()
}

func (db *Database) CreateSnippet(snippet Snippet) (int64, error) {
	stmt, err := db.conn.Prepare(`
		INSERT INTO snippets (title, description, tags, code)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(snippet.Title, snippet.Description, snippet.Tags, snippet.Code)
	if err != nil {
		return -1, err
	}

	return result.LastInsertId()
}

func (db *Database) GetSnippetByID(id int) (Snippet, error) {
	var snippet Snippet

	row := db.conn.QueryRow(`
		SELECT id, title, description, tags, code
		FROM snippets
		WHERE id = ?
	`, id)

	err := row.Scan(&snippet.ID, &snippet.Title, &snippet.Description, &snippet.Tags, &snippet.Code)
	if err != nil {
		return Snippet{}, err
	}
	return snippet, err
}

func (db *Database) SearchSnippets(keyword string) ([]Snippet, error) {
	query := fmt.Sprintf(`
		SELECT id, title, description, tags, code
		FROM snippets
		WHERE title LIKE '%%%s%%' OR description LIKE '%%%s%%' OR tags LIKE '%%%s%%'
	`, keyword, keyword, keyword)

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snippets []Snippet

	for rows.Next() {
		var snippet Snippet
		err := rows.Scan(&snippet.ID, &snippet.Title, &snippet.Description, &snippet.Tags, &snippet.Code)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		snippets = append(snippets, snippet)
	}

	return snippets, nil
}

func InitializeSchema(dbPath string) error {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Exec(`
		CREATE TABLE IF NOT EXISTS snippets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT,
			tags TEXT,
			code TEXT
		);
	`)

	return err
}
