// Package version includes the version information for installer.
package version

import "fmt"

var (
	// Raw is the string representation of the version. This will be replaced
	// with the calculated version at build time.
	// set in hack/build.sh
	Raw = "was not built correctly"

	// String is the human-friendly representation of the version.
	String = fmt.Sprintf("OpenShift Installer %s", Raw)

	// Commit is the commit hash from which the installer was built.
	// Set in hack/build.sh.
	Commit = ""
)
