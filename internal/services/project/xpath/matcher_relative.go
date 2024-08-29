package xpath

import (
	pathUtils "path"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
)

type MatcherRelative struct {
}

func NewMatcherRelative() *MatcherRelative {
	return &MatcherRelative{}
}

func (m *MatcherRelative) match(ctx *queryContext, query arch.FileQuery) ([]arch.PathDescriptor, error) {
	path := query.Path.(arch.PathRelative) // guaranteed by root composite
	path = arch.PathRelative(pathUtils.Join(string(query.WorkingDirectory), string(path)))

	descriptors := make([]arch.PathDescriptor, 0, 2)

	if query.Type == arch.FileMatchQueryTypeAll || query.Type == arch.FileMatchQueryTypeOnlyDirectories {
		if dir, found := ctx.index.directoryAt(path); found {
			descriptors = append(descriptors, dir)
		}
	}

	if query.Type == arch.FileMatchQueryTypeAll || query.Type == arch.FileMatchQueryTypeOnlyFiles {
		if file, found := ctx.index.fileAt(path); found {
			descriptors = append(descriptors, file)
		}
	}

	return descriptors, nil
}
