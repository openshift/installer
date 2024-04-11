/*
Copyright 2019 The OpenShift Authors.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TODO: these types should eventually be broken out, along with the actuator,
// to a separate repo.

// AzureProviderSpec contains the required information to create RBAC role
// bindings for Azure.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AzureProviderSpec struct {
	metav1.TypeMeta `json:",inline"`

	// RoleBindings contains a list of roles that should be associated with the minted credential.
	RoleBindings []RoleBinding `json:"roleBindings"`

	// Permissions is the list of Azure permissions required to create a more fine-grained custom role to
	// satisfy the CredentialsRequest.
	// The Permissions field may be provided in addition to RoleBindings. When both fields are specified,
	// the user-assigned managed identity will have union of permissions defined from both Permissions
	// and RoleBindings.
	// +optional
	Permissions []string `json:"permissions,omitempty"`

	// DataPermissions is the list of Azure data permissions required to create a more fine-grained custom
	// role to satisfy the CredentialsRequest.
	// The DataPermissions field may be provided in addition to RoleBindings. When both fields are specified,
	// the user-assigned managed identity will have union of permissions defined from both DataPermissions
	// and RoleBindings.
	// +optional
	DataPermissions []string `json:"dataPermissions,omitempty"`

	// The following fields are only required for Azure Workload Identity.
	// AzureClientID is the ID of the specific application you created in Azure
	// +optional
	AzureClientID string `json:"azureClientID,omitempty"`

	// AzureRegion is the geographic region of the Azure service.
	// +optional
	AzureRegion string `json:"azureRegion,omitempty"`

	// Each Azure subscription has an ID associated with it, as does the tenant to which a subscription belongs.
	// AzureSubscriptionID is the ID of the subscription.
	// +optional
	AzureSubscriptionID string `json:"azureSubscriptionID,omitempty"`

	// AzureTenantID is the ID of the tenant to which the subscription belongs.
	// +optional
	AzureTenantID string `json:"azureTenantID,omitempty"`
}

// RoleBinding models part of the Azure RBAC Role Binding
type RoleBinding struct {
	// Role defines a set of permissions that should be associated with the minted credential.
	Role string `json:"role"`
}

// AzureProviderStatus contains the status of the credentials request in Azure.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AzureProviderStatus struct {
	metav1.TypeMeta `json:",inline"`

	// ServicePrincipalName is the name of the service principal created in Azure for these credentials.
	ServicePrincipalName string `json:"name"`

	// AppID is the application id of the service principal created in Azure for these credentials.
	AppID string `json:"appID"`

	// SecretLastResourceVersion is the resource version of the secret resource
	// that was last synced. Used to determine if the object has changed and
	// requires a sync.
	SecretLastResourceVersion string `json:"secretLastResourceVersion"`
}
