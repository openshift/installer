/*
Copyright 2022 The Kubernetes Authors.

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

import clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

const (
	// ROSAControlPlaneReadyCondition condition reports on the successful reconciliation of ROSAControlPlane.
	ROSAControlPlaneReadyCondition clusterv1.ConditionType = "ROSAControlPlaneReady"

	// ROSAControlPlaneValidCondition condition reports whether ROSAControlPlane configuration is valid.
	ROSAControlPlaneValidCondition clusterv1.ConditionType = "ROSAControlPlaneValid"

	// ROSAControlPlaneUpgradingCondition condition reports whether ROSAControlPlane is upgrading or not.
	ROSAControlPlaneUpgradingCondition clusterv1.ConditionType = "ROSAControlPlaneUpgrading"

	// ExternalAuthConfiguredCondition condition reports whether external auth has beed correctly configured.
	ExternalAuthConfiguredCondition clusterv1.ConditionType = "ExternalAuthConfigured"

	// ROSARoleConfigReadyCondition condition reports whether the referenced RosaRoleConfig is ready.
	ROSARoleConfigReadyCondition clusterv1.ConditionType = "ROSARoleConfigReady"

	// ReconciliationFailedReason used to report reconciliation failures.
	ReconciliationFailedReason = "ReconciliationFailed"

	// ROSAControlPlaneDeletionFailedReason used to report failures while deleting ROSAControlPlane.
	ROSAControlPlaneDeletionFailedReason = "DeletionFailed"

	// ROSAControlPlaneInvalidConfigurationReason used to report invalid user input.
	ROSAControlPlaneInvalidConfigurationReason = "InvalidConfiguration"

	// ROSARoleConfigNotReadyReason used to report when referenced RosaRoleConfig is not ready.
	ROSARoleConfigNotReadyReason = "ROSARoleConfigNotReady"

	// ROSARoleConfigNotFoundReason used to report when referenced RosaRoleConfig is not found.
	ROSARoleConfigNotFoundReason = "ROSARoleConfigNotFound"
)
