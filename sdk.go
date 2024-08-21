package sdk

import (
	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/cfg"
	"github.com/fe3dback/go-arch-lint-sdk/internal/container"
	"github.com/fe3dback/go-arch-lint-sdk/mapping"
)

type (
	// SDK is root type for any other interaction with SDK
	// and all functions.
	//
	// Next steps:
	// - call SDK.Spec().FromXXX(..) to get or define "Arch specification"
	// - call SDK.Mapping($spec) to run "mapping" command
	// - call SDK.Check($spec) to run "check" command
	SDK struct {
		di *container.Container
	}
)

// NewSDK Creates new SDK library for specified project
// you need pass absolute projectDirectory where your GO project (go.mod) is located
func NewSDK(projectDirectory arch.PathAbsolute, opts ...CreateOptionsFn) *SDK {
	opt := &CreateOptions{
		usedContext: arch.UsedContextDefault,
		skipMissUse: false,
	}

	for _, mutate := range opts {
		mutate(opt)
	}

	return &SDK{
		di: container.NewContainer(
			projectDirectory,
			opt.usedContext,
			opt.skipMissUse,
		),
	}
}

// Spec returns config builder / fetcher
// Next you can use FromXXX method on it, to get proper Arch specification to work with.
func (sdk *SDK) Spec() *cfg.Definition {
	return sdk.di.ConfigDefinition()
}

// Mapping will return mapping information between project GO packages and
// components defined in spec, it useful for debugging purpose
// (ex: find out how glob path's (ex: "domain/*/repos/**") cover your go files
func (sdk *SDK) Mapping(spec arch.Spec, in mapping.In) (mapping.Out, error) {
	return sdk.di.OperationMapping().Execute(spec, in)
}
