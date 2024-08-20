package arch

type (
	// ProjectInfo contains basic info about specific GO project
	ProjectInfo struct {
		Directory PathAbsolute
		Module    GoModule
	}
)
