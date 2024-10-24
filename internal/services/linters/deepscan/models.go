package deepscan

type (
	Method struct {
		Name  string
		Gates []Gate
	}

	Gate struct {
		IsInterfaceType bool
	}
)
