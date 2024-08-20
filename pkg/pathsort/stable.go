package pathsort

import (
	"sort"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
)

func SortDescriptors(descriptors []arch.FileDescriptor) {
	tree := newNode("/")
	for _, descriptor := range descriptors {
		tree.append(&descriptor)
	}
	tree.sortLevels()

	newOrderList := tree.traversalDepthFirst()
	newOrderMap := make(map[arch.PathRelative]int)

	for order, dst := range newOrderList {
		newOrderMap[dst.PathRel] = order
	}

	sort.Slice(descriptors, func(i, j int) bool {
		orderI, orderJ := newOrderMap[descriptors[i].PathRel], newOrderMap[descriptors[j].PathRel]
		return orderI <= orderJ
	})
}
