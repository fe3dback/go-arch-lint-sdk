package arch

const (
	LinterIDOrphans          LinterID = "orphans"
	LinterIDComponentImports LinterID = "component_imports"
	LinterIDVendorImports    LinterID = "vendor_imports"
	LinterIDDeepScan         LinterID = "deepscan"
)

type (
	LinterID string

	Linter struct {
		ID                  LinterID `json:"ID"`
		Used                bool     `json:"Used"`
		Name                string   `json:"-"`
		Description         string   `json:"-"`
		EnableConditionHint string   `json:"-"`
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
	}

	LinterNoticeDetails struct {
		LinterID LinterID `json:"LinterID"`

		// exist when LinterID = arch.LinterIDOrphans
		LinterIDOrphan *LinterOrphanDetails `json:"LinterIDOrphan,omitempty"`

		// exist when LinterID = arch.LinterIDComponentImports
		LinterIDComponentImports *LinterImportDetails `json:"LinterIDComponentImports,omitempty"`

		// exist when LinterID = arch.LinterIDVendorImports
		LinterIDVendorImports *LinterImportDetails `json:"LinterIDVendorImports,omitempty"`

		// exist when LinterID = arch.LinterIDDeepScan
		LinterIDDeepscan *LinterDeepscanDetails `json:"LinterIDDeepscan,omitempty"`
	}

	LinterOrphanDetails struct {
		FileRelativePath PathRelative `json:"FileRelativePath"`
		FileAbsolutePath PathAbsolute `json:"FileAbsolutePath"`
	}

	LinterImportDetails struct {
		ComponentName      ComponentName `json:"ComponentName"`
		FileRelativePath   PathRelative  `json:"FileRelativePath"`
		FileAbsolutePath   PathAbsolute  `json:"FileAbsolutePath"`
		ResolvedImportName PathImport    `json:"ResolvedImportName"`
		Reference          Reference     `json:"Reference"`
	}

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
