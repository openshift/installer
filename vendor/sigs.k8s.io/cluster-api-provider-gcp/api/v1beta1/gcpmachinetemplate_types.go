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
)

// Architecture represents the CPU architecture of the node.
type Architecture string

const (
	// ArchitectureAmd64 represents the amd64 architecture.
	ArchitectureAmd64 Architecture = "amd64"
	// ArchitectureArm64 represents the arm64 architecture.
	ArchitectureArm64 Architecture = "arm64"
)

// NodeInfo defines the node information for this machine.
type NodeInfo struct {
	// Architecture defines the hardware architecture (e.g., amd64, arm64).
	// +kubebuilder:validation:Enum=amd64;arm64
	// +optional
	Architecture Architecture `json:"architecture,omitempty"`

	// OperatingSystem defines the operating system (e.g., linux, windows).
	// +kubebuilder:validation:Enum=linux;windows
	// +optional
	OperatingSystem corev1.OSName `json:"operatingSystem,omitempty"`
}

// GCPMachineTemplateStatus defines the observed state of a GCPMachineTemplate.
type GCPMachineTemplateStatus struct {
	// Capacity defines the resource capacity for this machine.
	// This value is used for autoscaling-from-zero operations as defined in:
	// https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20210310-opt-in-autoscaling-from-zero.md
	// +optional
	Capacity corev1.ResourceList `json:"capacity,omitempty"`

	// NodeInfo defines the node architecture and operating system.
	// This value is used for autoscaling-from-zero operations as defined in:
	// https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20210310-opt-in-autoscaling-from-zero.md
	// +optional
	NodeInfo *NodeInfo `json:"nodeInfo,omitempty"`
}

// GCPMachineTemplateSpec defines the desired state of GCPMachineTemplate.
type GCPMachineTemplateSpec struct {
	Template GCPMachineTemplateResource `json:"template"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=gcpmachinetemplates,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:subresource:status

// GCPMachineTemplate is the Schema for the gcpmachinetemplates API.
type GCPMachineTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GCPMachineTemplateSpec   `json:"spec,omitempty"`
	Status GCPMachineTemplateStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// GCPMachineTemplateList contains a list of GCPMachineTemplate.
type GCPMachineTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GCPMachineTemplate `json:"items"`
}

func init() {
	objectTypes = append(objectTypes, &GCPMachineTemplate{}, &GCPMachineTemplateList{})
}
