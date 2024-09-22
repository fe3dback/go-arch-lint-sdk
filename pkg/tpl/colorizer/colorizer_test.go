package colorizer_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fe3dback/go-arch-lint-sdk/pkg/tpl/colorizer"
)

const (
	sourceText     = "hello"
	sourceTextRed  = "\x1b[91mhello\x1b[0m"
	sourceTextBlue = "\x1b[94mhello\x1b[0m"
)

type in struct {
	useColors bool
	color     string
	text      string
}

func TestASCII_Colorize(t *testing.T) {
	tests := []struct {
		name string
		in   in
		out  string
	}{
		{
			name: "happy_no_colors",
			in: createIn("red", func(in *in) {
				in.useColors = false
			}),
			out: sourceText,
		},
		{
			name: "happy_red",
			in:   createIn("red"),
			out:  sourceTextRed,
		},
		{
			name: "happy_blue",
			in:   createIn("blue"),
			out:  sourceTextBlue,
		},
		{
			name: "unknown_color",
			in:   createIn("not-exist-color"),
			out:  sourceText,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := colorizer.New(tt.in.useColors, true)

			got := r.Colorize(tt.in.color, tt.in.text)
			require.Equal(t, tt.out, got)
		})
	}
}

func createIn(col string, mutators ...func(*in)) in {
	in := in{
		color:     col,
		text:      sourceText,
		useColors: true,
	}

	for _, mutate := range mutators {
		mutate(&in)
	}

	return in
}
