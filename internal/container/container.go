package container

import "github.com/fe3dback/go-arch-lint-sdk/arch"

type Container struct {
	projectDirectory arch.PathAbsolute
	usedContext      arch.UsedContext
	skipMissUsages   bool
	outputColors     bool
}

func NewContainer(
	projectDirectory arch.PathAbsolute,
	usedContext arch.UsedContext,
	skipMissUsages bool,
	outputColors bool,
) *Container {
	return &Container{
		projectDirectory: projectDirectory,
		usedContext:      usedContext,
		skipMissUsages:   skipMissUsages,
		outputColors:     outputColors,
	}
}
