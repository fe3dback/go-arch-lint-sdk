package arch

// todo: spec checking
// - pick file
// - find parent component
// - can check imports already
// - can check strictMode already
// - todo: some cache system (3 files: gate, injector, dependency)??

type (
	// Spec fully describes project configuration and all linter rules
	// is primary (main) structure for working with in all other code
	// Spec is indivisible struct, and it parts shouldn't be used separately
	Spec struct {
		Project          ProjectInfo
		WorkingDirectory Ref[PathRelative]
		Components       []SpecComponent
		Orphans          []SpecOrphan
	}

	// StructTag used in go source code for struct annotations (ex: "json", "db")
	StructTag string

	// ComponentName unique user-specified name of the component.
	// Component is alias for N go packages
	ComponentName string

	// VendorName unique user-specified name of the vendor library.
	// Vendor is alias for N specific imports (ex: [github.com/hello/world, go.example.com/my/package])
	VendorName string

	// SpecComponent fully describe one project component and it rules.
	// Component is alias for N go packages
	SpecComponent struct {
		Name                Ref[ComponentName]
		DefinitionComponent Reference // $.components.<NAME>
		DefinitionDeps      Reference // $.deps.<NAME>
		DeepScan            Ref[bool]
		StrictMode          Ref[bool]
		AllowAllProjectDeps Ref[bool]
		AllowAllVendorDeps  Ref[bool]
		AllowAllTags        Ref[bool]
		AllowedTags         RefSlice[StructTag]
		MayDependOn         RefSlice[ComponentName]
		CanUse              RefSlice[VendorName]
		MatchPatterns       RefSlice[PathRelativeGlob] // $.components.X.in
		MatchedFiles        []FileDescriptor           // all files matched by component "in" query
		MatchedPackages     []FileDescriptor           // all packages matched by component "in" query
		OwnedFiles          []FileDescriptor           // unique subset of MatchedFiles, belongs to this component (every file will belong only to single component)
		OwnedPackages       []FileDescriptor           // unique subset of MatchedPackages, belongs to this component (every package will belong only to single component)
	}

	// SpecOrphan describes project file that is not mapped to any component
	SpecOrphan struct {
		File FileDescriptor
	}
)
