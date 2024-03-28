// Package version includes the version information for installer.
package version

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

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

	// releaseArchitecture is the architecture of the release payload: multi, amd64, arm64, ppc64le, or s390x.
	// we don't know the releaseArchitecure by default "".
	releaseArchitecture = ""

	// defaultReleaseInfoPadded may be replaced in the binary with Release Metadata: Version that overrides defaultVersion as
	// a null-terminated string within the allowed character length. This allows a distributor to override the payload
	// location without having to rebuild the source.
	defaultVersionPadded = "\x00_RELEASE_VERSION_LOCATION_\x00XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX\x00"
	defaultVersionPrefix = "\x00_RELEASE_VERSION_LOCATION_\x00"
	defaultVersionLength = len(defaultVersionPadded)

	// releaseArchitecturesPadded may be replaced in the binary with Release Image Architecture(s): RELEASE_ARCHITECTURE that overrides releaseArchitecture as
	// a null-terminated string within the allowed character length. This allows a distributor to override the payload
	// location without having to rebuild the source.
	releaseArchitecturesPadded = "\x00_RELEASE_ARCHITECTURE_LOCATION_\x00XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX\x00"
	releaseArchitecturesPrefix = "\x00_RELEASE_ARCHITECTURE_LOCATION_\x00"
	releaseArchitecturesLength = len(releaseArchitecturesPadded)
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

// ReleaseArchitecture returns the release image cpu architecture version.
func ReleaseArchitecture() (string, error) {
	ri, okRI := os.LookupEnv("OPENSHIFT_INSTALL_RELEASE_IMAGE_OVERRIDE")
	if okRI {
		logrus.Warnf("Found override for release image (%s). Release Image Architecture is unknown", ri)
		return "unknown", nil
	}
	if strings.HasPrefix(releaseArchitecturesPadded, releaseArchitecturesPrefix) {
		logrus.Warn("Release Image Architecture not detected. Release Image Architecture is unknown")
		return "unknown", nil
	}
	nullTerminator := strings.IndexByte(releaseArchitecturesPadded, '\x00')
	if nullTerminator == -1 {
		// the binary has been altered, but we didn't find a null terminator within the release architecture constant which is an error
		return Raw, fmt.Errorf("release architecture location was replaced but without a null terminator before %d bytes", releaseArchitecturesLength)
	}
	if nullTerminator > len(releaseArchitecturesPadded) {
		// the binary has been altered, but the null terminator is *longer* than the constant encoded in the binary
		return Raw, fmt.Errorf("release architecture location contains no null-terminator and constant is corrupted")
	}
	releaseArchitecture = releaseArchitecturesPadded[:nullTerminator]
	if len(releaseArchitecture) == 0 {
		// the binary has been altered, but the replaced release architecture is empty which is incorrect
		return Raw, fmt.Errorf("release architecture was incorrectly replaced during extract")
	}
	return cleanArch(releaseArchitecture), nil
}

// cleanArch oc will embed linux/<arch> or multi (linux/<arch>) we want to clean this up so validation can more cleanly use this method.
// multi (linux/amd64) -> multi
// linux/amd64 -> amd64
// linux/<arch> -> <arch>.
func cleanArch(releaseArchitecture string) string {
	if strings.HasPrefix(releaseArchitecture, "multi") {
		return "multi"
	}
	// remove 'linux/', we just want <arch>
	return strings.ReplaceAll(releaseArchitecture, "linux/", "")
}

// DefaultArch returns the default release architecture
func DefaultArch() types.Architecture {
	return types.Architecture(defaultArch)
}
