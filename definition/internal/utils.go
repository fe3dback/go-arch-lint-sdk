package internal

import (
	"runtime"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
)

func getCaller(skip int) arch.Reference {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return arch.NewInvalidReference()
	}

	return arch.NewReference(arch.PathAbsolute(file), line, 0, "")
}

func GetParentCaller() arch.Reference {
	return getCaller(3)
}
