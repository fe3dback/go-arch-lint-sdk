package check

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/commands/check"
	"github.com/fe3dback/go-arch-lint-sdk/internal/models"
)

type Operation struct {
	linter linter
}

func NewOperation(
	linter linter,
) *Operation {
	return &Operation{
		linter: linter,
	}
}

func (o *Operation) Execute(spec arch.Spec, in check.In) (check.Out, error) {
	lintOptions := models.LintOptions{
		CheckSyntax: in.CheckSyntax,
	}

	linters, err := o.linter.Lint(spec, lintOptions)
	if err != nil {
		return check.Out{}, fmt.Errorf("linters failed: %w", err)
	}

	return check.Out{
		ModuleName:   spec.Project.Module,
		Linters:      linters,
		NoticesCount: calculateNotices(linters),
		OmittedCount: 0, // todo: limit
	}, nil
}

func calculateNotices(results []arch.LinterResult) int {
	count := 0

	for _, result := range results {
		count += len(result.Notices)
	}

	return count
}
