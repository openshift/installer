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

import (
	api_v1 "github.com/openshift-online/ocm-api-model/clientapi/addonsmgmt/v1"
)

// AddonParameterValueType represents the values of the 'addon_parameter_value_type' enumerated type.
type AddonParameterValueType = api_v1.AddonParameterValueType

const (
	// This value type enforces a valid CIDR value to be passed as parameter value
	AddonParameterValueTypeCIDR AddonParameterValueType = api_v1.AddonParameterValueTypeCIDR
	// This value type must be a valid boolean
	AddonParameterValueTypeBoolean AddonParameterValueType = api_v1.AddonParameterValueTypeBoolean
	// This value type must be a valid number, this includes integer and float type numbers
	AddonParameterValueTypeNumber AddonParameterValueType = api_v1.AddonParameterValueTypeNumber
	// This value must match a valid SKU resource in OCM
	AddonParameterValueTypeResource AddonParameterValueType = api_v1.AddonParameterValueTypeResource
	// This value must match a valid SKU resource in OCM and allows for validation of SKU resource in OCM
	AddonParameterValueTypeResourceRequirement AddonParameterValueType = api_v1.AddonParameterValueTypeResourceRequirement
	// This value type must be a valid string
	AddonParameterValueTypeString AddonParameterValueType = api_v1.AddonParameterValueTypeString
)
