package validator

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/internal/models"
)

type DepsComponentsValidator struct{}

func NewDepsComponentsValidator() *DepsComponentsValidator {
	return &DepsComponentsValidator{}
}

func (c *DepsComponentsValidator) Validate(ctx *validationContext) {
	ctx.conf.Dependencies.Map.Each(func(name arch.ComponentName, rules models.ConfigComponentDependencies, reference arch.Reference) {
		existComponents := make(map[arch.ComponentName]any)

		for _, anotherComponentName := range rules.MayDependOn {
			// check multiple usage
			if _, ok := existComponents[anotherComponentName.Value]; ok {
				ctx.AddNotice(
					fmt.Sprintf("Component '%s' dublicated in '%s' deps", anotherComponentName.Value, name),
					anotherComponentName.Ref,
				)
			}

			existComponents[anotherComponentName.Value] = struct{}{}

			// check is known
			if !ctx.IsKnownComponent(anotherComponentName.Value) {
				ctx.AddNotice(
					fmt.Sprintf("Component '%s' is not defined", anotherComponentName.Value),
					anotherComponentName.Ref,
				)
			}
		}
	})
}
