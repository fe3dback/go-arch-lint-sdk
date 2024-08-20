package yaml

import (
	"bytes"

	"github.com/goccy/go-yaml"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
)

type TransformContext struct {
	file   arch.PathAbsolute
	source []byte
}

func (tc *TransformContext) createReference(documentAstPath string) arch.Reference {
	astPath, err := yaml.PathString(documentAstPath)
	if err != nil {
		return arch.NewInvalidReference()
	}

	astNode, err := astPath.ReadNode(bytes.NewReader(tc.source))
	if err != nil {
		return arch.NewInvalidReference()
	}

	tok := astNode.GetToken()
	return arch.NewReference(
		tc.file,
		tok.Position.Line,
		tok.Position.Column,
		documentAstPath,
	)
}
