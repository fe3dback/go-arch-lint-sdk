package xpath

import (
	"path/filepath"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
)

type index struct {
	full           []arch.PathDescriptor
	files          map[arch.PathRelative]*arch.PathDescriptor
	directories    map[arch.PathRelative]*arch.PathDescriptor
	directoryFiles map[arch.PathRelative][]*arch.PathDescriptor
}

func newIndex() *index {
	return &index{
		full:           make([]arch.PathDescriptor, 0, 256),
		files:          make(map[arch.PathRelative]*arch.PathDescriptor, 256),
		directories:    make(map[arch.PathRelative]*arch.PathDescriptor, 64),
		directoryFiles: make(map[arch.PathRelative][]*arch.PathDescriptor, 64),
	}
}

func (ind *index) appendToIndex(path arch.PathRelative, src arch.PathDescriptor) {
	parent := arch.PathRelative(filepath.Dir(string(path)))

	ind.full = append(ind.full, src)
	descriptor := &ind.full[len(ind.full)-1]

	// add file to index
	if !descriptor.IsDir {
		ind.files[path] = descriptor
	}

	// create dirs index if not exist
	switch descriptor.IsDir {
	case true:
		ind.directories[path] = descriptor
		if _, exists := ind.directoryFiles[path]; !exists {
			ind.directoryFiles[path] = make([]*arch.PathDescriptor, 0, 8)
		}
	case false:
		if _, exists := ind.directoryFiles[parent]; !exists {
			ind.directoryFiles[parent] = make([]*arch.PathDescriptor, 0, 8)
		}
	}

	// add file to dir index
	if !descriptor.IsDir {
		ind.directoryFiles[parent] = append(ind.directoryFiles[parent], descriptor)
	}
}

func (ind *index) fileAt(path arch.PathRelative) (arch.PathDescriptor, bool) {
	dst, ok := ind.files[path]
	if !ok {
		return arch.PathDescriptor{}, false
	}

	if dst.IsDir {
		return arch.PathDescriptor{}, false
	}

	return *dst, true
}

func (ind *index) directoryAt(path arch.PathRelative) (arch.PathDescriptor, bool) {
	dst, ok := ind.directories[path]
	if !ok {
		return arch.PathDescriptor{}, false
	}

	if !dst.IsDir {
		return arch.PathDescriptor{}, false
	}

	return *dst, true
}

func (ind *index) each(fn func(arch.PathDescriptor)) {
	for _, descriptor := range ind.full {
		fn(descriptor)
	}
}
