package cfg

import (
	"fmt"
	"path/filepath"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
)

// FromDefaultFile will parse config from "{projectDirectory}/.go-arch-lint.yml"
// see also: FromRelativeFile, FromAbsoluteFile
func (def *Definition) FromDefaultFile() (arch.Spec, error) {
	return def.FromRelativeFile(".go-arch-lint.yml")
}

// FromRelativeFile will find and parse config file RELATIVE to your project directory
// you can also use FromDefaultFile() (for use default file ".go-arch-lint.yml")
func (def *Definition) FromRelativeFile(path arch.PathRelative) (arch.Spec, error) {
	return def.FromAbsoluteFile(
		arch.PathAbsolute(filepath.Join(string(def.projectPath), string(path))),
	)
}

// FromAbsoluteFile will find and parse config file in any directory, but all paths
// defined inside config will be related to project directory anyway
// see also: FromRelativeFile
func (def *Definition) FromAbsoluteFile(filePath arch.PathAbsolute) (arch.Spec, error) {
	config, err := def.reader.Read(filePath)
	if err != nil {
		return arch.Spec{}, fmt.Errorf("failed to read config at '%s': %w", filePath, err)
	}

	err = def.validator.Validate(config)
	if err != nil {
		return arch.Spec{}, fmt.Errorf("invalid config: %w", err)
	}

	spec, err := def.assembler.Assemble(config)
	if err != nil {
		return arch.Spec{}, fmt.Errorf("failed assemble spec: %w", err)
	}

	return spec, nil
}
