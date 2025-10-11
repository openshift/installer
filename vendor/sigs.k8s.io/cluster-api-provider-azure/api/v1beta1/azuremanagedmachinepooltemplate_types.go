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

// AzureManagedMachinePoolTemplateSpec defines the desired state of AzureManagedMachinePoolTemplate.
type AzureManagedMachinePoolTemplateSpec struct {
	Template AzureManagedMachinePoolTemplateResource `json:"template"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=azuremanagedmachinepooltemplates,scope=Namespaced,categories=cluster-api,shortName=ammpt
// +kubebuilder:storageversion
// +kubebuilder:deprecatedversion:warning="AzureManagedMachinePoolTemplate and the AzureManaged API are deprecated. Please migrate to infrastructure.cluster.x-k8s.io/v1beta1 AzureASOManagedMachinePoolTemplate and related AzureASOManaged resources instead."

// AzureManagedMachinePoolTemplate is the Schema for the AzureManagedMachinePoolTemplates API.
type AzureManagedMachinePoolTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec AzureManagedMachinePoolTemplateSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// AzureManagedMachinePoolTemplateList contains a list of AzureManagedMachinePoolTemplates.
type AzureManagedMachinePoolTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AzureManagedMachinePoolTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AzureManagedMachinePoolTemplate{}, &AzureManagedMachinePoolTemplateList{})
}

// AzureManagedMachinePoolTemplateResource describes the data needed to create an AzureManagedCluster from a template.
type AzureManagedMachinePoolTemplateResource struct {
	Spec AzureManagedMachinePoolTemplateResourceSpec `json:"spec"`
}
