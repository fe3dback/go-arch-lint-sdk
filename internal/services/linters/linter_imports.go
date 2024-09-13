package linters

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
)

type Imports struct {
}

type (
	importType string
)

const (
	importTypeStd     importType = "std"
	importTypeVendor  importType = "vendor"
	importTypeProject importType = "project"
)

func NewImports() *Imports {
	return &Imports{}
}

func (o *Imports) Information() arch.Linter {
	return arch.Linter{
		ID:                  arch.LinterIDComponentImports,
		Name:                "Imports",
		Description:         "Check that packages contain only allowed \"import\" statements",
		EnableConditionHint: "always on",
	}
}

func (o *Imports) IsSuitable(_ *arch.Spec) bool {
	return true
}

func (o *Imports) Lint(lCtx *lintContext) error {
	for ind := range lCtx.ro.spec.Components {
		o.checkComponent(lCtx, lCtx.ro.spec.Components[ind])
	}

	return nil
}

func (o *Imports) checkComponent(lCtx *lintContext, component *arch.SpecComponent) {
	for _, ownedPackage := range component.OwnedPackages {
		astPackages := lCtx.ro.packages[ownedPackage.PathRel]
		for _, astPackage := range astPackages {
			for _, astFile := range astPackage.Syntax {
				o.checkFile(lCtx, component, astFile)
			}
		}
	}
}

func (o *Imports) checkFile(lCtx *lintContext, component *arch.SpecComponent, astFile *ast.File) {
	for _, astImport := range astFile.Imports {
		importPath := arch.PathImport(strings.Trim(astImport.Path.Value, "\""))
		importRef := referenceFromAstToken(lCtx.ro.fileSet.Position(astImport.Pos()))
		importType := o.resolveImportType(lCtx, importPath)

		switch importType {
		case importTypeStd:
			// all std imports is always allowed
			continue
		case importTypeProject:
			if !o.isProjectImportAllowed(lCtx, component, importPath) {
				importOwnerName := o.findProjectImportPathOwner(lCtx, importPath)

				lCtx.state.AddNotice(arch.LinterNotice{
					Message: fmt.Sprintf("Component '%s' shouldn't depend on '%s' ('%s') in '%s:%d'",
						component.Name.Value,
						orDefault(importOwnerName, "unknown"),
						importPath,
						importRef.File,
						importRef.Line,
					),
					Reference: importRef,
					Details: arch.LinterNoticeDetails{
						LinterID: arch.LinterIDComponentImports,
						LinterIDComponentImports: &arch.LinterImportDetails{
							ComponentName:      component.Name.Value,
							FileRelativePath:   arch.PathRelative(strings.TrimPrefix(string(importRef.File), string(lCtx.ro.spec.Project.Directory))),
							FileAbsolutePath:   importRef.File,
							ResolvedImportName: importPath,
							Reference:          importRef,
						},
					},
				})
			}
		case importTypeVendor:
			if !o.isVendorImportAllowed(lCtx, component, importPath) {
				lCtx.state.AddNotice(arch.LinterNotice{
					Message: fmt.Sprintf("Component '%s' shouldn't use '%s' in '%s:%d'",
						component.Name.Value,
						importPath,
						importRef.File,
						importRef.Line,
					),
					Reference: importRef,
					Details: arch.LinterNoticeDetails{
						LinterID: arch.LinterIDVendorImports,
						LinterIDVendorImports: &arch.LinterImportDetails{
							ComponentName:      component.Name.Value,
							FileRelativePath:   arch.PathRelative(strings.TrimPrefix(string(importRef.File), string(lCtx.ro.spec.Project.Directory))),
							FileAbsolutePath:   importRef.File,
							ResolvedImportName: importPath,
							Reference:          importRef,
						},
					},
				})
			}
			continue
		}
	}
}

func (o *Imports) isProjectImportAllowed(lCtx *lintContext, component *arch.SpecComponent, importPath arch.PathImport) bool {
	if component.AllowAllProjectDeps.Value {
		return true
	}

	for _, dependComponentID := range component.MayDependOn {
		dependComponent := lCtx.ro.spec.Components[dependComponentID.Value]

		for _, dependOwnedPackage := range dependComponent.OwnedPackages {
			if dependOwnedPackage.Import == importPath {
				return true
			}
		}
	}

	return false
}

func (o *Imports) isVendorImportAllowed(lCtx *lintContext, component *arch.SpecComponent, importPath arch.PathImport) bool {
	if component.AllowAllVendorDeps.Value {
		return true
	}

	for _, dependVendorID := range component.CanUse {
		dependVendor := lCtx.ro.spec.Vendors[dependVendorID.Value]

		for _, vendorImport := range dependVendor.OwnedImports {
			glob := vendorImport.Value

			_ = glob
			// todo: if glob match (import) { return true }
		}
	}

	// todo: another linter ID:
	// todo: -this import %s matched by many vendors, but expected only one

	return false
}

func (o *Imports) findProjectImportPathOwner(lCtx *lintContext, importPath arch.PathImport) *arch.ComponentName {
	for _, component := range lCtx.ro.spec.Components {
		for _, ownedPackage := range component.OwnedPackages {
			if ownedPackage.Import == importPath {
				return &component.Name.Value
			}
		}
	}

	return nil
}

func (o *Imports) resolveImportType(lCtx *lintContext, importPath arch.PathImport) importType {
	if _, exist := lCtx.ro.stdPackageIDs[importPath]; exist {
		return importTypeStd
	}

	if strings.HasPrefix(string(importPath), string(lCtx.ro.spec.Project.Module)) {
		return importTypeProject
	}

	return importTypeVendor
}
