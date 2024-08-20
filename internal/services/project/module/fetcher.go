package module

import (
	"fmt"
	"os"
	"path"

	"golang.org/x/mod/modfile"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
)

type Fetcher struct {
	rootDirectory arch.PathAbsolute
}

func NewFetcher(
	rootDirectory arch.PathAbsolute,
) *Fetcher {
	return &Fetcher{
		rootDirectory: rootDirectory,
	}
}

func (f *Fetcher) Fetch() (arch.ProjectInfo, error) {
	modPath := arch.PathAbsolute(path.Join(string(f.rootDirectory), "go.mod"))
	modContent, err := os.ReadFile(string(modPath))
	if err != nil {
		return arch.ProjectInfo{}, fmt.Errorf("failed to read config file '%s': %w", modPath, err)
	}

	modData, err := modfile.ParseLax(string(modPath), modContent, nil)
	if err != nil {
		return arch.ProjectInfo{}, fmt.Errorf("failed parse go.mod file '%s': %w", modPath, err)
	}

	return arch.ProjectInfo{
		Directory: f.rootDirectory,
		Module:    arch.GoModule(modData.Module.Mod.Path),
	}, nil
}
