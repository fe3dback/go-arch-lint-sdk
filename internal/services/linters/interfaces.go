package linters

import (
	"github.com/fe3dback/go-arch-lint-sdk/arch"
)

type (
	// all methods will be called in goroutines
	// roSpec should not be changed anywhere
	// lCtx must be guarded when necessary
	linter interface {
		Information() arch.Linter
		IsSuitable(lCtx *lintContextReadOnly) bool
		Lint(lCtx *lintContext) error
	}
)
