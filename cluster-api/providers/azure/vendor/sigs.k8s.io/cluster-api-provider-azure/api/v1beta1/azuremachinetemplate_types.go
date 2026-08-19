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

package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
)

// Architecture represents the CPU architecture of the node.
// Its underlying type is a string and its value can be any of amd64, arm64.
type Architecture string

// Architecture constants.
const (
	// ArchitectureAmd64 represents the amd64 CPU architecture.
	ArchitectureAmd64 Architecture = "amd64"
	// ArchitectureArm64 represents the arm64 CPU architecture.
	ArchitectureArm64 Architecture = "arm64"
)

// OperatingSystem represents the operating system of the node.
// Its underlying type is a string and its value can be any of linux, windows.
type OperatingSystem string

// Operating system constants.
const (
	// OperatingSystemLinux represents the Linux operating system.
	OperatingSystemLinux OperatingSystem = "linux"
	// OperatingSystemWindows represents the Windows operating system.
	OperatingSystemWindows OperatingSystem = "windows"
)

// NodeInfo contains information about the node's architecture and operating system.
type NodeInfo struct {
	// Architecture is the CPU architecture of the node.
	// Its underlying type is a string and its value can be any of amd64, arm64.
	// +kubebuilder:validation:Enum=amd64;arm64
	// +optional
	Architecture Architecture `json:"architecture,omitempty"`
	// OperatingSystem is the operating system of the node.
	// Its underlying type is a string and its value can be any of linux, windows.
	// +kubebuilder:validation:Enum=linux;windows
	// +optional
	OperatingSystem OperatingSystem `json:"operatingSystem,omitempty"`
}

// AzureMachineTemplateStatus defines the observed state of AzureMachineTemplate.
type AzureMachineTemplateStatus struct {
	// Capacity defines the resource capacity for this machine.
	// This value is used for autoscaling from zero operations as defined in:
	// https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20210310-opt-in-autoscaling-from-zero.md
	// +optional
	Capacity corev1.ResourceList `json:"capacity,omitempty"`

	// NodeInfo contains information about the node's architecture and operating system.
	// This value is used for autoscaling from zero operations as defined in:
	// https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20210310-opt-in-autoscaling-from-zero.md
	// +optional
	NodeInfo *NodeInfo `json:"nodeInfo,omitempty"`
}

// AzureMachineTemplateSpec defines the desired state of AzureMachineTemplate.
type AzureMachineTemplateSpec struct {
	Template AzureMachineTemplateResource `json:"template"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=azuremachinetemplates,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:subresource:status

// AzureMachineTemplate is the Schema for the azuremachinetemplates API.
type AzureMachineTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AzureMachineTemplateSpec   `json:"spec,omitempty"`
	Status AzureMachineTemplateStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AzureMachineTemplateList contains a list of AzureMachineTemplates.
type AzureMachineTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AzureMachineTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AzureMachineTemplate{}, &AzureMachineTemplateList{})
}

// AzureMachineTemplateResource describes the data needed to create an AzureMachine from a template.
type AzureMachineTemplateResource struct {
	// +optional
	ObjectMeta clusterv1beta1.ObjectMeta `json:"metadata,omitempty"`
	// Spec is the specification of the desired behavior of the machine.
	Spec AzureMachineSpec `json:"spec"`
}
