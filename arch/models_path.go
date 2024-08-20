package arch

type (
	// PathRelative some relative to project directory path (ex: "internal/app.go")
	PathRelative string

	// PathRelativeRegExp some relative to project directory path contains regexp (ex: "^.*_test\\.go$")
	PathRelativeRegExp string

	// PathRelativeGlob some relative to project directory path with globs  (ex: "internal/*/repo/**/db.go")
	PathRelativeGlob string

	// PathAbsolute some absolute path (ex: "/home/user/admin/go/project/internal/app.go")
	PathAbsolute string

	// PathAbsoluteRegExp some absolute path contains regexp (ex: "^.*user\/admin*_\\.go$")
	PathAbsoluteRegExp string

	// PathAbsoluteGlob some absolute path with globs  (ex: "/home/user/admin/go/project/internal/*/repo/**/db.go")
	PathAbsoluteGlob string

	// PathImport is absolute golang import string (ex: "github.com/goccy/go-yaml")
	PathImport string

	// PathImportGlob is absolute golang import string with globs (ex: "oss.terrastruct.com/d2/*/libs/**/something")
	PathImportGlob string
)
