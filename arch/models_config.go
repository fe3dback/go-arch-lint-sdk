package arch

const (
	ConfigSettingsTagsEnumAll  ConfigSettingsTagsEnum = "All"
	ConfigSettingsTagsEnumNone ConfigSettingsTagsEnum = "None"
	ConfigSettingsTagsEnumList ConfigSettingsTagsEnum = "List"
)

type (
	// ConfigVersion primary arch-lint config version
	ConfigVersion int

	ConfigSettingsTagsEnum string
)
