package definition

import (
	"github.com/fe3dback/go-arch-lint-sdk/arch"
)

type (
	Definition struct {
		projectPath arch.PathAbsolute
		reader      reader
		validator   validator
		assembler   assembler
	}
)

func NewDefinition(
	projectPath arch.PathAbsolute,
	reader reader,
	validator validator,
	assembler assembler,
) *Definition {
	return &Definition{
		projectPath: projectPath,
		reader:      reader,
		validator:   validator,
		assembler:   assembler,
	}
}
