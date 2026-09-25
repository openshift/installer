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

import (
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
)

const (
	// CreateInfrastructureAnnotation is the name of an annotation that indicates if
	// Power VS infrastructure should be created as a part of cluster creation.
	CreateInfrastructureAnnotation = "powervs.cluster.x-k8s.io/create-infra"
)

// IBMPowerVSMachine's condition and corresponding reasons.
const (
	// IBMPowerVSMachineReadyCondition is true if the IBMPowerVSMachine's deletionTimestamp is not set, IBMPowerVSMachine's
	// IBMPowerVSMachineInstanceReadyCondition is true.
	IBMPowerVSMachineReadyCondition = clusterv1.ReadyCondition

	// IBMPowerVSMachineReadyReason surfaces when the IBMPowerVSMachine readiness criteria is met.
	IBMPowerVSMachineReadyReason = clusterv1.ReadyReason

	// IBMPowerVSMachineNotReadyReason surfaces when the IBMPowerVSMachine readiness criteria is not met.
	IBMPowerVSMachineNotReadyReason = clusterv1.NotReadyReason

	// IBMPowerVSMachineReadyUnknownReason surfaces when at least one IBMPowerVSMachine readiness criteria is unknown
	// and IBMPowerVSMachine readiness criteria is not met.
	IBMPowerVSMachineReadyUnknownReason = clusterv1.ReadyUnknownReason
)

// PowerVS instance related conditions and corresponding reasons (virtual machines).
const (
	// InstanceReadyCondition documents the status of the instance that is controlled
	// by the IBMPowerVSMachine.
	InstanceReadyCondition = "InstanceReady"

	// InstanceReadyReason surfaces when the instance that is controlled
	// by the IBMPowerVSMachine is ready.
	InstanceReadyReason = "InstanceReady"

	// InstanceNotReadyReason surfaces when the instance that is controlled
	// by the IBMPowerVSMachine is not ready.
	InstanceNotReadyReason = "InstanceNotReady"

	// InstanceProvisionFailedReason used for failures during instance provisioning.
	InstanceProvisionFailedReason = "InstanceProvisionFailed"

	// InstanceStoppedReason instance is in a stopped state.
	InstanceStoppedReason = "InstanceStopped"

	// InstanceErroredReason instance is in a errored state.
	InstanceErroredReason = "InstanceErrored"

	// InstanceStateUnknownReason used when the instance is in a unknown state.
	InstanceStateUnknownReason = "InstanceStateUnknown"

	// InstanceWaitingForClusterInfrastructureReadyReason documents the instance that is controller by
	// IBMPowerVSMachine waiting for the cluster infrastructure to be ready.
	InstanceWaitingForClusterInfrastructureReadyReason = clusterv1.WaitingForClusterInfrastructureReadyReason

	// InstanceWaitingForControlPlaneInitializedReason documents the instance that is controller by IBMPowerVSMachine waiting
	// for the control plane to be initialized.
	InstanceWaitingForControlPlaneInitializedReason = clusterv1.WaitingForControlPlaneInitializedReason

	// InstanceWaitingForBootstrapDataReason documents the instance that is controller by IBMPowerVSMachine waiting for the bootstrap
	// data to be ready.
	InstanceWaitingForBootstrapDataReason = clusterv1.WaitingForBootstrapDataReason

	// InstanceDeletingReason surfaces when the instance controller by IBMPowerVSMachine is deleting.
	InstanceDeletingReason = clusterv1.DeletingReason

	// InstanceLoadBalancerConfigurationFailedReason surfaces when configuring the instance IP to load balancer fails.
	InstanceLoadBalancerConfigurationFailedReason = "LoadBalancerConfigurationFailed"

	// InstanceWaitingForNetworkAddressReason surfaces when the instance that is controlled
	// by the IBMPowerVSMachine waiting for the machine network settings to be reported after machine being powered on.
	InstanceWaitingForNetworkAddressReason = "WaitingForNetworkAddress"

	// InstanceWaitingForImageReason surfaces when the instance that is controlled
	// by the IBMPowerVSMachine waiting for the Power VS image to be available in workspace.
	InstanceWaitingForImageReason = "WaitingForIBMImage"

	// InvalidMachineConfigurationReason used when the machine configuration is invalid.
	InvalidMachineConfigurationReason = "InvalidMachineConfiguration"
)

