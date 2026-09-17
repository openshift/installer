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

package v1beta3

// Conditions and condition Reasons for the CAPIBM v1beta2 object like IBMPowerVSCluster, IBMPowerVSMachine, IBMPowerVSImage.

import clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"

const (
	// ServiceInstanceReadyV1Beta2Condition reports on the successful reconciliation of a Power VS workspace.
	ServiceInstanceReadyV1Beta2Condition clusterv1.ConditionType = "ServiceInstanceReady"
	// ServiceInstanceReconciliationFailedV1Beta2Reason used when an error occurs during workspace reconciliation.
	ServiceInstanceReconciliationFailedV1Beta2Reason = "ServiceInstanceReconciliationFailed"

	// NetworkReadyV1Beta2Condition reports on the successful reconciliation of a Power VS network.
	NetworkReadyV1Beta2Condition clusterv1.ConditionType = "NetworkReady"
	// NetworkReconciliationFailedV1Beta2Reason used when an error occurs during network reconciliation.
	NetworkReconciliationFailedV1Beta2Reason = "NetworkReconciliationFailed"

	// VPCSecurityGroupReadyV1Beta2Condition reports on the successful reconciliation of a VPC.
	VPCSecurityGroupReadyV1Beta2Condition clusterv1.ConditionType = "VPCSecurityGroupReady"
	// VPCSecurityGroupReconciliationFailedV1Beta2Reason used when an error occurs during VPC reconciliation.
	VPCSecurityGroupReconciliationFailedV1Beta2Reason = "VPCSecurityGroupReconciliationFailed"

	// VPCReadyV1Beta2Condition reports on the successful reconciliation of a VPC.
	VPCReadyV1Beta2Condition clusterv1.ConditionType = "VPCReady"
	// VPCReconciliationFailedV1Beta2Reason used when an error occurs during VPC reconciliation.
	VPCReconciliationFailedV1Beta2Reason = "VPCReconciliationFailed"

	// VPCSubnetReadyV1Beta2Condition reports on the successful reconciliation of a VPC subnet.
	VPCSubnetReadyV1Beta2Condition clusterv1.ConditionType = "VPCSubnetReady"
	// VPCSubnetReconciliationFailedV1Beta2Reason used when an error occurs during VPC subnet reconciliation.
	VPCSubnetReconciliationFailedV1Beta2Reason = "VPCSubnetReconciliationFailed"

	// TransitGatewayReadyV1Beta2Condition reports on the successful reconciliation of a Power VS transit gateway.
	TransitGatewayReadyV1Beta2Condition clusterv1.ConditionType = "TransitGatewayReady"
	// TransitGatewayReconciliationFailedV1Beta2Reason used when an error occurs during transit gateway reconciliation.
	TransitGatewayReconciliationFailedV1Beta2Reason = "TransitGatewayReconciliationFailed"

	// LoadBalancerReadyV1Beta2Condition reports on the successful reconciliation of a Power VS network.
	LoadBalancerReadyV1Beta2Condition clusterv1.ConditionType = "LoadBalancerReady"
	// LoadBalancerReconciliationFailedV1Beta2Reason used when an error occurs during loadbalancer reconciliation.
	LoadBalancerReconciliationFailedV1Beta2Reason = "LoadBalancerReconciliationFailed"

	// COSInstanceReadyV1Beta2Condition reports on the successful reconciliation of a COS instance.
	COSInstanceReadyV1Beta2Condition clusterv1.ConditionType = "COSInstanceCreated"
	// COSInstanceReconciliationFailedV1Beta2Reason used when an error occurs during COS instance reconciliation.
	COSInstanceReconciliationFailedV1Beta2Reason = "COSInstanceCreationFailed"
)

