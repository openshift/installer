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

package v1alpha4

import capiv1 "sigs.k8s.io/cluster-api/api/v1beta1"

const (
	DeletionFailed = "DeletionFailed"
)

const (
	// ClusterCategoryCreatedCondition indicates the status of the category linked to the NutanixCluster
	ClusterCategoryCreatedCondition capiv1.ConditionType = "ClusterCategoryCreated"

	ClusterCategoryCreationFailed = "ClusterCategoryCreationFailed"
)

const (
	// PrismCentralClientCondition indicates the status of the client used to connect to Prism Central
	PrismCentralClientCondition capiv1.ConditionType = "PrismClientInit"

	PrismCentralClientInitializationFailed = "PrismClientInitFailed"
)

const (
	// VMProvisionedCondition shows the status of the VM provisioning process
	VMProvisionedCondition capiv1.ConditionType = "VMProvisioned"

	VMProvisionedTaskFailed = "FailedVMTask"

	// VMAddressesAssignedCondition shows the status of the process of assigning the VM addresses
	VMAddressesAssignedCondition capiv1.ConditionType = "VMAddressesAssigned"

	VMAddressesFailed = "VMAddressesFailed"
	VMBootTypeInvalid = "VMBootTypeInvalid"
)

const (
	// VMAddressesAssignedCondition shows the status of the process of assigning the VMs to a project
	ProjectAssignedCondition capiv1.ConditionType = "ProjectAssigned"

	ProjectAssignationFailed = "ProjectAssignationFailed"
)

const (
	// CredentialRefSecretOwnerSetCondition shows the status of setting the Owner
	CredentialRefSecretOwnerSetCondition capiv1.ConditionType = "CredentialRefSecretOwnerSet"

	CredentialRefSecretOwnerSetFailed = "CredentialRefSecretOwnerSetFailed"
)
