package models

import (
	"bytes"
	"fmt"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/pkg/stringutil"
)

// SDKError used for user-friendly error formatting
// this error will recursive parse wrapped errors for better display
type SDKError struct {
	internal         error
	projectDirectory arch.PathAbsolute
}

func NewSDKError(internal error, projectDirectory arch.PathAbsolute) SDKError {
	return SDKError{
		internal:         internal,
		projectDirectory: projectDirectory,
	}
}

func (e SDKError) Error() string {
	internalText := e.internalError()

	var buf bytes.Buffer
	buf.WriteString("Golang architectural linter SDK:\n")
	buf.WriteString("----------------------------\n")
	buf.WriteString(fmt.Sprintf("- project directory: %s\n", e.projectDirectory))
	buf.WriteString(fmt.Sprintf("- read more        : %s\n", "https://github.com/fe3dback/go-arch-lint-sdk"))
	buf.WriteString("----------------------------\n")
	buf.WriteString(stringutil.PrefixLines(internalText, "  |] "))
	buf.WriteString("\n")
	buf.WriteString("----------------------------\n\n")

	return buf.String()
}

func (e SDKError) internalError() string {
	return e.internal.Error()
}
