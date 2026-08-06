/*
Copyright 2022 Nutanix

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

import capiv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1" //nolint:staticcheck // suppress complaining on Deprecated package

const (
	Succeeded = "Succeeded"

	DeletionFailed = "DeletionFailed"

	VolumeGroupDetachFailed = "VolumeGroupDetachFailed"
)

// Conditions and Reasons related to failure domain
const (
	// FailureDomainSafeForDeletionCondition indicates whether the failure domain object is safe for deletion,
	// ie., when it is not used or referenced by other resources
	FailureDomainSafeForDeletionCondition capiv1beta1.ConditionType = "FailureDomainSafeForDeletion"

	// FailureDomainInUseReason indicates that the failure domain is used by
	// Machines and/or referenced by cluster
	FailureDomainInUseReason = "FailureDomainInUse"

	// NoFailureDomainsConfiguredCondition indicates no failure domains have been configured
	NoFailureDomainsConfiguredCondition capiv1beta1.ConditionType = "NoFailureDomainsConfigured"

	// FailureDomainsValidatedCondition indicates whether the failure domains are configured correctly or not.
	FailureDomainsValidatedCondition capiv1beta1.ConditionType = "FailureDomainsValidated"

	// FailureDomainsMisconfiguredReason (Severity=Warning) indicates that some of the failure domains
	// are misconfigured.
	FailureDomainsMisconfiguredReason = "FailureDomainsMisconfigured"

	// FailureDomainsValidatedCondition indicates that the failure domains are being validated.
	FailureDomainsValidationInProgressReason = "FailureDomainsValidationInProgress"
)

// Conditions and Reasons related to NutanixMetro and NutanixMetroSite
const (
	// MetroSafeForDeletionCondition indicates whether the NutanixMetro object is safe for deletion,
	// ie., when it is not used or referenced by other resources
	MetroSafeForDeletionCondition = "metroSafeForDeletion"

	// MetroInUseReason indicates that the NutanixMetro is used by other objects.
	MetroInUseReason = "metroInUse"

	// MetroNotInUseReason indicates that the NutanixMetro is not used by any objects.
	MetroNotInUseReason = "metroNotInUse"

	// MetroValidatedCondition indicates whether the NutanixMetro spec is configured correctly or not.
	MetroValidatedCondition = "metroValidated"

	// MetroMisconfiguredReason indicates that the NutanixMetro spec is misconfigured.
	MetroMisconfiguredReason = "metroMisconfigured"

	// MetroSpecValidReason indicates that the NutanixMetro spec is valid.
	MetroSpecValidReason = "metroSpecValid"

	// MetroSiteSafeForDeletionCondition indicates whether the NutanixMetroSite object is safe for deletion,
	// ie., when it is not used or referenced by other resources
	MetroSiteSafeForDeletionCondition = "metroSiteSafeForDeletion"

	// MetroSiteInUseReason indicates that the NutanixMetroSite is used by other objects.
	MetroSiteInUseReason = "metroSiteInUse"

	// MetroSiteNotInUseReason indicates that the NutanixMetroSite is not used by any objects.
	MetroSiteNotInUseReason = "metroSiteNotInUse"

	// MetroSiteValidatedCondition indicates whether the NutanixMetroSite spec is configured correctly or not.
	MetroSiteValidatedCondition = "metroSiteValidated"

	// MetroSiteMisconfiguredReason indicates that the NutanixMetroSite spec is misconfigured.
	MetroSiteMisconfiguredReason = "metroSiteMisconfigured"

	// MetroSiteSpecValidReason indicates that the NutanixMetroSite spec is configured correctly.
	MetroSiteSpecValidReason = "metroSiteSpecValid"
)

const (
	// ClusterCategoryCreatedCondition indicates the status of the category linked to the NutanixCluster
	ClusterCategoryCreatedCondition capiv1beta1.ConditionType = "ClusterCategoryCreated"

	ClusterCategoryCreationFailed = "ClusterCategoryCreationFailed"
)

const (
	// PrismCentralClientCondition indicates the status of the client used to connect to Prism Central
	PrismCentralClientCondition            capiv1beta1.ConditionType = "PrismClientInit"
	PrismCentralConvergedV4ClientCondition capiv1beta1.ConditionType = "PrismClientConvergedV4Init"

	PrismCentralClientInitializationFailed            = "PrismClientInitFailed"
	PrismCentralConvergedV4ClientInitializationFailed = "PrismClientConvergedV4InitFailed"
)

const (
	// VMProvisionedCondition shows the status of the VM provisioning process
	VMProvisionedCondition capiv1beta1.ConditionType = "VMProvisioned"

	VMProvisionedTaskFailed = "FailedVMTask"

	// VMAddressesAssignedCondition shows the status of the process of assigning the VM addresses
	VMAddressesAssignedCondition capiv1beta1.ConditionType = "VMAddressesAssigned"

	VMAddressesFailed             = "VMAddressesFailed"
	VMBootTypeInvalid             = "VMBootTypeInvalid"
	ClusterInfrastructureNotReady = "ClusterInfrastructureNotReady"
	BootstrapDataNotReady         = "BootstrapDataNotReady"
	ControlplaneNotInitialized    = "ControlplaneNotInitialized"
)

const (
	// VMAddressesAssignedCondition shows the status of the process of assigning the VMs to a project
	ProjectAssignedCondition capiv1beta1.ConditionType = "ProjectAssigned"

	ProjectAssignationFailed = "ProjectAssignationFailed"
)

// Conditions and Reasons related to NutanixVirtualHADomain.
const (
	// VHADomainPCResourcesValidatedCondition indicates whether the vHA domain PC resources
	// (categories, protection policy, and recovery plans) are valid and present in Prism Central.
	VHADomainPCResourcesValidatedCondition = "VHADomainPCResourcesValidated"

	// VHADomainPCResourcesValidReason indicates that the vHA domain PC resources are valid and present.
	VHADomainPCResourcesValidReason = "VHADomainPCResourcesValid"

	// VHADomainPCResourcesValidationFailedReason indicates that one or more vHA domain PC resources
	// are missing or could not be validated in Prism Central.
	VHADomainPCResourcesValidationFailedReason = "VHADomainPCResourcesValidationFailed"

	// VHADomainSafeForDeletionCondition indicates whether the NutanixVirtualHADomain object is safe
	// for deletion, ie., its owning Cluster/NutanixCluster are in deletion and its PC resources are cleaned up.
	VHADomainSafeForDeletionCondition = "VHADomainSafeForDeletion"

	// VHADomainOwnerNotInDeletion indicates that the vHA domain's owning Cluster/NutanixCluster
	// is not in deletion, so the vHA domain cannot be deleted yet.
	VHADomainOwnerNotInDeletion = "VHADomainOwnerNotInDeletion"

	// VHADomainFailedCleanup indicates that cleaning up the vHA domain's PC resources failed.
	VHADomainFailedCleanup = "VHADomainFailedCleanup"

	// VHADomainCleanupSucceeded indicates that the vHA domain's owning Cluster/NutanixCluster is in
	// deletion and its PC resources have been cleaned up, so the object is safe for deletion.
	VHADomainCleanupSucceeded = "VHADomainCleanupSucceeded"
)

const (
	// CredentialRefSecretOwnerSetCondition shows the status of setting the Owner
	CredentialRefSecretOwnerSetCondition capiv1beta1.ConditionType = "CredentialRefSecretOwnerSet"

	CredentialRefSecretOwnerSetFailed  = "CredentialRefSecretOwnerSetFailed"
	TrustBundleSecretOwnerSetCondition = "TrustBundleSecretOwnerSet"
	TrustBundleSecretOwnerSetFailed    = "TrustBundleSecretOwnerSetFailed"
)
