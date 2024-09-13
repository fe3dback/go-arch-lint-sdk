package definition

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/internal/models"
)

// fromConfig for internal use only. This method should be called
// only from top level factories
func (def *Definition) fromConfig(config models.Config) (arch.Spec, error) {
	err := def.validator.Validate(config)
	if err != nil {
		return arch.Spec{}, fmt.Errorf("invalid config: %w", err)
	}

	spec, err := def.assembler.Assemble(config)
	if err != nil {
		return arch.Spec{}, fmt.Errorf("failed assemble spec: %w", err)
	}

	return spec, nil
}
