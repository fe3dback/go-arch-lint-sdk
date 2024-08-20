package sdk

type (
	CreateOptions struct {
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
