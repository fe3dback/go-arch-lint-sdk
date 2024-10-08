package container

import (
	"github.com/fe3dback/go-arch-lint-sdk/internal/services/project/module"
	"github.com/fe3dback/go-arch-lint-sdk/internal/services/project/xpath"
)

func (c *Container) serviceProjectPathHelper() *xpath.Helper {
	return once(func() *xpath.Helper {
		return xpath.NewHelper(
			c.projectDirectory,
			c.serviceProjectPathFileScanner(),
			c.serviceProjectPathMatcherRelative(),
			c.serviceProjectPathMatcherAbsolute(),
			c.serviceProjectPathMatcherRelativeGlob(),
			nil, // todo
		)
	})
}

func (c *Container) serviceProjectPathFileScanner() *xpath.FileScanner {
	return once(xpath.NewFileScanner)
}

func (c *Container) serviceProjectPathMatcherRelative() *xpath.MatcherRelative {
	return once(xpath.NewMatcherRelative)
}

func (c *Container) serviceProjectPathMatcherRelativeGlob() *xpath.MatcherRelativeGlob {
	return once(xpath.NewMatcherRelativeGlob)
}

func (c *Container) serviceProjectPathMatcherAbsolute() *xpath.MatcherAbsolute {
	return once(func() *xpath.MatcherAbsolute {
		return xpath.NewMatcherAbsolute(
			c.serviceProjectPathMatcherRelative(),
		)
	})
}

func (c *Container) serviceProjectFetcher() *module.Fetcher {
	return once(func() *module.Fetcher {
		return module.NewFetcher(
			c.projectDirectory,
		)
	})
}
