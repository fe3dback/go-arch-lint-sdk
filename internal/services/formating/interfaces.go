package formating

import (
	"github.com/fe3dback/go-arch-lint-sdk/pkg/tpl/codeprinter"
)

type renderer interface {
	RegisterTemplate(id string, text []byte) error
	Render(id string, model any) (string, error)
}

type printer interface {
	Print(ref codeprinter.Reference, opts codeprinter.CodePrintOpts) (string, error)
}
