package tests

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/fe3dback/go-arch-lint-sdk"
	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/commands/check"
)

func TestCheck(t *testing.T) {
	// todo: delete this test

	archSDK, err := sdk.NewSDK(arch.PathAbsolute("/home/neo/code/fe3dback/linter/go-arch-lint/v4"))
	require.NoError(t, err)

	spec, err := archSDK.Spec().FromDefaultFile()
	require.NoError(t, err)

	out, err := archSDK.Check(spec, check.In{
		MaxWarnings: 32,
	})
	require.NoError(t, err)

	formattedOut, err := json.MarshalIndent(out, "", "  ")
	require.NoError(t, err)

	fmt.Println(string(formattedOut))
}
