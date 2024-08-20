package validator

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/internal/models"
)

func Test_CommonComponentsValidator(t *testing.T) {
	tests := []struct {
		name    string
		mutator func(config *models.Config)
		out     []string
	}{
		{
			name: "happy",
			mutator: func(config *models.Config) {
				config.CommonComponents = append(config.CommonComponents, nf(arch.ComponentName("my-cmp")))
			},
			out: []string{
				"Common component 'my-cmp' is not known",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vld := NewCommonComponentsValidator()

			in := createValidatorIn()
			vCtx := &validationContext{
				conf: in.conf,
			}

			tt.mutator(&vCtx.conf)
			vld.Validate(vCtx)

			wantNotices := make([]arch.Notice, 0, len(tt.out))
			for _, wantNoticeText := range tt.out {
				wantNotices = append(wantNotices, arch.Notice{
					Message:   wantNoticeText,
					Reference: arch.NewInvalidReference(),
				})
			}

			require.Equal(t, wantNotices, append(vCtx.notices, vCtx.missUsage...))
		})
	}
}
