package codeprinter

type (
	linesExtractor interface {
		ExtractLines(fileAbs string, from int, to int) ([]string, error)
	}
)
