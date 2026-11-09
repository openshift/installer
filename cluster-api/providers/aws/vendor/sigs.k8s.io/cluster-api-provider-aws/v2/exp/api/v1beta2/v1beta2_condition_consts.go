/*
Copyright 2026 The Kubernetes Authors.

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

package v1beta2

import clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"

// AWSMachinePool v1beta2 condition types.
const (
	// AWSMachinePoolReadyV1Beta2Condition defines the Ready condition type that summarizes the operational state of an AWSMachinePool.
	AWSMachinePoolReadyV1Beta2Condition = clusterv1.ReadyCondition

	// AWSMachinePoolASGReadyV1Beta2Condition reports on current status of the autoscaling group. Ready indicates the group is provisioned.
	AWSMachinePoolASGReadyV1Beta2Condition = "ASGReady"

	// AWSMachinePoolLaunchTemplateReadyV1Beta2Condition reports on the status of an AWSMachinePool's associated launch template.
	AWSMachinePoolLaunchTemplateReadyV1Beta2Condition = "LaunchTemplateReady"

	// AWSMachinePoolPreLaunchTemplateUpdateCheckV1Beta2Condition reports if all prerequisites are met for launch template update.
	AWSMachinePoolPreLaunchTemplateUpdateCheckV1Beta2Condition = "PreLaunchTemplateUpdateCheckSuccess"

	// AWSMachinePoolPostLaunchTemplateUpdateOperationV1Beta2Condition reports on successfully completed post launch template update operation.
	AWSMachinePoolPostLaunchTemplateUpdateOperationV1Beta2Condition = "PostLaunchTemplateUpdateOperationSuccess"

	// AWSMachinePoolInstanceRefreshStartedV1Beta2Condition reports on successfully starting instance refresh.
	AWSMachinePoolInstanceRefreshStartedV1Beta2Condition = "InstanceRefreshStarted"

	// AWSMachinePoolLifecycleHookReadyV1Beta2Condition reports on the status of the lifecycle hook.
	AWSMachinePoolLifecycleHookReadyV1Beta2Condition = "LifecycleHookReady"
)

// AWSMachinePool v1beta2 reason constants.
const (
	// AWSMachinePoolReadyV1Beta2Reason indicates the AWSMachinePool is ready.
	AWSMachinePoolReadyV1Beta2Reason = clusterv1.ReadyReason

	// AWSMachinePoolNotReadyV1Beta2Reason indicates the AWSMachinePool is not ready.
	AWSMachinePoolNotReadyV1Beta2Reason = clusterv1.NotReadyReason

	// AWSMachinePoolDeletingV1Beta2Reason indicates the AWSMachinePool is being deleted.
	AWSMachinePoolDeletingV1Beta2Reason = clusterv1.DeletingReason

	// AWSMachinePoolASGNotFoundV1Beta2Reason used when the autoscaling group couldn't be retrieved.
	AWSMachinePoolASGNotFoundV1Beta2Reason = "ASGNotFound"

	// AWSMachinePoolASGProvisionFailedV1Beta2Reason used for failures during autoscaling group provisioning.
	AWSMachinePoolASGProvisionFailedV1Beta2Reason = "ASGProvisionFailed"

	// AWSMachinePoolASGDeletionInProgressV1Beta2Reason used when the autoscaling group is in a deletion in progress state.
	AWSMachinePoolASGDeletionInProgressV1Beta2Reason = "ASGDeletionInProgress"

	// AWSMachinePoolLaunchTemplateNotFoundV1Beta2Reason used when an associated launch template can't be found.
	AWSMachinePoolLaunchTemplateNotFoundV1Beta2Reason = "LaunchTemplateNotFound"

	// AWSMachinePoolLaunchTemplateCreateFailedV1Beta2Reason used for failures during launch template creation.
	AWSMachinePoolLaunchTemplateCreateFailedV1Beta2Reason = "LaunchTemplateCreateFailed"

	// AWSMachinePoolLaunchTemplateReconcileFailedV1Beta2Reason used for failures during launch template reconciliation.
	AWSMachinePoolLaunchTemplateReconcileFailedV1Beta2Reason = "LaunchTemplateReconcileFailed"

	// AWSMachinePoolLaunchTemplateNitroEnclaveEdgeZoneV1Beta2Reason used when enclaveOptions is enabled but the pool
	// targets a Local Zone or Wavelength Zone, which does not support Nitro Enclaves.
	AWSMachinePoolLaunchTemplateNitroEnclaveEdgeZoneV1Beta2Reason = "NitroEnclaveEdgeZoneUnsupported"

	// AWSMachinePoolPreLaunchTemplateUpdateCheckFailedV1Beta2Reason used to report when not all prerequisites are met for launch template update.
	AWSMachinePoolPreLaunchTemplateUpdateCheckFailedV1Beta2Reason = "PreLaunchTemplateUpdateCheckFailed"

	// AWSMachinePoolPostLaunchTemplateUpdateOperationFailedV1Beta2Reason used to report when post launch template update operation failed.
	AWSMachinePoolPostLaunchTemplateUpdateOperationFailedV1Beta2Reason = "PostLaunchTemplateUpdateOperationFailed"

	// AWSMachinePoolInstanceRefreshNotReadyV1Beta2Reason used to report instance refresh is not initiated.
	// If there are instance refreshes that are in progress, then a new instance refresh request will fail.
	AWSMachinePoolInstanceRefreshNotReadyV1Beta2Reason = "InstanceRefreshNotReady"

	// AWSMachinePoolInstanceRefreshFailedV1Beta2Reason used to report when instance refresh is not initiated.
	AWSMachinePoolInstanceRefreshFailedV1Beta2Reason = "InstanceRefreshFailed"

	// AWSMachinePoolMachineCreationFailedV1Beta2Reason used to report when creating AWSMachines to represent ASG machines failed.
	AWSMachinePoolMachineCreationFailedV1Beta2Reason = "AWSMachineCreationFailed"

	// AWSMachinePoolMachineDeletionFailedV1Beta2Reason used to report when deleting AWSMachines failed.
	AWSMachinePoolMachineDeletionFailedV1Beta2Reason = "AWSMachineDeletionFailed"

	// AWSMachinePoolLifecycleHookCreationFailedV1Beta2Reason used for failures during lifecycle hook creation.
	AWSMachinePoolLifecycleHookCreationFailedV1Beta2Reason = "LifecycleHookCreationFailed"

	// AWSMachinePoolLifecycleHookUpdateFailedV1Beta2Reason used for failures during lifecycle hook update.
	AWSMachinePoolLifecycleHookUpdateFailedV1Beta2Reason = "LifecycleHookUpdateFailed"

	// AWSMachinePoolLifecycleHookDeletionFailedV1Beta2Reason used for failures during lifecycle hook deletion.
	AWSMachinePoolLifecycleHookDeletionFailedV1Beta2Reason = "LifecycleHookDeletionFailed"
)

// AWSManagedMachinePool v1beta2 condition types.
const (
	// AWSManagedMachinePoolReadyV1Beta2Condition defines the Ready condition type that summarizes the operational state of an AWSManagedMachinePool.
	AWSManagedMachinePoolReadyV1Beta2Condition = clusterv1.ReadyCondition

	// AWSManagedMachinePoolEKSNodegroupReadyV1Beta2Condition reports on the successful reconciliation of the EKS node group.
	AWSManagedMachinePoolEKSNodegroupReadyV1Beta2Condition = "EKSNodegroupReady"

	// AWSManagedMachinePoolIAMNodegroupRolesReadyV1Beta2Condition reports on the successful reconciliation of EKS node group IAM roles.
	AWSManagedMachinePoolIAMNodegroupRolesReadyV1Beta2Condition = "IAMNodegroupRolesReady"

	// AWSManagedMachinePoolLaunchTemplateReadyV1Beta2Condition reports on the status of the associated launch template.
	AWSManagedMachinePoolLaunchTemplateReadyV1Beta2Condition = "LaunchTemplateReady"
)

// AWSManagedMachinePool v1beta2 reason constants.
const (
	// AWSManagedMachinePoolReadyV1Beta2Reason indicates the AWSManagedMachinePool is ready.
	AWSManagedMachinePoolReadyV1Beta2Reason = clusterv1.ReadyReason

	// AWSManagedMachinePoolNotReadyV1Beta2Reason indicates the AWSManagedMachinePool is not ready.
	AWSManagedMachinePoolNotReadyV1Beta2Reason = clusterv1.NotReadyReason

	// AWSManagedMachinePoolDeletingV1Beta2Reason indicates the AWSManagedMachinePool is being deleted.
	AWSManagedMachinePoolDeletingV1Beta2Reason = clusterv1.DeletingReason

	// AWSManagedMachinePoolEKSNodegroupReconciliationFailedV1Beta2Reason used to report failures while reconciling EKS node group.
	AWSManagedMachinePoolEKSNodegroupReconciliationFailedV1Beta2Reason = "EKSNodegroupReconciliationFailed"

	// AWSManagedMachinePoolWaitingForEKSControlPlaneV1Beta2Reason used when the machine pool is waiting for
	// EKS control plane infrastructure to be ready before proceeding.
	AWSManagedMachinePoolWaitingForEKSControlPlaneV1Beta2Reason = "WaitingForEKSControlPlane"

	// AWSManagedMachinePoolIAMNodegroupRolesReconciliationFailedV1Beta2Reason used to report failures while reconciling EKS nodegroup IAM roles.
	AWSManagedMachinePoolIAMNodegroupRolesReconciliationFailedV1Beta2Reason = "IAMNodegroupRolesReconciliationFailed"

	// AWSManagedMachinePoolLaunchTemplateNotFoundV1Beta2Reason used when an associated launch template can't be found.
	AWSManagedMachinePoolLaunchTemplateNotFoundV1Beta2Reason = "LaunchTemplateNotFound"

	// AWSManagedMachinePoolLaunchTemplateCreateFailedV1Beta2Reason used for failures during launch template creation.
	AWSManagedMachinePoolLaunchTemplateCreateFailedV1Beta2Reason = "LaunchTemplateCreateFailed"

	// AWSManagedMachinePoolLaunchTemplateReconcileFailedV1Beta2Reason used for failures during launch template reconciliation.
	AWSManagedMachinePoolLaunchTemplateReconcileFailedV1Beta2Reason = "LaunchTemplateReconcileFailed"
)

// AWSFargateProfile v1beta2 condition types.
const (
	// AWSFargateProfileReadyV1Beta2Condition defines the Ready condition type that summarizes the operational state of an AWSFargateProfile.
	AWSFargateProfileReadyV1Beta2Condition = clusterv1.ReadyCondition

	// AWSFargateProfileEKSFargateProfileReadyV1Beta2Condition reports on the successful reconciliation of the EKS Fargate profile.
	AWSFargateProfileEKSFargateProfileReadyV1Beta2Condition = "EKSFargateProfileReady"

	// AWSFargateProfileEKSFargateCreatingV1Beta2Condition reports whether the Fargate profile is being created.
	AWSFargateProfileEKSFargateCreatingV1Beta2Condition = "EKSFargateCreating"

	// AWSFargateProfileEKSFargateDeletingV1Beta2Condition used to report that the profile is deleting.
	AWSFargateProfileEKSFargateDeletingV1Beta2Condition = "EKSFargateDeleting"

	// AWSFargateProfileIAMFargateRolesReadyV1Beta2Condition reports on the successful reconciliation of Fargate IAM roles.
	AWSFargateProfileIAMFargateRolesReadyV1Beta2Condition = "IAMFargateRolesReady"
)

// AWSFargateProfile v1beta2 reason constants.
const (
	// AWSFargateProfileReadyV1Beta2Reason indicates the AWSFargateProfile is ready.
	AWSFargateProfileReadyV1Beta2Reason = clusterv1.ReadyReason

	// AWSFargateProfileNotReadyV1Beta2Reason indicates the AWSFargateProfile is not ready.
	AWSFargateProfileNotReadyV1Beta2Reason = clusterv1.NotReadyReason

	// AWSFargateProfileDeletingV1Beta2Reason indicates the AWSFargateProfile is being deleted.
	AWSFargateProfileDeletingV1Beta2Reason = clusterv1.DeletingReason

	// AWSFargateProfileReconciliationFailedV1Beta2Reason used to report failures while reconciling EKS Fargate profile.
	AWSFargateProfileReconciliationFailedV1Beta2Reason = "EKSFargateReconciliationFailed"

	// AWSFargateProfileCreatingV1Beta2Reason used when the profile is creating.
	AWSFargateProfileCreatingV1Beta2Reason = "Creating"

	// AWSFargateProfileCreatedV1Beta2Reason used when the profile is created.
	AWSFargateProfileCreatedV1Beta2Reason = "Created"

	// AWSFargateProfileFailedV1Beta2Reason used when the profile failed.
	AWSFargateProfileFailedV1Beta2Reason = "Failed"

	// AWSFargateProfileDeletedV1Beta2Reason used when the profile is deleted.
	AWSFargateProfileDeletedV1Beta2Reason = "Deleted"

	// AWSFargateProfileIAMFargateRolesReconciliationFailedV1Beta2Reason used to report failures while reconciling Fargate IAM roles.
	AWSFargateProfileIAMFargateRolesReconciliationFailedV1Beta2Reason = "IAMFargateRolesReconciliationFailed"
)

// ROSAMachinePool v1beta2 condition types.
const (
	// ROSAMachinePoolReadyV1Beta2Condition defines the Ready condition type that summarizes the operational state of a ROSAMachinePool.
	ROSAMachinePoolReadyV1Beta2Condition = clusterv1.ReadyCondition

	// ROSAMachinePoolUpgradingV1Beta2Condition reports whether ROSAMachinePool is upgrading or not.
	ROSAMachinePoolUpgradingV1Beta2Condition = "ROSAMachinePoolUpgrading"
)

// ROSAMachinePool v1beta2 reason constants.
const (
	// ROSAMachinePoolReadyV1Beta2Reason indicates the ROSAMachinePool is ready.
	ROSAMachinePoolReadyV1Beta2Reason = clusterv1.ReadyReason

	// ROSAMachinePoolNotReadyV1Beta2Reason indicates the ROSAMachinePool is not ready.
	ROSAMachinePoolNotReadyV1Beta2Reason = clusterv1.NotReadyReason

	// ROSAMachinePoolDeletingV1Beta2Reason indicates the ROSAMachinePool is being deleted.
	ROSAMachinePoolDeletingV1Beta2Reason = clusterv1.DeletingReason

	// ROSAMachinePoolWaitingForROSAControlPlaneV1Beta2Reason used when the machine pool is waiting for
	// ROSA control plane infrastructure to be ready before proceeding.
	ROSAMachinePoolWaitingForROSAControlPlaneV1Beta2Reason = "WaitingForROSAControlPlane"

	// ROSAMachinePoolReconciliationFailedV1Beta2Reason used to report failures while reconciling ROSAMachinePool.
	ROSAMachinePoolReconciliationFailedV1Beta2Reason = "ReconciliationFailed"
)

// ROSACluster v1beta2 condition types.
const (
	// ROSAClusterReadyV1Beta2Condition defines the Ready condition type that summarizes the operational state of a ROSACluster.
	ROSAClusterReadyV1Beta2Condition = clusterv1.ReadyCondition
)

// ROSACluster v1beta2 reason constants.
const (
	// ROSAClusterReadyV1Beta2Reason indicates the ROSACluster is ready.
	ROSAClusterReadyV1Beta2Reason = clusterv1.ReadyReason

	// ROSAClusterNotReadyV1Beta2Reason indicates the ROSACluster is not ready.
	ROSAClusterNotReadyV1Beta2Reason = clusterv1.NotReadyReason
)

// ROSANetwork v1beta2 condition types.
const (
	// ROSANetworkReadyV1Beta2Condition defines the Ready condition type that summarizes the operational state of a ROSANetwork.
	ROSANetworkReadyV1Beta2Condition = clusterv1.ReadyCondition
)

// ROSANetwork v1beta2 reason constants.
const (
	// ROSANetworkReadyV1Beta2Reason indicates the ROSANetwork is ready.
	ROSANetworkReadyV1Beta2Reason = clusterv1.ReadyReason

	// ROSANetworkNotReadyV1Beta2Reason indicates the ROSANetwork is not ready.
	ROSANetworkNotReadyV1Beta2Reason = clusterv1.NotReadyReason

	// ROSANetworkCreatingV1Beta2Reason used when ROSANetwork is being created.
	ROSANetworkCreatingV1Beta2Reason = "Creating"

	// ROSANetworkCreatedV1Beta2Reason used when ROSANetwork is created.
	ROSANetworkCreatedV1Beta2Reason = "Created"

	// ROSANetworkFailedV1Beta2Reason used when ROSANetwork creation failed.
	ROSANetworkFailedV1Beta2Reason = "Failed"

	// ROSANetworkDeletingV1Beta2Reason indicates the ROSANetwork is being deleted.
	ROSANetworkDeletingV1Beta2Reason = clusterv1.DeletingReason

	// ROSANetworkDeletionFailedV1Beta2Reason used to report failures while deleting ROSANetwork.
	ROSANetworkDeletionFailedV1Beta2Reason = "DeletionFailed"
)

// ROSARoleConfig v1beta2 condition types.
const (
	// ROSARoleConfigReadyV1Beta2Condition defines the Ready condition type that summarizes the operational state of a ROSARoleConfig.
	ROSARoleConfigReadyV1Beta2Condition = clusterv1.ReadyCondition
)

// ROSARoleConfig v1beta2 reason constants.
const (
	// ROSARoleConfigReadyV1Beta2Reason indicates the ROSARoleConfig is ready.
	ROSARoleConfigReadyV1Beta2Reason = clusterv1.ReadyReason

	// ROSARoleConfigNotReadyV1Beta2Reason indicates the ROSARoleConfig is not ready.
	ROSARoleConfigNotReadyV1Beta2Reason = clusterv1.NotReadyReason

	// ROSARoleConfigDeletingV1Beta2Reason indicates the ROSARoleConfig is being deleted.
	ROSARoleConfigDeletingV1Beta2Reason = clusterv1.DeletingReason

	// ROSARoleConfigReconciliationFailedV1Beta2Reason used to report reconciliation failures.
	ROSARoleConfigReconciliationFailedV1Beta2Reason = "ReconciliationFailed"

	// ROSARoleConfigDeletionFailedV1Beta2Reason used to report failures while deleting ROSARoleConfig.
	ROSARoleConfigDeletionFailedV1Beta2Reason = "DeletionFailed"

	// ROSARoleConfigDeletionStartedV1Beta2Reason used to indicate that the deletion of ROSARoleConfig has started.
	ROSARoleConfigDeletionStartedV1Beta2Reason = "DeletionStarted"

	// ROSARoleConfigCreatedV1Beta2Reason used to indicate that the ROSARoleConfig has been created.
	ROSARoleConfigCreatedV1Beta2Reason = "Created"
)

// ROSAOCMRoleConfig v1beta2 condition types.
const (
	// ROSAOCMRoleConfigReadyV1Beta2Condition defines the Ready condition type that summarizes the operational state of a ROSAOCMRoleConfig.
	ROSAOCMRoleConfigReadyV1Beta2Condition = clusterv1.ReadyCondition
)

// ROSAOCMRoleConfig v1beta2 reason constants.
const (
	// ROSAOCMRoleConfigReadyV1Beta2Reason indicates the ROSAOCMRoleConfig is ready.
	ROSAOCMRoleConfigReadyV1Beta2Reason = clusterv1.ReadyReason

	// ROSAOCMRoleConfigNotReadyV1Beta2Reason indicates the ROSAOCMRoleConfig is not ready.
	ROSAOCMRoleConfigNotReadyV1Beta2Reason = clusterv1.NotReadyReason

	// ROSAOCMRoleConfigDeletingV1Beta2Reason indicates the ROSAOCMRoleConfig is being deleted.
	ROSAOCMRoleConfigDeletingV1Beta2Reason = clusterv1.DeletingReason

	// ROSAOCMRoleConfigDeletionStartedV1Beta2Reason used to indicate that the deletion of ROSAOCMRoleConfig has started.
	ROSAOCMRoleConfigDeletionStartedV1Beta2Reason = "DeletionStarted"

	// ROSAOCMRoleConfigReconciliationFailedV1Beta2Reason used to report reconciliation failures.
	ROSAOCMRoleConfigReconciliationFailedV1Beta2Reason = "ReconciliationFailed"

	// ROSAOCMRoleConfigDeletionFailedV1Beta2Reason used to report failures while deleting ROSAOCMRoleConfig.
	ROSAOCMRoleConfigDeletionFailedV1Beta2Reason = "DeletionFailed"

	// ROSAOCMRoleConfigCreatedV1Beta2Reason used to indicate that the ROSAOCMRoleConfig has been created.
	ROSAOCMRoleConfigCreatedV1Beta2Reason = "Created"

	// ROSAOCMRoleConfigLinkedV1Beta2Reason used to indicate that the OCM role has been linked to the organization.
	ROSAOCMRoleConfigLinkedV1Beta2Reason = "Linked"
)
