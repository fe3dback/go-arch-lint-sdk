package assembler

import (
	"github.com/fe3dback/go-arch-lint-sdk/arch"
)

type (
	projectInfoFetcher interface {
		Fetch() (arch.ProjectInfo, error)
	}

	pathHelper interface {
		FindProjectFiles(query arch.FileQuery) ([]arch.PathDescriptor, error)
	}
)
