package linters

import (
	"fmt"
	"go/token"
	"sync"

	"golang.org/x/sync/errgroup"
	"golang.org/x/tools/go/packages"

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
	roCtx, err := l.prepareLintContext(&spec)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare lint context: %w", err)
	}

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
				lintCtx := lintContext{
					ro: &roCtx,
					state: &lintContextMutable{
						mux:     &sync.Mutex{},
						notices: make([]arch.LinterNotice, 0, 16),
					},
				}

				err := linter.Lint(&lintCtx)
				if err != nil {
					return fmt.Errorf("linter %T failed: %w", linter, err)
				}

				result.Notices = lintCtx.state.Notices()
			}

			mux.Lock()
			results = append(results, result)
			mux.Unlock()

			return nil
		})
	}

	err = eg.Wait()
	if err != nil {
		return nil, fmt.Errorf("run linters failed: %w", err)
	}

	return results, nil
}

func (l *Root) prepareLintContext(spec *arch.Spec) (lintContextReadOnly, error) {
	fileSet := token.NewFileSet()

	parsedPackages, err := l.parsePackages(fileSet, spec)
	if err != nil {
		return lintContextReadOnly{}, fmt.Errorf("failed to parse packages: %w", err)
	}

	parsedStdPackagesIDs, err := l.parseStdPackageIDs(fileSet)
	if err != nil {
		return lintContextReadOnly{}, fmt.Errorf("failed to parse go std packages: %w", err)
	}

	lCtx := lintContextReadOnly{
		spec:          spec,
		fileSet:       fileSet,
		stdPackageIDs: parsedStdPackagesIDs,
		packages:      parsedPackages,
	}

	return lCtx, nil
}

func (l *Root) parseStdPackageIDs(fileSet *token.FileSet) (map[arch.PathImport]any, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName,
		Fset: fileSet,
	}

	stdPackages, err := packages.Load(cfg, "std")
	if err != nil {
		return nil, fmt.Errorf("failed to load standard packages: %w", err)
	}

	stdPackageIDs := make(map[arch.PathImport]any, 128)

	for _, stdPackage := range stdPackages {
		stdPackageIDs[arch.PathImport(stdPackage.ID)] = struct{}{}
	}

	return stdPackageIDs, nil
}

func (l *Root) parsePackages(fileSet *token.FileSet, spec *arch.Spec) (packagesMap, error) {
	const parseMode = packages.NeedName |
		packages.NeedFiles |
		packages.NeedTypes |
		packages.NeedSyntax |
		packages.NeedTypesInfo

	directoriesList := map[arch.PathAbsolute]*arch.PackageDescriptor{}
	for _, cmp := range spec.Components {
		for _, ownedDirectory := range cmp.OwnedPackages {
			directoriesList[ownedDirectory.PathAbs] = &ownedDirectory
		}
	}

	var wg errgroup.Group
	var mux sync.Mutex
	wg.SetLimit(16)

	packagesMap := make(packagesMap, len(directoriesList)*2)

	for directory, dst := range directoriesList {
		directory := directory
		dst := dst

		wg.Go(func() error {
			cfg := &packages.Config{
				Mode: parseMode,
				Fset: fileSet,
				Dir:  string(directory),
			}

			parsedPackages, err := packages.Load(cfg, string(directory))
			if err != nil {
				return fmt.Errorf("failed parse go source at '%s': %w", dst.PathRel, err)
			}

			mux.Lock()
			for _, goPackage := range parsedPackages {
				packagesMap[dst.PathRel] = append(packagesMap[dst.PathRel], goPackage)
			}
			mux.Unlock()

			return nil
		})
	}

	err := wg.Wait()
	if err != nil {
		return nil, fmt.Errorf("parse in goroutines failed: %w", err)
	}

	return packagesMap, nil
}
