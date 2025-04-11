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

// Package v1beta1 contains API types.
package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// VSphereResourceCPU defines Resource type CPU for VSphereMachines.
	VSphereResourceCPU corev1.ResourceName = "cpu"

	// VSphereResourceMemory defines Resource type memory for VSphereMachines.
	VSphereResourceMemory corev1.ResourceName = "memory"
)

// VSphereMachineTemplateSpec defines the desired state of VSphereMachineTemplate.
type VSphereMachineTemplateSpec struct {
	Template VSphereMachineTemplateResource `json:"template"`
}

// VSphereMachineTemplateStatus defines the observed state of VSphereMachineTemplate.
type VSphereMachineTemplateStatus struct {
	// Capacity defines the resource capacity for this VSphereMachineTemplate.
	// This value is used for autoscaling from zero operations as defined in:
	// https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20210310-opt-in-autoscaling-from-zero.md
	// +optional
	Capacity corev1.ResourceList `json:"capacity,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=vspheremachinetemplates,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:subresource:status

// VSphereMachineTemplate is the Schema for the vspheremachinetemplates API.
type VSphereMachineTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VSphereMachineTemplateSpec   `json:"spec,omitempty"`
	Status VSphereMachineTemplateStatus `json:"status,omitempty"`
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
