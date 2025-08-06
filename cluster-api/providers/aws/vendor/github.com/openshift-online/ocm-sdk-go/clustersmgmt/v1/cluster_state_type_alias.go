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

// ClusterState represents the values of the 'cluster_state' enumerated type.
type ClusterState = api_v1.ClusterState

const (
	// Error during installation.
	ClusterStateError ClusterState = api_v1.ClusterStateError
	// The cluster will consume marginal cloud provider infrastructure but will be counted for quota.
	ClusterStateHibernating ClusterState = api_v1.ClusterStateHibernating
	// The cluster is still being installed.
	ClusterStateInstalling ClusterState = api_v1.ClusterStateInstalling
	// The cluster is pending resources before being provisioned.
	ClusterStatePending ClusterState = api_v1.ClusterStatePending
	// The cluster is moving from 'Ready' state to 'Hibernating'.
	ClusterStatePoweringDown ClusterState = api_v1.ClusterStatePoweringDown
	// The cluster is ready to use.
	ClusterStateReady ClusterState = api_v1.ClusterStateReady
	// The cluster is moving from 'Hibernating' state to 'Ready'.
	ClusterStateResuming ClusterState = api_v1.ClusterStateResuming
	// The cluster is being uninstalled.
	ClusterStateUninstalling ClusterState = api_v1.ClusterStateUninstalling
	// The state of the cluster is unknown.
	ClusterStateUnknown ClusterState = api_v1.ClusterStateUnknown
	// The cluster is being updated.
	// This state is currently used only by aro hcp clusters.
	ClusterStateUpdating ClusterState = api_v1.ClusterStateUpdating
	// The cluster is validating user input.
	ClusterStateValidating ClusterState = api_v1.ClusterStateValidating
	// The cluster is waiting for user action.
	ClusterStateWaiting ClusterState = api_v1.ClusterStateWaiting
)
