package xpath

import (
	"fmt"
	pathUtils "path"
	"path/filepath"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
)

type MatcherAbsolute struct {
	matcherRelative typeMatcher
}

func NewMatcherAbsolute(
	matcherRelative typeMatcher,
) *MatcherAbsolute {
	return &MatcherAbsolute{
		matcherRelative: matcherRelative,
	}
}

func (m *MatcherAbsolute) match(ctx *queryContext, query arch.FileQuery) ([]arch.FileDescriptor, error) {
	path := query.Path.(arch.PathAbsolute) // guaranteed by root composite
	path = arch.PathAbsolute(pathUtils.Join(string(query.WorkingDirectory), string(path)))

	// take relative path
	relPathStr, err := filepath.Rel(string(ctx.projectDirectory), string(path))
	if err != nil {
		return nil, fmt.Errorf("failed get relative path from '%s': %w", path, err)
	}

	// pass to next matcher
	relPath := arch.PathRelative(relPathStr)
	relQuery := query
	relQuery.Path = relPath

	return m.matcherRelative.match(ctx, relQuery)
}
