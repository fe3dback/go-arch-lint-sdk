package codeprinter

import (
	"fmt"
	"os"
	"strings"
)

type ExtractorRaw struct{}

func NewExtractorRaw() *ExtractorRaw {
	return &ExtractorRaw{}
}

func (e *ExtractorRaw) ExtractLines(fileAbs string, from int, to int) ([]string, error) {
	data, err := os.ReadFile(fileAbs)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	return safeTakeLines(lines, from, to), nil
}
