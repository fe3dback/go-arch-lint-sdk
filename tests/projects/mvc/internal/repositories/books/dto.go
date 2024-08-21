package books

type (
	// this struct exist only for checking that
	// 'db' tags is allowed in repo component
	// but not allowed in other components
	bookAuthorDB struct {
		Name     string  `db:"name" json:"Name"`
		Lastname string  `db:"last_name" json:"Lastname"`
		Age      int     `db:"age" json:"Age"`
		Rating   float64 `db:"rating" json:"Rating"`
	}
)
