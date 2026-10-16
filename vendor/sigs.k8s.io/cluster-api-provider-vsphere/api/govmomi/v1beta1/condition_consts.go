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

// Conditions and condition Reasons for the VSphereCluster object.

const (
	// FailureDomainsAvailableCondition documents the status of the failure domains
	// associated to the VSphereCluster.
	FailureDomainsAvailableCondition clusterv1beta1.ConditionType = "FailureDomainsAvailable"

	// FailureDomainsSkippedReason (Severity=Info) documents that some of the failure domain statuses
	// associated to the VSphereCluster are reported as not ready.
	FailureDomainsSkippedReason = "FailureDomainsSkipped"

	// WaitingForFailureDomainStatusReason (Severity=Info) documents that some of the failure domains
	// associated to the VSphereCluster are not reporting the Ready status.
	// Instead of reporting a false ready status, these failure domains are still under the process of reconciling
	// and hence not yet reporting their status.
	WaitingForFailureDomainStatusReason = "WaitingForFailureDomainStatus"
)

// Conditions and condition Reasons for the VSphereMachine and the VSphereVM object.
//
// NOTE: VSphereMachine wraps a VMSphereVM, some we are using a unique set of conditions and reasons in order
// to ensure a consistent UX; differences between the two objects will be highlighted in the comments.

const (
	// VMProvisionedCondition documents the status of the provisioning of a VSphereMachine and its underlying VSphereVM.
	VMProvisionedCondition clusterv1beta1.ConditionType = "VMProvisioned"

	// WaitingForClusterInfrastructureReason (Severity=Info) documents a VSphereMachine waiting for the cluster
	// infrastructure to be ready before starting the provisioning process.
	//
	// NOTE: This reason does not apply to VSphereVM (this state happens before the VSphereVM is actually created).
	WaitingForClusterInfrastructureReason = "WaitingForClusterInfrastructure"

	// WaitingForBootstrapDataReason (Severity=Info) documents a VSphereMachine waiting for the bootstrap
	// script to be ready before starting the provisioning process.
	//
	// NOTE: This reason does not apply to VSphereVM (this state happens before the VSphereVM is actually created).
	WaitingForBootstrapDataReason = "WaitingForBootstrapData"

	// WaitingForStaticIPAllocationReason (Severity=Info) documents a VSphereVM waiting for the allocation of
	// a static IP address.
	WaitingForStaticIPAllocationReason = "WaitingForStaticIPAllocation"

	// WaitingForIPAllocationReason (Severity=Info) documents a VSphereVM waiting for the allocation of
	// an IP address.
	// This is used when the dhcp4 or dhcp6 for a VSphereVM is set and the VSphereVM is waiting for the
	// relevant IP address  to show up on the VM.
	WaitingForIPAllocationReason = "WaitingForIPAllocation"

	// CloningReason documents (Severity=Info) a VSphereMachine/VSphereVM currently executing the clone operation.
	CloningReason = "Cloning"

	// CloningFailedReason (Severity=Warning) documents a VSphereMachine/VSphereVM controller detecting
	// an error while provisioning; those kind of errors are usually transient and failed provisioning
	// are automatically re-tried by the controller.
	CloningFailedReason = "CloningFailed"

	// PoweringOnReason documents (Severity=Info) a VSphereMachine/VSphereVM currently executing the power on sequence.
	PoweringOnReason = "PoweringOn"

	// PoweringOnFailedReason (Severity=Warning) documents a VSphereMachine/VSphereVM controller detecting
	// an error while powering on; those kind of errors are usually transient and failed provisioning
	// are automatically re-tried by the controller.
	PoweringOnFailedReason = "PoweringOnFailed"

	// NotFoundByBIOSUUIDReason (Severity=Warning) documents a VSphereVM which can't be found by BIOS UUID.
	// Those kind of errors could be transient sometimes and failed VSphereVM are automatically
	// reconciled by the controller.
	NotFoundByBIOSUUIDReason = "NotFoundByBIOSUUID"

	// TaskFailure (Severity=Warning) documents a VSphereMachine/VSphere task failure; the reconcile look will automatically
	// retry the operation, but a user intervention might be required to fix the problem.
	TaskFailure = "TaskFailure"

	// WaitingForNetworkAddressesReason (Severity=Info) documents a VSphereMachine waiting for the machine network
	// settings to be reported after machine being powered on.
	//
	// NOTE: This reason does not apply to VSphereVM (this state happens after the VSphereVM is in ready state).
	WaitingForNetworkAddressesReason = "WaitingForNetworkAddresses"

	// TagsAttachmentFailedReason (Severity=Error) documents a VSphereMachine/VSphereVM tags attachment failure.
	TagsAttachmentFailedReason = "TagsAttachmentFailed"

	// PCIDevicesDetachedCondition documents the status of the attached PCI devices on the VSphereVM.
	// It is a negative condition to notify the user that the device(s) is no longer attached to
	// the underlying VM and would require manual intervention to fix the situation.
	//
	// NOTE: This condition does not apply to VSphereMachine.
	PCIDevicesDetachedCondition clusterv1beta1.ConditionType = "PCIDevicesDetached"

	// NotFoundReason (Severity=Warning) documents the VSphereVM not having the PCI device attached during VM startup.
	// This would indicate that the PCI devices were removed out of band by an external entity.
	NotFoundReason = "NotFound"
)

