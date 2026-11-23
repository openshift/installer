// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	// PoolMoIDs are the managed object ID of the vSphere ResourcePools
	// in an individual vSphere Zone. A zone may be comprised of
	// multiple ResourcePools.
	PoolMoIDs []string `json:"poolMoIDs,omitempty"`

	// +optional
	// FolderMoID is the managed object ID of the vSphere Folder for a
	// Namespace.
	FolderMoID string `json:"folderMoID,omitempty"`
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
}

// ZoneStatus defines the observed state of Zone.
type ZoneStatus struct {
	// +optional
	// Conditions describes the observed conditions of the Zone
	Conditions []metav1.Condition `json:"conditions,omitempty"`
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
	objectTypes = append(objectTypes, &Zone{}, &ZoneList{})
}
