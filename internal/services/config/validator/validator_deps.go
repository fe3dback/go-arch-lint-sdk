package validator

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/internal/models"
)

type DepsValidator struct{}

func NewDepsValidator() *DepsValidator {
	return &DepsValidator{}
}

func (c *DepsValidator) Validate(ctx *validationContext) {
	ctx.conf.Dependencies.Map.Each(func(name arch.ComponentName, rules models.ConfigComponentDependencies, reference arch.Reference) {
		if !ctx.IsKnownComponent(name) {
			ctx.AddNotice(
				fmt.Sprintf("Component '%s' in dependencies is not defined", name),
				reference,
			)
			return
		}

		if rules.AnyProjectDeps.Value && len(rules.MayDependOn) > 0 {
			ctx.AddMissUse(
				fmt.Sprintf("redundant: in component '%s': rule '%s' used with not empty '%s' list",
					name,
					xpathOr(rules.AnyProjectDeps.Ref.XPath, "anyProjectDeps"),
					xpathOr(rules.MayDependOn[0].Ref.XPath, "mayDependOn"),
				),
				rules.AnyProjectDeps.Ref,
			)
			return
		}

		if rules.AnyVendorDeps.Value && len(rules.CanUse) > 0 {
			ctx.AddMissUse(
				fmt.Sprintf("redundant: in component '%s': rule '%s' used with not empty '%s' list",
					name,
					xpathOr(rules.AnyVendorDeps.Ref.XPath, "anyVendorDeps"),
					xpathOr(rules.CanUse[0].Ref.XPath, "canUse"),
				),
				rules.AnyVendorDeps.Ref,
			)
			return
		}

		if len(rules.MayDependOn) == 0 && len(rules.CanUse) == 0 {
			if rules.AnyProjectDeps.Value {
				return
			}

			if rules.AnyVendorDeps.Value {
				return
			}

			if len(rules.CanContainTags.Values()) > 0 {
				return
			}

			ctx.AddNotice(
				fmt.Sprintf("In component '%s': no rules is defined (require at least one of [AnyProjectDeps, AnyVendorDeps, MayDependOn, CanUse, CanContainTags])", name),
				reference,
			)
			return
		}
	})
}
