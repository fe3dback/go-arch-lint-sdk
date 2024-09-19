package linters

import (
	"go/token"
	"strconv"
	"strings"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
)

func referenceFromAstToken(pos token.Position) arch.Reference {
	ref := arch.NewReference(
		arch.PathAbsolute(pos.Filename),
		pos.Line,
		pos.Column,
		"",
	)

	if pos.Line == 0 {
		ref.Valid = false
		ref.Line = 0

		return ref
	}

	return ref
}

// pos has format: "file:line:col" or "file:line" or "" or "-"
func referenceFromAstPos(pos string) arch.Reference {
	if pos == "" || pos == "-" {
		return arch.NewInvalidReference()
	}

	parts := strings.Split(pos, ":")
	switch len(parts) {
	case 3:
		file, line, col := parts[0], parts[1], parts[2]
		lineNum, _ := strconv.Atoi(line)
		colNum, _ := strconv.Atoi(col)

		return arch.NewReference(arch.PathAbsolute(file), lineNum, colNum, "")
	case 2:
		file, line := parts[0], parts[1]
		lineNum, _ := strconv.Atoi(line)

		return arch.NewReference(arch.PathAbsolute(file), lineNum, 0, "")
	default:
		return arch.NewInvalidReference()
	}
}

func orDefault[T any](value *T, defaultValue T) T {
	if value == nil {
		return defaultValue
	}

	return *value
}
