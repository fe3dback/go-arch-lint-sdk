package validator

import (
	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/pkg/tpl/codeprinter"
)

type internalValidator interface {
	Validate(ctx *validationContext)
}

type pathHelper interface {
	FindProjectFiles(query arch.FileQuery) ([]arch.PathDescriptor, error)
}

type printer interface {
	Print(ref codeprinter.Reference, opts codeprinter.CodePrintOpts) (string, error)
}
