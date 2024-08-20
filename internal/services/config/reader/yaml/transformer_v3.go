package yaml

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/internal/models"
)

//nolint:funlen
func transformV3(tCtx TransformContext, doc ModelV3) models.Config {
	return models.Config{
		Version:          arch.NewRef(arch.ConfigVersion(doc.Version), tCtx.createReference("$.version")),
		WorkingDirectory: arch.NewRef(arch.PathRelative(doc.WorkingDirectory), tCtx.createReference("$.workdir")),
		Settings: models.ConfigSettings{
			DeepScan: arch.NewRef(doc.Allow.DeepScan, tCtx.createReference("$.allow.deepScan")),
			Imports: models.ConfigSettingsImports{
				StrictMode:            arch.NewInvalidRef(false),
				AllowAnyVendorImports: arch.NewRef(doc.Allow.DepOnAnyVendor, tCtx.createReference("$.allow.depOnAnyVendor")),
			},
			Tags: models.ConfigSettingsTags{
				Allowed: arch.NewInvalidRef(arch.ConfigSettingsTagsEnumAll),
			},
		},
		Exclude: models.ConfigExclude{
			RelativeDirectories: sliceValuesAutoRef(tCtx, doc.ExcludeDirectories, "$.exclude", func(dir string) arch.PathRelative {
				return arch.PathRelative(dir)
			}),
			RelativeFiles: sliceValuesAutoRef(tCtx, doc.ExcludeFiles, "$.excludeFiles", func(file string) arch.PathRelativeRegExp {
				return arch.PathRelativeRegExp(file)
			}),
		},
		Components: models.ConfigComponents{
			Map: mapValuesAutoRef(tCtx, doc.Components, "$.components",
				func(tCtx TransformContext, name string, component ModelV3Component, refBasePath string) (arch.ComponentName, models.ConfigComponent) {
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
				func(tCtx TransformContext, name string, vendor ModelV3Vendor, refBasePath string) (arch.VendorName, models.ConfigVendor) {
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
			Map: mapValuesAutoRef(tCtx, doc.Dependencies, "$.deps",
				func(tCtx TransformContext, cmpName string, deps ModelV3ComponentDependencies, refBasePath string) (arch.ComponentName, models.ConfigComponentDependencies) {
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
						CanContainTags: []arch.Ref[arch.StructTag]{},
						DeepScan: models.ConfigOptional[arch.Ref[bool]]{
							Value:   arch.NewRef(deps.DeepScan.value, tCtx.createReference(fmt.Sprintf("%s.deepScan", refBasePath))),
							Defined: deps.DeepScan.defined,
						},
					}
				}),
		},
	}
}
