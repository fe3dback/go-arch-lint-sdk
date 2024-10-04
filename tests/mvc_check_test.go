package tests

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/fe3dback/go-arch-lint-sdk"
	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/commands/check"
)

func TestCheck(t *testing.T) {
	// todo: delete this test

	archSDK, err := sdk.NewSDK(arch.PathAbsolute("/home/neo/code/fe3dback/linter/go-arch-lint-sdk"))
	require.NoError(t, err)

	spec, err := archSDK.Spec().FromDefaultFile()
	require.NoError(t, err)

	out, err := archSDK.Check(spec, check.In{
		CheckSyntax: true,
		MaxWarnings: 32,
	})
	require.NoError(t, err)
	archSDK.Assert(t, out)
}
