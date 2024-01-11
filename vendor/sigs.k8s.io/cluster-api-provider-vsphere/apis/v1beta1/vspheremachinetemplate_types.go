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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VSphereMachineTemplateSpec defines the desired state of VSphereMachineTemplate
type VSphereMachineTemplateSpec struct {
	Template VSphereMachineTemplateResource `json:"template"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=vspheremachinetemplates,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion

// VSphereMachineTemplate is the Schema for the vspheremachinetemplates API
type VSphereMachineTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec VSphereMachineTemplateSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// VSphereMachineTemplateList contains a list of VSphereMachineTemplate
type VSphereMachineTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VSphereMachineTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VSphereMachineTemplate{}, &VSphereMachineTemplateList{})
}
