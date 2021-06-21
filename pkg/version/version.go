// Package version includes the version information for installer.
package version

import (
	"fmt"
	"strings"

	"github.com/openshift/installer/pkg/types"
)

// This file handles correctly identifying the default release version, which is expected to be
// replaced in the binary post-compile by the release name extracted from a payload. The expected modification is:
//
// 1. Extract a release binary from the installer image referenced within the release image
// 2. Identify the release name, add a NUL terminator byte (0x00) to the end, calculate length
// 3. Length must be less than the marker length
// 4. Search through the installer binary looking for `\x00_RELEASE_VERSION_LOCATION_\x00<PADDING_TO_LENGTH>`
//    where padding is the ASCII character X and length is the total length of the marker
// 5. Overwrite the beginning of the marker with the release name and a NUL terminator byte (0x00)

var (
	// Raw is the string representation of the version. This will be replaced
	// with the calculated version at build time.
	// set in hack/build.sh
	Raw = "was not built correctly"

	// Commit is the commit hash from which the installer was built.
	// Set in hack/build.sh.
	Commit = ""

	// defaultArch is the payload architecture for which the installer was built,
	// which even on Linux may not be the same as the architecture of the
	// installer binary itself.
	// Set in hack/build.sh.
	defaultArch = "amd64"

	// defaultReleaseInfoPadded may be replaced in the binary with Release Metadata: Version that overrides defaultVersion as
	// a null-terminated string within the allowed character length. This allows a distributor to override the payload
	// location without having to rebuild the source.
	defaultVersionPadded = "\x00_RELEASE_VERSION_LOCATION_\x00XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX\x00"
	defaultVersionPrefix = "\x00_RELEASE_VERSION_LOCATION_\x00"
	defaultVersionLength = len(defaultVersionPadded)
)

// String returns the human-friendly representation of the version.
func String() (string, error) {
	version, err := Version()
	return fmt.Sprintf("OpenShift Installer %s", version), err
}

// Version returns the installer/release version.
func Version() (string, error) {
	if strings.HasPrefix(defaultVersionPadded, defaultVersionPrefix) {
		return Raw, nil
	}
	nullTerminator := strings.IndexByte(defaultVersionPadded, '\x00')
	if nullTerminator == -1 {
		// the binary has been altered, but we didn't find a null terminator within the release name constant which is an error
		return Raw, fmt.Errorf("release name location was replaced but without a null terminator before %d bytes", defaultVersionLength)
	}
	if nullTerminator > len(defaultVersionPadded) {
		// the binary has been altered, but the null terminator is *longer* than the constant encoded in the binary
		return Raw, fmt.Errorf("release name location contains no null-terminator and constant is corrupted")
	}
	releaseName := defaultVersionPadded[:nullTerminator]
	if len(releaseName) == 0 {
		// the binary has been altered, but the replaced release name is empty which is incorrect
		// the oc binary will not be pinned to Release Metadata:Version
		return Raw, fmt.Errorf("release name was incorrectly replaced during extract")
	}
	return releaseName, nil
}

// DefaultArch returns the default release architecture
func DefaultArch() types.Architecture {
	return types.Architecture(defaultArch)
}
