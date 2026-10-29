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

package hub

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// +kubebuilder:validation:Enum=IndependentNonPersistent;IndependentPersistent;NonPersistent;Persistent

type VolumeDiskMode string

const (
	VolumeDiskModeIndependentNonPersistent VolumeDiskMode = "IndependentNonPersistent"
	VolumeDiskModeIndependentPersistent    VolumeDiskMode = "IndependentPersistent"
	VolumeDiskModeNonPersistent            VolumeDiskMode = "NonPersistent"
	VolumeDiskModePersistent               VolumeDiskMode = "Persistent"
)

// +kubebuilder:validation:Enum=MultiWriter;None

type VolumeSharingMode string

const (
	VolumeSharingModeMultiWriter VolumeSharingMode = "MultiWriter"
	VolumeSharingModeNone        VolumeSharingMode = "None"
)

// +kubebuilder:validation:Enum=OracleRAC;MicrosoftWSFC

type VolumeApplicationType string

const (
	VolumeApplicationTypeOracleRAC     VolumeApplicationType = "OracleRAC"
	VolumeApplicationTypeMicrosoftWSFC VolumeApplicationType = "MicrosoftWSFC"
)

// VirtualMachineVolume represents a named volume in a VM.
type VirtualMachineVolume struct {
	// Name represents the volume's name. Must be a DNS_LABEL and unique within
	// the VM.
	Name string `json:"name"`

	// VirtualMachineVolumeSource represents the location and type of a volume
	// to mount.
	VirtualMachineVolumeSource `json:",inline"`

	// +optional
	// +kubebuilder:default=true

	// Removable describes whether or not this volume may be removed from
	// spec.volumes. This is a safety measure that allows users to prevent the
	// accidental removal of disks from a VM that should not be removed, such as
	// the VM's boot disk(s).
	//
	// Users may change this value at any time.
	//
	// When deploying a VM, disks from the VM image are automatically marked as
	// removable=false in order to prevent the accidental removal of the disks
	// needed to boot the VM.
	//
	// Defaults to true.
	Removable *bool `json:"removable,omitempty"`

	// +optional

	// ApplicationType describes the type of application for which this volume
	// is intended to be used.
	//
	//   - OracleRAC      -- The volume is configured with
	//                       diskMode=IndependentPersistent and
	//                       sharingMode=MultiWriter and attached to the first
	//                       SCSI controller with an available slot and
	//                       sharingMode=None. If no such controller exists,
	//                       a new ParaVirtual SCSI controller will be created
	//                       with sharingMode=None long as there are currently
	//                       three or fewer SCSI controllers.
	//   - MicrosoftWSFC  -- The volume is configured with
	//                       diskMode=IndependentPersistent and attached to a
	//                       SCSI controller with sharingMode=Physical.
	//                       If no such controller exists, a new ParaVirtual
	//                       SCSI controller will be created with
	//                       sharingMode=Physical as long as there are currently
	//                       three or fewer SCSI controllers.
	ApplicationType VolumeApplicationType `json:"applicationType,omitempty"`

	// +optional

	// ControllerBusNumber describes the bus number of the controller to which
	// this volume should be attached.
	//
	// The bus number specifies a controller based on the value of the
	// controllerType field:
	//
	//   - IDE  -- spec.hardware.ideControllers
	//   - NVME -- spec.hardware.nvmeControllers
	//   - SATA -- spec.hardware.sataControllers
	//   - SCSI -- spec.hardware.scsiControllers
	//
	// If this and controllerType are both omitted, the volume will be attached
	// to the first available SCSI controller. If there is no SCSI controller
	// with an available slot, a new ParaVirtual SCSI controller will be added
	// as long as there are currently three or fewer SCSI controllers.
	//
	// If the specified controller has no available slots, the request will be
	// denied.
	ControllerBusNumber *int32 `json:"controllerBusNumber,omitempty"`

	// +optional

	// ControllerType describes the type of the controller to which this volume
	// should be attached.
	//
	// Please keep in mind the number of volumes supported by the different
	// types of controllers:
	//
	//   - IDE                -- 4 total (2 per controller)
	//   - NVME               -- 256 total (64 per controller)
	//   - SATA               -- 120 total (30 per controller)
	//   - SCSI (ParaVirtual) -- 256 total (64 per controller)
	//   - SCSI (BusLogic)    -- 60 total (15 per controller)
	//   - SCSI (LsiLogic)    -- 60 total (15 per controller)
	//   - SCSI (LsiLogicSAS) -- 60 total (15 per controller)
	//
	// Please note, the number of supported volumes per SCSI controller may seem
	// off, but remember that a SCSI controller occupies a slot on its own bus.
	// Thus even though a ParaVirtual SCSI controller supports 65 targets and
	// the other types of SCSI controllers support 16 targets, one of the
	// targets is occupied by the controller itself.
	//
	// Defaults to SCSI when controllerBusNumber is also omitted; otherwise the
	// default value is determined by the logic outlined in the description of
	// the controllerBusNumber field.
	ControllerType VirtualControllerType `json:"controllerType,omitempty"`

	// +optional

	// DiskMode describes the desired mode to use when attaching the volume.
	//
	// Please note, volumes attached as IndependentNonPersistent or
	// IndependentPersistent are not included in a VM's snapshots or backups.
	//
	// Also, any data written to volumes attached as IndependentNonPersistent or
	// NonPersistent will be discarded when the VM is powered off.
	//
	// When applicationType=OracleRAC or applicationType=MicrosoftWSFC, this
	// field defaults to IndependentPersistent.
	// Otherwise, defaults to Persistent.
	DiskMode VolumeDiskMode `json:"diskMode,omitempty"`

	// +optional

	// SharingMode describes the volume's desired sharing mode.
	//
	// When applicationType=OracleRAC, this field defaults to MultiWriter.
	// Otherwise, defaults to None.
	SharingMode VolumeSharingMode `json:"sharingMode,omitempty"`

	// UnitNumber describes the desired unit number for attaching the volume to
	// a storage controller.
	//
	// When omitted, the next available unit number of the selected controller
	// is used.
	//
	// This value must be unique for the controller referenced by the
	// controllerBusNumber and controllerType properties. If the value is
	// already used by another device, this volume will not be attached.
	//
	// Please note the value 7 is invalid if controllerType=SCSI as 7 is the
	// unit number of the SCSI controller on its own bus.
	UnitNumber *int32 `json:"unitNumber,omitempty"`
}

// VirtualMachineVolumeSource represents the source location of a volume to
// mount. Only one of its members may be specified.
type VirtualMachineVolumeSource struct {
	// +optional

	// PersistentVolumeClaim represents a reference to a PersistentVolumeClaim
	// in the same namespace.
	//
	// More information is available at
	// https://kubernetes.io/docs/concepts/storage/persistent-volumes#persistentvolumeclaims.
	PersistentVolumeClaim *PersistentVolumeClaimVolumeSource `json:"persistentVolumeClaim,omitempty"`
}

// PersistentVolumeClaimVolumeSource is a composite for the Kubernetes
// corev1.PersistentVolumeClaimVolumeSource and instance storage options.
type PersistentVolumeClaimVolumeSource struct {
	corev1.PersistentVolumeClaimVolumeSource `json:",inline" yaml:",inline"`

	// +optional

	// InstanceVolumeClaim is set if the PVC is backed by instance storage.
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
