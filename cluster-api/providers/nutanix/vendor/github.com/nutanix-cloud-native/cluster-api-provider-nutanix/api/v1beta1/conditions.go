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

import capiv1 "sigs.k8s.io/cluster-api/api/v1beta1"

const (
	DeletionFailed = "DeletionFailed"

	VolumeGroupDetachFailed = "VolumeGroupDetachFailed"
)

// Conditions and Reasons releated to failure domain
const (
	// FailureDomainSafeForDeletionCondition indicates whether the failure domain object is safe for deletion,
	// ie., when it is not used or referenced by other resources
	FailureDomainSafeForDeletionCondition capiv1.ConditionType = "FailureDomainSafeForDeletion"

	// FailureDomainInUseReason indicates that the failure domain is used by
	// Machines and/or referenced by cluster
	FailureDomainInUseReason = "FailureDomainInUse"

	// NoFailureDomainsConfiguredCondition indicates no failure domains have been configured
	NoFailureDomainsConfiguredCondition capiv1.ConditionType = "NoFailureDomainsConfigured"

	// FailureDomainsValidatedCondition indicates whether the failure domains are configured correctly or not.
	FailureDomainsValidatedCondition capiv1.ConditionType = "FailureDomainsValidated"

	// FailureDomainsMisconfiguredReason (Severity=Warning) indicates that some of the failure domains
	// are misconfigured.
	FailureDomainsMisconfiguredReason = "FailureDomainsMisconfigured"
)

const (
	// ClusterCategoryCreatedCondition indicates the status of the category linked to the NutanixCluster
	ClusterCategoryCreatedCondition capiv1.ConditionType = "ClusterCategoryCreated"

	ClusterCategoryCreationFailed = "ClusterCategoryCreationFailed"
)

const (
	// PrismCentralClientCondition indicates the status of the client used to connect to Prism Central
	PrismCentralClientCondition   capiv1.ConditionType = "PrismClientInit"
	PrismCentralV4ClientCondition capiv1.ConditionType = "PrismClientV4Init"

	PrismCentralClientInitializationFailed   = "PrismClientInitFailed"
	PrismCentralV4ClientInitializationFailed = "PrismClientV4InitFailed"
)

const (
	// VMProvisionedCondition shows the status of the VM provisioning process
	VMProvisionedCondition capiv1.ConditionType = "VMProvisioned"

	VMProvisionedTaskFailed = "FailedVMTask"

	// VMAddressesAssignedCondition shows the status of the process of assigning the VM addresses
	VMAddressesAssignedCondition capiv1.ConditionType = "VMAddressesAssigned"

	VMAddressesFailed             = "VMAddressesFailed"
	VMBootTypeInvalid             = "VMBootTypeInvalid"
	ClusterInfrastructureNotReady = "ClusterInfrastructureNotReady"
	BootstrapDataNotReady         = "BootstrapDataNotReady"
	ControlplaneNotInitialized    = "ControlplaneNotInitialized"
)

const (
	// VMAddressesAssignedCondition shows the status of the process of assigning the VMs to a project
	ProjectAssignedCondition capiv1.ConditionType = "ProjectAssigned"

	ProjectAssignationFailed = "ProjectAssignationFailed"
)

const (
	// CredentialRefSecretOwnerSetCondition shows the status of setting the Owner
	CredentialRefSecretOwnerSetCondition capiv1.ConditionType = "CredentialRefSecretOwnerSet"

	CredentialRefSecretOwnerSetFailed  = "CredentialRefSecretOwnerSetFailed"
	TrustBundleSecretOwnerSetCondition = "TrustBundleSecretOwnerSet"
	TrustBundleSecretOwnerSetFailed    = "TrustBundleSecretOwnerSetFailed"
)
