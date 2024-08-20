package validator

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/internal/models"
)

type TagsMissUseValidator struct {
}

func NewTagsMissUseValidator() *TagsMissUseValidator {
	return &TagsMissUseValidator{}
}

func (c *TagsMissUseValidator) Validate(ctx *validationContext) {
	c.validateTagsRedundantIfAllAllowed(ctx)
	c.validateTagsCollisionWithGlobal(ctx)
}

func (c *TagsMissUseValidator) validateTagsRedundantIfAllAllowed(ctx *validationContext) {
	setting := ctx.conf.Settings.Tags.Allowed
	if setting.Value != arch.ConfigSettingsTagsEnumAll {
		return
	}

	ctx.conf.Dependencies.Map.Each(func(name arch.ComponentName, rules models.ConfigComponentDependencies, _ arch.Reference) {
		for _, tag := range rules.CanContainTags {
			ctx.AddMissUse(
				fmt.Sprintf("redundant: component '%s' canContainTag '%s', but all tags in project is allowed in setting '%s = true' or by default",
					name,
					tag.Value,
					xpathOr(setting.Ref.XPath, "$.settings.structTags.allowed"),
				),
				tag.Ref,
			)
		}
	})
}

func (c *TagsMissUseValidator) validateTagsCollisionWithGlobal(ctx *validationContext) {
	setting := ctx.conf.Settings.Tags.Allowed
	if setting.Value != arch.ConfigSettingsTagsEnumList {
		return
	}

	allowedList := ctx.conf.Settings.Tags.AllowedList
	ctx.conf.Dependencies.Map.Each(func(name arch.ComponentName, rules models.ConfigComponentDependencies, _ arch.Reference) {
		for _, tag := range rules.CanContainTags {
			if !allowedList.Contains(tag) {
				continue
			}

			ctx.AddMissUse(
				fmt.Sprintf("redundant: component '%s' canContainTag '%s', but this tag allowed for all components in '%s'",
					name,
					tag.Value,
					xpathOr(setting.Ref.XPath, "$.settings.structTags.allowed"),
				),
				tag.Ref,
			)
		}
	})
}
