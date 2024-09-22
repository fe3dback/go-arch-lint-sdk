package formating

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
)

type NoticeFormatter struct {
	renderer renderer
}

func NewNoticeFormatter(renderer renderer) *NoticeFormatter {
	return &NoticeFormatter{
		renderer: renderer,
	}
}

func (nf *NoticeFormatter) Format(notice *arch.LinterNotice) error {
	id := notice.Details.LinterID

	err := nf.compileIfNecessary(id)
	if err != nil {
		return fmt.Errorf("failed compile template: %w", err)
	}

	out, err := nf.render(id, notice)
	if err != nil {
		return fmt.Errorf("failed render notice in template '%s': %w", id, err)
	}

	notice.Message = out
	return nil
}

func (nf *NoticeFormatter) extractID(linterID arch.LinterID) string {
	return fmt.Sprintf("notices-linter-%s", linterID)
}

func (nf *NoticeFormatter) compileIfNecessary(linterID arch.LinterID) error {
	id := nf.extractID(linterID)
	tplBytes, err := nf.getTemplateBytes(linterID)
	if err != nil {
		return fmt.Errorf("failed get template '%s' bytes: %w", id, err)
	}

	err = nf.renderer.RegisterTemplate(id, tplBytes)
	if err != nil {
		return fmt.Errorf("failed register template '%s': %w", id, err)
	}

	return nil
}

func (nf *NoticeFormatter) render(linterID arch.LinterID, data any) (string, error) {
	id := nf.extractID(linterID)
	return nf.renderer.Render(id, data)
}

func (nf *NoticeFormatter) getTemplateBytes(linterID arch.LinterID) ([]byte, error) {
	switch linterID {
	case arch.LinterIDSyntax:
		return tplNoticeSyntax, nil
	case arch.LinterIDOrphans:
		return tplNoticeOrphans, nil
	case arch.LinterIDImports:
		return tplNoticeImports, nil
	default:
		return nil, fmt.Errorf("formatting template for notice '%s' not defined", linterID)
	}
}