package tpl

import (
	"bytes"
	"fmt"
	"text/template"

	"golang.org/x/exp/maps"
)

const (
	fnColorize                  = "colorize"
	fnLines                     = "lines"
	fnTrimPrefix                = "trimPrefix"
	fnTrimSuffix                = "trimSuffix"
	fnTrimDef                   = "def"
	fnPadLeft                   = "padLeft"
	fnPadRight                  = "padRight"
	fnLinePrefix                = "linePrefix"
	fnLinePrefixExceptFirstLine = "linePrefixEFL"
	fnDir                       = "dir"
	fnPlus                      = "plus"
	fnMinus                     = "minus"
	fnConcat                    = "concat"
	fnIsMultiline               = "isMultiline"
)

type Renderer struct {
	asciiColorizer asciiColorizer
	templates      map[string]*template.Template
}

func NewRenderer(
	asciiColorizer asciiColorizer,
) *Renderer {
	return &Renderer{
		asciiColorizer: asciiColorizer,
		templates:      make(map[string]*template.Template, 16),
	}
}

func (r *Renderer) RegisterTemplate(id string, text []byte) error {
	if _, alreadyExist := r.templates[id]; alreadyExist {
		return nil
	}

	tpl, err := template.
		New(id).
		Funcs(map[string]interface{}{
			fnLines:                     r.asciiLines,
			fnColorize:                  r.asciiColorize,
			fnTrimPrefix:                r.asciiTrimPrefix,
			fnTrimSuffix:                r.asciiTrimSuffix,
			fnTrimDef:                   r.asciiDefaultValue,
			fnPadLeft:                   r.asciiPadLeft,
			fnPadRight:                  r.asciiPadRight,
			fnLinePrefix:                r.asciiLinePrefix,
			fnLinePrefixExceptFirstLine: r.asciiLinePrefixExceptFirstLine,
			fnDir:                       r.asciiPathDirectory,
			fnPlus:                      r.asciiPlus,
			fnMinus:                     r.asciiMinus,
			fnConcat:                    r.asciiConcat,
			fnIsMultiline:               r.asciiIsMultiline,
		}).
		Parse(
			preprocessRawASCIITemplate(string(text)),
		)
	if err != nil {
		return fmt.Errorf("failed to compile template '%s': %w", id, err)
	}

	r.templates[id] = tpl
	return nil
}

func (r *Renderer) Render(id string, model any) (string, error) {
	tpl, exist := r.templates[id]
	if !exist {
		return "", fmt.Errorf("template '%s' not exist. Found:[%#v]", id, maps.Keys(r.templates))
	}

	var buffer bytes.Buffer
	err := tpl.Execute(&buffer, model)
	if err != nil {
		return "", fmt.Errorf("failed to execute template '%s': %w", id, err)
	}

	return buffer.String(), nil
}
