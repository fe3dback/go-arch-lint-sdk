package validator

import (
	"github.com/fe3dback/go-arch-lint-sdk/arch"
)

type internalValidator interface {
	Validate(ctx *validationContext)
}

type pathHelper interface {
	FindProjectFiles(query arch.FileQuery) ([]arch.PathDescriptor, error)
}
