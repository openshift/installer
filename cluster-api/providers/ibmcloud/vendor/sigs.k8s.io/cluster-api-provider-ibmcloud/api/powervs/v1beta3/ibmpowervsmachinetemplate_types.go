/*
Copyright 2026 The Kubernetes Authors.

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

package v1beta3

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
)

func init() {
	objectTypes = append(objectTypes, &IBMPowerVSMachineTemplate{}, &IBMPowerVSMachineTemplateList{})
}

// IBMPowerVSMachineTemplateSpec defines the desired state of IBMPowerVSMachineTemplate.
type IBMPowerVSMachineTemplateSpec struct {
	// template is the IBMPowerVSMachineTemplateResource.
	// +required
	Template IBMPowerVSMachineTemplateResource `json:"template,omitempty,omitzero"`
}

// IBMPowerVSMachineTemplateStatus defines the observed state of IBMPowerVSMachineTemplate.
// +kubebuilder:validation:MinProperties=1
type IBMPowerVSMachineTemplateStatus struct {
	// capacity defines the resource capacity for this machine.
	// This value is used for autoscaling from zero operations as defined in:
	// https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20210310-opt-in-autoscaling-from-zero.md
	// +optional
	Capacity corev1.ResourceList `json:"capacity,omitempty"`
}

// +kubebuilder:subresource:status
// +kubebuilder:object:root=true
// +kubebuilder:storageversion

// IBMPowerVSMachineTemplate is the Schema for the ibmpowervsmachinetemplates API.
type IBMPowerVSMachineTemplate struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of IBMPowerVSMachineTemplate
	// +required
	Spec IBMPowerVSMachineTemplateSpec `json:"spec,omitempty,omitzero"`

	// status defines the observed state of IBMPowerVSMachineTemplate
	// +optional
	Status IBMPowerVSMachineTemplateStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// IBMPowerVSMachineTemplateList contains a list of IBMPowerVSMachineTemplate.
type IBMPowerVSMachineTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []IBMPowerVSMachineTemplate `json:"items"`
}

// IBMPowerVSMachineTemplateResource holds the IBMPowerVSMachine spec.
// +kubebuilder:validation:MinProperties=1
type IBMPowerVSMachineTemplateResource struct {
	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	ObjectMeta clusterv1.ObjectMeta `json:"metadata,omitempty,omitzero"`
	// spec is the IBMPowerVSMachineSpec.
	// +optional
	Spec IBMPowerVSMachineSpec `json:"spec,omitempty,omitzero"`
}
