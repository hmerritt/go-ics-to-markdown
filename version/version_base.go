package version

var (
	// The git commit that was compiled. This will be filled in by the compiler.
	GitCommit   string
	GitBranch   string
	GitDescribe string

	// The compilation date. This will be filled in by the compiler.
	BuildDate string

	Version           = "0.5.10"
	VersionPrerelease = ""
	VersionMetadata   = ""
)
