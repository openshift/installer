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

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

import (
	api_v1 "github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1"
)

// UpgradeType represents the values of the 'upgrade_type' enumerated type.
type UpgradeType = api_v1.UpgradeType

const (
	// Upgrade of OSD cluster, which will upgrade the cluster's control plane, and all the node pools.
	UpgradeTypeOSD UpgradeType = api_v1.UpgradeTypeOSD
	// Upgrade of an AddOn
	UpgradeTypeAddOn UpgradeType = api_v1.UpgradeTypeAddOn
	// Control plane upgrade, relevant only for hosted control plane clusters.
	UpgradeTypeControlPlane UpgradeType = api_v1.UpgradeTypeControlPlane
	// An upgrade required for security reasons.
	UpgradeTypeControlPlaneCVE UpgradeType = api_v1.UpgradeTypeControlPlaneCVE
	// Node pool upgrade, relevant only for hosted control plane clusters.
	UpgradeTypeNodePool UpgradeType = api_v1.UpgradeTypeNodePool
)
