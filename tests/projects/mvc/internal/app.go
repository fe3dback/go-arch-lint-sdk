package internal

import (
	"fmt"

	booksService "github.com/fe3dback/go-arch-lint-sdk/tests/projects/mvc/internal/domains/library/services/books/create"
	booksHandler "github.com/fe3dback/go-arch-lint-sdk/tests/projects/mvc/internal/handlers/books"
	"github.com/fe3dback/go-arch-lint-sdk/tests/projects/mvc/internal/models"
	booksRepo "github.com/fe3dback/go-arch-lint-sdk/tests/projects/mvc/internal/repositories/books"
)

func Execute() error {
	// di
	repo := booksRepo.NewRepository()
	service := booksService.NewService(repo)
	handler := booksHandler.NewHandler(
		service,
		// !IF[HANDLER_CALL_REPO_DIRECTLY]{
		repo,
		// !}
	)

	// execute logic
	book, err := handler.CreateBook(models.Author{
		ID:   444,
		Name: "Red Giant Media",
	}, "PHP for dummies, 6+")
	if err != nil {
		return fmt.Errorf("failed create book: %w", err)
	}

	fmt.Println(book.ID)
	return nil
}
