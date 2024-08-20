package xpath

import (
	"fmt"
	pathUtils "path"
	"strings"

	"github.com/gobwas/glob"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
)

type MatcherRelativeGlob struct{}

func NewMatcherRelativeGlob() *MatcherRelativeGlob {
	return &MatcherRelativeGlob{}
}

func (m *MatcherRelativeGlob) match(ctx *queryContext, query arch.FileQuery) ([]arch.FileDescriptor, error) {
	path := query.Path.(arch.PathRelativeGlob) // guaranteed by root composite
	path = arch.PathRelativeGlob(pathUtils.Join(string(query.WorkingDirectory), string(path)))

	var patternLast glob.Glob
	patternNormal, err := glob.Compile(string(path), '/')
	if err != nil {
		return nil, fmt.Errorf("failed compile glob matcher '%s': %w", path, err)
	}

	if strings.HasSuffix(string(path), "/**") {
		pathLast := strings.TrimSuffix(string(path), "/**")
		patternLast, err = glob.Compile(pathLast, '/')
		if err != nil {
			return nil, fmt.Errorf("failed compile glob matcher '%s': %w", pathLast, err)
		}
	}

	results := make([]arch.FileDescriptor, 0, 16)

	ctx.index.each(func(dsc arch.FileDescriptor) {
		if query.Type == arch.FileMatchQueryTypeOnlyDirectories && !dsc.IsDir {
			return
		}

		if query.Type == arch.FileMatchQueryTypeOnlyFiles && dsc.IsDir {
			return
		}

		matchedNormal := patternNormal.Match(string(dsc.PathRel))
		matchedLast := false

		if patternLast != nil {
			matchedLast = patternLast.Match(string(dsc.PathRel))
		}

		if !(matchedNormal || matchedLast) {
			return
		}

		results = append(results, dsc)
	})

	return results, nil
}
