package definition

import (
	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/definition/cfg"
	"github.com/fe3dback/go-arch-lint-sdk/definition/internal"
	"github.com/fe3dback/go-arch-lint-sdk/internal/models"
)

// FromCode will transform GO code struct into config
// and then build spec from this config
// for example:
//   - FromAbsoluteFile -> will read .yml             as `config` and build arch.Spec from it
//   - FromCode         -> will transform preConfig into `config` and build arch.Spec from it
func (def *Definition) FromCode(builders ...cfg.Factory) (arch.Spec, error) {
	conf := models.Config{
		Version: arch.NewRef(arch.ConfigVersion4, internal.GetParentCaller()),
		Components: models.ConfigComponents{
			Map: arch.NewRefMap[arch.ComponentName, models.ConfigComponent](16),
		},
	}

	for _, build := range builders {
		build(&conf)
	}

	return def.withUserFriendlyError(
		def.fromConfig(conf),
	)
}
