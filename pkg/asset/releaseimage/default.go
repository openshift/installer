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
// 3. Length must be less than the marker length
// 4. Search through the installer binary looking for `\x00_RELEASE_IMAGE_LOCATION_\x00<PADDING_TO_LENGTH>`
//    where padding is the ASCII character X and length is the total length of the marker
// 5. Overwrite the beginning of the marker with the release pullspec and a NUL terminator byte (0x00)
//
// On start the installer examines the constant and if it has been modified from the default the installer
// will use that image.

var (
	// defaultReleaseImageOriginal is the value served when defaultReleaseImagePadded is unmodified.
	defaultReleaseImageOriginal = "registry.ci.openshift.org/origin/release:4.9"
	// defaultReleaseImagePadded may be replaced in the binary with a pull spec that overrides defaultReleaseImage as
	// a null-terminated string within the allowed character length. This allows a distributor to override the payload
	// location without having to rebuild the source.
	defaultReleaseImagePadded = "\x00_RELEASE_IMAGE_LOCATION_\x00XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX\x00"
	defaultReleaseImagePrefix = "\x00_RELEASE_IMAGE_LOCATION_\x00"
	defaultReleaseImageLength = len(defaultReleaseImagePadded)
)

// Default abstracts how the binary loads the default release payload. We want to lock the binary
// to the pull spec of the payload we test it with, and since a payload contains an installer image we can't
// know that at build time. Instead, we make it possible to replace the release string after build via a
// known constant in the binary.
func Default() (string, error) {
	if strings.HasPrefix(defaultReleaseImagePadded, defaultReleaseImagePrefix) {
		// the defaultReleaseImagePadded constant hasn't been altered in the binary, fall back to the default
		return defaultReleaseImageOriginal, nil
	}
	nullTerminator := strings.IndexByte(defaultReleaseImagePadded, '\x00')
	if nullTerminator == -1 {
		// the binary has been altered, but we didn't find a null terminator within the constant which is an error
		return "", fmt.Errorf("release image location was replaced but without a null terminator before %d bytes", defaultReleaseImageLength)
	}
	if nullTerminator > len(defaultReleaseImagePadded) {
		// the binary has been altered, but the null terminator is *longer* than the constant encoded in the binary
		return "", fmt.Errorf("release image location contains no null-terminator and constant is corrupted")
	}
	pullspec := defaultReleaseImagePadded[:nullTerminator]
	if len(pullspec) == 0 {
		// the binary has been altered, but the replaced image is empty which is incorrect
		return "", fmt.Errorf("release image location is empty, this binary was incorrectly generated")
	}
	return pullspec, nil
}
