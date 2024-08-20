package container

import "github.com/fe3dback/go-arch-lint-sdk/cfg"

func (c *Container) ConfigDefinition() *cfg.Definition {
	return once(func() *cfg.Definition {
		return cfg.NewDefinition(
			c.projectDirectory,
			c.serviceConfigReader(),
			c.serviceConfigValidator(),
			c.serviceConfigAssembler(),
		)
	})
}
