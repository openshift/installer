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

import (
	capiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const (
	// InstanceProvisionFailedReason used for failures during instance provisioning.
	InstanceProvisionFailedReason = "InstanceProvisionFailed"
	// WaitingForClusterInfrastructureReason used when machine is waiting for cluster infrastructure to be ready before proceeding.
	WaitingForClusterInfrastructureReason = "WaitingForClusterInfrastructure"
	// WaitingForBootstrapDataReason used when machine is waiting for bootstrap data to be ready before proceeding.
	WaitingForBootstrapDataReason = "WaitingForBootstrapData"
)

const (
	// InstanceStoppedReason instance is in a stopped state.
	InstanceStoppedReason = "InstanceStopped"

	// InstanceErroredReason instance is in a errored state.
	InstanceErroredReason = "InstanceErrored"

	// InstanceNotReadyReason used when the instance is in a not ready state.
	InstanceNotReadyReason = "InstanceNotReady"

	// InstanceStateUnknownReason used when the instance is in a unknown state.
	InstanceStateUnknownReason = "InstanceStateUnknown"
)

const (
	// InstanceReadyCondition reports on current status of the instance. Ready indicates the instance is in a Running state.
	InstanceReadyCondition capiv1beta1.ConditionType = "InstanceReady"
)

const (
	// WaitingForIBMPowerVSImageReason used when machine is waiting for powervs image to be ready before proceeding.
	WaitingForIBMPowerVSImageReason = "WaitingForIBMPowerVSImage"
)

const (
	// ImageNotReadyReason used when the image is in a queued state.
	ImageNotReadyReason = "ImageNotReady"

	// ImageImportFailedReason used when the image import is failed.
	ImageImportFailedReason = "ImageImportFailed"

	// ImageReconciliationFailedReason used when an error occurs during VPC Custom Image reconciliation.
	ImageReconciliationFailedReason = "ImageReconciliationFailed"
)

const (
	// ImageReadyCondition reports on current status of the image. Ready indicates the image is in a active state.
	ImageReadyCondition capiv1beta1.ConditionType = "ImageReady"

	// ImageImportedCondition reports on current status of the image import job. Ready indicates the import job is finished.
	ImageImportedCondition capiv1beta1.ConditionType = "ImageImported"
)

const (
	// LoadBalancerNotReadyReason used when cluster is waiting for load balancer to be ready before proceeding.
	LoadBalancerNotReadyReason = "LoadBalancerNotReady"
)

const (
	// ServiceInstanceReadyCondition reports on the successful reconciliation of a Power VS workspace.
	ServiceInstanceReadyCondition capiv1beta1.ConditionType = "ServiceInstanceReady"
	// ServiceInstanceReconciliationFailedReason used when an error occurs during workspace reconciliation.
	ServiceInstanceReconciliationFailedReason = "ServiceInstanceReconciliationFailed"

	// NetworkReadyCondition reports on the successful reconciliation of a Power VS network.
	NetworkReadyCondition capiv1beta1.ConditionType = "NetworkReady"
	// NetworkReconciliationFailedReason used when an error occurs during network reconciliation.
	NetworkReconciliationFailedReason = "NetworkReconciliationFailed"

	// VPCSecurityGroupReadyCondition reports on the successful reconciliation of a VPC.
	VPCSecurityGroupReadyCondition capiv1beta1.ConditionType = "VPCSecurityGroupReady"
	// VPCSecurityGroupReconciliationFailedReason used when an error occurs during VPC reconciliation.
	VPCSecurityGroupReconciliationFailedReason = "VPCSecurityGroupReconciliationFailed"

	// VPCReadyCondition reports on the successful reconciliation of a VPC.
	VPCReadyCondition capiv1beta1.ConditionType = "VPCReady"
	// VPCReconciliationFailedReason used when an error occurs during VPC reconciliation.
	VPCReconciliationFailedReason = "VPCReconciliationFailed"

	// VPCSubnetReadyCondition reports on the successful reconciliation of a VPC subnet.
	VPCSubnetReadyCondition capiv1beta1.ConditionType = "VPCSubnetReady"
	// VPCSubnetReconciliationFailedReason used when an error occurs during VPC subnet reconciliation.
	VPCSubnetReconciliationFailedReason = "VPCSubnetReconciliationFailed"

	// TransitGatewayReadyCondition reports on the successful reconciliation of a Power VS transit gateway.
	TransitGatewayReadyCondition capiv1beta1.ConditionType = "TransitGatewayReady"
	// TransitGatewayReconciliationFailedReason used when an error occurs during transit gateway reconciliation.
	TransitGatewayReconciliationFailedReason = "TransitGatewayReconciliationFailed"

	// LoadBalancerReadyCondition reports on the successful reconciliation of a Power VS network.
	LoadBalancerReadyCondition capiv1beta1.ConditionType = "LoadBalancerReady"
	// LoadBalancerReconciliationFailedReason used when an error occurs during loadbalancer reconciliation.
	LoadBalancerReconciliationFailedReason = "LoadBalancerReconciliationFailed"

	// COSInstanceReadyCondition reports on the successful reconciliation of a COS instance.
	COSInstanceReadyCondition capiv1beta1.ConditionType = "COSInstanceCreated"
	// COSInstanceReconciliationFailedReason used when an error occurs during COS instance reconciliation.
	COSInstanceReconciliationFailedReason = "COSInstanceCreationFailed"
)

const (
	// CreateInfrastructureAnnotation is the name of an annotation that indicates if
	// Power VS infrastructure should be created as a part of cluster creation.
	CreateInfrastructureAnnotation = "powervs.cluster.x-k8s.io/create-infra"
)