// Conditions and Reasons related to utilizing a VSphereIdentity to make connections to a VCenter.
// Can currently be used by VSphereCluster and VSphereVM.
const (
	// VCenterAvailableCondition documents the connectivity with vcenter
	// for a given resource.
	VCenterAvailableCondition clusterv1beta1.ConditionType = "VCenterAvailable"

	// VCenterUnreachableReason (Severity=Error) documents a controller detecting
	// issues with VCenter reachability.
	VCenterUnreachableReason = "VCenterUnreachable"
)

const (
	// ClusterModulesAvailableCondition documents the availability of cluster modules for the VSphereCluster object.
	ClusterModulesAvailableCondition clusterv1beta1.ConditionType = "ClusterModulesAvailable"

	// MissingVCenterVersionReason (Severity=Warning) documents a controller detecting
	//  the scenario in which the vCenter version is not set in the status of the VSphereCluster object.
	MissingVCenterVersionReason = "MissingVCenterVersion"

	// VCenterVersionIncompatibleReason (Severity=Info) documents the case where the vCenter version of the
	// VSphereCluster object does not support cluster modules.
	VCenterVersionIncompatibleReason = "VCenterVersionIncompatible"

	// ClusterModuleSetupFailedReason (Severity=Warning) documents a controller detecting
	// issues when setting up anti-affinity constraints via cluster modules for objects
	// belonging to the cluster.
	ClusterModuleSetupFailedReason = "ClusterModuleSetupFailed"
)

const (
	// CredentialsAvailableCondidtion is used by VSphereClusterIdentity when a credential
	// secret is available and unused by other VSphereClusterIdentities.
	CredentialsAvailableCondidtion clusterv1beta1.ConditionType = "CredentialsAvailable"

	// SecretNotAvailableReason is used when the secret referenced by the VSphereClusterIdentity cannot be found.
	SecretNotAvailableReason = "SecretNotAvailable"

	// SecretOwnerReferenceFailedReason is used for errors while updating the owner reference of the secret.
	SecretOwnerReferenceFailedReason = "SecretOwnerReferenceFailed"

	// SecretAlreadyInUseReason is used when another VSphereClusterIdentity is using the secret.
	SecretAlreadyInUseReason = "SecretInUse"
)

