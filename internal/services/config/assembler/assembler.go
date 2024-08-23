package assembler

import (
	"fmt"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/internal/models"
	"github.com/fe3dback/go-arch-lint-sdk/pkg/pathsort"
)

type Assembler struct {
	projectInfoFetcher projectInfoFetcher
	pathHelper         pathHelper
}

func NewAssembler(projectInfoFetcher projectInfoFetcher, pathHelper pathHelper) *Assembler {
	return &Assembler{
		projectInfoFetcher: projectInfoFetcher,
		pathHelper:         pathHelper,
	}
}

// nolint
func (a *Assembler) Assemble(conf models.Config) (spec arch.Spec, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed assemble: %v\n%s", r, debug.Stack())
		}
	}()

	projectInfo, err := a.projectInfoFetcher.Fetch()
	if err != nil {
		return arch.Spec{}, fmt.Errorf("failed fetch project info: %w", err)
	}

	components := make([]*arch.SpecComponent, 0, conf.Dependencies.Map.Len())
	mappedFiles := make(map[arch.PathRelative]any, 128)

	conf.Components.Map.Each(func(name arch.ComponentName, definition models.ConfigComponent, definitionRef arch.Reference) {

		rules, rulesRef, exist := conf.Dependencies.Map.Get(name)
		if !exist {
			// defaults for rules
			rules = models.ConfigComponentDependencies{}
			rulesRef = arch.NewInvalidReference()
		}

		var deepScan arch.Ref[bool]
		if rules.DeepScan.Defined {
			deepScan = rules.DeepScan.Value
		} else {
			deepScan = conf.Settings.DeepScan
		}

		tagsAllowedAll, tagsAllowedWhiteList := a.figureOutAllowedStructTags(&conf, &rules)

		matchedFiles, err := a.findOwnedFiles(&conf, definition)
		if err != nil {
			panic(fmt.Errorf("failed finding owned files by component '%s': %w", name, err))
		}

		components = append(components, &arch.SpecComponent{
			Name:                arch.NewRef(name, definitionRef),
			DefinitionComponent: definitionRef,
			DefinitionDeps:      rulesRef,
			DeepScan:            deepScan,
			StrictMode:          conf.Settings.Imports.StrictMode,
			AllowAllProjectDeps: rules.AnyProjectDeps,
			AllowAllVendorDeps:  rules.AnyVendorDeps,
			AllowAllTags:        tagsAllowedAll,
			AllowedTags:         tagsAllowedWhiteList,
			MayDependOn:         rules.MayDependOn,
			CanUse:              rules.CanUse,
			MatchPatterns:       definition.In,
			MatchedFiles:        matchedFiles,
		})

		for _, dst := range matchedFiles {
			mappedFiles[dst.PathRel] = struct{}{}
		}
	})

	// copy matched files to owned files (but each file owned only by one component)
	a.calculateFilesOwnage(components)

	// find matched/owned packages from files
	for _, component := range components {
		component.MatchedPackages = a.extractUniquePackages(component.MatchedFiles)
		component.OwnedPackages = a.extractUniquePackages(component.OwnedFiles)
	}

	// find orphan files
	files, err := a.pathHelper.FindProjectFiles(arch.FileQuery{
		Path:               arch.PathRelativeGlob("*/**"),
		WorkingDirectory:   conf.WorkingDirectory.Value,
		Type:               arch.FileMatchQueryTypeOnlyFiles,
		ExcludeDirectories: conf.Exclude.RelativeDirectories.Values(),
		ExcludeRegexp:      conf.Exclude.RelativeFiles.Values(),
		Extensions:         []string{"go"},
	})
	if err != nil {
		return arch.Spec{}, fmt.Errorf("failed finding project files: %w", err)
	}

	orphanFiles := make([]arch.SpecOrphan, 0, 32)
	for _, projectFile := range files {
		if _, mapped := mappedFiles[projectFile.PathRel]; mapped {
			continue
		}

		orphanFiles = append(orphanFiles, arch.SpecOrphan{
			File: projectFile,
		})
	}

	// sort paths
	for _, component := range components {
		pathsort.SortDescriptors(component.MatchedFiles)
		pathsort.SortDescriptors(component.MatchedPackages)
		pathsort.SortDescriptors(component.OwnedFiles)
		pathsort.SortDescriptors(component.OwnedPackages)
	}

	// finalize
	resultComponents := make([]arch.SpecComponent, 0, len(components))
	for _, component := range components {
		resultComponents = append(resultComponents, *component)
	}

	return arch.Spec{
		Project:          projectInfo,
		WorkingDirectory: conf.WorkingDirectory,
		Components:       resultComponents,
		Orphans:          orphanFiles,
	}, nil
}

func (a *Assembler) figureOutAllowedStructTags(conf *models.Config, rules *models.ConfigComponentDependencies) (arch.Ref[bool], arch.RefSlice[arch.StructTag]) {
	if conf.Settings.Tags.Allowed.Value == arch.ConfigSettingsTagsEnumAll {
		return arch.NewRef(true, conf.Settings.Tags.Allowed.Ref), nil
	}

	globalTags := conf.Settings.Tags.AllowedList
	localTags := rules.CanContainTags

	allowedList := make(arch.RefSlice[arch.StructTag], 0, len(globalTags)+len(localTags))
	allowedList = append(allowedList, globalTags...)
	allowedList = append(allowedList, localTags...)

	return arch.NewRef(false, conf.Settings.Tags.Allowed.Ref), allowedList
}

func (a *Assembler) findOwnedFiles(conf *models.Config, component models.ConfigComponent) ([]arch.FileDescriptor, error) {
	list := make([]arch.FileDescriptor, 0, 32)

	for _, globPath := range component.In {
		// convert directory glob to file scope glob
		fileGlob := globPath.Value

		if !strings.HasSuffix(string(fileGlob), "/**") {
			// rules:
			// "app"            | match only app itself (directory), but no files inside
			// "app/*"          | match all files inside app itself, but no directory and subdirs (and subdirs files)
			// "app/**"         | match all files inside app and all subdirs (will recursive files on any level)

			// convert "app" -> "app/*", because we want to find files inside this directory
			fileGlob = arch.PathRelativeGlob(fmt.Sprintf("%s/*", fileGlob))
		}

		files, err := a.pathHelper.FindProjectFiles(arch.FileQuery{
			Path:               fileGlob,
			WorkingDirectory:   conf.WorkingDirectory.Value,
			Type:               arch.FileMatchQueryTypeOnlyFiles,
			ExcludeDirectories: conf.Exclude.RelativeDirectories.Values(),
			ExcludeRegexp:      conf.Exclude.RelativeFiles.Values(),
			Extensions:         []string{"go"},
		})

		if err != nil {
			return nil, fmt.Errorf("matching glob path failed '%v': %w", globPath.Value, err)
		}

		list = append(list, files...)
	}

	return list, nil
}

func (a *Assembler) extractUniquePackages(files []arch.FileDescriptor) []arch.FileDescriptor {
	packages := make([]arch.FileDescriptor, 0, len(files))
	unique := make(map[arch.PathRelative]any)

	for _, file := range files {
		packagePathRelative := arch.PathRelative(filepath.Dir(string(file.PathRel)))
		packagePathAbsolute := arch.PathAbsolute(filepath.Dir(string(file.PathAbs)))

		if _, ok := unique[packagePathRelative]; ok {
			continue
		}

		unique[packagePathRelative] = struct{}{}
		packages = append(packages, arch.FileDescriptor{
			PathRel: packagePathRelative,
			PathAbs: packagePathAbsolute,
			IsDir:   true,
		})
	}

	return packages
}
