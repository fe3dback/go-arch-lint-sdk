package definition

import (
	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/internal/models"
)

func (def *Definition) withUserFriendlyError(spec arch.Spec, err error) (arch.Spec, error) {
	// happy path
	if err == nil {
		return spec, nil
	}

	// no not wrap anything if we in CLI context
	if def.usedContext == arch.UsedContextCLI {
		return spec, err
	}

	// already wrapped, skip
	if _, ok := err.(models.SDKError); ok {
		return arch.Spec{}, err
	}

	// todo: add auto preview (if error is referenced. ex: config has syntax problem)
	// wrap
	return arch.Spec{}, models.NewSDKError(err, def.projectPath)
}
