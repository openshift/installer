package releaseimage

import (
	"fmt"
	"strings"
)

// This file handles correctly identifying the default release image, which is expected to be
// replaced in the binary post-compile after being extracted from a payload. The expected modification is:
//
// 1. Extract a release binary from the installer image referenced within the release image
// 2. Identify the release image pull spec, add a NUL terminator byte (0x00) to the end, calculate length
// 3. Length must be less than 300 bytes
// 4. Search through the installer binary looking for `\x00_RELEASE_IMAGE_LOCATION_\x00<PADDING_TO_LENGTH>`
//    where padding is the ASCII character X and length is the total length of the image
// 5. Overwrite that chunk of the bytes if found, otherwise return error.
//
// On start the installer examines the constant and if it has been modified from the default the installer
// will use that image.

var (
	// defaultReleaseImageOriginal is the value served when defaultReleaseInfoPadded is unmodified.
	defaultReleaseImageOriginal = "registry.svc.ci.openshift.org/origin/release:4.2"
	// defaultReleaseNameOriginal is the value served when defaultReleaseInfoPadded is unmodified.
	defaultReleaseNameOriginal = "4.2"
	// defaultReleaseInfoPadded may be replaced in the binary with a pull spec that overrides defaultReleaseInfo as
	// a null-terminated string within the allowed character length. This allows a distributor to override the payload
	// location without having to rebuild the source.
	defaultReleaseInfoPadded = "\x00_RELEASE_IMAGE_LOCATION_\x00XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX\x00"
	defaultReleaseInfoPrefix = "\x00_RELEASE_IMAGE_LOCATION_\x00"
	defaultReleaseInfoLength = len(defaultReleaseInfoPadded)
)

// Default abstracts how the binary loads the default release payload. We want to lock the binary
// to the pull spec of the payload we test it with, and since a payload contains an installer image we can't
// know that at build time. Instead, we make it possible to replace the release string after build via a
// known constant in the binary.
func Default() (string, string, error) {
	releaseName := defaultReleaseNameOriginal
	pullspec := defaultReleaseImageOriginal
	if strings.HasPrefix(defaultReleaseInfoPadded, defaultReleaseInfoPrefix) {
		// the defaultReleaseInfoPadded constant hasn't been altered in the binary, fall back to the default
		return pullspec, releaseName, nil
	}
	// look for release image info
	// when extracting from payload 'oc adm release extract --tools', info is injected into binary in form of
	// pullspec,releaseName
	nullTerminator := strings.IndexByte(defaultReleaseInfoPadded, '\x00')
	if nullTerminator == -1 {
		// the binary has been altered, but we didn't find a null terminator within the constant which is an error
		return "", "", fmt.Errorf("release image location was replaced but without a null terminator before %d bytes", defaultReleaseInfoLength)
	}
	if nullTerminator > len(defaultReleaseInfoPadded) {
		// the binary has been altered, but the null terminator is *longer* than the constant encoded in the binary
		return "", "", fmt.Errorf("release image location contains no null-terminator and constant is corrupted")
	}
	injectedSlice := strings.Split(defaultReleaseInfoPadded[:nullTerminator], ",")
	// TODO: check len == 2 after oc #88 change merges
	if len(injectedSlice) < 1 {
		return "", "", fmt.Errorf("unexpected release image info, this binary was incorrectly generated")
	}
	pullspec = injectedSlice[0]
	if len(pullspec) == 0 {
		// the binary has been altered, but the replaced image is empty which is incorrect
		return "", "", fmt.Errorf("release image location is empty, this binary was incorrectly generated")
	}
	// now get short form to report for version
	// TODO: remove this conditional when oc #88 merges
	if len(injectedSlice) > 1 {
		releaseName = injectedSlice[1]
	}
	return pullspec, releaseName, nil
}
