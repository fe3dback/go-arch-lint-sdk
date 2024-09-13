package arch

const (
	ConfigSettingsTagsEnumAll  ConfigSettingsTagsEnum = "All"
	ConfigSettingsTagsEnumNone ConfigSettingsTagsEnum = "None"
	ConfigSettingsTagsEnumList ConfigSettingsTagsEnum = "List"
)

const (
	ConfigVersion4 = ConfigVersion(4)
)

type (
	// ConfigVersion primary arch-lint config version
	ConfigVersion int

	ConfigSettingsTagsEnum string
)
