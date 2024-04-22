/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package version

import (
	"fmt"
	"runtime"
	"strings"

	"k8s.io/apimachinery/pkg/version"
)

var (
	// commitFromGit is a constant representing the source version that
	// generated this build. It should be set during build via -ldflags.
	commitFromGit string
	// versionFromGit is a constant representing the version tag that
	// generated this build. It should be set during build via -ldflags.
	versionFromGit = "unknown"
	// major version
	majorFromGit string
	// minor version
	minorFromGit string
	// build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
	buildDate string
	// state of git tree, either "clean" or "dirty"
	gitTreeState string
)

// Get returns the overall codebase version. It's for detecting
// what code a binary was built from.
func Get() version.Info {
	return version.Info{
		Major:        majorFromGit,
		Minor:        minorFromGit,
		GitCommit:    commitFromGit,
		GitVersion:   versionFromGit,
		GitTreeState: gitTreeState,
		BuildDate:    buildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// This file handles correctly identifying the default release version, which is expected to be
// replaced in the binary post-compile by the release name extracted from a payload. The expected modification is:
//
// 1. Extract a release binary that contains a cli image
// 2. Identify the release name, add a NUL terminator byte (0x00) to the end, calculate length
// 3. Length must be less than the marker length
// 4. Search through the cli binary for `\x00_RELEASE_VERSION_LOCATION_\x00<PADDING_TO_LENGTH>`
//    where padding is the ASCII character X and length is the total length of the marker
// 5. Overwrite the beginning of the marker with the release name and a NUL terminator byte (0x00)

var (
	// defaultReleaseInfoPadded may be replaced in the binary with Release Metadata: Version that overrides defaultVersion as
	// a null-terminated string within the allowed character length. This allows a distributor to override the payload
	// location without having to rebuild the source.
	defaultVersionPadded = "\x00_RELEASE_VERSION_LOCATION_\x00XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX\x00"
	defaultVersionPrefix = "\x00_RELEASE_VERSION_LOCATION_\x00"
	defaultVersionLength = len(defaultVersionPadded)
)

// ExtractVersion() abstracts how the binary loads the default release payload. We want to lock the binary
// to version of the payload we test it with, and since a payload contains the cli image we can't
// know that at build time. Instead, we make it possible to replace the version string after build via a
// known constant in the binary.  When extracting oc from release image, 'oc version' reports the payload version.
func ExtractVersion() (version.Info, string, error) {
	if strings.HasPrefix(defaultVersionPadded, defaultVersionPrefix) {
		return Get(), "", nil
	}
	nullTerminator := strings.IndexByte(defaultVersionPadded, '\x00')
	if nullTerminator == -1 {
		// the binary has been altered, but we didn't find a null terminator within the release name constant which is an error
		return version.Info{}, "", fmt.Errorf("release name location was replaced but without a null terminator before %d bytes", defaultVersionLength)
	}
	if nullTerminator > len(defaultVersionPadded) {
		// the binary has been altered, but the null terminator is *longer* than the constant encoded in the binary
		return version.Info{}, "", fmt.Errorf("release name location contains no null-terminator and constant is corrupted")
	}
	releaseName := defaultVersionPadded[:nullTerminator]
	if len(releaseName) == 0 {
		// the binary has been altered, but the replaced release name is empty which is incorrect
		// the oc binary will not be pinned to Release Metadata:Version
		return version.Info{}, "", fmt.Errorf("oc version was incorrectly replaced during extract")
	}
	return Get(), releaseName, nil
}
