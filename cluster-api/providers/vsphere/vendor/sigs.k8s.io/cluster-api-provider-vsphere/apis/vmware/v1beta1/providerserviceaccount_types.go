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

//nolint:godot
package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ProviderServiceAccountSpec defines the desired state of ProviderServiceAccount.
type ProviderServiceAccountSpec struct {
	// Ref specifies the reference to the VSphereCluster for which the ProviderServiceAccount needs to be realized.
	Ref *corev1.ObjectReference `json:"ref"`

	// Rules specifies the privileges that need to be granted to the service account.
	Rules []rbacv1.PolicyRule `json:"rules"`

	// TargetNamespace is the namespace in the target cluster where the secret containing the generated service account
	// token needs to be created.
	TargetNamespace string `json:"targetNamespace"`

	// TargetSecretName is the name of the secret in the target cluster that contains the generated service account
	// token.
	TargetSecretName string `json:"targetSecretName"`
}

// ProviderServiceAccountStatus defines the observed state of ProviderServiceAccount.
type ProviderServiceAccountStatus struct {
	Ready    bool   `json:"ready,omitempty"`
	ErrorMsg string `json:"errorMsg,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=providerserviceaccounts,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="VSphereCluster",type=string,JSONPath=.spec.ref.name
// +kubebuilder:printcolumn:name="TargetNamespace",type=string,JSONPath=.spec.targetNamespace
// +kubebuilder:printcolumn:name="TargetSecretName",type=string,JSONPath=.spec.targetSecretName
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// ProviderServiceAccount is the schema for the ProviderServiceAccount API.
type ProviderServiceAccount struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ProviderServiceAccountSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// ProviderServiceAccountList contains a list of ProviderServiceAccount
type ProviderServiceAccountList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProviderServiceAccount `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ProviderServiceAccount{}, &ProviderServiceAccountList{})
}
