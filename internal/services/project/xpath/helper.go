package xpath

import (
	"fmt"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
)

type Helper struct {
	indexActual   bool
	fileScanner   fileScanner
	matchers      map[string]typeMatcher
	queryCtx      queryContext
	cachedRegExps map[string]*regexp.Regexp
}

func NewHelper(
	projectDirectory arch.PathAbsolute,
	fileScanner fileScanner,
	matcherRelative typeMatcher,
	matcherAbsolute typeMatcher,
	matcherGlobRelative typeMatcher,
	matcherGlobAbsolute typeMatcher,
) *Helper {
	rootDirectory, err := filepath.Abs(string(projectDirectory))
	if err != nil {
		panic(fmt.Errorf("failed get working directory: %w", err))
	}

	srv := &Helper{
		indexActual: false,
		fileScanner: fileScanner,
		queryCtx:    newQueryContext(arch.PathAbsolute(rootDirectory)),
		matchers: map[string]typeMatcher{
			getType(arch.PathRelative("/")):     matcherRelative,
			getType(arch.PathAbsolute("/")):     matcherAbsolute,
			getType(arch.PathRelativeGlob("/")): matcherGlobRelative,
			getType(arch.PathAbsoluteGlob("/")): matcherGlobAbsolute,
		},
		cachedRegExps: make(map[string]*regexp.Regexp, 4),
	}

	return srv
}

func (h *Helper) reindexProjectFilesIfNecessary() error {
	if h.indexActual {
		return nil
	}
	h.indexActual = true

	return h.fileScanner.Scan(string(h.queryCtx.projectDirectory), func(path string, isDir bool) error {
		relativePathStr, err := filepath.Rel(string(h.queryCtx.projectDirectory), path)
		if err != nil {
			return fmt.Errorf("failed getting relative path '%q' from '%q': %w", path, h.queryCtx.projectDirectory, err)
		}

		relativePath := arch.PathRelative(relativePathStr)
		extLower := strings.ToLower(filepath.Ext(path))
		extLower = strings.TrimLeft(extLower, ".")

		h.queryCtx.index.appendToIndex(relativePath, arch.PathDescriptor{
			PathRel:   relativePath,
			PathAbs:   arch.PathAbsolute(path),
			IsDir:     isDir,
			Extension: extLower,
		})

		return nil
	})
}

func (h *Helper) FindProjectFiles(query arch.FileQuery) ([]arch.PathDescriptor, error) {
	err := h.reindexProjectFilesIfNecessary()
	if err != nil {
		return nil, fmt.Errorf("failed build files index from project directory '%s': %w", h.queryCtx.projectDirectory, err)
	}

	pathType := getType(query.Path)
	matcher, exist := h.matchers[pathType]
	if !exist {
		return nil, fmt.Errorf("unknown matcher type %s", pathType)
	}

	if matcher == nil {
		return nil, fmt.Errorf("NIL matcher registered for type %s", pathType)
	}

	// match
	found, err := matcher.match(&h.queryCtx, query)
	if err != nil {
		return nil, fmt.Errorf("failed match files by path '%s/%s': %w", query.WorkingDirectory, query.Path, err)
	}

	// filter
	result := make([]arch.PathDescriptor, 0, len(found))
	for _, dst := range found {
		suitable, err := h.isSuitable(dst, &query)
		if err != nil {
			return nil, fmt.Errorf("failed check file name '%s': %w", dst.PathRel, err)
		}

		if !suitable {
			continue
		}

		result = append(result, dst)
	}

	// sort
	sort.Slice(result, func(i, j int) bool {
		return result[i].PathRel < result[j].PathRel
	})

	return result, nil
}

//nolint:funlen
func (h *Helper) isSuitable(dst arch.PathDescriptor, query *arch.FileQuery) (bool, error) {
	// only directories
	if dst.IsDir && !(query.Type == arch.FileMatchQueryTypeAll || query.Type == arch.FileMatchQueryTypeOnlyDirectories) {
		return false, nil
	}

	// only files
	if !dst.IsDir && !(query.Type == arch.FileMatchQueryTypeAll || query.Type == arch.FileMatchQueryTypeOnlyFiles) {
		return false, nil
	}

	// find dir
	dstDirectory := dst.PathRel
	if !dst.IsDir {
		dstDirectory = arch.PathRelative(filepath.Dir(string(dst.PathRel)))
	}

	// exclude by directory
	if len(query.ExcludeDirectories) > 0 {
		excludedDirs := make([]arch.PathRelative, 0, len(query.ExcludeDirectories))
		for _, excDirectory := range query.ExcludeDirectories {
			excludedDirs = append(excludedDirs, arch.PathRelative(filepath.Join(string(query.WorkingDirectory), string(excDirectory))))
		}

		if slices.Contains(excludedDirs, dstDirectory) {
			return false, nil
		}
	}

	// exclude by file path
	if len(query.ExcludeFiles) > 0 {
		excludedFiles := make([]arch.PathRelative, 0, len(query.ExcludeFiles))
		for _, excFile := range query.ExcludeFiles {
			excludedFiles = append(excludedFiles, arch.PathRelative(filepath.Join(string(query.WorkingDirectory), string(excFile))))
		}

		if !dst.IsDir && slices.Contains(excludedFiles, dst.PathRel) {
			return false, nil
		}
	}

	// exclude by regexp
	for _, pattern := range query.ExcludeRegexp {
		reg, err := h.compileRegexp(string(pattern))
		if err != nil {
			return false, fmt.Errorf("failed compile regular expression '%s': %w", pattern, err)
		}

		if reg.MatchString("/" + string(dst.PathRel)) {
			return false, nil
		}
	}

	// exclude by file ext
	if len(query.Extensions) > 0 {
		if !dst.IsDir && !slices.Contains(query.Extensions, dst.Extension) {
			return false, nil
		}
	}

	// ok
	return true, nil
}

func (h *Helper) compileRegexp(expr string) (*regexp.Regexp, error) {
	if _, ok := h.cachedRegExps[expr]; !ok {
		regular, err := regexp.Compile(expr)
		if err != nil {
			return nil, fmt.Errorf("invalid regular expression '%s': %w", expr, err)
		}

		h.cachedRegExps[expr] = regular
	}

	return h.cachedRegExps[expr], nil
}
