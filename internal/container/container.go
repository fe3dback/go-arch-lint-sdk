package container

import "github.com/fe3dback/go-arch-lint-sdk/arch"

type Container struct {
	projectDirectory arch.PathAbsolute
	skipMissUsages   bool
}

func NewContainer(projectDirectory arch.PathAbsolute, skipMissUsages bool) *Container {
	return &Container{
		projectDirectory: projectDirectory,
		skipMissUsages:   skipMissUsages,
	}
}
