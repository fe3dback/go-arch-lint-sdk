package linters

import (
	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/internal/services/linters/deepscan"
)

type Deepscan struct {
	ast *deepscan.ProjectAstFetcher
}

func NewDeepscan(ast *deepscan.ProjectAstFetcher) *Deepscan {
	return &Deepscan{
		ast: ast,
	}
}

func (o *Deepscan) Information() arch.Linter {
	return arch.Linter{
		ID:          arch.LinterIDDeepScan,
		Name:        "Deepscan",
		Description: "Check actual code injections through interfaces",
	}
}

func (o *Deepscan) IsSuitable(_ *lintContextReadOnly) bool {
	return false
}

func (o *Deepscan) Lint(lCtx *lintContext) error {
	componentIDs := o.findActiveComponentIDs(lCtx)
	if len(componentIDs) <= 0 {
		return nil
	}

	for _, componentID := range componentIDs {
		methods := o.findPublicMethods(lCtx, componentID)

		for _, method := range methods {
			for _, gate := range method.Gates {
				if !gate.IsInterfaceType {
					continue
				}

				// found method with interface gate(argument)
				// we need to check this, possible some bad thing
				// is injected through this interface into component

				// todo: check all ast packages
				// todo: if package contain import to any of owned packages in this component
				// todo: we can check decl and found exactly injector of this method
				// injectors := o.findInjectors(componentID, method)
			}
		}
	}

	return nil
}

// findActiveComponentIDs will extract all components from spec, with deepscan=true
func (o *Deepscan) findActiveComponentIDs(lCtx *lintContext) []arch.ComponentName {
	list := make([]arch.ComponentName, 0, len(lCtx.ro.spec.Components))

	for cmpID := range lCtx.ro.spec.Components {
		cmp := lCtx.ro.spec.Components[cmpID]
		if !cmp.DeepScan.Value {
			continue
		}

		list = append(list, cmpID)
	}

	return list
}

// findPublicMethods will find all public methods (functions) in all packages
// that belongs to this component. Our main target is find
// methods like `books.NewRepository(db)`.
// this useful for checking gates [example: `db`], and verify
// that component `cmp` can depend on owner of `db`
func (o *Deepscan) findPublicMethods(lCtx *lintContext, cmpID arch.ComponentName) []deepscan.Method {
	cmp := lCtx.ro.spec.Components[cmpID]
	list := make([]deepscan.Method, 0, 4)

	for _, ownedPackage := range cmp.OwnedPackages {
		for _, astPackage := range lCtx.ro.packages[ownedPackage.PathRel] {
			list = append(list, o.ast.FindPublicMethods(astPackage)...)
		}
	}

	return list
}
