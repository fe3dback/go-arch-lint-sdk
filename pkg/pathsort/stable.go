package pathsort

import (
	"sort"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
)

func SortFileTree[T any](list []T, getRel func(*T) (relPath arch.PathRelative, isDir bool)) {
	tree := newNode("/")
	for _, elem := range list {
		pathRel, isDir := getRel(&elem)
		tree.append(value{pathRel: pathRel, isDir: isDir})
	}
	tree.sortLevels()

	newOrderList := tree.traversalDepthFirst()
	newOrderMap := make(map[arch.PathRelative]int)

	for order, dst := range newOrderList {
		newOrderMap[dst.pathRel] = order
	}

	sort.Slice(list, func(i, j int) bool {
		pathRelI, _ := getRel(&list[i])
		pathRelJ, _ := getRel(&list[j])

		orderI, orderJ := newOrderMap[pathRelI], newOrderMap[pathRelJ]
		return orderI <= orderJ
	})
}
