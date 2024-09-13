package codeprinter

type CodePrintMode string

const (
	CodePrintModeOneLine CodePrintMode = "oneLine"
	CodePrintModeExtend  CodePrintMode = "extend"
)

type CodePrintOpts struct {
	Borders     bool
	LineNumbers bool
	Arrows      bool
	ColumnArrow bool
	Highlight   bool
	Mode        CodePrintMode
}
