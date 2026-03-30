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

package v1alpha1 // github.com/openshift-online/ocm-api-model/clientapi/arohcp/v1alpha1

// UpgradePolicyStateValue represents the values of the 'upgrade_policy_state_value' enumerated type.
type UpgradePolicyStateValue string

const (
	// Upgrade got cancelled (temporary state - the policy will get removed).
	UpgradePolicyStateValueCancelled UpgradePolicyStateValue = "cancelled"
	// Upgrade completed (temporary state - the policy will be removed in case of
	// manual upgrade, or move back to pending in case of automatic upgrade)
	UpgradePolicyStateValueCompleted UpgradePolicyStateValue = "completed"
	// Upgrade is taking longer than expected
	UpgradePolicyStateValueDelayed UpgradePolicyStateValue = "delayed"
	// Upgrade failed
	UpgradePolicyStateValueFailed UpgradePolicyStateValue = "failed"
	// Upgrade policy set but an upgrade wasn't scheduled yet
	UpgradePolicyStateValuePending UpgradePolicyStateValue = "pending"
	// Upgrade policy set and was scheduled
	UpgradePolicyStateValueScheduled UpgradePolicyStateValue = "scheduled"
	// Upgrade started
	UpgradePolicyStateValueStarted UpgradePolicyStateValue = "started"
)
