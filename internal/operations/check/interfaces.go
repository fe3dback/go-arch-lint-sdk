package check

import (
	"github.com/fe3dback/go-arch-lint-sdk/arch"
)

type (
	linter interface {
		Lint(spec arch.Spec) ([]arch.LinterResult, error)
	}
)
