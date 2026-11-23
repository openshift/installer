// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package v1alpha6

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type VirtualMachineReservedProfileSpec struct {
	// +required

	// Name is the name of the reserved profile ID.
	Name string `json:"name"`
}

type VirtualMachineReservedProfileZoneStatus struct {
	// +required

	// Name is the name of the zone.
	Name string `json:"name"`

	// +optional
	// +kubebuilder:validation:Minimum=0

	// ReservedSlots describes the number of slots reserved for VMs that use
	// this reserved profile.
	ReservedSlots int32 `json:"reservedSlots,omitempty"`

	// +optional
	// +listType=map
	// +listMapKey=type

	// Conditions describes the observed conditions of the
	// VirtualMachineReservedProfile in the zone.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// VirtualMachineReservedProfileStatus defines the observed state of
// VirtualMachineReservedProfile.
type VirtualMachineReservedProfileStatus struct {
	// +optional
	// +listType=map
	// +listMapKey=type

	// Conditions describes the observed conditions of the
	// VirtualMachineReservedProfile.
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// +optional
	// +listType=map
	// +listMapKey=name

	// Zones describes the observed status of the class in one or more zones.
	Zones []VirtualMachineReservedProfileZoneStatus `json:"zones,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:storageversion
// +kubebuilder:subresource:status

// VirtualMachineReservedProfile is the schema for the
// virtualmachinereservedprofiles API and
// represents the desired state and observed status of a
// virtualmachinereservedprofiles resource.
type VirtualMachineReservedProfile struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualMachineReservedProfileSpec   `json:"spec,omitempty"`
	Status VirtualMachineReservedProfileStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// VirtualMachineReservedProfileList contains a list of
// VirtualMachineReservedProfile.
type VirtualMachineReservedProfileList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualMachineReservedProfile `json:"items"`
}

func init() {
	objectTypes = append(
		objectTypes,
		&VirtualMachineReservedProfile{},
		&VirtualMachineReservedProfileList{})
}
