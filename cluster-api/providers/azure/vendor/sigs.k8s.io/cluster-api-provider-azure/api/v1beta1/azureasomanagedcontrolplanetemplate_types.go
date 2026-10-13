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

// AzureASOManagedControlPlaneTemplateSpec defines the desired state of AzureASOManagedControlPlane.
type AzureASOManagedControlPlaneTemplateSpec struct {
	Template AzureASOManagedControlPlaneResource `json:"template"`
}

// AzureASOManagedControlPlaneResource defines the templated resource.
type AzureASOManagedControlPlaneResource struct {
	Spec AzureASOManagedControlPlaneTemplateResourceSpec `json:"spec,omitempty"`
}

// AzureASOManagedControlPlaneTemplateResourceSpec defines the desired state of the templated resource.
type AzureASOManagedControlPlaneTemplateResourceSpec struct {
	// Version is the Kubernetes version of the control plane. It fulfills Cluster API's control plane
	// provider contract.
	//+optional
	Version string `json:"version,omitempty"`

	// Resources are embedded ASO resources to be managed by this resource.
	//+optional
	Resources []runtime.RawExtension `json:"resources,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:storageversion

// AzureASOManagedControlPlaneTemplate is the Schema for the azureasomanagedcontrolplanetemplates API.
type AzureASOManagedControlPlaneTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec AzureASOManagedControlPlaneTemplateSpec `json:"spec,omitempty"`
}

//+kubebuilder:object:root=true

// AzureASOManagedControlPlaneTemplateList contains a list of AzureASOManagedControlPlaneTemplate.
type AzureASOManagedControlPlaneTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AzureASOManagedControlPlaneTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AzureASOManagedControlPlaneTemplate{}, &AzureASOManagedControlPlaneTemplateList{})
}
