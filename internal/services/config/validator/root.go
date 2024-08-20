package validator

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/internal/models"
)

type Root struct {
	skipMissuse bool
	validators  []internalValidator
}

func NewRoot(skipMissuse bool, validators ...internalValidator) *Root {
	return &Root{
		skipMissuse: skipMissuse,
		validators:  validators,
	}
}

func (v *Root) Validate(config models.Config) error {
	ctx := &validationContext{
		conf:             config,
		notices:          make([]arch.Notice, 0, 16),
		currentValidator: "main",
	}

	for _, validator := range v.validators {
		ctx.currentValidator = fmt.Sprintf("%T", validator)
		validator.Validate(ctx)

		if ctx.critical {
			break
		}
	}

	if len(ctx.notices) > 0 {
		return arch.NewErrorWithNotices(
			"Config validator find some notices",
			ctx.notices,
		)
	}

	if !v.skipMissuse && len(ctx.missUsage) > 0 {
		return arch.NewErrorWithNotices(
			"Config validator find miss usages. You can hide this message by adding '--skip-missuse' flag",
			ctx.missUsage,
		)
	}

	return nil
}