// IBMPowerVSImage's Ready condition and corresponding reasons.
const (
	// IBMPowerVSImageReadyCondition is true if the IBMPowerVSImage's deletionTimestamp is not set, IBMPowerVSImage's IBMPowerVSImageReadyCondition is true.
	IBMPowerVSImageReadyCondition = clusterv1.ReadyCondition

	// IBMPowerVSImageReadyReason surfaces when the IBMPowerVSImage readiness criteria is met.
	IBMPowerVSImageReadyReason = clusterv1.ReadyReason

	// IBMPowerVSImageNotReadyReason surfaces when the IBMPowerVSImage readiness criteria is not met.
	IBMPowerVSImageNotReadyReason = clusterv1.NotReadyReason

	// IBMPowerVSImageReadyUnknownReason surfaces when at least one of the IBMPowerVSImage readiness criteria is unknown
	// and none of the IBMPowerVSImage readiness criteria is met.
	IBMPowerVSImageReadyUnknownReason = clusterv1.ReadyUnknownReason

	// IBMPowerVSImageImportFailedReason used when the image import is failed.
	IBMPowerVSImageImportFailedReason = "ImageImportFailed"
)

const (
	// ServiceInstanceReadyCondition reports on the successful reconciliation of a Power VS workspace.
	ServiceInstanceReadyCondition = "ServiceInstanceReady"
	// ServiceInstanceReconciliationFailedReason used when an error occurs during workspace reconciliation.
	ServiceInstanceReconciliationFailedReason = "ServiceInstanceReconciliationFailed"

	// NetworkReadyCondition reports on the successful reconciliation of a Power VS network.
	NetworkReadyCondition = "NetworkReady"
	// NetworkReconciliationFailedReason used when an error occurs during network reconciliation.
	NetworkReconciliationFailedReason = "NetworkReconciliationFailed"

	// VPCSecurityGroupReadyCondition reports on the successful reconciliation of a VPC.
	VPCSecurityGroupReadyCondition = "VPCSecurityGroupReady"
	// VPCSecurityGroupReconciliationFailedReason used when an error occurs during VPC reconciliation.
	VPCSecurityGroupReconciliationFailedReason = "VPCSecurityGroupReconciliationFailed"

	// VPCReadyCondition reports on the successful reconciliation of a VPC.
	VPCReadyCondition = "VPCReady"
	// VPCReconciliationFailedReason used when an error occurs during VPC reconciliation.
	VPCReconciliationFailedReason = "VPCReconciliationFailed"

	// VPCSubnetReadyCondition reports on the successful reconciliation of a VPC subnet.
	VPCSubnetReadyCondition = "VPCSubnetReady"
	// VPCSubnetReconciliationFailedReason used when an error occurs during VPC subnet reconciliation.
	VPCSubnetReconciliationFailedReason = "VPCSubnetReconciliationFailed"

	// TransitGatewayReadyCondition reports on the successful reconciliation of a Power VS transit gateway.
	TransitGatewayReadyCondition = "TransitGatewayReady"
	// TransitGatewayReconciliationFailedReason used when an error occurs during transit gateway reconciliation.
	TransitGatewayReconciliationFailedReason = "TransitGatewayReconciliationFailed"

	// LoadBalancerReadyCondition reports on the successful reconciliation of a Power VS network.
	LoadBalancerReadyCondition = "LoadBalancerReady"
	// LoadBalancerReconciliationFailedReason used when an error occurs during loadbalancer reconciliation.
	LoadBalancerReconciliationFailedReason = "LoadBalancerReconciliationFailed"

	// COSInstanceReadyCondition reports on the successful reconciliation of a COS instance.
	COSInstanceReadyCondition = "COSInstanceCreated"
	// COSInstanceReconciliationFailedReason used when an error occurs during COS instance reconciliation.
	COSInstanceReconciliationFailedReason = "COSInstanceCreationFailed"
)

// IBMPowerVSCluster's Ready condition and corresponding reasons.
const (
	// IBMPowerVSClusterReadyCondition is true if the IBMPowerVSCluster's deletionTimestamp is not set, IBMPowerVSCluster's
	// FailureDomainsReady, VCenterAvailable and ClusterModulesReady conditions are true.
	IBMPowerVSClusterReadyCondition = clusterv1.ReadyCondition

	// IBMPowerVSClusterReadyReason surfaces when the IBMPowerVSCluster readiness criteria is met.
	IBMPowerVSClusterReadyReason = clusterv1.ReadyReason

	// IBMPowerVSClusterNotReadyReason surfaces when the IBMPowerVSCluster readiness criteria is not met.
	IBMPowerVSClusterNotReadyReason = clusterv1.NotReadyReason

	// IBMPowerVSImageDeletingReason surfaces when the image is in deleting state.
	IBMPowerVSImageDeletingReason = clusterv1.DeletingReason

	// IBMPowerVSClusterReadyUnknownReason surfaces when at least one of the IBMPowerVSCluster readiness criteria is unknown
	// and none of the IBMPowerVSCluster readiness criteria is met.
	IBMPowerVSClusterReadyUnknownReason = clusterv1.ReadyUnknownReason

	// IBMPowerVSImageQueuedReason used when the image is in queued state.
	IBMPowerVSImageQueuedReason = "ImageQueued"
)

