package linters

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
)

type Orphans struct {
}

func NewOrphans() *Orphans {
	return &Orphans{}
}

func (o *Orphans) Information() arch.Linter {
	return arch.Linter{
		ID:                  arch.LinterIDOrphans,
		Name:                "Orphans",
		Description:         "Will return list of all files that not mapped to any components specified in config",
		EnableConditionHint: "always on",
	}
}

func (o *Orphans) IsSuitable(_ *arch.Spec) bool {
	return true
}

func (o *Orphans) Lint(_ *lintContext, spec *arch.Spec) ([]arch.LinterNotice, error) {
	notices := make([]arch.LinterNotice, 0, len(spec.Orphans))

	for _, orphan := range spec.Orphans {
		notices = append(notices, arch.LinterNotice{
			Message:   fmt.Sprintf("File '%s' not attached to any component in archfile", orphan.File.PathRel),
			Reference: arch.NewFileReference(orphan.File.PathAbs),
			Details: arch.LinterNoticeDetails{
				LinterID: arch.LinterIDOrphans,
				LinterIDOrphan: &arch.LinterOrphanDetails{
					FileRelativePath: orphan.File.PathRel,
					FileAbsolutePath: orphan.File.PathAbs,
				},
			},
		})
	}

	return notices, nil
}
