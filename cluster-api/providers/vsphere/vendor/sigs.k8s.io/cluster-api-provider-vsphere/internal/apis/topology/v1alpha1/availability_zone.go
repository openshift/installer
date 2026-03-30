// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// AvailabiltyZoneConditionInUse indiciates that the AvailabilityZone is in
	// use.
	//
	//	See Zone.Status for more details.
	AvailabiltyZoneConditionInUse = "AvailabiltyZoneInUse"
)

// NamespaceInfo contains identifying information about the vSphere resources
// used to represent a Kubernetes namespace on individual vSphere Zones.
type NamespaceInfo struct {
	// PoolMoId is the managed object ID of the vSphere ResourcePool for a
	// Namespace on an individual vSphere Cluster.
	PoolMoId string `json:"poolMoId,omitempty"`

	// PoolMoIDs are the managed object ID of the vSphere ResourcePools for a
	// Namespace in an individual vSphere Zone. A zone may be comprised of
	// multiple ResourcePools.
	PoolMoIDs []string `json:"poolMoIDs,omitempty"`

	// FolderMoId is the managed object ID of the vSphere Folder for a
	// Namespace. Folders are global and not per-vSphere Cluster, but the
	// FolderMoId is stored here, alongside the PoolMoId for convenience.
	FolderMoId string `json:"folderMoId,omitempty"`
}

// SystemInfo contains identifying information about the vSphere resources
// holding or otherwise used by a Supervisor's system-managed resources in an
// individual vSphere Zone.
type SystemInfo struct {
	// PoolMoIDs are the managed object IDs of the vSphere ResourcePools for
	// this Supervisor's system-managed resources within a zone. A zone may be
	// comprised of multiple ResourcePools.
	// +optional
	PoolMoIDs []string `json:"poolMoIDs,omitempty"`

	// FolderMoID is the managed object ID of the vSphere Folder for this
	// Supervisor's system-managed resources. Folders are global and not scoped to
	// a zone or vSphere cluster, but the FolderMoID is stored here for
	// convenience.
	// +optional
	FolderMoID string `json:"folderMoID,omitempty"`
}

// AvailabilityZoneSpec defines the desired state of AvailabilityZone.
type AvailabilityZoneSpec struct {
	// ClusterComputeResourceMoId is the managed object ID of the vSphere
	// ClusterComputeResource represented by this availability zone.
	ClusterComputeResourceMoId string `json:"clusterComputeResourceMoId,omitempty"`

	// ClusterComputeResourceMoIDs are the managed object IDs of the vSphere
	// ClusterComputeResources represented by this availability zone.
	ClusterComputeResourceMoIDs []string `json:"clusterComputeResourceMoIDs,omitempty"`

	// SystemInfo holds identifying information about vSphere resource grouping
	// objects used by the Supervisor's system objects. These are typically
	// the top Supervisor-level "Namespaces" Resource Pools and Folder.
	// +optional
	SystemInfo *SystemInfo `json:"systemInfo,omitempty"`

	// Namespaces is a map that enables querying information about the vSphere
	// objects that make up a Kubernetes Namespace based on its name.
	Namespaces map[string]NamespaceInfo `json:"namespaces,omitempty"`
}

// AvailabilityZoneStatus defines the observed state of AvailabilityZone.
type AvailabilityZoneStatus struct {
	// +optional
	// Conditions describes the observed conditions of the AvailabilityZone
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// +optional
	// MarkedForRemoval describes if the AvailabiltyZone is marked for removal.
	MarkedForRemoval bool `json:"markedForRemoval,omitempty"`
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
