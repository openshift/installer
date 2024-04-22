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

// TODO: these types should eventually be broken out, along with the actuator, to a separate repo.

// GCPProviderSpec contains the required information to create a service account with policy bindings in GCP.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type GCPProviderSpec struct {
	metav1.TypeMeta `json:",inline"`
	// PredefinedRoles is the list of GCP pre-defined roles
	// that the CredentialsRequest requires.
	PredefinedRoles []string `json:"predefinedRoles"`
	// Permissions is the list of GCP permissions required to create a more fine-grained custom role to
	// satisfy the CredentialsRequest.
	// The Permissions field may be provided in addition to PredefinedRoles. When both fields are specified,
	// the service account will have union of permissions defined from both Permissions and PredefinedRoles.
	// +optional
	Permissions []string `json:"permissions,omitempty"`
	// SkipServiceCheck can be set to true to skip the check whether the requested roles or permissions
	// have the necessary services enabled
	// +optional
	SkipServiceCheck bool `json:"skipServiceCheck,omitempty"`
}

// GCPProviderStatus contains the status of the GCP credentials request.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type GCPProviderStatus struct {
	metav1.TypeMeta `json:",inline"`
	// ServiceAccountID is the ID of the service account created in GCP for the requested credentials.
	ServiceAccountID string `json:"serviceAccountID"`
	// RoleID is the ID of the custom role created in GCP for the requested permissions apart from
	// permissions granted by the pre-defined roles.
	// RoleID is set by the Cloud Credential Operator controllers and should not be set manually.
	// +optional
	RoleID string `json:"roleID,omitempty"`
}
