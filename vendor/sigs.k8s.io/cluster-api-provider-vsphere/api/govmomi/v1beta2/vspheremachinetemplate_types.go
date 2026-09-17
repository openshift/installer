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

package v1beta2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VSphereMachineTemplateSpec defines the desired state of VSphereMachineTemplate.
type VSphereMachineTemplateSpec struct {
	// template defines the desired state of VSphereMachineTemplate.
	// +required
	Template VSphereMachineTemplateResource `json:"template,omitzero"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=vspheremachinetemplates,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="ClusterClass",type="string",JSONPath=`.metadata.ownerReferences[?(@.kind=="ClusterClass")].name`,description="Name of the ClusterClass owning this template"
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=`.metadata.ownerReferences[?(@.kind=="Cluster")].name`,description="Name of the Cluster owning this template"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of VSphereMachineTemplate"

// VSphereMachineTemplate is the Schema for the vspheremachinetemplates API.
type VSphereMachineTemplate struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is the desired state of VSphereMachineTemplate.
	// +optional
	Spec VSphereMachineTemplateSpec `json:"spec,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// VSphereMachineTemplateList contains a list of VSphereMachineTemplate.
type VSphereMachineTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VSphereMachineTemplate `json:"items"`
}

func init() {
	objectTypes = append(objectTypes, &VSphereMachineTemplate{}, &VSphereMachineTemplateList{})
}
