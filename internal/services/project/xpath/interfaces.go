package xpath

import (
	"github.com/fe3dback/go-arch-lint-sdk/arch"
)

type (
	typeMatcher interface {
		match(ctx *queryContext, query arch.FileQuery) ([]arch.FileDescriptor, error)
	}

	fileScanner interface {
		Scan(scanDirectory string, fn func(path string, isDir bool) error) error
	}
)
