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

// AWSManagedControlPlane v1beta2 condition types.
const (
	// AWSManagedControlPlaneReadyV1Beta2Condition defines the Ready condition type that summarizes the operational state of an AWSManagedControlPlane.
	AWSManagedControlPlaneReadyV1Beta2Condition = clusterv1.ReadyCondition

	// AWSManagedControlPlaneEKSControlPlaneReadyV1Beta2Condition reports on the successful reconciliation of the EKS control plane.
	AWSManagedControlPlaneEKSControlPlaneReadyV1Beta2Condition = "EKSControlPlaneReady"

	// AWSManagedControlPlaneEKSControlPlaneCreatingV1Beta2Condition reports whether the EKS control plane is being created.
	AWSManagedControlPlaneEKSControlPlaneCreatingV1Beta2Condition = "EKSControlPlaneCreating"

	// AWSManagedControlPlaneEKSControlPlaneUpdatingV1Beta2Condition reports whether the EKS control plane is being updated.
	AWSManagedControlPlaneEKSControlPlaneUpdatingV1Beta2Condition = "EKSControlPlaneUpdating"

	// AWSManagedControlPlaneIAMControlPlaneRolesReadyV1Beta2Condition reports on the successful reconciliation of EKS control plane IAM roles.
	AWSManagedControlPlaneIAMControlPlaneRolesReadyV1Beta2Condition = "IAMControlPlaneRolesReady"

	// AWSManagedControlPlaneIAMAuthenticatorConfiguredV1Beta2Condition reports on the successful reconciliation of the aws-iam-authenticator config.
	AWSManagedControlPlaneIAMAuthenticatorConfiguredV1Beta2Condition = "IAMAuthenticatorConfigured"

	// AWSManagedControlPlaneEKSAddonsConfiguredV1Beta2Condition reports on the successful reconciliation of EKS addons.
	AWSManagedControlPlaneEKSAddonsConfiguredV1Beta2Condition = "EKSAddonsConfigured"

	// AWSManagedControlPlaneEKSIdentityProviderConfiguredV1Beta2Condition reports on the successful association of identity provider config.
	AWSManagedControlPlaneEKSIdentityProviderConfiguredV1Beta2Condition = "EKSIdentityProviderConfigured"

	// AWSManagedControlPlaneEKSPodIdentityAssociationConfiguredV1Beta2Condition reports on the successful reconciliation of EKS pod identity associations.
	AWSManagedControlPlaneEKSPodIdentityAssociationConfiguredV1Beta2Condition = "EKSPodIdentityAssociationConfigured"
)

// AWSManagedControlPlane v1beta2 reason constants.
const (
	// AWSManagedControlPlaneReadyV1Beta2Reason indicates the AWSManagedControlPlane is ready.
	AWSManagedControlPlaneReadyV1Beta2Reason = clusterv1.ReadyReason

	// AWSManagedControlPlaneNotReadyV1Beta2Reason indicates the AWSManagedControlPlane is not ready.
	AWSManagedControlPlaneNotReadyV1Beta2Reason = clusterv1.NotReadyReason

	// AWSManagedControlPlaneDeletingV1Beta2Reason indicates the AWSManagedControlPlane is being deleted.
	AWSManagedControlPlaneDeletingV1Beta2Reason = clusterv1.DeletingReason

	// AWSManagedControlPlaneEKSControlPlaneReconciliationFailedV1Beta2Reason used to report failures while reconciling EKS control plane.
	AWSManagedControlPlaneEKSControlPlaneReconciliationFailedV1Beta2Reason = "EKSControlPlaneReconciliationFailed"

	// AWSManagedControlPlaneIAMControlPlaneRolesReconciliationFailedV1Beta2Reason used to report failures while reconciling EKS control plane IAM roles.
	AWSManagedControlPlaneIAMControlPlaneRolesReconciliationFailedV1Beta2Reason = "IAMControlPlaneRolesReconciliationFailed"

	// AWSManagedControlPlaneIAMAuthenticatorConfigurationFailedV1Beta2Reason used to report failures while reconciling the aws-iam-authenticator config.
	AWSManagedControlPlaneIAMAuthenticatorConfigurationFailedV1Beta2Reason = "IAMAuthenticatorConfigurationFailed"

	// AWSManagedControlPlaneEKSAddonsConfiguredFailedV1Beta2Reason used to report failures while reconciling the EKS addons.
	AWSManagedControlPlaneEKSAddonsConfiguredFailedV1Beta2Reason = "EKSAddonsConfiguredFailed"

	// AWSManagedControlPlaneEKSIdentityProviderConfiguredFailedV1Beta2Reason used to report failures while reconciling the identity provider config association.
	AWSManagedControlPlaneEKSIdentityProviderConfiguredFailedV1Beta2Reason = "EKSIdentityProviderConfiguredFailed"

	// AWSManagedControlPlaneEKSPodIdentityAssociationConfigurationFailedV1Beta2Reason used to report failures while reconciling the EKS pod identity associations.
	AWSManagedControlPlaneEKSPodIdentityAssociationConfigurationFailedV1Beta2Reason = "EKSPodIdentityAssociationConfigurationFailed"
)