// Power VS instance related conditions and corresponding reasons (virtual machines).
const (
	// InstanceReadyV1Beta2Condition documents the status of the instance that is controlled
	// by the IBMPowerVSMachine.
	InstanceReadyV1Beta2Condition clusterv1.ConditionType = "InstanceReady"

	// InstanceNotReadyV1Beta2Reason surfaces when the instance that is controlled
	// by the IBMPowerVSMachine is not ready.
	InstanceNotReadyV1Beta2Reason = "InstanceNotReady"

	// InstanceProvisionFailedV1Beta2Reason used for failures during instance provisioning.
	InstanceProvisionFailedV1Beta2Reason = "InstanceProvisionFailed"

	// InstanceStoppedV1Beta2Reason instance is in a stopped state.
	InstanceStoppedV1Beta2Reason = "InstanceStopped"

	// InstanceErroredV1Beta2Reason instance is in a errored state.
	InstanceErroredV1Beta2Reason = "InstanceErrored"

	// InstanceStateUnknownV1Beta2Reason used when the instance is in a unknown state.
	InstanceStateUnknownV1Beta2Reason = "InstanceStateUnknown"

	// InstanceWaitingForClusterInfrastructureReadyV1Beta2Reason documents the instance that is controller by
	// IBMPowerVSMachine waiting for the cluster infrastructure to be ready.
	InstanceWaitingForClusterInfrastructureReadyV1Beta2Reason = clusterv1.WaitingForClusterInfrastructureReadyReason

	// InstanceWaitingForControlPlaneInitializedV1Beta2Reason documents the instance that is controller by IBMPowerVSMachine waiting
	// for the control plane to be initialized.
	InstanceWaitingForControlPlaneInitializedV1Beta2Reason = clusterv1.WaitingForControlPlaneInitializedReason

	// InstanceWaitingForBootstrapDataV1Beta2Reason documents the instance that is controller by IBMPowerVSMachine waiting for the bootstrap
	// data to be ready.
	InstanceWaitingForBootstrapDataV1Beta2Reason = clusterv1.WaitingForBootstrapDataReason

	// InstanceLoadBalancerConfigurationFailedV1Beta2Reason surfaces when configuring the instance IP to load balancer fails.
	InstanceLoadBalancerConfigurationFailedV1Beta2Reason = "LoadBalancerConfigurationFailed"

	// InstanceWaitingForNetworkAddressV1Beta2Reason surfaces when the instance that is controlled
	// by the IBMPowerVSMachine waiting for the machine network settings to be reported after machine being powered on.
	InstanceWaitingForNetworkAddressV1Beta2Reason = "WaitingForNetworkAddress"

	// InstanceWaitingForImageV1Beta2Reason surfaces when the instance that is controlled
	// by the IBMPowerVSMachine waiting for the Power VS image to be available in workspace.
	InstanceWaitingForImageV1Beta2Reason = "WaitingForIBMImage"
)

// PowerVS Image related conditions and corresponding reasons.
const (
	// ImageReadyV1Beta2Condition reports on current status of the image. Ready indicates the image is in a active state.
	ImageReadyV1Beta2Condition clusterv1.ConditionType = "ImageReady"

	// ImageImportedV1Beta2Condition reports on current status of the image import job. Ready indicates the import job is finished.
	ImageImportedV1Beta2Condition clusterv1.ConditionType = "ImageImported"

	// ImageNotReadyV1Beta2Reason used when the image is not ready.
	ImageNotReadyV1Beta2Reason = "ImageNotReady"

	// ImageImportFailedV1Beta2Reason used when the image import is failed.
	ImageImportFailedV1Beta2Reason = "ImageImportFailed"

	// ImageStateUnknownV1Beta2Reason used when the image is in an unknown state.
	ImageStateUnknownV1Beta2Reason = "ImageStateUnknown"
)

const (
	// DeletingV1Beta2Reason surfaces when an object is deleting because the
	// DeletionTimestamp is set. This reason is used if none of the more specific reasons apply.
	DeletingV1Beta2Reason = "Deleting"

	// InternalErrorV1Beta2Reason surfaces unexpected errors reporting by controllers.
	// In most cases, it will be required to look at controllers logs to properly triage those issues.
	InternalErrorV1Beta2Reason = "InternalError"
)
