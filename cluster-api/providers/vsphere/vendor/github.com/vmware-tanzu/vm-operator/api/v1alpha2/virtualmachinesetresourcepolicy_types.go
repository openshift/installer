// Copyright (c) 2023 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ResourcePoolSpec defines a Logical Grouping of workloads that share resource
// policies.
type ResourcePoolSpec struct {
	// Name describes the name of the ResourcePool grouping.
	// +optional
	Name string `json:"name,omitempty"`

	// Reservations describes the guaranteed resources reserved for the
	// ResourcePool.
	// +optional
	Reservations VirtualMachineResourceSpec `json:"reservations,omitempty"`

	// Limits describes the limit to resources available to the ResourcePool.
	// +optional
	Limits VirtualMachineResourceSpec `json:"limits,omitempty"`
}

// VirtualMachineSetResourcePolicySpec defines the desired state of
// VirtualMachineSetResourcePolicy.
type VirtualMachineSetResourcePolicySpec struct {
	ResourcePool        ResourcePoolSpec `json:"resourcePool,omitempty"`
	Folder              string           `json:"folder,omitempty"`
	ClusterModuleGroups []string         `json:"clusterModuleGroups,omitempty"`
}

// VirtualMachineSetResourcePolicyStatus defines the observed state of
// VirtualMachineSetResourcePolicy.
type VirtualMachineSetResourcePolicyStatus struct {
	ClusterModules []VSphereClusterModuleStatus `json:"clustermodules,omitempty"`
}

// VSphereClusterModuleStatus describes the observed state of a vSphere
// cluster module.
type VSphereClusterModuleStatus struct {
	GroupName   string `json:"groupName"`
	ModuleUuid  string `json:"moduleUUID"` //nolint:revive,stylecheck
	ClusterMoID string `json:"clusterMoID"`
}

// +kubebuilder:object:root=true
// +kubebuilder:storageversion:false
// +kubebuilder:subresource:status

// VirtualMachineSetResourcePolicy is the Schema for the virtualmachinesetresourcepolicies API.
type VirtualMachineSetResourcePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualMachineSetResourcePolicySpec   `json:"spec,omitempty"`
	Status VirtualMachineSetResourcePolicyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// VirtualMachineSetResourcePolicyList contains a list of VirtualMachineSetResourcePolicy.
type VirtualMachineSetResourcePolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualMachineSetResourcePolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VirtualMachineSetResourcePolicy{}, &VirtualMachineSetResourcePolicyList{})
}
