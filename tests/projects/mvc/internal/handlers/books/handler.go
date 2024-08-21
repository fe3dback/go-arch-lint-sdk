package books

import "github.com/fe3dback/go-arch-lint-sdk/tests/projects/mvc/internal/models"

type (
	booksService interface {
		CreateBook(author models.Author, title string) (models.Book, error)
	}

	// !IF[HANDLER_CALL_REPO_DIRECTLY]{
	booksRepository interface {
		UpsertAuthor(author models.Author) error
	}
	// !}
)

type Handler struct {
	booksService booksService
	// !IF[HANDLER_CALL_REPO_DIRECTLY]{
	booksRepository booksRepository
	// !}
}

func NewHandler(
	booksService booksService,
	// !IF[HANDLER_CALL_REPO_DIRECTLY]{
	booksRepository booksRepository,
	// !}
) *Handler {
	return &Handler{
		booksService: booksService,
		// !IF[HANDLER_CALL_REPO_DIRECTLY]{
		booksRepository: booksRepository,
		// !}
	}
}
