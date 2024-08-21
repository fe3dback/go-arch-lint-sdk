package books

import (
	"errors"
	"fmt"

	"github.com/fe3dback/go-arch-lint-sdk/tests/projects/mvc/internal/models"
)

func (h *Handler) CreateBook(author models.Author, title string) (models.Book, error) {
	// input validation
	if title == "" {
		return models.Book{}, errors.New("title is required")
	}

	if author.Name == "" {
		return models.Book{}, errors.New("author name is required")
	}

	if author.ID == 0 {
		return models.Book{}, errors.New("author id is required")
	}

	// call business service
	book, err := h.booksService.CreateBook(author, title)
	if err != nil {
		return models.Book{}, fmt.Errorf("error creating book: %w", err)
	}

	// !IF[HANDLER_CALL_REPO_DIRECTLY]{
	err = h.booksRepository.UpsertAuthor(author)
	if err != nil {
		return models.Book{}, fmt.Errorf("error updating author: %w", err)
	}
	// !}

	return book, nil
}
