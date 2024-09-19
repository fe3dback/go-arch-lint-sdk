package linters

import (
	"fmt"

	"golang.org/x/tools/go/packages"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
)

type Syntax struct {
}

func NewSyntax() *Syntax {
	return &Syntax{}
}

func (o *Syntax) Information() arch.Linter {
	return arch.Linter{
		ID:                  arch.LinterIDSyntax,
		Name:                "Go Syntax",
		Description:         "Check that go files has correct go AST code (if not, checker can guarantee correct work of other arch-linters)",
		EnableConditionHint: "always on",
	}
}

func (o *Syntax) IsSuitable(lCtx *lintContextReadOnly) bool {
	return lCtx.options.CheckSyntax
}

func (o *Syntax) Lint(lCtx *lintContext) error {
	for ind := range lCtx.ro.spec.Components {
		o.checkComponent(lCtx, lCtx.ro.spec.Components[ind])
	}

	return nil
}

func (o *Syntax) checkComponent(lCtx *lintContext, component *arch.SpecComponent) {
	for _, ownedPackage := range component.OwnedPackages {
		astPackages := lCtx.ro.packages[ownedPackage.PathRel]
		for _, astPackage := range astPackages {
			if len(astPackage.Errors) == 0 {
				continue
			}

			for _, err := range astPackage.Errors {
				if err.Kind != packages.TypeError {
					continue
				}

				lCtx.state.AddNotice(arch.LinterNotice{
					Message: fmt.Sprintf("Failed parse component '%s' package '%s' at '%s': invalid syntax: %v",
						component.Name.Value,
						astPackage.Name,
						ownedPackage.PathRel,
						err.Msg,
					),
					Reference: referenceFromAstPos(err.Pos),
					Details: arch.LinterNoticeDetails{
						LinterID: arch.LinterIDSyntax,
					},
				})
			}

		}
	}
}
