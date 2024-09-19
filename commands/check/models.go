package check

import "github.com/fe3dback/go-arch-lint-sdk/arch"

type (
	In struct {
		// Check go AST syntax? (default: true)
		CheckSyntax bool

		// Output notices limit
		MaxWarnings int
	}

	Out struct {
		ModuleName   arch.GoModule       `json:"ModuleName"`
		NoticesCount int                 `json:"NoticesCount"`
		OmittedCount int                 `json:"OmittedCount"`
		Linters      []arch.LinterResult `json:"Linters"`
	}
)
