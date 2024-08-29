package arch

// PathDescriptor hold meta information for specific file path
// PathDescriptor is not synced with OS anyway and can contain old (irrelevant) information
type PathDescriptor struct {
	PathRel   PathRelative // relative to (projectDirectory + workingDirectory)
	PathAbs   PathAbsolute
	IsDir     bool
	Extension string // in lowercase
}

// PackageDescriptor is extension for PathDescriptor but can describe
// only go packages. Also, this struct contain import path for this go package
type PackageDescriptor struct {
	PathDescriptor
	Import PathImport
}

type FileMatchQueryType string

const (
	FileMatchQueryTypeAll             FileMatchQueryType = "all"
	FileMatchQueryTypeOnlyFiles       FileMatchQueryType = "files"
	FileMatchQueryTypeOnlyDirectories FileMatchQueryType = "directories"
)

type FileQuery struct {
	Path               any          // support models.PathXXX types
	WorkingDirectory   PathRelative // fill be prepended to Path
	Type               FileMatchQueryType
	ExcludeDirectories []PathRelative
	ExcludeFiles       []PathRelative
	ExcludeRegexp      []PathRelativeRegExp
	Extensions         []string // without dot, example: [js, go, jpg]. Nil = no filter
}
