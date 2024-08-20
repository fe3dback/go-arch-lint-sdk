package yaml_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fe3dback/go-arch-lint-sdk/internal/services/config/reader/yaml"
	testUtils "github.com/fe3dback/go-arch-lint-sdk/internal/services/config/reader/yaml/tests"
)

//go:generate go run ./tests/gen/gen_stub.go
// readme: this command will do:
// - read all files in ./test/*.yml
// - parse it with current parser code
// - output parsed DTO to ./test/*_parsed.go
// - this DTO should be checked by human, when changed
// - auto below test will compare latest transformer with stub file

func TestReader_Parse(t *testing.T) {
	tests := []struct {
		testConfig string
		wantError  string
	}{
		{testConfig: "version_below_supported"},
		{testConfig: "version_above_supported"},
		{testConfig: "syntax_problem_elem"},
		{testConfig: "syntax_problem_sys"},
		{testConfig: "3_min"},
		{testConfig: "3_full"},
		{testConfig: "3_deepscan"},
		{testConfig: "4_min"},
		{testConfig: "4_full"},
		{testConfig: "4_tags_list"},
	}
	for _, tt := range tests {
		reader := yaml.NewReader()

		t.Run(tt.testConfig, func(t *testing.T) {
			sourceCode, err := os.ReadFile(fmt.Sprintf("./tests/data/%s.yml", tt.testConfig))
			require.NoError(t, err)

			conf, err := reader.Parse("/conf.yml", sourceCode)

			if tt.wantError != "" {
				require.EqualError(t, err, tt.wantError)
			} else {
				require.NoError(t, err)
			}

			got := testUtils.Dump(conf)
			wantWithHeader, err := os.ReadFile(fmt.Sprintf("./tests/data/%s_parsed.gold", tt.testConfig))
			require.NoError(t, err)

			want := trimGoldHeader(string(wantWithHeader))

			require.Equal(t, want, got)
		})
	}
}

func trimGoldHeader(content string) string {
	lines := strings.Split(content, "\n")
	if len(lines) < 4 {
		return "unexpected file OR check 'trimGoldHeader' function"
	}

	return strings.Join(lines[3:], "\n")
}
