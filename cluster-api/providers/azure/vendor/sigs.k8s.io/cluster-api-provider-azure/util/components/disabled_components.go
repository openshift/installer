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

package components

import (
	"slices"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
)

// IsValidDisableComponent validates if the provided value is a valid disable component by checking if the value exists
// in the infrav1.ValidDisableableComponents map.
func IsValidDisableComponent(value string) bool {
	_, ok := infrav1.ValidDisableableComponents[infrav1.DisableComponent(value)]
	return ok
}

// IsComponentDisabled checks if the provided component is in the list of disabled components.
func IsComponentDisabled(disabledComponents []string, component infrav1.DisableComponent) bool {
	return slices.Contains(disabledComponents, string(component))
}
