package pathsort

import (
	"fmt"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
)

func dsc(relPath string) arch.FileDescriptor {
	ext := filepath.Ext(relPath)
	isDir := ext == ""

	return arch.FileDescriptor{
		PathRel:   arch.PathRelative(relPath),
		PathAbs:   arch.PathAbsolute("/project/" + relPath),
		IsDir:     isDir,
		Extension: ext,
	}
}

func TestSortDescriptors(t *testing.T) {
	want := []arch.FileDescriptor{
		dsc("conf/assembler"),
		dsc("conf/assembler/tests"),
		dsc("conf/assembler/tests/some.txt"),
		dsc("conf/assembler/aaa.go"),
		dsc("conf/assembler/c.go"),
		dsc("conf/assembler/c.txt"),
		dsc("conf/assembler/ggg.go"),
		dsc("conf/reader/yaml/tests/data/3_full.yml"),
		dsc("conf/reader/yaml/tests/utils.go"),
		dsc("conf/reader/yaml/reader.go"),
		dsc("conf/reader/yaml/utils.go"),
		dsc("conf/reader/interfaces.go"),
		dsc("conf/reader/reader.go"),
		dsc("conf/validator/ctx.go"),
		dsc("conf/validator/root.go"),
		dsc("project/reader"),
		dsc("project/reader/reader.go"),
		dsc("project/validator/ctx.go"),
		dsc("project/validator/root.go"),
	}

	in := make([]arch.FileDescriptor, len(want))
	copy(in, want)

	sort.Slice(in, func(_, _ int) bool {
		return rand.Int31n(100) < 50
	})

	got := make([]arch.FileDescriptor, len(in))
	copy(got, in)

	SortDescriptors(got)

	if !assert.Equal(t, want, got) {
		printSlice("in", in)
		printSlice("got", got)
		printSlice("want", want)
	}
}

func printSlice(name string, list []arch.FileDescriptor) {
	fmt.Printf("%s:\n", name)

	for _, value := range list {
		fmt.Printf("- %s:\n", value.PathRel)
	}

	fmt.Printf("\n\n")
}
