package container

import "github.com/fe3dback/go-arch-lint-sdk/internal/operations/mapping"

func (c *Container) OperationMapping() *mapping.Operation {
	return once(func() *mapping.Operation {
		return mapping.NewOperation()
	})
}
