package reader

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/internal/models"
)

type Reader struct {
	yamlReader fileReader
}

func NewReader(
	yamlReader fileReader,
) *Reader {
	return &Reader{
		yamlReader: yamlReader,
	}
}

func (r *Reader) Read(path arch.PathAbsolute) (models.Config, error) {
	conf, err := r.readFile(path)
	if err != nil {
		return models.Config{}, err
	}

	if len(conf.SyntaxProblems) > 0 {
		var err error

		for _, problem := range conf.SyntaxProblems {
			err = errors.Join(err, arch.NewReferencedError(errors.New(problem.Value), problem.Ref))
		}

		return models.Config{}, fmt.Errorf("found problems in config: %w", err)
	}

	return conf, nil
}

func (r *Reader) readFile(path arch.PathAbsolute) (models.Config, error) {
	ext := filepath.Ext(string(path))

	switch ext {
	case ".yml", ".yaml":
		return r.yamlReader.ReadFile(path)
	default:
		return models.Config{}, fmt.Errorf("unknown config file '%s' ext: %s", path, ext)
	}
}
