package deepscan

import (
	"go/ast"
	"golang.org/x/tools/go/packages"
)

type ProjectAstFetcher struct {
}

func NewProjectAstFetcher() *ProjectAstFetcher {
	return &ProjectAstFetcher{}
}

// FindPublicMethods will find all public methods (functions) in all packages
// that belongs to this component. Our main target is find
// methods like `books.NewRepository(db)`.
// this useful for checking gates [example: `db`], and verify
// that component `cmp` can depend on owner of `db`
func (o *ProjectAstFetcher) FindPublicMethods(astPackage *packages.Package) []Method {
	list := make([]Method, 0, 4)

	for _, astFile := range astPackage.Syntax {
		list = append(list, o.extractPublicMethodsFromAstFile(astFile)...)
	}

	return list
}

// extractPublicMethodsFromAstFile will return all public methods (functions)
// containing in this go file and not attached to any structs
// example: `func New(a, b) {}`
func (o *ProjectAstFetcher) extractPublicMethodsFromAstFile(astFile *ast.File) []Method {
	// todo:
	//for _, decl := range astFile.Decls {
	//
	//}
	return nil
}
