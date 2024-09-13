package mapping

import (
	"path/filepath"
	"sort"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/commands/mapping"
)

const unknownComponentID = "[not attached]"

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
	flatList := o.buildList(spec)
	components := make(map[string][]string, 32)

	for _, elem := range flatList {
		cmpPackages, exist := components[elem.ComponentName]
		if !exist {
			cmpPackages = make([]string, 0, 8)
		}

		cmpPackages = append(cmpPackages, elem.Package)
		components[elem.ComponentName] = cmpPackages
	}

	list := make([]mapping.OutGrouped, 0, 32)

	for componentName, ownedPackages := range components {
		group := mapping.OutGrouped{
			ComponentName:  componentName,
			Packages:       make([]string, 0, len(ownedPackages)),
			ComponentExist: componentName != unknownComponentID,
		}

		for _, ownedPackage := range ownedPackages {
			group.Packages = append(group.Packages, ownedPackage)
		}

		list = append(list, group)
	}

	// add empty components
	for _, component := range spec.Components {
		if len(component.OwnedPackages) > 0 {
			continue
		}

		list = append(list, mapping.OutGrouped{
			ComponentName:  string(component.Name.Value),
			Packages:       make([]string, 0),
			ComponentExist: true,
		})
	}

	sort.SliceStable(list, func(i, j int) bool {
		if list[i].ComponentName != list[j].ComponentName {
			return list[i].ComponentName <= list[j].ComponentName
		}

		return i < j
	})

	return list
}

func (o *Operation) buildList(spec arch.Spec) []mapping.OutList {

	list := make([]mapping.OutList, 0, 128)

	for _, component := range spec.Components {
		for _, ownedPackage := range component.OwnedPackages {
			list = append(list, mapping.OutList{
				Package:        string(ownedPackage.PathAbs),
				ComponentName:  string(component.Name.Value),
				ComponentExist: true,
			})
		}
	}

	orphanProcessedPackages := make(map[string]any, 4)

	for _, orphan := range spec.Orphans {
		orphanPackage := filepath.Dir(string(orphan.File.PathAbs))

		// filter unique only
		if _, ok := orphanProcessedPackages[orphanPackage]; ok {
			continue
		}

		orphanProcessedPackages[orphanPackage] = struct{}{}

		// add to list
		list = append(list, mapping.OutList{
			Package:        orphanPackage,
			ComponentName:  unknownComponentID,
			ComponentExist: false,
		})
	}

	sort.SliceStable(list, func(i, j int) bool {
		if list[i].ComponentName != list[j].ComponentName {
			return list[i].ComponentName <= list[j].ComponentName
		}

		return i < j
	})

	return list
}
