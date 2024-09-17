package tests

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/fe3dback/go-arch-lint-sdk"
	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/definition/cfg"
)

func TestScratch(t *testing.T) {
	projectID := "mvc"

	rootDirectory, err := filepath.Abs(fmt.Sprintf("./projects/%s", projectID))
	require.NoError(t, err)

	archSDK, err := sdk.NewSDK(arch.PathAbsolute(rootDirectory))
	require.NoError(t, err)

	const (
		cmpApp          = "app"
		cmpHandlers     = "handlers"
		cmpServices     = "services"
		cmpRepositories = "repositories"
		cmpModels       = "models"
	)

	// todo: alt Spec().FromYAML("version: 4\n ...")

	spec, err := archSDK.Spec().FromCode(
		// settings
		cfg.WorkingDirectory("internal"),
		// components
		cfg.Component(cmpApp, "."),
		cfg.Component(cmpHandlers, "handlers/*"),
		cfg.Component(cmpServices, "domains/*/services/**"),
		cfg.Component(cmpRepositories, "repositories/**"),
		cfg.Component(cmpModels, "models"),
		// common
		cfg.CommonComponents(cmpModels),
		// todo: other config props
	)
	require.NoError(t, err)

	_ = spec
}
