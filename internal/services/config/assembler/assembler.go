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

func (a *Assembler) Assemble(conf models.Config) (spec arch.Spec, err error) {
	projectInfo, err := a.projectInfoFetcher.Fetch()
	if err != nil {
		return arch.Spec{}, fmt.Errorf("failed fetch project info: %w", err)
	}

	components, err := a.assembleComponents(&conf, &projectInfo)
	if err != nil {
		return arch.Spec{}, fmt.Errorf("failed assemble components: %w", err)
	}

	vendors, err := a.assembleVendors(&conf)
	if err != nil {
		return arch.Spec{}, fmt.Errorf("failed assemble vendors: %w", err)
	}

	orphans, err := a.assembleOrphans(&conf, components)
	if err != nil {
		return spec, fmt.Errorf("failed assembling orphans: %w", err)
	}

	return arch.Spec{
		Project:          projectInfo,
		WorkingDirectory: conf.WorkingDirectory,
		Components:       components,
		Vendors:          vendors,
		Orphans:          orphans,
	}, nil
}

func (a *Assembler) assembleComponents(conf *models.Config, project *arch.ProjectInfo) (result arch.Components, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed assemble components: %v\n%s", r, debug.Stack())
		}
	}()

	components := make(arch.Components, conf.Components.Map.Len())

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

		tagsAllowedAll, tagsAllowedWhiteList := a.figureOutAllowedStructTags(conf, &rules)

		matchedFiles, err := a.findOwnedFiles(conf, definition)
		if err != nil {
			panic(fmt.Errorf("failed finding owned files by component '%s': %w", name, err))
		}

		components[name] = &arch.SpecComponent{
			Name:                arch.NewRef(name, definitionRef),
			DefinitionComponent: definitionRef,
			DefinitionDeps:      rulesRef,
			DeepScan:            deepScan,
			StrictMode:          conf.Settings.Imports.StrictMode,
			AllowAllProjectDeps: rules.AnyProjectDeps,
			AllowAllVendorDeps:  rules.AnyVendorDeps,
			AllowAllTags:        tagsAllowedAll,
			AllowedTags:         tagsAllowedWhiteList,
			MayDependOn:         append(conf.CommonComponents, rules.MayDependOn...),
			CanUse:              append(conf.CommonVendors, rules.CanUse...),
			MatchPatterns:       definition.In,
			MatchedFiles:        matchedFiles,
		}
	})

	// copy matched files to owned files (but each file owned only by one component)
	a.calculateFilesOwnage(components)

	// find matched/owned packages from files
	for _, component := range components {
		component.MatchedPackages = a.extractUniquePackages(project, component.MatchedFiles)
		component.OwnedPackages = a.extractUniquePackages(project, component.OwnedFiles)
	}

	// sort paths
	providerPathFn := func(p *arch.PathDescriptor) (relPath arch.PathRelative, isDir bool) {
		return p.PathRel, p.IsDir
	}
	providerPackageFn := func(p *arch.PackageDescriptor) (relPath arch.PathRelative, isDir bool) {
		return p.PathRel, p.IsDir
	}

	for _, component := range components {
		pathsort.SortFileTree(component.MatchedFiles, providerPathFn)
		pathsort.SortFileTree(component.MatchedPackages, providerPackageFn)
		pathsort.SortFileTree(component.OwnedFiles, providerPathFn)
		pathsort.SortFileTree(component.OwnedPackages, providerPackageFn)
	}

	return components, nil
}

func (a *Assembler) assembleVendors(conf *models.Config) (arch.Vendors, error) {
	vendors := make(arch.Vendors, conf.Vendors.Map.Len())

	conf.Vendors.Map.Each(func(name arch.VendorName, vendor models.ConfigVendor, definitionRef arch.Reference) {
		cleanedIn := make(arch.RefSlice[arch.PathImportGlob], 0, len(vendor.In))
		for _, path := range vendor.In {
			path.Value = arch.PathImportGlob(strings.TrimRight(string(path.Value), "/"))
			cleanedIn = append(cleanedIn, path)
		}

		vendors[name] = &arch.SpecVendor{
			Name:         arch.NewRef(name, definitionRef),
			Definition:   definitionRef,
			OwnedImports: cleanedIn,
		}
	})

	return vendors, nil
}

func (a *Assembler) assembleOrphans(conf *models.Config, components arch.Components) ([]arch.SpecOrphan, error) {
	mappedFiles := make(map[arch.PathRelative]any, 128)
	for _, component := range components {
		for _, dst := range component.MatchedFiles {
			mappedFiles[dst.PathRel] = struct{}{}
		}
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
		return nil, fmt.Errorf("failed finding project files: %w", err)
	}

	list := make([]arch.SpecOrphan, 0, 32)
	for _, projectFile := range files {
		if _, mapped := mappedFiles[projectFile.PathRel]; mapped {
			continue
		}

		list = append(list, arch.SpecOrphan{
			File: projectFile,
		})
	}

	pathsort.SortFileTree(list, func(a *arch.SpecOrphan) (relPath arch.PathRelative, isDir bool) {
		return a.File.PathRel, false
	})

	return list, nil
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

func (a *Assembler) findOwnedFiles(conf *models.Config, component models.ConfigComponent) ([]arch.PathDescriptor, error) {
	list := make([]arch.PathDescriptor, 0, 32)

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

func (a *Assembler) extractUniquePackages(project *arch.ProjectInfo, files []arch.PathDescriptor) []arch.PackageDescriptor {
	packages := make([]arch.PackageDescriptor, 0, len(files))
	unique := make(map[arch.PathRelative]any)

	for _, file := range files {
		packagePathRelative := arch.PathRelative(filepath.Dir(string(file.PathRel)))
		packagePathAbsolute := arch.PathAbsolute(filepath.Dir(string(file.PathAbs)))

		if _, ok := unique[packagePathRelative]; ok {
			continue
		}

		unique[packagePathRelative] = struct{}{}

		pathDescriptor := arch.PathDescriptor{
			PathRel: packagePathRelative,
			PathAbs: packagePathAbsolute,
			IsDir:   true,
		}
		packages = append(packages, arch.PackageDescriptor{
			PathDescriptor: pathDescriptor,
			Import:         a.transformToImportPath(project, pathDescriptor),
		})
	}

	return packages
}

func (a *Assembler) transformToImportPath(project *arch.ProjectInfo, ownedPackage arch.PathDescriptor) arch.PathImport {
	return arch.PathImport(string(project.Module) + "/" + string(ownedPackage.PathRel))
}
