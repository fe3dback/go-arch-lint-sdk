package validator

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/internal/models"
	"github.com/fe3dback/go-arch-lint-sdk/pkg/tpl/codeprinter"
)

type Root struct {
	usedContext arch.UsedContext
	skipMissuse bool
	printer     printer
	validators  []internalValidator
}

func NewRoot(
	usedContext arch.UsedContext,
	skipMissuse bool,
	printer printer,
	validators ...internalValidator,
) *Root {
	return &Root{
		usedContext: usedContext,
		skipMissuse: skipMissuse,
		printer:     printer,
		validators:  validators,
	}
}

func (v *Root) Validate(config models.Config) error {
	ctx := &validationContext{
		conf:             config,
		notices:          make([]arch.Notice, 0, 16),
		currentValidator: "main",
	}

	for _, validator := range v.validators {
		ctx.currentValidator = fmt.Sprintf("%T", validator)
		validator.Validate(ctx)

		if ctx.critical {
			break
		}
	}

	for ind := range ctx.notices {
		notice := &ctx.notices[ind]
		shortPreview := len(ctx.notices) >= 10

		err := v.renderPreview(notice, shortPreview)
		if err != nil {
			return err
		}
	}

	if len(ctx.notices) > 0 {
		return arch.NewErrorWithNotices(
			"Config validator find some notices",
			ctx.notices,
			v.usedContext != arch.UsedContextCLI,
		)
	}

	if !v.skipMissuse && len(ctx.missUsage) > 0 {
		return arch.NewErrorWithNotices(
			"Config validator find miss usages. You can hide this message by adding '--skip-missuse' flag",
			ctx.missUsage,
			v.usedContext != arch.UsedContextCLI,
		)
	}

	return nil
}

func (v *Root) renderPreview(notice *arch.Notice, shortPreview bool) error {
	if !notice.Reference.Valid {
		return nil
	}

	internalRef := codeprinter.Reference{
		File:   string(notice.Reference.File),
		Line:   notice.Reference.Line,
		Column: notice.Reference.Column,
		Valid:  notice.Reference.Valid,
	}

	mode := codeprinter.CodePrintModeExtend
	if shortPreview {
		mode = codeprinter.CodePrintModeOneLine
	}

	preview, err := v.printer.Print(internalRef, codeprinter.CodePrintOpts{
		LineNumbers: true,
		Arrows:      true,
		Mode:        mode,
	})
	if err != nil {
		return fmt.Errorf("error while rendering notice preview: %w", err)
	}

	notice.CodePreview = preview

	return nil
}
