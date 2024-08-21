package arch

type UsedContext int

const (
	// UsedContextDefault should be used by default, except rare cased described in UsedContextCLI
	UsedContextDefault UsedContext = 0

	// UsedContextCLI Marks that SDK used from CLI program
	// this affects small details, e.g. validation errors will not be pasted into error,
	// but will be formatted separately in stdout.
	//
	// DO NOT USE anywhere, except go-arch-lint CLI program, or analog
	UsedContextCLI UsedContext = 1
)
