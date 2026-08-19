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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OpenStackCredentialSecretReference references a Secret containing OpenStack credentials.
type OpenStackCredentialSecretReference struct {
	// Name of the Secret which contains a `clouds.yaml` key (and optionally `cacert`).
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Namespace where the Secret resides.
	// +kubebuilder:validation:Required
	Namespace string `json:"namespace"`
}

// OpenStackClusterIdentitySpec defines the desired state for an OpenStackClusterIdentity.
type OpenStackClusterIdentitySpec struct {
	// SecretRef references the credentials Secret containing a `clouds.yaml` file.
	// +kubebuilder:validation:Required
	SecretRef OpenStackCredentialSecretReference `json:"secretRef"`

	// NamespaceSelector limits which namespaces may use this identity. If nil, all namespaces are allowed.
	// +optional
	NamespaceSelector *metav1.LabelSelector `json:"namespaceSelector,omitempty"`
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=openstackclusteridentities,scope=Cluster,categories=cluster-api,shortName=osci

// OpenStackClusterIdentity is a cluster-scoped identity that centralizes OpenStack credentials.
type OpenStackClusterIdentity struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec OpenStackClusterIdentitySpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// OpenStackClusterIdentityList contains a list of OpenStackClusterIdentity.
type OpenStackClusterIdentityList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenStackClusterIdentity `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OpenStackClusterIdentity{}, &OpenStackClusterIdentityList{})
}
