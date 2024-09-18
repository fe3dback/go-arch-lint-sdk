package xpath

import (
	"strings"

	"github.com/gobwas/glob"
)

var (
	globMatchersCache = make(map[string]glob.Glob)
)

func init() {
	globMatchersCache = make(map[string]glob.Glob, 32)
}

func IsGlobMatched(pattern string, value string) bool {
	if getCompiledGlobMatcher(pattern).Match(value) {
		return true
	}

	if strings.HasSuffix(pattern, "/**") {
		patternLast := strings.TrimSuffix(pattern, "/**")
		if getCompiledGlobMatcher(patternLast).Match(value) {
			return true
		}
	}

	return false
}

func getCompiledGlobMatcher(path string) glob.Glob {
	if cached, exist := globMatchersCache[path]; exist {
		return cached
	}

	// err guaranteed can`t be here, because every importGlobPath already checked on validation stage
	matcher, _ := glob.Compile(path, '/')
	globMatchersCache[path] = matcher

	return matcher
}
