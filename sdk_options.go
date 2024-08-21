package sdk

import "github.com/fe3dback/go-arch-lint-sdk/arch"

type (
	CreateOptions struct {
		usedContext arch.UsedContext
		skipMissUse bool
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

// WithUsedContext sets the logic context in which the SDK runs
// In 99.9% you don`t need set this option
func WithUsedContext(usedContext arch.UsedContext) CreateOptionsFn {
	return func(opt *CreateOptions) {
		opt.usedContext = usedContext
	}
}
