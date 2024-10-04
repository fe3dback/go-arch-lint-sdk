package codeprinter

type (
	// Reference is copy of arch.Reference
	// it required because /pkg should not depend on SDK models
	Reference struct {
		File   string
		Line   int
		Column int
		Valid  bool
	}
)