const (
	// WorkspaceReadyCondition reports on the successful reconciliation of a PowerVS workspace.
	WorkspaceReadyCondition = "WorkspaceReady"

	// WorkspaceReadyReason surfaces when the PowerVS workspace is ready.
	WorkspaceReadyReason = clusterv1.ReadyReason

	// WorkspaceNotReadyReason surfaces when PowerVS workspace is not ready.
	WorkspaceNotReadyReason = clusterv1.NotReadyReason

	// WorkspaceDeletingReason surfaces when the PowerVS workspace is being deleted.
	WorkspaceDeletingReason = clusterv1.DeletingReason

	// NetworkReadyReason surfaces when PowerVS workspace is ready.
	NetworkReadyReason = clusterv1.ReadyReason

	// NetworkNotReadyReason surfaces when the PowerVS network is not ready.
	NetworkNotReadyReason = clusterv1.NotReadyReason

	// NetworkDeletingReason surfaces when the PowerVS network is being deleted.
	NetworkDeletingReason = clusterv1.DeletingReason

	// VPCReadyReason surfaces when the VPC is ready.
	VPCReadyReason = clusterv1.ReadyReason

	// VPCNotReadyReason surfaces when VPC is not ready.
	VPCNotReadyReason = clusterv1.NotReadyReason

	// VPCDeletingReason surfaces when the VPC is being deleted.
	VPCDeletingReason = clusterv1.DeletingReason

	// VPCSubnetReadyReason surfaces when the VPC subnet is ready.
	VPCSubnetReadyReason = clusterv1.ReadyReason

	// VPCSubnetNotReadyReason surfaces when VPC subnet is not ready.
	VPCSubnetNotReadyReason = clusterv1.NotReadyReason

	// VPCSubnetDeletingReason surfaces when the VPC subnet is being deleted.
	VPCSubnetDeletingReason = clusterv1.DeletingReason

	// VPCSecurityGroupReadyReason surfaces when the VPC security group is ready.
	VPCSecurityGroupReadyReason = clusterv1.ReadyReason

	// VPCSecurityGroupNotReadyCondition surfaces when VPC security group is not ready.
	VPCSecurityGroupNotReadyCondition = clusterv1.NotReadyReason

	// VPCSecurityGroupDeletingReason surfaces when the VPC security group is being deleted.
	VPCSecurityGroupDeletingReason = clusterv1.DeletingReason

	// TransitGatewayReadyReason surfaces when the transit gateway is ready.
	TransitGatewayReadyReason = clusterv1.ReadyReason

	// TransitGatewayNotReadyReason surfaces when the transit gateway is not ready.
	TransitGatewayNotReadyReason = clusterv1.NotReadyReason

	// TransitGatewayDeletingReason surfaces when the transit gateway is being deleted.
	TransitGatewayDeletingReason = clusterv1.DeletingReason

	// VPCLoadBalancerReadyCondition reports on the successful reconciliation of a VPC LoadBalancer.
	VPCLoadBalancerReadyCondition = "LoadBalancerReady"

	// VPCLoadBalancerReadyReason surfaces when the VPC LoadBalancer is ready.
	VPCLoadBalancerReadyReason = clusterv1.ReadyReason

	// VPCLoadBalancerNotReadyReason surfaces when VPC LoadBalancer is not ready.
	VPCLoadBalancerNotReadyReason = clusterv1.NotReadyReason

	// VPCLoadBalancerDeletingReason surfaces when the VPC LoadBalancer is being deleted.
	VPCLoadBalancerDeletingReason = clusterv1.DeletingReason

	// COSInstanceReadyReason surfaces when the COS instance is ready.
	COSInstanceReadyReason = clusterv1.ReadyReason

	// COSInstanceNotReadyReason surfaces when the COS instance is not ready.
	COSInstanceNotReadyReason = clusterv1.NotReadyReason

	// COSInstanceDeletingReason surfaces when the COS instance is being deleted.
	COSInstanceDeletingReason = clusterv1.DeletingReason
)
