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

// ClusterState represents the values of the 'cluster_state' enumerated type.
type ClusterState string

const (
	// Error during installation.
	ClusterStateError ClusterState = "error"
	// The cluster will consume marginal cloud provider infrastructure but will be counted for quota.
	ClusterStateHibernating ClusterState = "hibernating"
	// The cluster is still being installed.
	ClusterStateInstalling ClusterState = "installing"
	// The cluster is pending resources before being provisioned.
	ClusterStatePending ClusterState = "pending"
	// The cluster is moving from 'Ready' state to 'Hibernating'.
	ClusterStatePoweringDown ClusterState = "powering_down"
	// The cluster is ready to use.
	ClusterStateReady ClusterState = "ready"
	// The cluster is moving from 'Hibernating' state to 'Ready'.
	ClusterStateResuming ClusterState = "resuming"
	// The cluster is being uninstalled.
	ClusterStateUninstalling ClusterState = "uninstalling"
	// The state of the cluster is unknown.
	ClusterStateUnknown ClusterState = "unknown"
	// The cluster is being updated.
	// This state is currently used only by aro hcp clusters.
	ClusterStateUpdating ClusterState = "updating"
	// The cluster is validating user input.
	ClusterStateValidating ClusterState = "validating"
	// The cluster is waiting for user action.
	ClusterStateWaiting ClusterState = "waiting"
)
