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

// Package v1beta2 contains API types.
package v1beta2

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

// Architecture represents the CPU architecture of the node.
// Its underlying type is a string and its value can be any of amd64, arm64, s390x, ppc64le.
// +kubebuilder:validation:Enum=amd64;arm64;s390x;ppc64le
// +enum
type Architecture string

const (
	// ArchitectureAmd64 is the AMD64 architecture.
	ArchitectureAmd64 Architecture = "amd64"
	// ArchitectureArm64 is the ARM64 architecture.
	ArchitectureArm64 Architecture = "arm64"
	// ArchitectureS390x is the S390X architecture.
	ArchitectureS390x Architecture = "s390x"
	// ArchitecturePpc64le is the PPC64LE architecture.
	ArchitecturePpc64le Architecture = "ppc64le"
)

const (
	// VMwareSystemOSArchPropertyKey is the key for the architecture property in the
	// ClusterVirtualMachineImage's vmwareSystemProperties. This is defined by VM Operator.
	VMwareSystemOSArchPropertyKey = "vmware-system.tkr.os-arch"
	// VMwareSystemOSTypePropertyKey is the key for the operating system type property in the
	// ClusterVirtualMachineImage's vmwareSystemProperties. This is defined by VM Operator.
	VMwareSystemOSTypePropertyKey = "vmware-system.tkr.os-type"
)

// OperatingSystem represents the operating system of the node.
// Its underlying type is a string and its value can be any of linux, windows.
// +kubebuilder:validation:Enum=linux;windows
// +enum
type OperatingSystem string

const (
	// OperatingSystemLinux is the Linux operating system.
	OperatingSystemLinux OperatingSystem = "linux"
	// OperatingSystemWindows is the Windows operating system.
	OperatingSystemWindows OperatingSystem = "windows"
)

// VSphereMachineTemplateSpec defines the desired state of VSphereMachineTemplate.
type VSphereMachineTemplateSpec struct {
	// template defines the desired state of VSphereMachineTemplate.
	// +required
	Template VSphereMachineTemplateResource `json:"template,omitempty,omitzero"`
}

// VSphereMachineTemplateStatus defines the observed state of VSphereMachineTemplate.
// +kubebuilder:validation:MinProperties=1
type VSphereMachineTemplateStatus struct {
	// capacity defines the resource capacity for this VSphereMachineTemplate.
	// This value is used for autoscaling from zero operations as defined in:
	// https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20210310-opt-in-autoscaling-from-zero.md
	// +optional
	Capacity corev1.ResourceList `json:"capacity,omitempty,omitzero"`
	// nodeInfo defines the node's architecture and operating system.
	// This value is used for autoscaling from zero operations as defined in:
	// https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20210310-opt-in-autoscaling-from-zero.md#implementation-detailsnotesconstraints
	// +optional
	NodeInfo NodeInfo `json:"nodeInfo,omitempty,omitzero"`
}

// NodeInfo contains information about the node's architecture and operating system.
// +kubebuilder:validation:MinProperties=1
type NodeInfo struct {
	// architecture is the CPU architecture of the node.
	// Its underlying type is a string and its value can be any of amd64, arm64, s390x, ppc64le.
	// +optional
	Architecture Architecture `json:"architecture,omitempty"`
	// operatingSystem is a string representing the operating system of the node.
	// This may be a string like 'linux' or 'windows'.
	// +optional
	OperatingSystem OperatingSystem `json:"operatingSystem,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=vspheremachinetemplates,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
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

	// status is the observed state of VSphereMachineTemplate.
	// +optional
	Status VSphereMachineTemplateStatus `json:"status,omitempty,omitzero"`
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
