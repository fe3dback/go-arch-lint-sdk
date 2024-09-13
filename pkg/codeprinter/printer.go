package codeprinter

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Printer struct {
	rawExtractor         linesExtractor
	highlightedExtractor linesExtractor
}

func NewPrinter(
	rawExtractor linesExtractor,
	highlightedExtractor linesExtractor,
) *Printer {
	return &Printer{
		rawExtractor:         rawExtractor,
		highlightedExtractor: highlightedExtractor,
	}
}

func (p *Printer) Print(ref Reference, opts CodePrintOpts) (string, error) {
	if !ref.Valid {
		return "", nil
	}

	maxLines, err := fileLinesCount(ref.File)
	if err != nil {
		return "", fmt.Errorf("failed count file lines: %w", err)
	}

	region := calculateRegion(ref, opts, maxLines)
	lines, err := p.extractLines(opts, region)
	if err != nil {
		return "", fmt.Errorf("failed extract lines from '%s' [%d-%d]: %w", ref.File, region.from, region.to, err)
	}

	if opts.LineNumbers {
		lines = addLinesNumbers(lines, region)
	}

	if opts.Arrows {
		lines = addArrows(lines, region, opts.ColumnArrow)
	}

	if opts.Borders {
		lines = addBorders(lines, region)
	}

	return strings.Join(lines, "\n"), nil
}

func (p *Printer) extractLines(opts CodePrintOpts, region area) ([]string, error) {
	if opts.Highlight {
		return p.highlightedExtractor.ExtractLines(region.ref.File, region.from, region.to)
	}

	return p.rawExtractor.ExtractLines(region.ref.File, region.from, region.to)
}

func fileLinesCount(pathAbs string) (int, error) {
	data, err := os.ReadFile(pathAbs)
	if err != nil {
		return 0, fmt.Errorf("failed to read file '%s': %w", pathAbs, err)
	}

	return len(strings.Split(string(data), "\n")), nil
}

func calculateRegion(ref Reference, opts CodePrintOpts, maxLines int) area {
	line := clamp(ref.Line, 1, maxLines)
	from, to := line, line

	if opts.Mode == CodePrintModeExtend {
		from -= 1
		to += 2
	}

	from = max(from, 1)
	to = min(to, maxLines)

	return area{
		ref:  ref,
		from: from,
		to:   to,
	}
}

func addLinesNumbers(lines []string, region area) []string {
	width := len(fmt.Sprintf("%d", region.to))

	region.lineNumbers(func(ind int, line int, isReferenced bool) {
		paddedNumber := padLeft(width, " ", strconv.FormatInt(int64(line), 10))

		lines[ind] = fmt.Sprintf("%s | %s", paddedNumber, lines[ind])
	})

	return lines
}

func padLeft(width int, padStr string, value string) string {
	padCountInt := 1 + ((width - len(padStr)) / len(padStr))
	retStr := strings.Repeat(padStr, padCountInt) + value
	return retStr[(len(retStr) - width):]
}

func addArrows(lines []string, region area, cursor bool) []string {
	result := make([]string, 0, len(lines)+1)
	lineNumberOffset := strings.IndexByte(lines[0], '|') + 3

	region.lineNumbers(func(ind int, line int, isReferenced bool) {
		if isReferenced {
			column := clamp(region.ref.Column, 1, len(lines[ind]))

			result = append(result, "> "+lines[ind])
			if cursor {
				result = append(result, strings.Repeat(" ", lineNumberOffset+column)+"^")
			}
		} else {
			result = append(result, "  "+lines[ind])
		}
	})

	return result
}

func addBorders(lines []string, region area) []string {
	result := make([]string, 0, len(lines)+2)

	name := filepath.Base(region.ref.File)
	width := len(name)

	result = append(result, name)
	result = append(result, strings.Repeat("~", width))
	result = append(result, lines...)

	return result
}
