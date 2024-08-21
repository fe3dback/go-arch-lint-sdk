package arch

// FileDescriptor hold meta information for specific file path
// FileDescriptor is not synced with OS anyway and can contain old (irrelevant) information
type FileDescriptor struct {
	PathRel   PathRelative // relative to (projectDirectory + workingDirectory)
	PathAbs   PathAbsolute
	IsDir     bool
	Extension string // in lowercase
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