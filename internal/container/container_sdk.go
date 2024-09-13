package container

import "github.com/fe3dback/go-arch-lint-sdk/definition"

func (c *Container) ConfigDefinition() *definition.Definition {
	return once(func() *definition.Definition {
		return definition.NewDefinition(
			c.projectDirectory,
			c.serviceConfigReader(),
			c.serviceConfigValidator(),
			c.serviceConfigAssembler(),
		)
	})
}
