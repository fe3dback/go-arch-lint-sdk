package sdk

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/commands/check"
	"github.com/fe3dback/go-arch-lint-sdk/commands/mapping"
	"github.com/fe3dback/go-arch-lint-sdk/definition"
	"github.com/fe3dback/go-arch-lint-sdk/internal/container"
	"github.com/fe3dback/go-arch-lint-sdk/internal/models"
)

type (
	// SDK is root type for any other interaction with SDK
	// and all functions.
	SDK struct {
		projectDirectory arch.PathAbsolute
		usedContext      arch.UsedContext
		di               *container.Container
	}
)

// NewSDK Creates new SDK library for specified project
// you need pass absolute projectDirectory where your GO project (go.mod) is located
//
// Next steps:
// - call SDK.Spec().FromXXX(..) to get or define "Arch specification"
// - call SDK.Mapping($spec) to run "mapping" command
// - call SDK.Check($spec) to run "check" command
func NewSDK(projectDirectory arch.PathAbsolute, opts ...CreateOptionsFn) (*SDK, error) {
	if !filepath.IsAbs(string(projectDirectory)) {
		cleanedRoot, err := filepath.Abs(string(projectDirectory))
		if err != nil {
			return nil, fmt.Errorf("invalid project directory '%s': %w", projectDirectory, err)
		}

		projectDirectory = arch.PathAbsolute(cleanedRoot)
	}

	opt := &CreateOptions{
		usedContext:  arch.UsedContextDefault,
		outputColors: true,
		skipMissUse:  false,
	}

	for _, mutate := range opts {
		mutate(opt)
	}

	return &SDK{
		projectDirectory: projectDirectory,
		usedContext:      opt.usedContext,
		di: container.NewContainer(
			projectDirectory,
			opt.usedContext,
			opt.skipMissUse,
			opt.outputColors,
		),
	}, nil
}

// Spec returns config builder / fetcher
// Next you can use FromXXX method on it, to get proper Arch specification to work with.
func (sdk *SDK) Spec() *definition.Definition {
	return sdk.di.ConfigDefinition()
}

// Mapping will return mapping information between project GO packages and
// components defined in spec, it useful for debugging purpose
// (ex: find out how glob path's (ex: "domain/*/repos/**") cover your go files
func (sdk *SDK) Mapping(spec arch.Spec, in mapping.In) (mapping.Out, error) {
	out, err := sdk.di.OperationMapping().Execute(spec, in)
	err = sdk.wrapErrWithFriendlyHelp(err)

	return out, err
}

// Check will run all configured arch linters and return all found notices
func (sdk *SDK) Check(spec arch.Spec, in check.In) (check.Out, error) {
	out, err := sdk.di.OperationCheck().Execute(spec, in)
	err = sdk.wrapErrWithFriendlyHelp(err)

	return out, err
}

func (sdk *SDK) Assert(t *testing.T, result check.Out) {
	// todo: output full formatted text (like in CLI)

	if result.NoticesCount == 0 {
		t.Logf("Arch linter OK: project '%s' corresponds to defined architecture", sdk.projectDirectory)
		return
	}

	for _, linter := range result.Linters {
		for ind, notice := range linter.Notices {
			// todo: or user normal test ID
			t.Run(fmt.Sprintf("Linter-%s-%d", linter.Linter.ID, ind), func(t *testing.T) {
				t.Errorf("\n\n%s\n\n%s\n\n", notice.Message, notice.ReferencePreview)
			})
		}
	}
}

func (sdk *SDK) wrapErrWithFriendlyHelp(err error) error {
	// happy path
	if err == nil {
		return nil
	}

	// no not wrap anything if we in CLI context
	if sdk.usedContext == arch.UsedContextCLI {
		return err
	}

	// already wrapped, skip
	if _, ok := err.(models.SDKError); ok {
		return err
	}

	// wrap
	return models.NewSDKError(err, sdk.projectDirectory)
}
