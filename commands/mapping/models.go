package mapping

const (
	SchemeGrouped Scheme = "grouped"
	SchemeList    Scheme = "list"
)

var Schemes = []string{
	SchemeList,
	SchemeGrouped,
}

type (
	Scheme = string

	In struct {
		Scheme Scheme
	}

	Out struct {
		ProjectDirectory string       `json:"ProjectDirectory"`
		ModuleName       string       `json:"ModuleName"`
		MappingGrouped   []OutGrouped `json:"MappingGrouped"`
		MappingList      []OutList    `json:"MappingList"`
		Scheme           Scheme       `json:"-"`
	}

	OutGrouped struct {
		ComponentName string   `json:"ComponentName"`
		Packages      []string `json:"Packages"`
	}

	OutList struct {
		Package       string `json:"Package"`
		ComponentName string `json:"ComponentName"`
	}
)
