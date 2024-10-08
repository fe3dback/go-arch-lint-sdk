package validator

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/internal/models"
)

type DepsVendorsValidator struct{}

func NewDepsVendorsValidator() *DepsVendorsValidator {
	return &DepsVendorsValidator{}
}

func (c *DepsVendorsValidator) Validate(ctx *validationContext) {
	ctx.conf.Dependencies.Map.Each(func(name arch.ComponentName, rules models.ConfigComponentDependencies, reference arch.Reference) {
		existVendors := make(map[arch.VendorName]any)

		for _, vendorName := range rules.CanUse {
			// check multiple usage
			if _, ok := existVendors[vendorName.Value]; ok {
				ctx.AddNotice(
					fmt.Sprintf("Vendor '%s' dublicated in '%s' deps", vendorName.Value, name),
					vendorName.Ref,
				)
			}

			existVendors[vendorName.Value] = struct{}{}

			// check is known
			if !ctx.IsKnownVendor(vendorName.Value) {
				ctx.AddNotice(
					fmt.Sprintf("Vendor '%s' is not defined", vendorName.Value),
					vendorName.Ref,
				)
			}
		}
	})
}
