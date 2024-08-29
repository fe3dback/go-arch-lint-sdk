package linters

import (
	"go/token"

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

func orDefault[T any](value *T, defaultValue T) T {
	if value == nil {
		return defaultValue
	}

	return *value
}
