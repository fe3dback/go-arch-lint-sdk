package arch

import (
	"os"

	"github.com/muesli/termenv"
)

type (
	UsedContext int

	TerminalColorEnv string
)

const (
	// UsedContextDefault should be used by default, except rare cased described in UsedContextCLI
	UsedContextDefault UsedContext = 0

	// UsedContextCLI Marks that SDK used from CLI program
	// this affects small details, e.g. validation errors will not be pasted into error,
	// but will be formatted separately in stdout.
	//
	// DO NOT USE anywhere, except go-arch-lint CLI program, or analog
	UsedContextCLI UsedContext = 1
)

const (

	// TerminalColorEnvBlackAndWhite only basic text (black/white) without any formatting
	TerminalColorEnvBlackAndWhite TerminalColorEnv = "ASCII"

	// TerminalColorEnvColored mode affect all output and printing.
	// this mode will try to turn on colors (if terminal/emulator/std env is support it)
	TerminalColorEnvColored TerminalColorEnv = "Colored"
)

func DetectColorProfile(useColors bool) TerminalColorEnv {
	dsds := os.Environ()
	_ = dsds

	// forced colors off
	if !useColors {
		return TerminalColorEnvBlackAndWhite
	}

	// detect env mode
	prof := termenv.ColorProfile()
	if prof == termenv.Ascii {
		return TerminalColorEnvBlackAndWhite
	}

	// allow basic colors
	return TerminalColorEnvColored
}
