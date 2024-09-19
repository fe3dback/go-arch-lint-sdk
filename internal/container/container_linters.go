package container

import (
	"github.com/fe3dback/go-arch-lint-sdk/internal/services/linters"
)

func (c *Container) lintersRoot() *linters.Root {
	return once(func() *linters.Root {
		return linters.NewRoot(
			c.lintersSyntax(),
			c.lintersOrphans(),
			c.lintersImports(),
		)
	})
}

func (c *Container) lintersSyntax() *linters.Syntax {
	return once(linters.NewSyntax)
}

func (c *Container) lintersOrphans() *linters.Orphans {
	return once(linters.NewOrphans)
}

func (c *Container) lintersImports() *linters.Imports {
	return once(linters.NewImports)
}
