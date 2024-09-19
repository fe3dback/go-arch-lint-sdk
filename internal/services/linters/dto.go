package linters

import (
	"go/token"
	"sync"

	"golang.org/x/tools/go/packages"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/internal/models"
)

type (
	packagesMap = map[arch.PathRelative][]*packages.Package

	lintContext struct {
		ro    *lintContextReadOnly
		state *lintContextMutable
	}

	// lintContext is READ ONLY and shouldn't been modified inside linters
	// all linters will work in goroutines and can read all data in parallel
	lintContextReadOnly struct {
		// global options
		options models.LintOptions

		// go ast fileSet. All parse function will use this fileSet
		// so, it can be used for reference resolving after parsing
		fileSet *token.FileSet

		// contain map of packages.Package[ID]
		// of all known go std packages for current version of GO binary
		// available in user host
		stdPackageIDs map[arch.PathImport]any

		// every directory can contain multiple go packages (ex: linters, linters_test)
		packages packagesMap

		// arch spec reference
		spec *arch.Spec
	}

	// lintContextMutable is safe to use, because it unique for each linter
	lintContextMutable struct {
		mux     *sync.Mutex
		notices []arch.LinterNotice
	}
)

func (lcs *lintContextMutable) AddNotice(notice arch.LinterNotice) {
	lcs.mux.Lock()
	defer lcs.mux.Unlock()

	lcs.notices = append(lcs.notices, notice)
}

func (lcs *lintContextMutable) Notices() []arch.LinterNotice {
	lcs.mux.Lock()
	defer lcs.mux.Unlock()

	newSlice := make([]arch.LinterNotice, len(lcs.notices))
	copy(newSlice, lcs.notices)

	return newSlice
}
