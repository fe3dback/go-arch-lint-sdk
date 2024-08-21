package models

type (
	BookID string

	AuthorID int64

	Author struct {
		ID   AuthorID
		Name string
	}

	Book struct {
		ID     BookID
		Author Author
		Title  string
	}
)
