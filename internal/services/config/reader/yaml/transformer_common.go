package yaml

import (
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/internal/models"
)

func transformFromSyntaxError(tCtx TransformContext, err error) models.Config {
	errText := yaml.FormatError(err, false, true)
	errLines := strings.Split(errText, "\n")
	plainErr := errLines[0]

	ref := extractReferenceFromError(tCtx, errLines[1:])

	return models.Config{
		SyntaxProblems: []arch.Ref[string]{
			arch.NewRef(plainErr, ref),
		},
	}
}

func transform(tCtx TransformContext, doc any) (models.Config, error) {
	switch typed := doc.(type) {
	case *ModelV3:
		return transformV3(tCtx, *typed), nil
	case *ModelV4:
		return transformV4(tCtx, *typed), nil
	}

	return models.Config{}, fmt.Errorf("unknown document type %T", doc)
}
