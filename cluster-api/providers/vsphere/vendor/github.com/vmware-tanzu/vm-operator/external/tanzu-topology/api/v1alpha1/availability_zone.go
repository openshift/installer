// Copyright (c) 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NamespaceInfo contains identifying information about the vSphere resources
// used to represent a Kuberentes namespace on individual vSphere Clusters.
type NamespaceInfo struct {
	// PoolMoId is the managed object ID of the vSphere ResourcePool for a
	// Namespace on an individual vSphere Cluster.
	PoolMoId string `json:"poolMoId,omitempty"`

	// FolderMoId is the managed object ID of the vSphere Folder for a
	// Namespace. Folders are global and not per-vSphere Cluster, but the
	// FolderMoId is stored here, alongside the PoolMoId for convenience.
	FolderMoId string `json:"folderMoId,omitempty"`
}

// AvailabilityZoneSpec defines the desired state of AvailabilityZone.
type AvailabilityZoneSpec struct {
	// ClusterComputeResourceMoId is the managed object ID of the vSphere
	// ClusterComputeResource represented by this availability zone.
	ClusterComputeResourceMoId string `json:"clusterComputeResourceMoId,omitempty"`

	// Namespaces is a map that enables querying information about the vSphere
	// objects that make up a Kubernetes Namespace based on its name.
	Namespaces map[string]NamespaceInfo `json:"namespaces,omitempty"`
}

// AvailabilityZoneStatus defines the observed state of AvailabilityZone.
type AvailabilityZoneStatus struct {
}

// AvailabilityZone is the schema for the AvailabilityZone resource for the
// vSphere topology API.
//
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=availabilityzones,scope=Cluster,shortName=az
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
type AvailabilityZone struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AvailabilityZoneSpec   `json:"spec,omitempty"`
	Status AvailabilityZoneStatus `json:"status,omitempty"`
}

// AvailabilityZoneList contains a list of AvailabilityZone resources.
//
// +kubebuilder:object:root=true
type AvailabilityZoneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AvailabilityZone `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AvailabilityZone{}, &AvailabilityZoneList{})
}
