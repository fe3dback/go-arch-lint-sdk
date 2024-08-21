package create

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint-sdk/tests/projects/mvc/internal/models"
)

type (
	booksRepository interface {
		UpsertAuthor(author models.Author) error
		CreateBook(author models.Author, title string) (models.Book, error)
	}
)

type Service struct {
	booksRepository booksRepository
}

func NewService(
	booksRepository booksRepository,
) *Service {
	return &Service{
		booksRepository: booksRepository,
	}
}

func (s *Service) CreateBook(author models.Author, title string) (models.Book, error) {
	// some business logic
	if title == "Hello World" {
		return models.Book{}, fmt.Errorf("book name is too common")
	}

	// store data in repo
	err := s.booksRepository.UpsertAuthor(author)
	if err != nil {
		return models.Book{}, fmt.Errorf("failed to upsert author: %w", err)
	}

	book, err := s.booksRepository.CreateBook(author, title)
	if err != nil {
		return models.Book{}, fmt.Errorf("failed to create book: %w", err)
	}

	return book, nil
}
