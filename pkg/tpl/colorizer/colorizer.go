package colorizer

import (
	"fmt"
	"os"

	"github.com/muesli/termenv"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
)

var palette = map[string]string{
	"red":     "#d62d20",
	"green":   "#008744",
	"blue":    "#0057e7",
	"yellow":  "#ffa700",
	"magenta": "#d62976",
	"cyan":    "#2691a5",
	"gray":    "#777777",
}

type Colorizer struct {
	colorEnv arch.TerminalColorEnv
	env      *termenv.Output
}

func New(colorEnv arch.TerminalColorEnv) *Colorizer {
	// detect color profile from TTY
	prof := termenv.ColorProfile()

	// reset if not allowed
	if colorEnv == arch.TerminalColorEnvBlackAndWhite {
		prof = termenv.Ascii
	}

	// force override from ascii to ansi
	// but if color profile better than ansi - will stay as-is
	if colorEnv == arch.TerminalColorEnvColored && prof == termenv.Ascii {
		// force turn on colors
		prof = termenv.ANSI
	}

	return &Colorizer{
		colorEnv: colorEnv,
		env:      termenv.NewOutput(os.Stdout, termenv.WithProfile(prof)),
	}
}

func (c *Colorizer) Colorize(color string, text string) string {
	if c.colorEnv == arch.TerminalColorEnvBlackAndWhite {
		return text
	}

	hex, ok := palette[color]
	if !ok {
		return text
	}

	return c.env.
		String(text).
		Foreground(c.env.Color(hex)).
		String()
}

func (c *Colorizer) wrapWithQuotes(text string) string {
	return fmt.Sprintf("\"%s\"", text)
}
