package codeprinter_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fe3dback/go-arch-lint-sdk/pkg/codeprinter"
)

const (
	//
	modeGenerateGold = "gen"
	modeVerify       = "verify"
)

// will regenerate all *.golden files in tests
// you need to verify all generated files before commit
// and change mode back to "verify"
const mode = modeVerify

func TestPrinter_Print(t *testing.T) {
	type ref struct {
		testFile string
		line     int
		column   int
	}

	matrix := map[string]codeprinter.CodePrintOpts{
		"one_line": {
			Borders:     false,
			LineNumbers: false,
			Arrows:      false,
			Highlight:   false,
			Mode:        codeprinter.CodePrintModeOneLine,
		},
		"b0_n0_a1_h0_mOL": {
			Borders:     false,
			LineNumbers: false,
			Arrows:      true,
			Highlight:   false,
			Mode:        codeprinter.CodePrintModeOneLine,
		},
		"b0_n1_a1_h0_mOL": {
			Borders:     false,
			LineNumbers: true,
			Arrows:      true,
			Highlight:   false,
			Mode:        codeprinter.CodePrintModeOneLine,
		},
		"b0_n1_a0_h0_mE": {
			Borders:     false,
			LineNumbers: true,
			Arrows:      false,
			Highlight:   false,
			Mode:        codeprinter.CodePrintModeExtend,
		},
		"b0_n1_a1_h0_mE": {
			Borders:     false,
			LineNumbers: true,
			Arrows:      true,
			Highlight:   false,
			Mode:        codeprinter.CodePrintModeExtend,
		},
		"b0_n1_a1_h1_mE": {
			Borders:     false,
			LineNumbers: true,
			Arrows:      true,
			Highlight:   true,
			Mode:        codeprinter.CodePrintModeExtend,
		},
		"full": {
			Borders:     true,
			LineNumbers: true,
			Arrows:      true,
			Highlight:   false,
			Mode:        codeprinter.CodePrintModeExtend,
		},
		"full_colored": {
			Borders:     true,
			LineNumbers: true,
			Arrows:      true,
			Highlight:   true,
			Mode:        codeprinter.CodePrintModeExtend,
		},
	}

	tests := []struct {
		group   string
		name    string
		ref     ref
		wantErr string
	}{
		{
			group: "yaml",
			name:  "arr_start",
			ref:   ref{testFile: "bigconf.yml", line: 10, column: 9},
		},
		{
			group: "yaml",
			name:  "first_line",
			ref:   ref{testFile: "bigconf.yml", line: 1, column: 0},
		},
		{
			group: "yaml",
			name:  "above_max",
			ref:   ref{testFile: "bigconf.yml", line: 9000, column: 4000},
		},
		{
			group: "go",
			name:  "time_value",
			ref:   ref{testFile: "some_code.go", line: 17, column: 16},
		},
		{
			group: "go",
			name:  "below_zero",
			ref:   ref{testFile: "some_code.go", line: -15, column: 3},
		},
		{
			group: "go",
			name:  "strange_column",
			ref:   ref{testFile: "some_code.go", line: 12, column: 5000},
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s_%s", tt.group, tt.name), func(t *testing.T) {
			p := codeprinter.NewPrinter(
				codeprinter.NewExtractorRaw(),
				codeprinter.NewExtractorHL(),
			)

			pathSrc := filepath.Clean(fmt.Sprintf("./tests/%s", tt.ref.testFile))
			srcReference := codeprinter.Reference{
				File:   pathSrc,
				Line:   tt.ref.line,
				Column: tt.ref.column,
				Valid:  true,
			}

			for variantName, opts := range matrix {
				dirName := strings.ReplaceAll(tt.ref.testFile, ".", "_")
				pathDst := filepath.Clean(fmt.Sprintf("./tests/%s/%s/%s.golden", dirName, tt.name, variantName))

				got, err := p.Print(srcReference, opts)
				require.NoError(t, err)

				switch mode {
				case modeGenerateGold:
					err = os.MkdirAll(filepath.Dir(pathDst), os.ModePerm)
					require.NoError(t, err)

					err = os.WriteFile(pathDst, []byte(got), 0600)
					require.NoError(t, err)
				case modeVerify:
					want, err := os.ReadFile(pathDst)
					require.NoError(t, err)
					require.Equal(t, string(want), got)

					t.Logf("\nout:\n--\n%s\n--\n", got)
				}
			}
		})
	}
}
