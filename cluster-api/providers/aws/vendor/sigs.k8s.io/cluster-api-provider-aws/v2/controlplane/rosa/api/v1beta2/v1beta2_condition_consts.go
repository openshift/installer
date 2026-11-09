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

// ROSAControlPlane v1beta2 condition types.
const (
	// ROSAControlPlaneReadyV1Beta2Condition defines the Ready condition type that summarizes the operational state of a ROSAControlPlane.
	ROSAControlPlaneReadyV1Beta2Condition = clusterv1.ReadyCondition

	// ROSAControlPlaneValidV1Beta2Condition reports whether ROSAControlPlane configuration is valid.
	ROSAControlPlaneValidV1Beta2Condition = "ROSAControlPlaneValid"

	// ROSAControlPlaneUpgradingV1Beta2Condition reports whether ROSAControlPlane is upgrading or not.
	ROSAControlPlaneUpgradingV1Beta2Condition = "ROSAControlPlaneUpgrading"

	// ROSAControlPlaneExternalAuthConfiguredV1Beta2Condition reports whether external auth has been correctly configured.
	ROSAControlPlaneExternalAuthConfiguredV1Beta2Condition = "ExternalAuthConfigured"

	// ROSAControlPlaneRoleConfigReadyV1Beta2Condition reports whether the referenced RosaRoleConfig is ready.
	ROSAControlPlaneRoleConfigReadyV1Beta2Condition = "ROSARoleConfigReady"
)

// ROSAControlPlane v1beta2 reason constants.
const (
	// ROSAControlPlaneReadyV1Beta2Reason indicates the ROSAControlPlane is ready.
	ROSAControlPlaneReadyV1Beta2Reason = clusterv1.ReadyReason

	// ROSAControlPlaneNotReadyV1Beta2Reason indicates the ROSAControlPlane is not ready.
	ROSAControlPlaneNotReadyV1Beta2Reason = clusterv1.NotReadyReason

	// ROSAControlPlaneDeletingV1Beta2Reason indicates the ROSAControlPlane is being deleted.
	ROSAControlPlaneDeletingV1Beta2Reason = clusterv1.DeletingReason

	// ROSAControlPlaneReconciliationFailedV1Beta2Reason used to report reconciliation failures.
	ROSAControlPlaneReconciliationFailedV1Beta2Reason = "ReconciliationFailed"

	// ROSAControlPlaneDeletionFailedV1Beta2Reason used to report failures while deleting ROSAControlPlane.
	ROSAControlPlaneDeletionFailedV1Beta2Reason = "DeletionFailed"

	// ROSAControlPlaneInvalidConfigurationV1Beta2Reason used to report invalid user input.
	ROSAControlPlaneInvalidConfigurationV1Beta2Reason = "InvalidConfiguration"

	// ROSAControlPlaneRoleConfigNotReadyV1Beta2Reason used to report when referenced RosaRoleConfig is not ready.
	ROSAControlPlaneRoleConfigNotReadyV1Beta2Reason = "ROSARoleConfigNotReady"

	// ROSAControlPlaneRoleConfigNotFoundV1Beta2Reason used to report when referenced RosaRoleConfig is not found.
	ROSAControlPlaneRoleConfigNotFoundV1Beta2Reason = "ROSARoleConfigNotFound"
)
