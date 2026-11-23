// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package v1alpha2

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// VirtualMachineVolumeProvisioningMode is the type used to express the
// desired or observed provisioning mode for a virtual machine disk.
//
// +kubebuilder:validation:Enum=Thin;Thick;ThickEagerZero
type VirtualMachineVolumeProvisioningMode string

const (
	VirtualMachineVolumeProvisioningModeThin           VirtualMachineVolumeProvisioningMode = "Thin"
	VirtualMachineVolumeProvisioningModeThick          VirtualMachineVolumeProvisioningMode = "Thick"
	VirtualMachineVolumeProvisioningModeThickEagerZero VirtualMachineVolumeProvisioningMode = "ThickEagerZero"
)

// VirtualMachineVolume represents a named volume in a VM.
type VirtualMachineVolume struct {
	// Name represents the volume's name. Must be a DNS_LABEL and unique within
	// the VM.
	Name string `json:"name"`

	// VirtualMachineVolumeSource represents the location and type of a volume
	// to mount.
	VirtualMachineVolumeSource `json:",inline"`
}

// VirtualMachineVolumeSource represents the source location of a volume to
// mount. Only one of its members may be specified.
type VirtualMachineVolumeSource struct {
	// PersistentVolumeClaim represents a reference to a PersistentVolumeClaim
	// in the same namespace.
	//
	// More information is available at
	// https://kubernetes.io/docs/concepts/storage/persistent-volumes#persistentvolumeclaims.
	//
	// +optional
	PersistentVolumeClaim *PersistentVolumeClaimVolumeSource `json:"persistentVolumeClaim,omitempty"`
}

// PersistentVolumeClaimVolumeSource is a composite for the Kubernetes
// corev1.PersistentVolumeClaimVolumeSource and instance storage options.
type PersistentVolumeClaimVolumeSource struct {
	corev1.PersistentVolumeClaimVolumeSource `json:",inline" yaml:",inline"`

	// InstanceVolumeClaim is set if the PVC is backed by instance storage.
	// +optional
	InstanceVolumeClaim *InstanceVolumeClaimVolumeSource `json:"instanceVolumeClaim,omitempty"`
}

// InstanceVolumeClaimVolumeSource contains information about the instance
// storage volume claimed as a PVC.
type InstanceVolumeClaimVolumeSource struct {
	// StorageClass is the name of the Kubernetes StorageClass that provides
	// the backing storage for this instance storage volume.
	StorageClass string `json:"storageClass"`

	// Size is the size of the requested instance storage volume.
	Size resource.Quantity `json:"size"`
}

// VirtualMachineVolumeStatus defines the observed state of a
// VirtualMachineVolume instance.
type VirtualMachineVolumeStatus struct {
	// Name is the name of the attached volume.
	Name string `json:"name"`

	// Attached represents whether a volume has been successfully attached to
	// the VirtualMachine or not.
	// +optional
	Attached bool `json:"attached,omitempty"`

	// DiskUUID represents the underlying virtual disk UUID and is present when
	// attachment succeeds.
	// +optional
	DiskUUID string `json:"diskUUID,omitempty"`

	// Error represents the last error seen when attaching or detaching a
	// volume.  Error will be empty if attachment succeeds.
	// +optional
	Error string `json:"error,omitempty"`
}
