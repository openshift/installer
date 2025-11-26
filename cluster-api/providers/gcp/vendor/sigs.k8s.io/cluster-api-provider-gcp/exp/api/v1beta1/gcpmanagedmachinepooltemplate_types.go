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
)

// GCPManagedMachinePoolTemplateSpec defines the desired state of GCPManagedMachinePoolTemplate.
type GCPManagedMachinePoolTemplateSpec struct {
	Template GCPManagedMachinePoolTemplateResource `json:"template"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=gcpmanagedmachinepooltemplates,scope=Namespaced,categories=cluster-api,shortName=ammpt
// +kubebuilder:storageversion

// GCPManagedMachinePoolTemplate is the Schema for the GCPManagedMachinePoolTemplates API.
type GCPManagedMachinePoolTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec GCPManagedMachinePoolTemplateSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// GCPManagedMachinePoolTemplateList contains a list of GCPManagedMachinePoolTemplates.
type GCPManagedMachinePoolTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GCPManagedMachinePoolTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GCPManagedMachinePoolTemplate{}, &GCPManagedMachinePoolTemplateList{})
}

// GCPManagedMachinePoolTemplateResource describes the data needed to create an GCPManagedCluster from a template.
type GCPManagedMachinePoolTemplateResource struct {
	Spec GCPManagedMachinePoolTemplateResourceSpec `json:"spec"`
}
