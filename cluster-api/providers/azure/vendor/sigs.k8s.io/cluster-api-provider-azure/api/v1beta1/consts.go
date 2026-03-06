/*
Copyright 2021 The Kubernetes Authors.

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

import clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"

// AzureCluster Conditions and Reasons.
const (
	// NetworkInfrastructureReadyCondition reports of current status of cluster infrastructure.
	NetworkInfrastructureReadyCondition clusterv1beta1.ConditionType = "NetworkInfrastructureReady"
	// NamespaceNotAllowedByIdentity used to indicate cluster in a namespace not allowed by identity.
	NamespaceNotAllowedByIdentity = "NamespaceNotAllowedByIdentity"
)

// AzureMachine Conditions and Reasons.
const (
	// VMRunningCondition reports on current status of the Azure VM.
	VMRunningCondition clusterv1beta1.ConditionType = "VMRunning"
	// VMIdentitiesReadyCondition reports on the readiness of the Azure VM identities.
	VMIdentitiesReadyCondition clusterv1beta1.ConditionType = "VMIdentitiesReady"
	// VMCreatingReason used when the vm creation is in progress.
	VMCreatingReason = "VMCreating"
	// VMUpdatingReason used when the vm updating is in progress.
	VMUpdatingReason = "VMUpdating"
	// VMDeletingReason used when the vm is in a deleting state.
	VMDeletingReason = "VMDeleting"
	// VMProvisionFailedReason used for failures during vm provisioning.
	VMProvisionFailedReason = "VMProvisionFailed"
	// UserAssignedIdentityMissingReason used for failures when a user-assigned identity is missing.
	UserAssignedIdentityMissingReason = "UserAssignedIdentityMissing"
	// WaitingForClusterInfrastructureReason used when machine is waiting for cluster infrastructure to be ready before proceeding.
	WaitingForClusterInfrastructureReason = "WaitingForClusterInfrastructure"
	// WaitingForBootstrapDataReason used when machine is waiting for bootstrap data to be ready before proceeding.
	WaitingForBootstrapDataReason = "WaitingForBootstrapData"
	// BootstrapSucceededCondition reports the result of the execution of the bootstrap data on the machine.
	BootstrapSucceededCondition clusterv1beta1.ConditionType = "BootstrapSucceeded"
	// BootstrapInProgressReason is used to indicate the bootstrap data has not finished executing.
	BootstrapInProgressReason = "BootstrapInProgress"
	// BootstrapFailedReason is used to indicate the bootstrap process ran into an error.
	BootstrapFailedReason = "BootstrapFailed"
)

// AzureMachinePool Conditions and Reasons.
const (
	// ScaleSetRunningCondition reports on current status of the Azure Scale Set.
	ScaleSetRunningCondition clusterv1beta1.ConditionType = "ScaleSetRunning"
	// ScaleSetCreatingReason used when the scale set creation is in progress.
	ScaleSetCreatingReason = "ScaleSetCreating"
	// ScaleSetUpdatingReason used when the scale set updating is in progress.
	ScaleSetUpdatingReason = "ScaleSetUpdating"
	// ScaleSetDeletingReason used when the scale set is in a deleting state.
	ScaleSetDeletingReason = "ScaleSetDeleting"
	// ScaleSetProvisionFailedReason used for failures during scale set provisioning.
	ScaleSetProvisionFailedReason = "ScaleSetProvisionFailed"

	// ScaleSetDesiredReplicasCondition reports on the scaling state of the machine pool.
	ScaleSetDesiredReplicasCondition clusterv1beta1.ConditionType = "ScaleSetDesiredReplicas"
	// ScaleSetScaleUpReason describes the machine pool scaling up.
	ScaleSetScaleUpReason = "ScaleSetScalingUp"
	// ScaleSetScaleDownReason describes the machine pool scaling down.
	ScaleSetScaleDownReason = "ScaleSetScalingDown"

	// ScaleSetModelUpdatedCondition reports on the model state of the pool.
	ScaleSetModelUpdatedCondition clusterv1beta1.ConditionType = "ScaleSetModelUpdated"
	// ScaleSetModelOutOfDateReason describes the machine pool model being out of date.
	ScaleSetModelOutOfDateReason = "ScaleSetModelOutOfDate"
)

// AzureManagedCluster Conditions and Reasons.
const (
	// ManagedClusterRunningCondition means the AKS cluster exists and is in a running state.
	ManagedClusterRunningCondition clusterv1beta1.ConditionType = "ManagedClusterRunning"
	// AgentPoolsReadyCondition means the AKS agent pools exist and are ready to be used.
	AgentPoolsReadyCondition clusterv1beta1.ConditionType = "AgentPoolsReady"
	// AzureResourceAvailableCondition means the AKS cluster is healthy according to Azure's Resource Health API.
	AzureResourceAvailableCondition clusterv1beta1.ConditionType = "AzureResourceAvailable"
)

// Azure Services Conditions and Reasons.
const (
	// ResourceGroupReadyCondition means the resource group exists and is ready to be used.
	ResourceGroupReadyCondition clusterv1beta1.ConditionType = "ResourceGroupReady"
	// VNetReadyCondition means the virtual network exists and is ready to be used.
	VNetReadyCondition clusterv1beta1.ConditionType = "VNetReady"
	// VnetPeeringReadyCondition means the virtual network peerings exist and are ready to be used.
	VnetPeeringReadyCondition clusterv1beta1.ConditionType = "VnetPeeringReady"
	// SecurityGroupsReadyCondition means the security groups exist and are ready to be used.
	SecurityGroupsReadyCondition clusterv1beta1.ConditionType = "SecurityGroupsReady"
	// RouteTablesReadyCondition means the route tables exist and are ready to be used.
	RouteTablesReadyCondition clusterv1beta1.ConditionType = "RouteTablesReady"
	// PublicIPsReadyCondition means the public IPs exist and are ready to be used.
	PublicIPsReadyCondition clusterv1beta1.ConditionType = "PublicIPsReady"
	// NATGatewaysReadyCondition means the NAT gateways exist and are ready to be used.
	NATGatewaysReadyCondition clusterv1beta1.ConditionType = "NATGatewaysReady"
	// SubnetsReadyCondition means the subnets exist and are ready to be used.
	SubnetsReadyCondition clusterv1beta1.ConditionType = "SubnetsReady"
	// LoadBalancersReadyCondition means the load balancers exist and are ready to be used.
	LoadBalancersReadyCondition clusterv1beta1.ConditionType = "LoadBalancersReady"
	// PrivateDNSZoneReadyCondition means the private DNS zone exists and is ready to be used.
	PrivateDNSZoneReadyCondition clusterv1beta1.ConditionType = "PrivateDNSZoneReady"
	// PrivateDNSLinkReadyCondition means the private DNS links exist and are ready to be used.
	PrivateDNSLinkReadyCondition clusterv1beta1.ConditionType = "PrivateDNSLinkReady"
	// PrivateDNSRecordReadyCondition means the private DNS records exist and are ready to be used.
	PrivateDNSRecordReadyCondition clusterv1beta1.ConditionType = "PrivateDNSRecordReady"
	// BastionHostReadyCondition means the bastion host exists and is ready to be used.
	BastionHostReadyCondition clusterv1beta1.ConditionType = "BastionHostReady"
	// InboundNATRulesReadyCondition means the inbound NAT rules exist and are ready to be used.
	InboundNATRulesReadyCondition clusterv1beta1.ConditionType = "InboundNATRulesReady"
	// AvailabilitySetReadyCondition means the availability set exists and is ready to be used.
	AvailabilitySetReadyCondition clusterv1beta1.ConditionType = "AvailabilitySetReady"
	// RoleAssignmentReadyCondition means the role assignment exists and is ready to be used.
	RoleAssignmentReadyCondition clusterv1beta1.ConditionType = "RoleAssignmentReady"
	// DisksReadyCondition means the disks exist and are ready to be used.
	DisksReadyCondition clusterv1beta1.ConditionType = "DisksReady"
	// NetworkInterfaceReadyCondition means the network interfaces exist and are ready to be used.
	NetworkInterfaceReadyCondition clusterv1beta1.ConditionType = "NetworkInterfacesReady"
	// PrivateEndpointsReadyCondition means the private endpoints exist and are ready to be used.
	PrivateEndpointsReadyCondition clusterv1beta1.ConditionType = "PrivateEndpointsReady"
	// FleetReadyCondition means the Fleet exists and is ready to be used.
	FleetReadyCondition clusterv1beta1.ConditionType = "FleetReady"
	// AKSExtensionsReadyCondition means the AKS Extensions exist and are ready to be used.
	AKSExtensionsReadyCondition clusterv1beta1.ConditionType = "AKSExtensionsReady"

	// CreatingReason means the resource is being created.
	CreatingReason = "Creating"
	// FailedReason means the resource failed to be created.
	FailedReason = "Failed"
	// DeletingReason means the resource is being deleted.
	DeletingReason = "Deleting"
	// DeletedReason means the resource was deleted.
	DeletedReason = "Deleted"
	// DeletionFailedReason means the resource failed to be deleted.
	DeletionFailedReason = "DeletionFailed"
	// UpdatingReason means the resource is being updated.
	UpdatingReason = "Updating"
)

const (
	// LinuxOS is Linux OS value for OSDisk.OSType.
	LinuxOS = "Linux"
	// WindowsOS is Windows OS value for OSDisk.OSType.
	WindowsOS = "Windows"
)

const (
	// OwnedByClusterLabelKey communicates CAPZ's ownership of an ASO resource
	// independently of its ownership of the underlying Azure resource. The
	// value for the label is the CAPI Cluster Name.
	//
	// Deprecated: OwnerReferences now determine ownership.
	OwnedByClusterLabelKey = NameAzureProviderPrefix + string(ResourceLifecycleOwned)
)

const (
	// AzureNetworkPluginName is the name of the Azure network plugin.
	AzureNetworkPluginName = "azure"
)

const (
	// AzureClusterKind indicates the kind of an AzureCluster.
	AzureClusterKind = "AzureCluster"
	// AzureClusterTemplateKind indicates the kind of an AzureClusterTemplate.
	AzureClusterTemplateKind = "AzureClusterTemplate"
	// AzureMachineKind indicates the kind of an AzureMachine.
	AzureMachineKind = "AzureMachine"
	// AzureMachineTemplateKind indicates the kind of an AzureMachineTemplate.
	AzureMachineTemplateKind = "AzureMachineTemplate"
	// AzureMachinePoolKind indicates the kind of an AzureMachinePool.
	AzureMachinePoolKind = "AzureMachinePool"
	// AzureManagedMachinePoolKind indicates the kind of an AzureManagedMachinePool.
	AzureManagedMachinePoolKind = "AzureManagedMachinePool"
	// AzureManagedClusterKind indicates the kind of an AzureManagedCluster.
	AzureManagedClusterKind = "AzureManagedCluster"
	// AzureManagedControlPlaneKind indicates the kind of an AzureManagedControlPlane.
	AzureManagedControlPlaneKind = "AzureManagedControlPlane"
	// AzureManagedControlPlaneTemplateKind indicates the kind of an AzureManagedControlPlaneTemplate.
	AzureManagedControlPlaneTemplateKind = "AzureManagedControlPlaneTemplate"
	// AzureManagedMachinePoolTemplateKind indicates the kind of an AzureManagedMachinePoolTemplate.
	AzureManagedMachinePoolTemplateKind = "AzureManagedMachinePoolTemplate"
	// AzureClusterIdentityKind indicates the kind of an AzureClusterIdentity.
	AzureClusterIdentityKind = "AzureClusterIdentity"
)
