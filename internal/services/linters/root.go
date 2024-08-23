package linters

import (
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
)

// todo: imp
//var linters = []check.Linter{
//	{
//		ID:   check.LinterIDComponentImports,
//		Name: "Base: component imports",
//		Hint: "always on",
//	},
//	{
//		ID:   check.LinterIDVendorImports,
//		Name: "Advanced: vendor imports",
//		Hint: "switch 'allow.depOnAnyVendor = false' (or delete) to on",
//	},
//	{
//		ID:   check.LinterIDDeepScan,
//		Name: "Advanced: method calls and dependency injections",
//		Hint: "switch 'allow.depOnAnyVendor = false' (or delete) to on",
//	},
//}

type Root struct {
	linters []linter
}

func NewRoot(linters ...linter) *Root {
	return &Root{
		linters: linters,
	}
}

func (l *Root) Lint(spec arch.Spec) ([]arch.LinterResult, error) {
	var mux sync.Mutex

	results := make([]arch.LinterResult, 0, len(l.linters))
	lCtx := lintContext{}

	eg := new(errgroup.Group)

	for _, linter := range l.linters {
		linter := linter

		eg.Go(func() error {
			info := linter.Information()
			info.Used = linter.IsSuitable(&spec) // spec is RO

			result := arch.LinterResult{
				Linter: info,
			}

			if info.Used {
				// spec is RO
				// lCtx must be guarded inside linter

				notices, err := linter.Lint(&lCtx, &spec)
				if err != nil {
					return fmt.Errorf("linter %T failed: %w", linter, err)
				}

				result.Notices = notices
			}

			mux.Lock()
			results = append(results, result)
			mux.Unlock()

			return nil
		})
	}

	err := eg.Wait()
	if err != nil {
		return nil, fmt.Errorf("run linters failed: %w", err)
	}

	return results, nil
}
