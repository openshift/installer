/*
Copyright 2025 The Kubernetes Authors.

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

package controllers

import (
	"fmt"
	"strings"
)

// TODO: EKS and ROSA are excluded from this list for the time being because they are behind feature gates.
// They should be added to this list when they graduate.
const (
	Unmanaged = "unmanaged"
)

// disabledControllers tracks which controllers are disabled.
// A value of `false` (default) means a controller is _enabled_.
var disabledControllers = map[string]bool{
	Unmanaged: false,
}

var notValidErr = "%q is not a valid controller name"

// IsDisabled checks if a controller is disabled.
// If the name provided is not in the map, this will return 'false'.
func IsDisabled(name string) bool {
	return disabledControllers[name]
}

// GetValidNames returns a list of controller names that are valid to disable.
func GetValidNames() []string {
	ret := make([]string, 0, len(disabledControllers))
	for name := range disabledControllers {
		ret = append(ret, name)
	}
	return ret
}

// ValidateNamesAndDisable validates a list of controller names against the known set, and disables valid names.
func ValidateNamesAndDisable(names []string) error {
	// This list is not de-deduplicated, so in someone could specify valid names multiple times.
	for _, n := range names {
		// Make sure we're doing a case-insensitive comaparison
		name := strings.ToLower(n)
		if _, ok := disabledControllers[name]; !ok {
			return fmt.Errorf(notValidErr, name)
		}
		disabledControllers[name] = true
	}
	return nil
}
