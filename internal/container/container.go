package container

import (
	"github.com/fe3dback/go-arch-lint-sdk/arch"
)

type Container struct {
	projectDirectory arch.PathAbsolute
	usedContext      arch.UsedContext
	terminalColorEnv arch.TerminalColorEnv
	skipMissUsages   bool
}

func NewContainer(
	projectDirectory arch.PathAbsolute,
	usedContext arch.UsedContext,
	skipMissUsages bool,
	useColors bool,
) *Container {
	return &Container{
		projectDirectory: projectDirectory,
		usedContext:      usedContext,
		terminalColorEnv: arch.DetectColorProfile(useColors),
		skipMissUsages:   skipMissUsages,
	}
}
