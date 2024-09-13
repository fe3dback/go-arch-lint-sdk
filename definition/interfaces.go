package definition

import (
	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/internal/models"
)

type (
	reader interface {
		Read(path arch.PathAbsolute) (models.Config, error)
	}

	validator interface {
		Validate(config models.Config) error
	}

	assembler interface {
		Assemble(conf models.Config) (arch.Spec, error)
	}
)
