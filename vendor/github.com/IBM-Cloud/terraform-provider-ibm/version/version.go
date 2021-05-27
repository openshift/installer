package version

import (
	"github.com/hashicorp/go-version"
)

// Version is the current provider main version
const Version = "1.25.0"

// GitCommit is the git commit that was compiled. This will be filled in by the compiler.
var GitCommit string

//VersionPrerelease is the marker for version. If this is "" (empty string)
// then it means that it is a final release. Otherwise, this is a pre-release
// such as "dev" (in development), "beta", "rc1", etc.
var VersionPrerelease = ""

// SemVersion is an instance of version.Version. This has the secondary
// benefit of verifying during tests and init time that our version is a
// proper semantic version, which should always be the case.
var SemVersion = version.Must(version.NewVersion(Version))
