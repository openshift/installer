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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// AzureASOManagedClusterTemplateSpec defines the desired state of AzureASOManagedClusterTemplate.
type AzureASOManagedClusterTemplateSpec struct {
	Template AzureASOManagedClusterTemplateResource `json:"template"`
}

// AzureASOManagedClusterTemplateResource defines the templated resource.
type AzureASOManagedClusterTemplateResource struct {
	Spec AzureASOManagedClusterTemplateResourceSpec `json:"spec,omitempty"`
}

// AzureASOManagedClusterTemplateResourceSpec defines the desired state of the templated resource.
type AzureASOManagedClusterTemplateResourceSpec struct {
	// Resources are embedded ASO resources to be managed by this resource.
	//+optional
	Resources []runtime.RawExtension `json:"resources,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:storageversion

// AzureASOManagedClusterTemplate is the Schema for the azureasomanagedclustertemplates API.
type AzureASOManagedClusterTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec AzureASOManagedClusterTemplateSpec `json:"spec,omitempty"`
}

//+kubebuilder:object:root=true

// AzureASOManagedClusterTemplateList contains a list of AzureASOManagedClusterTemplate.
type AzureASOManagedClusterTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AzureASOManagedClusterTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AzureASOManagedClusterTemplate{}, &AzureASOManagedClusterTemplateList{})
}
