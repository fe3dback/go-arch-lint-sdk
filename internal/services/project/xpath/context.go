package xpath

import (
	"github.com/fe3dback/go-arch-lint-sdk/arch"
)

type queryContext struct {
	projectDirectory arch.PathAbsolute
	index            *index
}

func newQueryContext(projectDirectory arch.PathAbsolute) queryContext {
	return queryContext{
		projectDirectory: projectDirectory,
		index:            newIndex(),
	}
}
