package check

import (
	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/internal/models"
)

type (
	linter interface {
		Lint(spec arch.Spec, options models.LintOptions) ([]arch.LinterResult, error)
	}
)
