package validator

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/internal/models"
)

type validationContext struct {
	conf             models.Config
	notices          []arch.Notice // critical (will block linter for working)
	missUsage        []arch.Notice // optional (will block linter, if arg --skip-missuse is not set)
	critical         bool
	currentValidator string // for error prefixing
}

func (vc *validationContext) SkipOtherValidators() {
	vc.critical = true
}

func (vc *validationContext) AddNotice(message string, ref arch.Reference) {
	vc.notices = append(vc.notices, arch.Notice{
		Message:   fmt.Sprintf("%s: %s", vc.currentValidator, message),
		Reference: ref,
	})
}

func (vc *validationContext) AddMissUse(message string, ref arch.Reference) {
	vc.missUsage = append(vc.missUsage, arch.Notice{
		Message:   fmt.Sprintf("%s: %s", vc.currentValidator, message),
		Reference: ref,
	})
}

func (vc *validationContext) IsKnownComponent(name arch.ComponentName) bool {
	return vc.conf.Components.Map.Has(name)
}

func (vc *validationContext) IsKnownVendor(name arch.VendorName) bool {
	return vc.conf.Vendors.Map.Has(name)
}
