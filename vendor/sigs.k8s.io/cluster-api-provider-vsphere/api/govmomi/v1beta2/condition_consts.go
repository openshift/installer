/*
Copyright 2025 The Kubernetes Authors.

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
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
)

// Conditions and condition Reasons for the VSphereCluster object.

const (
	// FailureDomainsAvailableV1Beta1Condition documents the status of the failure domains
	// associated to the VSphereCluster.
	FailureDomainsAvailableV1Beta1Condition clusterv1.ConditionType = "FailureDomainsAvailable"

	// FailureDomainsSkippedV1Beta1Reason (Severity=Info) documents that some of the failure domain statuses
	// associated to the VSphereCluster are reported as not ready.
	FailureDomainsSkippedV1Beta1Reason = "FailureDomainsSkipped"

	// WaitingForFailureDomainStatusV1Beta1Reason (Severity=Info) documents that some of the failure domains
	// associated to the VSphereCluster are not reporting the Ready status.
	// Instead of reporting a false ready status, these failure domains are still under the process of reconciling
	// and hence not yet reporting their status.
	WaitingForFailureDomainStatusV1Beta1Reason = "WaitingForFailureDomainStatus"
)

// Conditions and condition Reasons for the VSphereMachine and the VSphereVM object.
//
// NOTE: VSphereMachine wraps a VMSphereVM, some we are using a unique set of conditions and reasons in order
// to ensure a consistent UX; differences between the two objects will be highlighted in the comments.

const (
	// VMProvisionedV1Beta1Condition documents the status of the provisioning of a VSphereMachine and its underlying VSphereVM.
	VMProvisionedV1Beta1Condition clusterv1.ConditionType = "VMProvisioned"

	// WaitingForClusterInfrastructureV1Beta1Reason (Severity=Info) documents a VSphereMachine waiting for the cluster
	// infrastructure to be ready before starting the provisioning process.
	//
	// NOTE: This reason does not apply to VSphereVM (this state happens before the VSphereVM is actually created).
	WaitingForClusterInfrastructureV1Beta1Reason = "WaitingForClusterInfrastructure"

	// WaitingForBootstrapDataV1Beta1Reason (Severity=Info) documents a VSphereMachine waiting for the bootstrap
	// script to be ready before starting the provisioning process.
	//
	// NOTE: This reason does not apply to VSphereVM (this state happens before the VSphereVM is actually created).
	WaitingForBootstrapDataV1Beta1Reason = "WaitingForBootstrapData"

	// WaitingForStaticIPAllocationV1Beta1Reason (Severity=Info) documents a VSphereVM waiting for the allocation of
	// a static IP address.
	WaitingForStaticIPAllocationV1Beta1Reason = "WaitingForStaticIPAllocation"

	// WaitingForIPAllocationV1Beta1Reason (Severity=Info) documents a VSphereVM waiting for the allocation of
	// an IP address.
	// This is used when the dhcp4 or dhcp6 for a VSphereVM is set and the VSphereVM is waiting for the
	// relevant IP address  to show up on the VM.
	WaitingForIPAllocationV1Beta1Reason = "WaitingForIPAllocation"

	// CloningV1Beta1Reason documents (Severity=Info) a VSphereMachine/VSphereVM currently executing the clone operation.
	CloningV1Beta1Reason = "Cloning"

	// CloningFailedV1Beta1Reason (Severity=Warning) documents a VSphereMachine/VSphereVM controller detecting
	// an error while provisioning; those kind of errors are usually transient and failed provisioning
	// are automatically re-tried by the controller.
	CloningFailedV1Beta1Reason = "CloningFailed"

	// PoweringOnV1Beta1Reason documents (Severity=Info) a VSphereMachine/VSphereVM currently executing the power on sequence.
	PoweringOnV1Beta1Reason = "PoweringOn"

	// PoweringOnFailedV1Beta1Reason (Severity=Warning) documents a VSphereMachine/VSphereVM controller detecting
	// an error while powering on; those kind of errors are usually transient and failed provisioning
	// are automatically re-tried by the controller.
	PoweringOnFailedV1Beta1Reason = "PoweringOnFailed"

	// NotFoundByBIOSUUIDV1Beta1Reason (Severity=Warning) documents a VSphereVM which can't be found by BIOS UUID.
	// Those kind of errors could be transient sometimes and failed VSphereVM are automatically
	// reconciled by the controller.
	NotFoundByBIOSUUIDV1Beta1Reason = "NotFoundByBIOSUUID"

	// TaskFailure (Severity=Warning) documents a VSphereMachine/VSphere task failure; the reconcile look will automatically
	// retry the operation, but a user intervention might be required to fix the problem.
	TaskFailure = "TaskFailure"

	// WaitingForNetworkAddressesV1Beta1Reason (Severity=Info) documents a VSphereMachine waiting for the machine network
	// settings to be reported after machine being powered on.
	//
	// NOTE: This reason does not apply to VSphereVM (this state happens after the VSphereVM is in ready state).
	WaitingForNetworkAddressesV1Beta1Reason = "WaitingForNetworkAddresses"

	// TagsAttachmentFailedV1Beta1Reason (Severity=Error) documents a VSphereMachine/VSphereVM tags attachment failure.
	TagsAttachmentFailedV1Beta1Reason = "TagsAttachmentFailed"

	// PCIDevicesDetachedV1Beta1Condition documents the status of the attached PCI devices on the VSphereVM.
	// It is a negative condition to notify the user that the device(s) is no longer attached to
	// the underlying VM and would require manual intervention to fix the situation.
	//
	// NOTE: This condition does not apply to VSphereMachine.
	PCIDevicesDetachedV1Beta1Condition clusterv1.ConditionType = "PCIDevicesDetached"

	// NotFoundV1Beta1Reason (Severity=Warning) documents the VSphereVM not having the PCI device attached during VM startup.
	// This would indicate that the PCI devices were removed out of band by an external entity.
	NotFoundV1Beta1Reason = "NotFound"
)

// Conditions and Reasons related to utilizing a VSphereIdentity to make connections to a VCenter.
// Can currently be used by VSphereCluster and VSphereVM.
const (
	// VCenterAvailableV1Beta1Condition documents the connectivity with vcenter
	// for a given resource.
	VCenterAvailableV1Beta1Condition clusterv1.ConditionType = "VCenterAvailable"

	// VCenterUnreachableV1Beta1Reason (Severity=Error) documents a controller detecting
	// issues with VCenter reachability.
	VCenterUnreachableV1Beta1Reason = "VCenterUnreachable"
)

const (
	// ClusterModulesAvailableV1Beta1Condition documents the availability of cluster modules for the VSphereCluster object.
	ClusterModulesAvailableV1Beta1Condition clusterv1.ConditionType = "ClusterModulesAvailable"

	// MissingVCenterVersionV1Beta1Reason (Severity=Warning) documents a controller detecting
	//  the scenario in which the vCenter version is not set in the status of the VSphereCluster object.
	MissingVCenterVersionV1Beta1Reason = "MissingVCenterVersion"

	// VCenterVersionIncompatibleV1Beta1Reason (Severity=Info) documents the case where the vCenter version of the
	// VSphereCluster object does not support cluster modules.
	VCenterVersionIncompatibleV1Beta1Reason = "VCenterVersionIncompatible"

	// ClusterModuleSetupFailedV1Beta1Reason (Severity=Warning) documents a controller detecting
	// issues when setting up anti-affinity constraints via cluster modules for objects
	// belonging to the cluster.
	ClusterModuleSetupFailedV1Beta1Reason = "ClusterModuleSetupFailed"
)

const (
	// CredentialsAvailableV1Beta1Condition is used by VSphereClusterIdentity when a credential
	// secret is available and unused by other VSphereClusterIdentities.
	CredentialsAvailableV1Beta1Condition clusterv1.ConditionType = "CredentialsAvailable"

	// SecretNotAvailableV1Beta1Reason is used when the secret referenced by the VSphereClusterIdentity cannot be found.
	SecretNotAvailableV1Beta1Reason = "SecretNotAvailable"

	// SecretOwnerReferenceFailedV1Beta1Reason is used for errors while updating the owner reference of the secret.
	SecretOwnerReferenceFailedV1Beta1Reason = "SecretOwnerReferenceFailed"

	// SecretAlreadyInUseV1Beta1Reason is used when another VSphereClusterIdentity is using the secret.
	SecretAlreadyInUseV1Beta1Reason = "SecretInUse"
)

const (
	// PlacementConstraintMetV1Beta1Condition documents whether the placement constraint is configured correctly or not.
	PlacementConstraintMetV1Beta1Condition clusterv1.ConditionType = "PlacementConstraintMet"

	// ResourcePoolNotFoundV1Beta1Reason (Severity=Error) documents that the resource pool in the placement constraint
	// associated to the VSphereDeploymentZone is misconfigured.
	ResourcePoolNotFoundV1Beta1Reason = "ResourcePoolNotFound"

	// FolderNotFoundV1Beta1Reason (Severity=Error) documents that the folder in the placement constraint
	// associated to the VSphereDeploymentZone is misconfigured.
	FolderNotFoundV1Beta1Reason = "FolderNotFound"
)

const (
	// VSphereFailureDomainValidatedV1Beta1Condition documents whether the failure domain for the deployment zone is configured correctly or not.
	VSphereFailureDomainValidatedV1Beta1Condition clusterv1.ConditionType = "VSphereFailureDomainValidated"

	// RegionMisconfiguredV1Beta1Reason (Severity=Error) documents that the region for the Failure Domain associated to
	// the VSphereDeploymentZone is misconfigured.
	RegionMisconfiguredV1Beta1Reason = "FailureDomainRegionMisconfigured"

	// ZoneMisconfiguredV1Beta1Reason (Severity=Error) documents that the zone for the Failure Domain associated to
	// the VSphereDeploymentZone is misconfigured.
	ZoneMisconfiguredV1Beta1Reason = "FailureDomainZoneMisconfigured"

	// ComputeClusterNotFoundV1Beta1Reason (Severity=Error) documents that the Compute Cluster for the Failure Domain
	// associated to the VSphereDeploymentZone cannot be found.
	ComputeClusterNotFoundV1Beta1Reason = "ComputeClusterNotFound"

	// HostsMisconfiguredV1Beta1Reason (Severity=Error) documents that the VM & Host Group details for the Failure Domain
	// associated to the VSphereDeploymentZone are misconfigured.
	HostsMisconfiguredV1Beta1Reason = "HostsMisconfigured"

	// HostsAffinityMisconfiguredV1Beta1Reason (Severity=Warning) documents that the VM & Host Group affinity rule for the FailureDomain is disabled.
	HostsAffinityMisconfiguredV1Beta1Reason = "HostsAffinityMisconfigured"

	// NetworkNotFoundV1Beta1Reason (Severity=Error) documents that the networks in the topology for the Failure Domain
	// associated to the VSphereDeploymentZone are misconfigured.
	NetworkNotFoundV1Beta1Reason = "NetworkNotFound"

	// DatastoreNotFoundV1Beta1Reason (Severity=Error) documents that the datastore in the topology for the Failure Domain
	// associated to the VSphereDeploymentZone is misconfigured.
	DatastoreNotFoundV1Beta1Reason = "DatastoreNotFound"
)

const (
	// IPAddressClaimedV1Beta1Condition documents the status of claiming an IP address
	// from an IPAM provider.
	IPAddressClaimedV1Beta1Condition clusterv1.ConditionType = "IPAddressClaimed"

	// IPAddressClaimsBeingCreatedV1Beta1Reason (Severity=Info) documents that claims for the
	// IP addresses required by the VSphereVM are being created.
	IPAddressClaimsBeingCreatedV1Beta1Reason = "IPAddressClaimsBeingCreated"

	// WaitingForIPAddressV1Beta1Reason (Severity=Info) documents that the VSphereVM is
	// currently waiting for an IP address to be provisioned.
	WaitingForIPAddressV1Beta1Reason = "WaitingForIPAddress"

	// IPAddressInvalidV1Beta1Reason (Severity=Error) documents that the IP address
	// provided by the IPAM provider is not valid.
	IPAddressInvalidV1Beta1Reason = "IPAddressInvalid"

	// IPAddressClaimNotFoundV1Beta1Reason (Severity=Error) documents that the IPAddressClaim
	// cannot be found.
	IPAddressClaimNotFoundV1Beta1Reason = "IPAddressClaimNotFound"
)

const (
	// GuestSoftPowerOffSucceededV1Beta1Condition documents the status of performing guest initiated
	// graceful shutdown.
	GuestSoftPowerOffSucceededV1Beta1Condition clusterv1.ConditionType = "GuestSoftPowerOffSucceeded"

	// GuestSoftPowerOffInProgressV1Beta1Reason (Severity=Info) documents that the guest receives
	// a graceful shutdown request.
	GuestSoftPowerOffInProgressV1Beta1Reason = "GuestSoftPowerOffInProgress"

	// GuestSoftPowerOffFailedV1Beta1Reason (Severity=Warning) documents that the graceful
	// shutdown request fails.
	GuestSoftPowerOffFailedV1Beta1Reason = "GuestSoftPowerOffFailed"
)
