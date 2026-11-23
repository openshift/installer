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

package v1beta2

import (
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ProviderServiceAccountSpec defines the desired state of ProviderServiceAccount.
type ProviderServiceAccountSpec struct {
	// ref specifies the reference to the VSphereCluster for which the ProviderServiceAccount needs to be realized.
	// +required
	Ref *corev1.ObjectReference `json:"ref,omitempty"`

	// rules specifies the privileges that need to be granted to the service account.
	// +required
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=1024
	Rules []rbacv1.PolicyRule `json:"rules,omitempty"`

	// targetNamespace is the namespace in the target cluster where the secret containing the generated service account
	// token needs to be created.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	TargetNamespace string `json:"targetNamespace,omitempty"`

	// targetSecretName is the name of the secret in the target cluster that contains the generated service account
	// token.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	TargetSecretName string `json:"targetSecretName,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=providerserviceaccounts,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:deprecatedversion:warning="ProviderServiceAccount is deprecated and is going to be removed in one of the upcoming CAPV releases"
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="VSphereCluster",type=string,JSONPath=.spec.ref.name
// +kubebuilder:printcolumn:name="TargetNamespace",type=string,JSONPath=.spec.targetNamespace
// +kubebuilder:printcolumn:name="TargetSecretName",type=string,JSONPath=.spec.targetSecretName
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// ProviderServiceAccount is the schema for the ProviderServiceAccount API.
//
// Deprecated: ProviderServiceAccount is deprecated and is going to be removed in one of the upcoming CAPV releases.
type ProviderServiceAccount struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is the desired state of ProviderServiceAccount.
	// +required
	Spec ProviderServiceAccountSpec `json:"spec,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// ProviderServiceAccountList contains a list of ProviderServiceAccount.
type ProviderServiceAccountList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProviderServiceAccount `json:"items"`
}

func init() {
	objectTypes = append(objectTypes, &ProviderServiceAccount{}, &ProviderServiceAccountList{})
}
