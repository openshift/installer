/*
Copyright The Kubernetes Authors.
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

package util

import "strings"

type NormalizedArch string

const (
	// ArchitectureAmd64 is the normalized architecture name for amd64.
	ArchitectureAmd64 NormalizedArch = "amd64"
	// ArchitectureArm64 is the normalized architecture name for arm64.
	ArchitectureArm64 NormalizedArch = "arm64"
)

// machineTypePrefixArchitectureMap contains a map of (machineTypePrefix, architecture) tuples
var machineTypePrefixArchitectureMap = map[string]NormalizedArch{
	"c4a": ArchitectureArm64,
	"t2a": ArchitectureArm64,
}

// CPUArchitecture gets a machineType string parameter and returns the architecture for the machineType, if it is known
// and stored in the machineTypePrefixArchitectureMap. Otherwise, it returns amd64.
func CPUArchitecture(machineType string) NormalizedArch {
	prefix, _, _ := strings.Cut(machineType, "-")
	if arch, ok := machineTypePrefixArchitectureMap[prefix]; ok {
		return arch
	}
	// Fallback to Amd64 for any unknown machine types prefixes
	return ArchitectureAmd64
}
