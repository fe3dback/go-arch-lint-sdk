package container

import (
	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/pkg/tpl"
	"github.com/fe3dback/go-arch-lint-sdk/pkg/tpl/codeprinter"
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
		return colorizer.New(c.terminalColorEnv)
	})
}

func (c *Container) serviceCodePrinter() *codeprinter.Printer {
	return once(func() *codeprinter.Printer {
		return codeprinter.NewPrinter(
			codeprinter.NewExtractorRaw(),
			codeprinter.NewExtractorHL(),
			c.terminalColorEnv == arch.TerminalColorEnvColored,
		)
	})
}
