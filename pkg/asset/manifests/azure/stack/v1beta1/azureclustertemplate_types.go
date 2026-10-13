/*
Copyright 2022 The Kubernetes Authors.

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

// AzureClusterTemplateSpec defines the desired state of AzureClusterTemplate.
type AzureClusterTemplateSpec struct {
	Template AzureClusterTemplateResource `json:"template"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=azureclustertemplates,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion

// AzureClusterTemplate is the Schema for the azureclustertemplates API.
type AzureClusterTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec AzureClusterTemplateSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// AzureClusterTemplateList contains a list of AzureClusterTemplate.
type AzureClusterTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AzureClusterTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AzureClusterTemplate{}, &AzureClusterTemplateList{})
}

// AzureClusterTemplateResource describes the data needed to create an AzureCluster from a template.
type AzureClusterTemplateResource struct {
	Spec AzureClusterTemplateResourceSpec `json:"spec"`
}
