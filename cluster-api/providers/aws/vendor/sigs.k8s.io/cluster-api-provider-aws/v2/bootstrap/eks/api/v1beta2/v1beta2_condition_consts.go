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

// EKSConfig v1beta2 condition types.
const (
	// EKSConfigReadyV1Beta2Condition defines the Ready condition type that summarizes the operational state of an EKSConfig.
	EKSConfigReadyV1Beta2Condition = clusterv1.ReadyCondition

	// EKSConfigDataSecretAvailableV1Beta2Condition reports on the status of the bootstrap secret generation process.
	//
	// NOTE: When the DataSecret generation starts the process completes immediately and within the
	// same reconciliation, so the user will always see a transition from Wait to Generated without having
	// evidence that BootstrapSecret generation is started/in progress.
	EKSConfigDataSecretAvailableV1Beta2Condition = "DataSecretAvailable"
)

// EKSConfig v1beta2 reason constants.
const (
	// EKSConfigReadyV1Beta2Reason indicates the EKSConfig is ready.
	EKSConfigReadyV1Beta2Reason = clusterv1.ReadyReason

	// EKSConfigNotReadyV1Beta2Reason indicates the EKSConfig is not ready.
	EKSConfigNotReadyV1Beta2Reason = clusterv1.NotReadyReason

	// EKSConfigDeletingV1Beta2Reason indicates the EKSConfig is being deleted.
	EKSConfigDeletingV1Beta2Reason = clusterv1.DeletingReason

	// EKSConfigDataSecretGenerationFailedV1Beta2Reason indicates an error while generating a data secret;
	// those kind of errors are usually due to misconfigurations and user intervention is required to get them fixed.
	EKSConfigDataSecretGenerationFailedV1Beta2Reason = "DataSecretGenerationFailed"

	// EKSConfigWaitingForClusterInfrastructureV1Beta2Reason indicates waiting for cluster infrastructure to be ready.
	//
	// NOTE: Having the cluster infrastructure ready is a pre-condition for starting to create machines;
	// the EKSConfig controller ensures this pre-condition is satisfied.
	EKSConfigWaitingForClusterInfrastructureV1Beta2Reason = clusterv1.WaitingForClusterInfrastructureReadyReason

	// EKSConfigWaitingForControlPlaneInitializationV1Beta2Reason indicates waiting for the control plane to be initialized.
	//
	// NOTE: This is a pre-condition for starting to create machines;
	// the EKSConfig controller ensures this pre-condition is satisfied.
	EKSConfigWaitingForControlPlaneInitializationV1Beta2Reason = clusterv1.WaitingForControlPlaneInitializedReason
)

// NodeadmConfig v1beta2 condition types.
const (
	// NodeadmConfigReadyV1Beta2Condition defines the Ready condition type that summarizes the operational state of a NodeadmConfig.
	NodeadmConfigReadyV1Beta2Condition = clusterv1.ReadyCondition

	// NodeadmConfigDataSecretAvailableV1Beta2Condition reports on the status of the bootstrap secret generation process.
	//
	// NOTE: When the DataSecret generation starts the process completes immediately and within the
	// same reconciliation, so the user will always see a transition from Wait to Generated without having
	// evidence that BootstrapSecret generation is started/in progress.
	NodeadmConfigDataSecretAvailableV1Beta2Condition = "DataSecretAvailable"
)

// NodeadmConfig v1beta2 reason constants.
const (
	// NodeadmConfigReadyV1Beta2Reason indicates the NodeadmConfig is ready.
	NodeadmConfigReadyV1Beta2Reason = clusterv1.ReadyReason

	// NodeadmConfigNotReadyV1Beta2Reason indicates the NodeadmConfig is not ready.
	NodeadmConfigNotReadyV1Beta2Reason = clusterv1.NotReadyReason

	// NodeadmConfigDeletingV1Beta2Reason indicates the NodeadmConfig is being deleted.
	NodeadmConfigDeletingV1Beta2Reason = clusterv1.DeletingReason

	// NodeadmConfigDataSecretGenerationFailedV1Beta2Reason indicates an error while generating a data secret;
	// those kind of errors are usually due to misconfigurations and user intervention is required to get them fixed.
	NodeadmConfigDataSecretGenerationFailedV1Beta2Reason = "DataSecretGenerationFailed"

	// NodeadmConfigWaitingForClusterInfrastructureV1Beta2Reason indicates waiting for cluster infrastructure to be ready.
	//
	// NOTE: Having the cluster infrastructure ready is a pre-condition for starting to create machines;
	// the NodeadmConfig controller ensures this pre-condition is satisfied.
	NodeadmConfigWaitingForClusterInfrastructureV1Beta2Reason = clusterv1.WaitingForClusterInfrastructureReadyReason

	// NodeadmConfigWaitingForControlPlaneInitializationV1Beta2Reason indicates waiting for the control plane to be initialized.
	//
	// NOTE: This is a pre-condition for starting to create machines;
	// the NodeadmConfig controller ensures this pre-condition is satisfied.
	NodeadmConfigWaitingForControlPlaneInitializationV1Beta2Reason = clusterv1.WaitingForControlPlaneInitializedReason
)
