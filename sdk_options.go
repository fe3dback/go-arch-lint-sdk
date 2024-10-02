package sdk

import "github.com/fe3dback/go-arch-lint-sdk/arch"

type (
	CreateOptions struct {
		usedContext  arch.UsedContext
		outputColors bool
		skipMissUse  bool
	}

	CreateOptionsFn func(opt *CreateOptions)
)

// WithSkipMissUse if set to true:
// - will skip not critical notices in config validation
func WithSkipMissUse(skipMissUse bool) CreateOptionsFn {
	return func(opt *CreateOptions) {
		opt.skipMissUse = skipMissUse
	}
}

// WithOutputColors can be used for override default=true mode (ansi)
//
//	true  = ansi  (use colors)
//	false = ascii (only plain text)
func WithOutputColors(outputColors bool) CreateOptionsFn {
	return func(opt *CreateOptions) {
		opt.outputColors = outputColors
	}
}

// WithUsedContext sets the logic context in which the SDK runs
// In 99.9% you don`t need set this option
func WithUsedContext(usedContext arch.UsedContext) CreateOptionsFn {
	return func(opt *CreateOptions) {
		opt.usedContext = usedContext
	}
}
