package linters

import (
	"fmt"
	"go/token"
	"sort"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
	"golang.org/x/tools/go/packages"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/internal/models"
)

type Root struct {
	linters []linter
}

func NewRoot(linters ...linter) *Root {
	return &Root{
		linters: linters,
	}
}

func (l *Root) Lint(spec arch.Spec, options models.LintOptions) ([]arch.LinterResult, error) {
	var mux sync.Mutex

	results := make([]arch.LinterResult, 0, len(l.linters))
	roCtx, err := l.prepareLintContext(&spec, options)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare lint context: %w", err)
	}

	eg := new(errgroup.Group)

	for _, linter := range l.linters {
		linter := linter

		eg.Go(func() error {
			info := linter.Information()
			info.Used = linter.IsSuitable(&roCtx) // spec is RO

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

	sort.SliceStable(results, func(i, j int) bool {
		lntI, lntJ := results[i], results[j]
		orderI, orderJ := arch.LintersSortOrder[lntI.Linter.ID], arch.LintersSortOrder[lntJ.Linter.ID]

		return orderI < orderJ
	})

	return results, nil
}

func (l *Root) prepareLintContext(spec *arch.Spec, options models.LintOptions) (lintContextReadOnly, error) {
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
		options:       options,
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
				Dir:  string(spec.Project.Directory),
			}

			parsedPackages, err := packages.Load(cfg, string(directory))
			if err != nil {
				return fmt.Errorf("failed parse go source at '%s': %w", dst.PathRel, err)
			}

			if len(parsedPackages) == 0 {
				return fmt.Errorf("directory '%s' not contain any go packages", directory)
			}

			for _, goPackage := range parsedPackages {
				for _, err := range goPackage.Errors {
					if err.Kind == packages.ListError {
						help := ""
						if strings.Contains(err.Msg, "does not contain package") {
							help = " (maybe 'go.work' used and this module is not listed?)"
						}

						return fmt.Errorf("failed load package at '%s': %v%s",
							directory,
							err.Error(),
							help,
						)
					}
				}

				mux.Lock()
				packagesMap[dst.PathRel] = append(packagesMap[dst.PathRel], goPackage)
				mux.Unlock()
			}

			return nil
		})
	}

	err := wg.Wait()
	if err != nil {
		return nil, fmt.Errorf("parse in goroutines failed: %w", err)
	}

	return packagesMap, nil
}
