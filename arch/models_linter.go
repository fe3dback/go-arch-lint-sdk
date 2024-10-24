package arch

const (
	LinterIDSyntax   LinterID = "syntax"
	LinterIDOrphans  LinterID = "orphans"
	LinterIDImports  LinterID = "imports"
	LinterIDDeepScan LinterID = "deepscan"
)

var LintersSortOrder = map[LinterID]int{
	LinterIDSyntax:   10,
	LinterIDOrphans:  20,
	LinterIDImports:  100,
	LinterIDDeepScan: 200,
}

var on = struct{}{}
var LintersWithPreview = map[LinterID]any{
	LinterIDSyntax:   on,
	LinterIDImports:  on,
	LinterIDDeepScan: on,
}

type (
	LinterID string

	Linter struct {
		ID          LinterID `json:"ID"`
		Used        bool     `json:"Used"`
		Name        string   `json:"-"`
		Description string   `json:"-"`
	}
)

type (
	LinterResult struct {
		Linter  Linter         `json:"Linter"`
		Notices []LinterNotice `json:"Notices"`
	}

	LinterNotice struct {
		Message   string              `json:"Message"`
		Reference Reference           `json:"Reference"`
		Details   LinterNoticeDetails `json:"Details"`
		Preview   string              `json:"-"`
	}

	LinterNoticeDetails struct {
		LinterID LinterID `json:"LinterID"`

		// exist when LinterID = arch.LinterIDOrphans
		LinterIDSyntax *LinterSyntaxDetails `json:"LinterIDSyntax,omitempty"`

		// exist when LinterID = arch.LinterIDOrphans
		LinterIDOrphan *LinterOrphanDetails `json:"LinterIDOrphan,omitempty"`

		// exist when LinterID = arch.LinterIDComponentImports
		LinterIDImports *LinterImportDetails `json:"LinterIDImports,omitempty"`

		// exist when LinterID = arch.LinterIDDeepScan
		LinterIDDeepscan *LinterDeepscanDetails `json:"LinterIDDeepscan,omitempty"`
	}

	LinterSyntaxDetails struct {
		ComponentName ComponentName `json:"ComponentName"`
		GoPackageName string        `json:"GoPackageName"`
		GoPackagePath PathRelative  `json:"GoPackagePath"`
		SyntaxError   string        `json:"SyntaxError"`
	}

	LinterOrphanDetails struct {
		FileRelativePath PathRelative `json:"FileRelativePath"`
		FileAbsolutePath PathAbsolute `json:"FileAbsolutePath"`
	}

	LinterImportDetails struct {
		ComponentName      ComponentName                 `json:"ComponentName"`
		TargetType         LinterImportDetailsTargetType `json:"TargetType"`
		TargetName         string                        `json:"TargetName"`    // Owner of ResolvedImportName (component or vendor)
		TargetDefined      bool                          `json:"TargetDefined"` // true if Target is known component or vendor in config
		FileRelativePath   PathRelative                  `json:"FileRelativePath"`
		FileAbsolutePath   PathAbsolute                  `json:"FileAbsolutePath"`
		ResolvedImportName PathImport                    `json:"ResolvedImportName"`
		Reference          Reference                     `json:"Reference"`
	}

	LinterImportDetailsTargetType string

	LinterDeepscanDetails struct {
		Gate       LinterDeepscanGate       `json:"Gate"`
		Dependency LinterDeepscanDependency `json:"Dependency"`
		Target     LinterDeepscanTarget     `json:"Target"`
	}

	LinterDeepscanGate struct {
		ComponentName ComponentName `json:"ComponentName"` // operations
		MethodName    string        `json:"MethodName"`    // NewOperation
		Definition    Reference     `json:"Definition"`    // internal/glue/code/line_count.go:54
		RelativePath  PathRelative  `json:"-"`             // internal/glue/code/line_count.go:54
	}

	LinterDeepscanDependency struct {
		ComponentName     ComponentName `json:"ComponentName"` // repository
		Name              string        `json:"Name"`          // micro.ViewRepository
		InjectionAST      string        `json:"InjectionAST"`  // c.provideMicroViewRepository()
		Injection         Reference     `json:"Injection"`     // internal/app/internal/container/cmd_mapping.go:15
		InjectionPath     PathRelative  `json:"-"`             // internal/app/internal/container/cmd_mapping.go:15
		SourceCodePreview []byte        `json:"-"`
	}

	LinterDeepscanTarget struct {
		Definition   Reference    `json:"Definition"`
		RelativePath PathRelative `json:"-"` // internal/app/internal/container/cmd_mapping.go:15
	}
)

const (
	LinterImportDetailsTargetTypeComponent LinterImportDetailsTargetType = "Component"
	LinterImportDetailsTargetTypeVendor    LinterImportDetailsTargetType = "Vendor"
)
