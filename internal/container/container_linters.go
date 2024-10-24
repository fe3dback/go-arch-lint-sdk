package container

import (
	"github.com/fe3dback/go-arch-lint-sdk/internal/services/linters"
	"github.com/fe3dback/go-arch-lint-sdk/internal/services/linters/deepscan"
)

func (c *Container) lintersRoot() *linters.Root {
	return once(func() *linters.Root {
		return linters.NewRoot(
			c.serviceNoticeFormatter(),
			c.lintersSyntax(),
			c.lintersOrphans(),
			c.lintersImports(),
			c.lintersDeepscan(),
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

func (c *Container) lintersDeepscan() *linters.Deepscan {
	return once(func() *linters.Deepscan {
		return linters.NewDeepscan(
			deepscan.NewProjectAstFetcher(),
		)
	})
}
