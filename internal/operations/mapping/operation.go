package mapping

import (
	"sort"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/commands/mapping"
)

type Operation struct{}

func NewOperation() *Operation {
	return &Operation{}
}

func (o *Operation) Execute(spec arch.Spec, in mapping.In) (mapping.Out, error) {
	return mapping.Out{
		ProjectDirectory: string(spec.Project.Directory),
		ModuleName:       string(spec.Project.Module),
		MappingGrouped:   o.buildGrouped(spec),
		MappingList:      o.buildList(spec),
		Scheme:           in.Scheme,
	}, nil
}

func (o *Operation) buildGrouped(spec arch.Spec) []mapping.OutGrouped {
	list := make([]mapping.OutGrouped, 0, len(spec.Components))

	for _, component := range spec.Components {
		group := mapping.OutGrouped{
			ComponentName: string(component.Name.Value),
			Packages:      make([]string, 0, len(component.OwnedPackages)),
		}

		for _, pkg := range component.OwnedPackages {
			group.Packages = append(group.Packages, string(pkg.PathAbs))
		}

		list = append(list, group)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].ComponentName <= list[j].ComponentName
	})

	return list
}

func (o *Operation) buildList(spec arch.Spec) []mapping.OutList {
	list := make([]mapping.OutList, 0, 128)

	for _, component := range spec.Components {
		for _, ownedPackage := range component.OwnedPackages {
			list = append(list, mapping.OutList{
				Package:       string(ownedPackage.PathAbs),
				ComponentName: string(component.Name.Value),
			})
		}
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].ComponentName <= list[j].ComponentName
	})

	return list
}
