/*
Copyright 2023 The Kubernetes Authors.

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

package v1beta1

import (
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
)

const (
	// ReadyCondition is our overall Ready condition.
	ReadyCondition string = "Ready"

	// GKEControlPlaneReadyCondition condition reports on the successful reconciliation of GKE control plane.
	GKEControlPlaneReadyCondition clusterv1beta1.ConditionType = "GKEControlPlaneReady"
	// GKEControlPlaneCreatingCondition condition reports on whether the GKE control plane is creating.
	GKEControlPlaneCreatingCondition clusterv1beta1.ConditionType = "GKEControlPlaneCreating"
	// GKEControlPlaneUpdatingCondition condition reports on whether the GKE control plane is updating.
	GKEControlPlaneUpdatingCondition clusterv1beta1.ConditionType = "GKEControlPlaneUpdating"
	// GKEControlPlaneDeletingCondition condition reports on whether the GKE control plane is deleting.
	GKEControlPlaneDeletingCondition clusterv1beta1.ConditionType = "GKEControlPlaneDeleting"

	// GKEControlPlaneCreatingReason used to report GKE control plane being created.
	GKEControlPlaneCreatingReason = "GKEControlPlaneCreating"
	// GKEControlPlaneCreatedReason used to report GKE control plane is created.
	GKEControlPlaneCreatedReason = "GKEControlPlaneCreated"
	// GKEControlPlaneUpdatedReason used to report GKE control plane is updated.
	GKEControlPlaneUpdatedReason = "GKEControlPlaneUpdated"
	// GKEControlPlaneDeletingReason used to report GKE control plane being deleted.
	GKEControlPlaneDeletingReason = "GKEControlPlaneDeleting"
	// GKEControlPlaneDeletedReason used to report GKE control plane is deleted.
	GKEControlPlaneDeletedReason = "GKEControlPlaneDeleted"
	// GKEControlPlaneErrorReason used to report GKE control plane is in error state.
	GKEControlPlaneErrorReason = "GKEControlPlaneError"
	// GKEControlPlaneReconciliationFailedReason used to report failures while reconciling GKE control plane.
	GKEControlPlaneReconciliationFailedReason = "GKEControlPlaneReconciliationFailed"
	// GKEControlPlaneRequiresAtLeastOneNodePoolReason used to report that no node pool is specified for the GKE control plane.
	GKEControlPlaneRequiresAtLeastOneNodePoolReason = "GKEControlPlaneRequiresAtLeastOneNodePool"

	// GKEMachinePoolReadyCondition condition reports on the successful reconciliation of GKE node pool.
	GKEMachinePoolReadyCondition clusterv1beta1.ConditionType = "GKEMachinePoolReady"
	// GKEMachinePoolCreatingCondition condition reports on whether the GKE node pool is creating.
	GKEMachinePoolCreatingCondition clusterv1beta1.ConditionType = "GKEMachinePoolCreating"
	// GKEMachinePoolUpdatingCondition condition reports on whether the GKE node pool is updating.
	GKEMachinePoolUpdatingCondition clusterv1beta1.ConditionType = "GKEMachinePoolUpdating"
	// GKEMachinePoolDeletingCondition condition reports on whether the GKE node pool is deleting.
	GKEMachinePoolDeletingCondition clusterv1beta1.ConditionType = "GKEMachinePoolDeleting"

	// WaitingForGKEControlPlaneReason used when the machine pool is waiting for GKE control plane infrastructure to be ready before proceeding.
	WaitingForGKEControlPlaneReason = "WaitingForGKEControlPlane"
	// GKEMachinePoolCreatingReason used to report GKE node pool being created.
	GKEMachinePoolCreatingReason = "GKEMachinePoolCreating"
	// GKEMachinePoolCreatedReason used to report GKE node pool is created.
	GKEMachinePoolCreatedReason = "GKEMachinePoolCreated"
	// GKEMachinePoolUpdatedReason used to report GKE node pool is updated.
	GKEMachinePoolUpdatedReason = "GKEMachinePoolUpdated"
	// GKEMachinePoolDeletingReason used to report GKE node pool being deleted.
	GKEMachinePoolDeletingReason = "GKEMachinePoolDeleting"
	// GKEMachinePoolDeletedReason used to report GKE node pool is deleted.
	GKEMachinePoolDeletedReason = "GKEMachinePoolDeleted"
	// GKEMachinePoolErrorReason used to report GKE node pool is in error state.
	GKEMachinePoolErrorReason = "GKEMachinePoolError"
	// GKEMachinePoolReconciliationFailedReason used to report failures while reconciling GKE node pool.
	GKEMachinePoolReconciliationFailedReason = "GKEMachinePoolReconciliationFailed"

	// MIGReadyCondition reports on current status of the managed instance group. Ready indicates the group is provisioned.
	MIGReadyCondition clusterv1.ConditionType = "ManagedInstanceGroupReady"
	// MIGProvisionFailedReason used for failures during managed instance group provisioning.
	MIGProvisionFailedReason = "ManagedInstanceGroupProvisionFailed"
	// MIGDeletionInProgress MIG is in a deletion in progress state.
	MIGDeletionInProgress = "ManagedInstanceGroupDeletionInProgress"

	// InstanceTemplateReadyCondition represents the status of an AWSMachinePool's associated Launch Template.
	InstanceTemplateReadyCondition clusterv1.ConditionType = "InstanceTemplateReady"
	// InstanceTemplateReconcileFailedReason used for failures during Launch Template reconciliation.
	InstanceTemplateReconcileFailedReason = "InstanceTemplateReconcileFailed"
)
