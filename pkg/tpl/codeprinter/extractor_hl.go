package codeprinter

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

func init() {
	linterColorTheme, err := chroma.NewStyle("linter", map[chroma.TokenType]string{
		chroma.Literal: "#ffff00",
		chroma.String:  "#00ff00",
		chroma.NameTag: "#00ffff",
	})
	_ = err

	styles.Register(linterColorTheme)
}

type ExtractorHL struct {
}

func NewExtractorHL() *ExtractorHL {
	return &ExtractorHL{}
}

func (e *ExtractorHL) ExtractLines(fileAbs string, from int, to int) ([]string, error) {
	lexer := lexers.Match(fileAbs)

	if lexer == nil {
		lexer = lexers.Fallback
	}

	style := styles.Get("linter")
	if style == nil {
		style = styles.Fallback
	}

	formatter := formatters.Get("terminal8")
	if formatter == nil {
		formatter = formatters.Fallback
	}

	data, err := os.ReadFile(string(fileAbs))
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}

	content := getContents(data, lexer, style, formatter)
	lines := strings.Split(content, "\n")

	return safeTakeLines(lines, from, to), nil
}

func getContents(data []byte, lexer chroma.Lexer, style *chroma.Style, formatter chroma.Formatter) string {
	content := string(data)
	iterator, err := lexer.Tokenise(nil, content)
	if err != nil {
		// fallback to raw output
		return content
	}

	var buff bytes.Buffer
	err = formatter.Format(&buff, style, iterator)
	if err != nil {
		// fallback to raw output
		return content
	}

	return buff.String()
}
