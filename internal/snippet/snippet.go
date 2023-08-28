package snippet

import (
	"github.com/Prixix/CodeScribe/pkg/database"
)

type Manager struct {
	db *database.Database
}

func NewManager(db *database.Database) *Manager {
	return &Manager{db: db}
}

func (m *Manager) CreateSnippet(title, description, tags, code string) error {
	snippet := database.Snippet{
		Title:       title,
		Description: description,
		Tags:        tags,
		Code:        code,
	}

	_, err := m.db.CreateSnippet(snippet)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) GetSnippetByID(id int) (database.Snippet, error) {
	return m.db.GetSnippetByID(id)
}

func (m *Manager) SearchSnippets(keyword string) ([]database.Snippet, error) {
	return m.db.SearchSnippets(keyword)
}

func (m *Manager) GetAllSnippets() ([]database.Snippet, error) {
	var snippets []database.Snippet
	return snippets, m.db.GetAllSnippets(&snippets)
}

func InitializeSchema(dbPath string) error {
	return database.InitializeSchema(dbPath)
}
