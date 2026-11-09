/*
Copyright 2022 The Kubernetes Authors.

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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
)

// Architecture represents the CPU architecture of the node.
// Its underlying type is a string and its value can be any of amd64, arm64.
type Architecture string

// Architecture constants.
const (
	ArchitectureAmd64 Architecture = "amd64"
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

// AWSMachineTemplateStatus defines a status for an AWSMachineTemplate.
type AWSMachineTemplateStatus struct {
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

	// Conditions defines current service state of the AWSMachineTemplate.
	// +optional
	Conditions clusterv1beta1.Conditions `json:"conditions,omitempty"`

	// v1beta2 groups all the fields that will be added or modified in AWSMachineTemplate's status with the V1Beta2 version.
	// +optional
	V1Beta2 *AWSMachineTemplateV1Beta2Status `json:"v1beta2,omitempty"`
}

// AWSMachineTemplateV1Beta2Status groups all the fields that will be added or modified in AWSMachineTemplate with the V1Beta2 version.
// See https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more context.
type AWSMachineTemplateV1Beta2Status struct {
	// conditions represents the observations of an AWSMachineTemplate's current state.
	// +optional
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=32
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// AWSMachineTemplateSpec defines the desired state of AWSMachineTemplate.
type AWSMachineTemplateSpec struct {
	Template AWSMachineTemplateResource `json:"template"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=awsmachinetemplates,scope=Namespaced,categories=cluster-api,shortName=awsmt
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +k8s:defaulter-gen=true

// AWSMachineTemplate is the schema for the Amazon EC2 Machine Templates API.
type AWSMachineTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AWSMachineTemplateSpec   `json:"spec,omitempty"`
	Status AWSMachineTemplateStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AWSMachineTemplateList contains a list of AWSMachineTemplate.
type AWSMachineTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSMachineTemplate `json:"items"`
}

// AWSMachineTemplateResource describes the data needed to create am AWSMachine from a template.
type AWSMachineTemplateResource struct {
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	ObjectMeta clusterv1beta1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the specification of the desired behavior of the machine.
	Spec AWSMachineSpec `json:"spec"`
}

// GetConditions returns the observations of the operational state of the AWSMachineTemplate resource.
func (r *AWSMachineTemplate) GetConditions() clusterv1beta1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the underlying service state of the AWSMachineTemplate to the predescribed clusterv1beta1.Conditions.
func (r *AWSMachineTemplate) SetConditions(conditions clusterv1beta1.Conditions) {
	r.Status.Conditions = conditions
}

// GetV1Beta2Conditions returns the set of conditions for this object.
func (r *AWSMachineTemplate) GetV1Beta2Conditions() []metav1.Condition {
	if r.Status.V1Beta2 == nil {
		return nil
	}
	return r.Status.V1Beta2.Conditions
}

// SetV1Beta2Conditions sets conditions for an API object.
func (r *AWSMachineTemplate) SetV1Beta2Conditions(conditions []metav1.Condition) {
	if r.Status.V1Beta2 == nil {
		r.Status.V1Beta2 = &AWSMachineTemplateV1Beta2Status{}
	}
	r.Status.V1Beta2.Conditions = conditions
}

func init() {
	SchemeBuilder.Register(&AWSMachineTemplate{}, &AWSMachineTemplateList{})
}
