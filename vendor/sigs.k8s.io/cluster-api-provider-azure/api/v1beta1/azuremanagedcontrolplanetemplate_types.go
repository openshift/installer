/*
Copyright 2023 The Kubernetes Authors.

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

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AzureManagedControlPlaneTemplateSpec defines the desired state of AzureManagedControlPlaneTemplate.
type AzureManagedControlPlaneTemplateSpec struct {
	Template AzureManagedControlPlaneTemplateResource `json:"template"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=azuremanagedcontrolplanetemplates,scope=Namespaced,categories=cluster-api,shortName=amcpt
// +kubebuilder:storageversion
// +kubebuilder:deprecatedversion:warning="AzureManagedControlPlaneTemplate and the AzureManaged API are deprecated. Please migrate to infrastructure.cluster.x-k8s.io/v1beta1 AzureASOManagedControlPlaneTemplate and related AzureASOManaged resources instead."

// AzureManagedControlPlaneTemplate is the Schema for the AzureManagedControlPlaneTemplates API.
type AzureManagedControlPlaneTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec AzureManagedControlPlaneTemplateSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// AzureManagedControlPlaneTemplateList contains a list of AzureManagedControlPlaneTemplates.
type AzureManagedControlPlaneTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AzureManagedControlPlaneTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AzureManagedControlPlaneTemplate{}, &AzureManagedControlPlaneTemplateList{})
}

// AzureManagedControlPlaneTemplateResource describes the data needed to create an AzureManagedCluster from a template.
type AzureManagedControlPlaneTemplateResource struct {
	Spec AzureManagedControlPlaneTemplateResourceSpec `json:"spec"`
}
