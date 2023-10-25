/*
Copyright (c) 2020 Red Hat, Inc.

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

// IMPORTANT: This file has been generated automatically, refrain from modifying it manually as all
// your changes will be lost when the file is generated again.

package v1 // github.com/openshift-online/ocm-sdk-go/addonsmgmt/v1

// AddonRequirementResource represents the values of the 'addon_requirement_resource' enumerated type.
type AddonRequirementResource string

const (
	// This requirement resource data will be validated and checked against addon installation resources
	AddonRequirementResourceAddon AddonRequirementResource = "addon"
	// This requirement resource data will be validated and checked against cluster resources
	AddonRequirementResourceCluster AddonRequirementResource = "cluster"
	// This requirement resource data will be validated and checked against machine pool resources
	AddonRequirementResourceMachinePool AddonRequirementResource = "machine_pool"
)
