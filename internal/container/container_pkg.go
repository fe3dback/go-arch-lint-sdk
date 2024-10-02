package container

import (
	"github.com/fe3dback/go-arch-lint-sdk/pkg/tpl"
	"github.com/fe3dback/go-arch-lint-sdk/pkg/tpl/colorizer"
)

func (c *Container) serviceRenderer() *tpl.Renderer {
	return once(func() *tpl.Renderer {
		return tpl.NewRenderer(
			c.serviceColorizer(),
		)
	})
}

func (c *Container) serviceColorizer() *colorizer.Colorizer {
	return once(func() *colorizer.Colorizer {
		return colorizer.New(c.outputColors, false)
	})
}
