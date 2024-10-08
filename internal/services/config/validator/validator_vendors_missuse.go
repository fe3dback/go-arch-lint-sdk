package validator

import (
	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/internal/models"
)

type VendorsMissUseValidator struct {
}

func NewVendorsMissUseValidator() *VendorsMissUseValidator {
	return &VendorsMissUseValidator{}
}

func (c *VendorsMissUseValidator) Validate(ctx *validationContext) {
	if !ctx.conf.Settings.Imports.AllowAnyVendorImports.Value {
		return
	}

	msg := "there is no point in vendors, since a global flag 'settings.imports.AllowAnyVendorImports' is enabled. (linter will not check vendors)"
	for _, cmnVendor := range ctx.conf.CommonVendors {
		ctx.AddMissUse(
			msg,
			cmnVendor.Ref,
		)
		break
	}

	first := true
	ctx.conf.Vendors.Map.Each(func(_ arch.VendorName, _ models.ConfigVendor, reference arch.Reference) {
		if !first {
			return
		}
		first = false

		ctx.AddMissUse(
			msg,
			reference,
		)
	})

	ctx.conf.Dependencies.Map.Each(func(_ arch.ComponentName, rules models.ConfigComponentDependencies, _ arch.Reference) {
		for _, vendorName := range rules.CanUse {
			ctx.AddMissUse(
				msg,
				vendorName.Ref,
			)
			break
		}
	})
}
