package container

import (
	"github.com/fe3dback/go-arch-lint-sdk/internal/services/config/assembler"
	"github.com/fe3dback/go-arch-lint-sdk/internal/services/config/reader"
	"github.com/fe3dback/go-arch-lint-sdk/internal/services/config/reader/yaml"
	"github.com/fe3dback/go-arch-lint-sdk/internal/services/config/validator"
)

func (c *Container) serviceConfigReader() *reader.Reader {
	return once(func() *reader.Reader {
		return reader.NewReader(
			c.serviceConfigReaderYAML(),
		)
	})
}

func (c *Container) serviceConfigReaderYAML() *yaml.Reader {
	return once(yaml.NewReader)
}

func (c *Container) serviceConfigValidator() *validator.Root {
	return once(func() *validator.Root {
		return validator.NewRoot(
			c.usedContext,
			c.skipMissUsages,
			c.serviceCodePrinter(),
			c.serviceConfigValidatorWorkdir(),
			c.serviceConfigValidatorExcludedFiles(),
			c.serviceConfigValidatorCmnComponents(),
			c.serviceConfigValidatorCmnVendors(),
			c.serviceConfigValidatorComponents(),
			c.serviceConfigValidatorVendors(),
			c.serviceConfigValidatorDeps(),
			c.serviceConfigValidatorDepsComponents(),
			c.serviceConfigValidatorDepsVendors(),
			c.serviceConfigValidatorCommonCollisionMissuse(),
			c.serviceConfigValidatorVendorsMissuse(),
			c.serviceConfigValidatorTagsMissuse(),
		)
	})
}

func (c *Container) serviceConfigValidatorWorkdir() *validator.WorkdirValidator {
	return once(func() *validator.WorkdirValidator {
		return validator.NewWorkdirValidator(
			c.serviceProjectPathHelper(),
		)
	})
}

func (c *Container) serviceConfigValidatorExcludedFiles() *validator.ExcludedFilesValidator {
	return once(validator.NewExcludedFilesValidator)
}

func (c *Container) serviceConfigValidatorCmnComponents() *validator.CommonComponentsValidator {
	return once(validator.NewCommonComponentsValidator)
}

func (c *Container) serviceConfigValidatorCmnVendors() *validator.CommonVendorsValidator {
	return once(validator.NewCommonVendorsValidator)
}

func (c *Container) serviceConfigValidatorComponents() *validator.ComponentsValidator {
	return once(func() *validator.ComponentsValidator {
		return validator.NewComponentsValidator(
			c.serviceProjectPathHelper(),
		)
	})
}

func (c *Container) serviceConfigValidatorVendors() *validator.VendorsValidator {
	return once(func() *validator.VendorsValidator {
		return validator.NewVendorsValidator(
			c.serviceProjectPathHelper(),
		)
	})
}

func (c *Container) serviceConfigValidatorDeps() *validator.DepsValidator {
	return once(validator.NewDepsValidator)
}

func (c *Container) serviceConfigValidatorDepsComponents() *validator.DepsComponentsValidator {
	return once(validator.NewDepsComponentsValidator)
}

func (c *Container) serviceConfigValidatorDepsVendors() *validator.DepsVendorsValidator {
	return once(validator.NewDepsVendorsValidator)
}

func (c *Container) serviceConfigValidatorVendorsMissuse() *validator.VendorsMissUseValidator {
	return once(validator.NewVendorsMissUseValidator)
}

func (c *Container) serviceConfigValidatorCommonCollisionMissuse() *validator.VendorsCommonCollisionMissUseValidator {
	return once(validator.NewCommonCollisionMissUseValidator)
}

func (c *Container) serviceConfigValidatorTagsMissuse() *validator.TagsMissUseValidator {
	return once(validator.NewTagsMissUseValidator)
}

func (c *Container) serviceConfigAssembler() *assembler.Assembler {
	return once(func() *assembler.Assembler {
		return assembler.NewAssembler(
			c.serviceProjectFetcher(),
			c.serviceProjectPathHelper(),
		)
	})
}
