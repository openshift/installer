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

// AzureASOManagedMachinePoolTemplateSpec defines the desired state of AzureASOManagedMachinePoolTemplate.
type AzureASOManagedMachinePoolTemplateSpec struct {
	Template AzureASOManagedControlPlaneResource `json:"template"`
}

// AzureASOManagedMachinePoolResource defines the templated resource.
type AzureASOManagedMachinePoolResource struct {
	Spec AzureASOManagedMachinePoolTemplateResourceSpec `json:"spec,omitempty"`
}

// AzureASOManagedMachinePoolTemplateResourceSpec defines the desired state of the templated resource.
type AzureASOManagedMachinePoolTemplateResourceSpec struct {
	// ProviderIDList is the list of cloud provider IDs for the instances. It fulfills Cluster API's machine
	// pool infrastructure provider contract.
	ProviderIDList []string `json:"providerIDList,omitempty"`

	// Resources are embedded ASO resources to be managed by this resource.
	//+optional
	Resources []runtime.RawExtension `json:"resources,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:storageversion

// AzureASOManagedMachinePoolTemplate is the Schema for the azureasomanagedmachinepooltemplates API.
type AzureASOManagedMachinePoolTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec AzureASOManagedMachinePoolTemplateSpec `json:"spec,omitempty"`
}

//+kubebuilder:object:root=true

// AzureASOManagedMachinePoolTemplateList contains a list of AzureASOManagedMachinePoolTemplate.
type AzureASOManagedMachinePoolTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AzureASOManagedMachinePoolTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AzureASOManagedMachinePoolTemplate{}, &AzureASOManagedMachinePoolTemplateList{})
}
