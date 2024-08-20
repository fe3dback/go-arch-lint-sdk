package validator

import (
	"github.com/fe3dback/go-arch-lint-sdk/arch"
	"github.com/fe3dback/go-arch-lint-sdk/internal/models"
)

const (
	cmpContainer  = arch.ComponentName("container")
	cmpHandler    = arch.ComponentName("handler")
	cmpUsecase    = arch.ComponentName("usecase")
	cmpRepository = arch.ComponentName("repository")
	cmpModels     = arch.ComponentName("models")

	vendorDB      = arch.VendorName("db")
	vendorTracing = arch.VendorName("tracing")
)

type validatorIn struct {
	conf models.Config
}

// shortcut for models.NewInvalidRef
// for better code readability
func nf[T any](value T) arch.Ref[T] {
	return arch.NewInvalidRef(value)
}

//nolint:funlen
func createValidatorIn(mutators ...func(*validatorIn)) validatorIn {
	in := validatorIn{
		conf: models.Config{
			Version:          nf(arch.ConfigVersion(4)),
			WorkingDirectory: nf(arch.PathRelative("internal")),
			Settings: models.ConfigSettings{
				DeepScan: nf(false),
				Imports: models.ConfigSettingsImports{
					StrictMode:            nf(false),
					AllowAnyVendorImports: nf(false),
				},
				Tags: models.ConfigSettingsTags{
					Allowed: nf(arch.ConfigSettingsTagsEnumAll),
				},
			},
			Components: models.ConfigComponents{
				Map: arch.NewRefMapFrom(map[arch.ComponentName]arch.Ref[models.ConfigComponent]{
					cmpContainer: nf(models.ConfigComponent{
						In: []arch.Ref[arch.PathRelativeGlob]{
							nf(arch.PathRelativeGlob("app/internal/container")),
							nf(arch.PathRelativeGlob("plugins/*/container")),
						},
					}),
					cmpHandler: nf(models.ConfigComponent{
						In: []arch.Ref[arch.PathRelativeGlob]{
							nf(arch.PathRelativeGlob("handlers")),
						},
					}),
					cmpUsecase: nf(models.ConfigComponent{
						In: []arch.Ref[arch.PathRelativeGlob]{
							nf(arch.PathRelativeGlob("app/*/business/**")),
						},
					}),
					cmpRepository: nf(models.ConfigComponent{
						In: []arch.Ref[arch.PathRelativeGlob]{
							nf(arch.PathRelativeGlob("app/*/repo")),
						},
					}),
					cmpModels: nf(models.ConfigComponent{
						In: []arch.Ref[arch.PathRelativeGlob]{
							nf(arch.PathRelativeGlob("models/**")),
						},
					}),
				}),
			},
			Vendors: models.ConfigVendors{
				Map: arch.NewRefMapFrom(map[arch.VendorName]arch.Ref[models.ConfigVendor]{
					vendorDB: nf(models.ConfigVendor{
						In: []arch.Ref[arch.PathImportGlob]{
							nf(arch.PathImportGlob("github.com/fe3dback/orm")),
							nf(arch.PathImportGlob("github.com/fe3dback/libs/**/db/*")),
							nf(arch.PathImportGlob("github.com/fe3dback/transactional")),
						},
					}),
					vendorTracing: nf(models.ConfigVendor{
						In: []arch.Ref[arch.PathImportGlob]{
							nf(arch.PathImportGlob("io.org.example.com/telemetry/*/tracing")),
						},
					}),
				}),
			},
			CommonComponents: []arch.Ref[arch.ComponentName]{
				nf(cmpModels),
			},
			CommonVendors: []arch.Ref[arch.VendorName]{
				nf(vendorTracing),
			},
			Dependencies: models.ConfigDependencies{
				Map: arch.NewRefMapFrom(map[arch.ComponentName]arch.Ref[models.ConfigComponentDependencies]{
					cmpContainer: nf(models.ConfigComponentDependencies{
						AnyVendorDeps:  nf(true),
						AnyProjectDeps: nf(true),
					}),
					cmpHandler: nf(models.ConfigComponentDependencies{
						MayDependOn: []arch.Ref[arch.ComponentName]{
							nf(cmpUsecase),
						},
					}),
					cmpUsecase: nf(models.ConfigComponentDependencies{
						MayDependOn: []arch.Ref[arch.ComponentName]{
							nf(cmpUsecase),
							nf(cmpRepository),
						},
					}),
					cmpRepository: nf(models.ConfigComponentDependencies{
						CanUse: []arch.Ref[arch.VendorName]{
							nf(vendorDB),
						},
						CanContainTags: []arch.Ref[arch.StructTag]{
							nf(arch.StructTag("db")),
						},
					}),
				}),
			},
		},
	}

	for _, mutate := range mutators {
		mutate(&in)
	}

	return in
}
