package models

import "github.com/fe3dback/go-arch-lint-sdk/arch"

type (
	Config struct {
		SyntaxProblems arch.RefSlice[string]

		Version          arch.Ref[arch.ConfigVersion]
		WorkingDirectory arch.Ref[arch.PathRelative]
		Settings         ConfigSettings
		Exclude          ConfigExclude
		Components       ConfigComponents
		Vendors          ConfigVendors
		CommonComponents arch.RefSlice[arch.ComponentName]
		CommonVendors    arch.RefSlice[arch.VendorName]
		Dependencies     ConfigDependencies
	}

	ConfigSettings struct {
		DeepScan arch.Ref[bool]
		Imports  ConfigSettingsImports
		Tags     ConfigSettingsTags
	}

	ConfigSettingsImports struct {
		StrictMode            arch.Ref[bool]
		AllowAnyVendorImports arch.Ref[bool]
	}

	ConfigSettingsTags struct {
		Allowed     arch.Ref[arch.ConfigSettingsTagsEnum]
		AllowedList arch.RefSlice[arch.StructTag]
	}

	ConfigExclude struct {
		RelativeDirectories arch.RefSlice[arch.PathRelative]
		RelativeFiles       arch.RefSlice[arch.PathRelativeRegExp]
	}

	ConfigComponents struct {
		Map arch.RefMap[arch.ComponentName, ConfigComponent]
	}

	ConfigComponent struct {
		In arch.RefSlice[arch.PathRelativeGlob]
	}

	ConfigVendors struct {
		Map arch.RefMap[arch.VendorName, ConfigVendor]
	}

	ConfigVendor struct {
		In arch.RefSlice[arch.PathImportGlob]
	}

	ConfigDependencies struct {
		Map arch.RefMap[arch.ComponentName, ConfigComponentDependencies]
	}

	ConfigComponentDependencies struct {
		MayDependOn    arch.RefSlice[arch.ComponentName]
		CanUse         arch.RefSlice[arch.VendorName]
		AnyVendorDeps  arch.Ref[bool]
		AnyProjectDeps arch.Ref[bool]
		CanContainTags arch.RefSlice[arch.StructTag]
		DeepScan       ConfigOptional[arch.Ref[bool]]
	}

	ConfigOptional[T any] struct {
		Value   T
		Defined bool
	}
)
