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

package v1beta1

import (
	"strings"
)

// NormalizeVersion normalizes the Kubernetes version string to include a "v" prefix.
func NormalizeVersion(version string) string {
	if version != "" && !strings.HasPrefix(version, "v") {
		normalizedVersion := "v" + version
		version = normalizedVersion
	}
	return version
}

// SetZeroPointerDefault sets the default value for a pointer to a value for any comparable type.
func SetZeroPointerDefault[T comparable](field *T, value T) {
	if field == nil {
		// shouldn't happen with proper use
		return
	}
	var zero T
	if *field == zero {
		*field = value
	}
}
