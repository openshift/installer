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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1

// UpgradeType represents the values of the 'upgrade_type' enumerated type.
type UpgradeType string

const (
	// Upgrade of OSD cluster, which will upgrade the cluster's control plane, and all the node pools.
	UpgradeTypeOSD UpgradeType = "OSD"
	// Upgrade of an AddOn
	UpgradeTypeAddOn UpgradeType = "ADDON"
	// Control plane upgrade, relevant only for hosted control plane clusters.
	UpgradeTypeControlPlane UpgradeType = "ControlPlane"
	// An upgrade required for security reasons.
	UpgradeTypeControlPlaneCVE UpgradeType = "ControlPlaneCVE"
	// Node pool upgrade, relevant only for hosted control plane clusters.
	UpgradeTypeNodePool UpgradeType = "NodePool"
)
