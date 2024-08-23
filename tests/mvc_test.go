package tests

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/fe3dback/go-arch-lint-sdk"
	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/commands/mapping"
)

// If you set this to TRUE:
// - test will not compare gold result with actual
// - test will WRITE(update) gold files from actual result
const modeUpdateGold = false

func TestMVC(t *testing.T) {
	projectID := "mvc"

	rootDirectory, err := filepath.Abs(fmt.Sprintf("./projects/%s", projectID))
	require.NoError(t, err)

	archSDK := sdk.NewSDK(arch.PathAbsolute(rootDirectory))

	spec, err := archSDK.Spec().FromDefaultFile()
	require.NoError(t, err)

	type testCase struct {
		in        func(archSDK *sdk.SDK) (wantData any, err error)
		wantError string
	}

	testCases := map[string]testCase{
		"mapping_list": {
			in: func(archSDK *sdk.SDK) (any, error) {
				return archSDK.Mapping(spec, mapping.In{
					Scheme: mapping.SchemeList,
				})
			},
		},
		"mapping_grouped": {
			in: func(archSDK *sdk.SDK) (any, error) {
				return archSDK.Mapping(spec, mapping.In{
					Scheme: mapping.SchemeGrouped,
				})
			},
		},
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			outData, err := tt.in(archSDK)

			if tt.wantError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.wantError)
			}

			outBytes, err := json.MarshalIndent(outData, "", "  ")
			require.NoError(t, err)

			outBytes = replaceRoots(rootDirectory, outBytes)
			goldFile := filepath.Join("./data", projectID, fmt.Sprintf("%s.gold.json", name))

			if modeUpdateGold {
				err := os.WriteFile(goldFile, outBytes, 0644)
				require.NoError(t, err)
			} else {
				wantBytes, err := os.ReadFile(goldFile)
				require.NoError(t, err)

				require.Equal(t, string(wantBytes), string(outBytes))
			}
		})
	}
}

func replaceRoots(rootDirectory string, src []byte) []byte {
	data := string(src)
	data = strings.ReplaceAll(data, rootDirectory, "/root")

	return []byte(data)
}
