package yaml

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/internal/models"
)

//nolint:funlen
func transformV4(tCtx TransformContext, doc ModelV4) models.Config {
	return models.Config{
		Version:          arch.NewRef(arch.ConfigVersion(doc.Version), tCtx.createReference("$.version")),
		WorkingDirectory: arch.NewRef(arch.PathRelative(doc.WorkingDirectory), tCtx.createReference("$.workingDirectory")),
		Settings: models.ConfigSettings{
			DeepScan: arch.NewInvalidRef(true),
			Imports: models.ConfigSettingsImports{
				StrictMode:            arch.NewRef(doc.Settings.Imports.StrictMode, tCtx.createReference("$.settings.imports.strictMode")),
				AllowAnyVendorImports: arch.NewRef(doc.Settings.Imports.AllowAnyVendorImports, tCtx.createReference("$.settings.imports.allowAnyVendorImports")),
			},
			Tags: transformV4SettingsTags(tCtx, doc.Settings.Tags),
		},
		Exclude: models.ConfigExclude{
			RelativeDirectories: sliceValuesAutoRef(tCtx, doc.Exclude.RelativeDirectories, "$.exclude.directories", func(dir string) arch.PathRelative {
				return arch.PathRelative(dir)
			}),
			RelativeFiles: sliceValuesAutoRef(tCtx, doc.Exclude.RelativeFiles, "$.exclude.files", func(file string) arch.PathRelativeRegExp {
				return arch.PathRelativeRegExp(file)
			}),
		},
		Components: models.ConfigComponents{
			Map: mapValuesAutoRef(tCtx, doc.Components, "$.components",
				func(tCtx TransformContext, name string, component ModelV4Component, refBasePath string) (arch.ComponentName, models.ConfigComponent) {
					return arch.ComponentName(name), models.ConfigComponent{
						In: sliceValuesAutoRef(tCtx, component.In, fmt.Sprintf("%s.in", refBasePath),
							func(value string) arch.PathRelativeGlob {
								return arch.PathRelativeGlob(value)
							}),
					}
				}),
		},
		Vendors: models.ConfigVendors{
			Map: mapValuesAutoRef(tCtx, doc.Vendors, "$.vendors",
				func(tCtx TransformContext, name string, vendor ModelV4Vendor, refBasePath string) (arch.VendorName, models.ConfigVendor) {
					return arch.VendorName(name), models.ConfigVendor{
						In: sliceValuesAutoRef(tCtx, vendor.In, fmt.Sprintf("%s.in", refBasePath),
							func(value string) arch.PathImportGlob {
								return arch.PathImportGlob(value)
							}),
					}
				}),
		},
		CommonComponents: sliceValuesAutoRef(tCtx, doc.CommonComponents, "$.commonComponents", func(v string) arch.ComponentName {
			return arch.ComponentName(v)
		}),
		CommonVendors: sliceValuesAutoRef(tCtx, doc.CommonVendors, "$.commonVendors", func(v string) arch.VendorName {
			return arch.VendorName(v)
		}),
		Dependencies: models.ConfigDependencies{
			Map: mapValuesAutoRef(tCtx, doc.Dependencies, "$.dependencies",
				func(tCtx TransformContext, cmpName string, deps ModelV4ComponentDependencies, refBasePath string) (arch.ComponentName, models.ConfigComponentDependencies) {
					return arch.ComponentName(cmpName), models.ConfigComponentDependencies{
						MayDependOn: sliceValuesAutoRef(tCtx, deps.MayDependOn, fmt.Sprintf("%s.mayDependOn", refBasePath),
							func(anotherCmpName string) arch.ComponentName {
								return arch.ComponentName(anotherCmpName)
							}),
						CanUse: sliceValuesAutoRef(tCtx, deps.CanUse, fmt.Sprintf("%s.canUse", refBasePath),
							func(vendorName string) arch.VendorName {
								return arch.VendorName(vendorName)
							}),
						AnyVendorDeps:  arch.NewRef(deps.AnyVendorDeps, tCtx.createReference(fmt.Sprintf("%s.anyVendorDeps", refBasePath))),
						AnyProjectDeps: arch.NewRef(deps.AnyProjectDeps, tCtx.createReference(fmt.Sprintf("%s.anyProjectDeps", refBasePath))),
						CanContainTags: sliceValuesAutoRef(tCtx, deps.CanContainTags, fmt.Sprintf("%s.canContainTags", refBasePath),
							func(tag string) arch.StructTag {
								return arch.StructTag(tag)
							}),
						DeepScan: models.ConfigOptional[arch.Ref[bool]]{
							Value:   arch.NewInvalidRef(false),
							Defined: false,
						},
					}
				}),
		},
	}
}

func transformV4SettingsTags(tCtx TransformContext, tags ModelV4SettingsTags) models.ConfigSettingsTags {
	refBasePath := "$.settings.structTags.allowed"
	ref := tCtx.createReference(refBasePath)

	if !tags.Allowed.defined {
		return models.ConfigSettingsTags{
			Allowed: arch.NewRef(arch.ConfigSettingsTagsEnumAll, ref),
		}
	}

	if len(tags.Allowed.value) == 0 {
		return models.ConfigSettingsTags{
			Allowed: arch.NewRef(arch.ConfigSettingsTagsEnumNone, ref),
		}
	}

	if tags.Allowed.value[0] == "true" {
		return models.ConfigSettingsTags{
			Allowed: arch.NewRef(arch.ConfigSettingsTagsEnumAll, ref),
		}
	}

	if tags.Allowed.value[0] == "false" {
		return models.ConfigSettingsTags{
			Allowed: arch.NewRef(arch.ConfigSettingsTagsEnumNone, ref),
		}
	}

	return models.ConfigSettingsTags{
		Allowed: arch.NewRef(arch.ConfigSettingsTagsEnumList, ref),
		AllowedList: sliceValuesAutoRef(tCtx, tags.Allowed.value, refBasePath,
			func(value string) arch.StructTag {
				return arch.StructTag(value)
			}),
	}
}
