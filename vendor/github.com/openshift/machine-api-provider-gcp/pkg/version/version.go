package version

import (
	"fmt"
	"strings"

	"github.com/blang/semver"
)

var (
	// Raw is the string representation of the version. This will be replaced
	// with the calculated version at build time.
	Raw = "v0.0.0-was-not-built-properly"

	// Version is semver representation of the version.
	Version = semver.MustParse(strings.TrimLeft(Raw, "v"))

	// String is the human-friendly representation of the version.
	String = fmt.Sprintf("ClusterAPIProviderGCP %s", Raw)
)
