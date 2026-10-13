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

// AddonStatusConditionType represents the values of the 'addon_status_condition_type' enumerated type.
type AddonStatusConditionType = api_v1.AddonStatusConditionType

const (
	//
	AddonStatusConditionTypeAvailable AddonStatusConditionType = api_v1.AddonStatusConditionTypeAvailable
	//
	AddonStatusConditionTypeDegraded AddonStatusConditionType = api_v1.AddonStatusConditionTypeDegraded
	//
	AddonStatusConditionTypeDeleteTimeout AddonStatusConditionType = api_v1.AddonStatusConditionTypeDeleteTimeout
	//
	AddonStatusConditionTypeHealthy AddonStatusConditionType = api_v1.AddonStatusConditionTypeHealthy
	//
	AddonStatusConditionTypeInstalled AddonStatusConditionType = api_v1.AddonStatusConditionTypeInstalled
	//
	AddonStatusConditionTypePaused AddonStatusConditionType = api_v1.AddonStatusConditionTypePaused
	//
	AddonStatusConditionTypeReadyToBeDeleted AddonStatusConditionType = api_v1.AddonStatusConditionTypeReadyToBeDeleted
	//
	AddonStatusConditionTypeUpgradeStarted AddonStatusConditionType = api_v1.AddonStatusConditionTypeUpgradeStarted
	//
	AddonStatusConditionTypeUpgradeSucceeded AddonStatusConditionType = api_v1.AddonStatusConditionTypeUpgradeSucceeded
)
