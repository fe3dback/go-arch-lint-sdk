package reader

import (
	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/internal/models"
)

type (
	fileReader interface {
		ReadFile(path arch.PathAbsolute) (models.Config, error)
	}
)
