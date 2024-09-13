package cfg

import (
	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/definition/internal"
	"github.com/fe3dback/go-arch-lint-sdk/internal/models"
)

type (
	Factory func(*models.Config)
)

func WorkingDirectory(relative arch.PathRelative) Factory {
	ref := internal.GetParentCaller()
	return func(config *models.Config) {
		config.WorkingDirectory = arch.NewRef(relative, ref)
	}
}

func Component(name arch.ComponentName, in ...arch.PathRelativeGlob) Factory {
	ref := internal.GetParentCaller()
	return func(config *models.Config) {
		inList := make(arch.RefSlice[arch.PathRelativeGlob], len(in))
		for ind, value := range in {
			inList[ind] = arch.NewRef(value, ref)
		}

		config.Components.Map.Set(
			name,
			models.ConfigComponent{
				In: inList,
			},
			ref,
		)
	}
}

func CommonComponents(list ...arch.ComponentName) Factory {
	ref := internal.GetParentCaller()
	return func(config *models.Config) {
		refList := make(arch.RefSlice[arch.ComponentName], len(list))
		for ind, name := range list {
			refList[ind] = arch.NewRef(name, ref)
		}

		config.CommonComponents = append(config.CommonComponents, refList...)
	}
}
