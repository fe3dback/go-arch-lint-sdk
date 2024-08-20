package cfg

import (
	"fmt"
	"path/filepath"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
)

// FromFile will find and parse config file RELATIVE to your project directory
// you can also use FromDefaultFile() (for use default file ".go-arch-lint.yml")
func (def *Definition) FromFile(path arch.PathRelative) (arch.Spec, error) {
	filePath := arch.PathAbsolute(filepath.Join(string(def.projectPath), string(path)))

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

func (def *Definition) FromDefaultFile() (arch.Spec, error) {
	return def.FromFile(".go-arch-lint.yml")
}
