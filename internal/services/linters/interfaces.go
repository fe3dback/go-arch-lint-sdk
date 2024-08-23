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
		IsSuitable(roSpec *arch.Spec) bool
		Lint(lCtx *lintContext, roSpec *arch.Spec) ([]arch.LinterNotice, error)
	}
)
