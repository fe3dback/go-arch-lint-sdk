package container

import (
	"github.com/fe3dback/go-arch-lint-sdk/internal/operations/check"
	"github.com/fe3dback/go-arch-lint-sdk/internal/operations/mapping"
)

func (c *Container) OperationMapping() *mapping.Operation {
	return once(func() *mapping.Operation {
		return mapping.NewOperation()
	})
}

func (c *Container) OperationCheck() *check.Operation {
	return once(func() *check.Operation {
		return check.NewOperation(
			c.lintersRoot(),
		)
	})
}