const (
	// PlacementConstraintMetCondition documents whether the placement constraint is configured correctly or not.
	PlacementConstraintMetCondition clusterv1beta1.ConditionType = "PlacementConstraintMet"

	// ResourcePoolNotFoundReason (Severity=Error) documents that the resource pool in the placement constraint
	// associated to the VSphereDeploymentZone is misconfigured.
	ResourcePoolNotFoundReason = "ResourcePoolNotFound"

	// FolderNotFoundReason (Severity=Error) documents that the folder in the placement constraint
	// associated to the VSphereDeploymentZone is misconfigured.
	FolderNotFoundReason = "FolderNotFound"
)

const (
	// VSphereFailureDomainValidatedCondition documents whether the failure domain for the deployment zone is configured correctly or not.
	VSphereFailureDomainValidatedCondition clusterv1beta1.ConditionType = "VSphereFailureDomainValidated"

	// RegionMisconfiguredReason (Severity=Error) documents that the region for the Failure Domain associated to
	// the VSphereDeploymentZone is misconfigured.
	RegionMisconfiguredReason = "FailureDomainRegionMisconfigured"

	// ZoneMisconfiguredReason (Severity=Error) documents that the zone for the Failure Domain associated to
	// the VSphereDeploymentZone is misconfigured.
	ZoneMisconfiguredReason = "FailureDomainZoneMisconfigured"

	// ComputeClusterNotFoundReason (Severity=Error) documents that the Compute Cluster for the Failure Domain
	// associated to the VSphereDeploymentZone cannot be found.
	ComputeClusterNotFoundReason = "ComputeClusterNotFound"

	// HostsMisconfiguredReason (Severity=Error) documents that the VM & Host Group details for the Failure Domain
	// associated to the VSphereDeploymentZone are misconfigured.
	HostsMisconfiguredReason = "HostsMisconfigured"

	// HostsAffinityMisconfiguredReason (Severity=Warning) documents that the VM & Host Group affinity rule for the FailureDomain is disabled.
	HostsAffinityMisconfiguredReason = "HostsAffinityMisconfigured"

	// NetworkNotFoundReason (Severity=Error) documents that the networks in the topology for the Failure Domain
	// associated to the VSphereDeploymentZone are misconfigured.
	NetworkNotFoundReason = "NetworkNotFound"

	// DatastoreNotFoundReason (Severity=Error) documents that the datastore in the topology for the Failure Domain
	// associated to the VSphereDeploymentZone is misconfigured.
	DatastoreNotFoundReason = "DatastoreNotFound"
)

const (
	// IPAddressClaimedCondition documents the status of claiming an IP address
	// from an IPAM provider.
	IPAddressClaimedCondition clusterv1beta1.ConditionType = "IPAddressClaimed"

	// IPAddressClaimsBeingCreatedReason (Severity=Info) documents that claims for the
	// IP addresses required by the VSphereVM are being created.
	IPAddressClaimsBeingCreatedReason = "IPAddressClaimsBeingCreated"

	// WaitingForIPAddressReason (Severity=Info) documents that the VSphereVM is
	// currently waiting for an IP address to be provisioned.
	WaitingForIPAddressReason = "WaitingForIPAddress"

	// IPAddressInvalidReason (Severity=Error) documents that the IP address
	// provided by the IPAM provider is not valid.
	IPAddressInvalidReason = "IPAddressInvalid"

	// IPAddressClaimNotFoundReason (Severity=Error) documents that the IPAddressClaim
	// cannot be found.
	IPAddressClaimNotFoundReason = "IPAddressClaimNotFound"
)

const (
	// GuestSoftPowerOffSucceededCondition documents the status of performing guest initiated
	// graceful shutdown.
	GuestSoftPowerOffSucceededCondition clusterv1beta1.ConditionType = "GuestSoftPowerOffSucceeded"

	// GuestSoftPowerOffInProgressReason (Severity=Info) documents that the guest receives
	// a graceful shutdown request.
	GuestSoftPowerOffInProgressReason = "GuestSoftPowerOffInProgress"

	// GuestSoftPowerOffFailedReason (Severity=Warning) documents that the graceful
	// shutdown request fails.
	GuestSoftPowerOffFailedReason = "GuestSoftPowerOffFailed"
)
