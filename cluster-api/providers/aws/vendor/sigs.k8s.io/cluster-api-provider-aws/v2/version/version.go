/*
Copyright 2018 The Kubernetes Authors.

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

// Package version provides the version of the manager.
package version

import (
	"fmt"
	"runtime"

	"github.com/aws/aws-sdk-go-v2/aws"
)

var (
	gitMajor     string // major version, always numeric
	gitMinor     string // minor version, numeric possibly followed by "+"
	gitVersion   string // semantic version, derived by build scripts
	gitCommit    string // sha1 from git, output of $(git rev-parse HEAD)
	gitTreeState string // state of git tree, either "clean" or "dirty"
	buildDate    string // build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
)

// Info defines the version.
type Info struct {
	Major         string `json:"major,omitempty"`
	Minor         string `json:"minor,omitempty"`
	GitVersion    string `json:"gitVersion,omitempty"`
	GitCommit     string `json:"gitCommit,omitempty"`
	GitTreeState  string `json:"gitTreeState,omitempty"`
	BuildDate     string `json:"buildDate,omitempty"`
	GoVersion     string `json:"goVersion,omitempty"`
	AwsSdkVersion string `json:"awsSdkVersion,omitempty"`
	Compiler      string `json:"compiler,omitempty"`
	Platform      string `json:"platform,omitempty"`
}

// Get returns metadata and information regarding the version.
func Get() Info {
	return Info{
		Major:         gitMajor,
		Minor:         gitMinor,
		GitVersion:    gitVersion,
		GitCommit:     gitCommit,
		GitTreeState:  gitTreeState,
		BuildDate:     buildDate,
		GoVersion:     runtime.Version(),
		AwsSdkVersion: fmt.Sprintf("v%s", aws.SDKVersion),
		Compiler:      runtime.Compiler,
		Platform:      fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// String returns info as a human-friendly version string.
func (info Info) String() string {
	return info.GitVersion
}
