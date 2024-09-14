package validator

import (
	"fmt"

	"github.com/gobwas/glob"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/internal/models"
)

type VendorsValidator struct {
	pathHelper pathHelper
}

func NewVendorsValidator(
	pathHelper pathHelper,
) *VendorsValidator {
	return &VendorsValidator{
		pathHelper: pathHelper,
	}
}

func (c *VendorsValidator) Validate(ctx *validationContext) {
	ctx.conf.Vendors.Map.Each(func(_ arch.VendorName, vendor models.ConfigVendor, _ arch.Reference) {
		for _, pathGlob := range vendor.In {
			_, err := glob.Compile(string(pathGlob.Value), '/')
			if err != nil {
				ctx.AddNotice(
					fmt.Sprintf("invalid glob path '%s': %v", pathGlob.Value, err),
					pathGlob.Ref,
				)
				return
			}
		}
	})
}
