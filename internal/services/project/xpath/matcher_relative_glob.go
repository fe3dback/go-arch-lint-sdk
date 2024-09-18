package xpath

import (
	pathUtils "path"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
)

type MatcherRelativeGlob struct{}

func NewMatcherRelativeGlob() *MatcherRelativeGlob {
	return &MatcherRelativeGlob{}
}

func (m *MatcherRelativeGlob) match(ctx *queryContext, query arch.FileQuery) ([]arch.PathDescriptor, error) {
	pattern := query.Path.(arch.PathRelativeGlob) // guaranteed by root composite
	pattern = arch.PathRelativeGlob(pathUtils.Join(string(query.WorkingDirectory), string(pattern)))

	results := make([]arch.PathDescriptor, 0, 16)

	ctx.index.each(func(dsc arch.PathDescriptor) {
		if query.Type == arch.FileMatchQueryTypeOnlyDirectories && !dsc.IsDir {
			return
		}

		if query.Type == arch.FileMatchQueryTypeOnlyFiles && dsc.IsDir {
			return
		}

		if !IsGlobMatched(string(pattern), string(dsc.PathRel)) {
			return
		}

		results = append(results, dsc)
	})

	return results, nil
}
