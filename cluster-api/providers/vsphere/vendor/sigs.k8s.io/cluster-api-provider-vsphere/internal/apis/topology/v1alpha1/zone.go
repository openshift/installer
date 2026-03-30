// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// ZoneConditionPersistentVolumeClaimsExist indicates that PVCs exist in the
	// Zone.
	ZoneConditionPersistentVolumeClaimsExist = "ZonePersistentVolumeClaimsExist"

	// ZoneConditionPodsExist indicates that Pods exist in the Zone.
	ZoneConditionPodsExist = "ZonePodsExist"

	// ZoneConditionVirtualMachinesExist indicates that VirtualMachines exist
	// in the Zone.
	ZoneConditionVirtualMachinesExist = "ZoneVirtualMachinesExist"

	// ZoneConditionNamespaceResourcePoolReconciled indicates that the
	// Namespace ResourcePool has been reconciled.
	ZoneConditionNamespaceResourcePoolReconciled = "ZoneNamespaceResourcePoolReconciled"
)

// AvailabilityZoneReference describes a reference to the cluster scoped
// AvailabilityZone object.
type AvailabilityZoneReference struct {
	// APIVersion defines the versioned schema of this reference to the cluster scoped
	// AvailabilityZone object.
	APIVersion string `json:"apiVersion"`

	// Name is the name of  the cluster scoped AvailabilityZone to refer to.
	Name string `json:"name"`
}

// VSphereEntityInfo contains the managed object IDs associated with
// a vSphere entity
type VSphereEntityInfo struct {
	// +optional
	// ClusterMoIDs are the managed object IDs of the vSphere Clusters in an
	// individual vSphere Zone. A zone may be comprised of multiple Clusters.
	ClusterMoIDs []string `json:"clusterMoIDs,omitempty"`

	// +optional
	// PoolMoIDs are the managed object ID of the vSphere ResourcePools
	// in an individual vSphere Zone. A zone may be comprised of
	// multiple ResourcePools.
	PoolMoIDs []string `json:"poolMoIDs,omitempty"`

	// +optional
	// FolderMoID is the managed object ID of the vSphere Folder for a
	// Namespace.
	FolderMoID string `json:"folderMoID,omitempty"`
}

// VirtualMachineClassAllocationInfo describes the definition of allocations
// for Virtual Machines of a given class.
type VirtualMachineClassAllocationInfo struct {
	// +optional
	// Identifier of the Virtual Machine class used for allocation.
	ReservedVMClass string `json:"reservedVmClass,omitempty"`

	// +optional
	// Number of instances of given Virtual Machine class.
	Count int64 `json:"count,omitempty"`
}

// ZoneSpec contains identifying information about the
// vSphere resources used to represent a Kubernetes namespace on individual
// vSphere Zones.
type ZoneSpec struct {
	// Namespace contains ResourcePool and folder moIDs to represent the namespace
	Namespace VSphereEntityInfo `json:"namespace,omitempty"`

	// VSpherePods contains ResourcePool and folder moIDs to represent vSpherePods
	// entity within the namespace
	VSpherePods VSphereEntityInfo `json:"vSpherePods,omitempty"`

	// ManagedVMs contains ResourcePool and folder moIDs to represent managedVMs
	// entity within the namespace
	ManagedVMs VSphereEntityInfo `json:"managedVMs,omitempty"`

	// Zone is a reference to the cluster scoped AvailabilityZone this
	// Zone is derived from.
	Zone AvailabilityZoneReference `json:"availabilityZoneReference"`

	// +optional
	// Guaranteed number of reserved Virtual Machine class instances that are available for the
	// namespace in this zone.
	VirtualMachineReservations []VirtualMachineClassAllocationInfo `json:"virtualMachineReservations,omitempty"`

	// +optional
	// CPU limit (in megahertz) for the namespace in this zone in addition to the limits specified as part of
	// reserved Virtual Machine classes.
	CPULimitMHz int64 `json:"cpuLimitMHz,omitempty"`

	// +optional
	// CPU reservation (in megahertz) for the namespace in this zone, for VMs
	// that are not using reserved Virtual Machine class instances.
	CPUReservationMHz int64 `json:"cpuReservationMHz,omitempty"`

	// +optional
	// Memory limit (in mebibytes) for the namespace in this zone in addition
	// to the limits specified as part of reserved Virtual Machine classes.
	MemoryLimitMiB int64 `json:"memoryLimitMiB,omitempty"`

	// +optional
	// Memory reservation (in mebibytes) for the namespace in this zone, for
	// VMs that are not using reserved Virtual Machine class instances.
	MemoryReservationMiB int64 `json:"memoryReservationMiB,omitempty"`

	// +optional
	// Determines whether workloads that don't use a reserved Virtual Machine class
	// instance can use a DirectPath device.
	DisallowUnreservedDirectPathUsage bool `json:"disallowUnreservedDirectPathUsage,omitempty"`

	// +optional
	// AllowedClusterComputeResourceMoIDs are the managed object IDs of the vSphere
	// ClusterComputeResources in this vSphere Zone on which workloads in this Supervisor
	// Namespace can be placed on. If empty, all the vSphere Clusters in the vSphere Zone are
	// candidates to place the workloads in this vSphere Namespace.
	AllowedClusterComputeResourceMoIDs []string `json:"allowedClusterComputeResourceMoIDs,omitempty"`
}

// ZoneStatus defines the observed state of Zone.
type ZoneStatus struct {
	// +optional
	// Conditions describes the observed conditions of the Zone
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// +optional
	// MarkedForRemoval describes if the Zone is marked for removal from the
	// Namespace.
	MarkedForRemoval bool `json:"markedForRemoval,omitempty"`
}

// Zone is the schema for the Zone resource for the vSphere topology API.
//
// A Zone is the zone the k8s namespace is confined to. That is workloads will
// be limited to the Zones in the namespace.  For more information about
// availability zones, refer to:
// https://kubernetes.io/docs/setup/best-practices/multiple-zones/
//
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=zones,scope=Namespaced
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
type Zone struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ZoneSpec   `json:"spec,omitempty"`
	Status ZoneStatus `json:"status,omitempty"`
}

// ZoneList contains a list of Zone resources.
//
// +kubebuilder:object:root=true
type ZoneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Zone `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Zone{}, &ZoneList{})
}
