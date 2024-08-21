package books

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint-sdk/tests/projects/mvc/internal/models"
)

type Repository struct {
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) UpsertAuthor(author models.Author) error {
	// SQL select author ...
	foundRow := bookAuthorDB{
		Name:     author.Name,
		Lastname: "Some",
		Age:      37,
		Rating:   0.84,
	}

	// some nonsense check
	if foundRow.Age >= 100 {
		return fmt.Errorf("age is too big")
	}

	// SQL update book ...

	return nil
}

func (r *Repository) CreateBook(author models.Author, title string) (models.Book, error) {
	// SQL insert ...

	return models.Book{
		ID:     "123",
		Author: author,
		Title:  title,
	}, nil
}
