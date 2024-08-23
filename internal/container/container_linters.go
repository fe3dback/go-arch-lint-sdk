package container

import (
	"github.com/fe3dback/go-arch-lint-sdk/internal/services/linters"
)

func (c *Container) lintersRoot() *linters.Root {
	return once(func() *linters.Root {
		return linters.NewRoot(
			c.lintersOrphans(),
		)
	})
}

func (c *Container) lintersOrphans() *linters.Orphans {
	return once(func() *linters.Orphans {
		return linters.NewOrphans()
	})
}
