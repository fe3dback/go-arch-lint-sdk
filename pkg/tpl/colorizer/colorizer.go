package colorizer

import (
	"fmt"
	"os"

	"github.com/muesli/termenv"
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
	useColors bool
	env       *termenv.Output
}

func New(useColors bool, forceASCII bool) *Colorizer {
	var env *termenv.Output

	if !forceASCII {
		// detect profile automatically
		env = termenv.NewOutput(os.Stdout)
	} else {
		env = termenv.NewOutput(os.Stdout, termenv.WithProfile(termenv.ANSI))
	}

	return &Colorizer{
		useColors: useColors,
		env:       env,
	}
}

func (c *Colorizer) Colorize(color string, text string) string {
	if !c.useColors {
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
