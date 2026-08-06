/*
Copyright 2026 Nutanix

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

package feature

import "k8s.io/component-base/featuregate"

const (
	// DefaultToPlaceholderImageName controls whether the NutanixMachineTemplate
	// defaulting webhook injects a placeholder into image.name when the image
	// type is "name" but the name field is unset or empty.
	//
	// This is a brownfield compatibility feature: objects created before the CEL
	// validation rules were added may have the image type set without a value.
	// Enabling this gate allows those objects to pass admission without weakening
	// the API contract globally.
	//
	// owner: @thunderboltsid
	// alpha: v1.5
	DefaultToPlaceholderImageName featuregate.Feature = "DefaultToPlaceholderImageName"

	// DefaultToPlaceholderImageUUID controls whether the NutanixMachineTemplate
	// defaulting webhook injects a placeholder into image.uuid when the image
	// type is "uuid" but the uuid field is unset or empty.
	//
	// This is a brownfield compatibility feature: objects created before the CEL
	// validation rules were added may have the image type set without a value.
	// Enabling this gate allows those objects to pass admission without weakening
	// the API contract globally.
	//
	// owner: @thunderboltsid
	// alpha: v1.5
	DefaultToPlaceholderImageUUID featuregate.Feature = "DefaultToPlaceholderImageUUID"
)

// defaultFeatureGates returns all known feature gates and their default specs.
// To add a new feature, define a constant above and add it here.
func defaultFeatureGates() map[featuregate.Feature]featuregate.FeatureSpec {
	return map[featuregate.Feature]featuregate.FeatureSpec{
		DefaultToPlaceholderImageName: {Default: false, PreRelease: featuregate.Alpha},
		DefaultToPlaceholderImageUUID: {Default: false, PreRelease: featuregate.Alpha},
	}
}
