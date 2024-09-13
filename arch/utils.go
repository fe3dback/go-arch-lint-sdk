package arch

import "github.com/fe3dback/go-arch-lint-sdk/pkg/codeprinter"

func printCode(ref Reference) (string, bool) {
	if !ref.Valid {
		return "", false
	}

	printer := codeprinter.NewPrinter(
		codeprinter.NewExtractorRaw(),
		codeprinter.NewExtractorHL(),
	)

	internalRef := codeprinter.Reference{
		File:   string(ref.File),
		Line:   ref.Line,
		Column: ref.Column,
		Valid:  ref.Valid,
	}

	preview, err := printer.Print(internalRef, codeprinter.CodePrintOpts{
		LineNumbers: true,
		Arrows:      true,
		Mode:        codeprinter.CodePrintModeExtend,
	})
	if err != nil {
		return "", false
	}

	return preview, true
}
