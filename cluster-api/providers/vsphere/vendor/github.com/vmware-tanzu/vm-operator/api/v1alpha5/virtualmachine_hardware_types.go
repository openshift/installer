// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package v1alpha5

import "k8s.io/apimachinery/pkg/api/resource"

// +kubebuilder:validation:Enum=IDE;NVME;SCSI;SATA

type VirtualControllerType string

const (
	VirtualControllerTypeIDE  VirtualControllerType = "IDE"
	VirtualControllerTypeNVME VirtualControllerType = "NVME"
	VirtualControllerTypeSCSI VirtualControllerType = "SCSI"
	VirtualControllerTypeSATA VirtualControllerType = "SATA"
)

// MaxCount returns the maximum number of controllers per VM.
func (t VirtualControllerType) MaxCount() int32 {
	switch t {
	case VirtualControllerTypeIDE:
		return 2
	case VirtualControllerTypeNVME,
		VirtualControllerTypeSATA,
		VirtualControllerTypeSCSI:
		return 4
	}
	return 0
}

type VirtualControllerSharingMode string

const (
	VirtualControllerSharingModeNone     VirtualControllerSharingMode = "None"
	VirtualControllerSharingModePhysical VirtualControllerSharingMode = "Physical"
	VirtualControllerSharingModeVirtual  VirtualControllerSharingMode = "Virtual"
)

// +kubebuilder:validation:Enum=ParaVirtual;BusLogic;LsiLogic;LsiLogicSAS

type SCSIControllerType string

const (
	SCSIControllerTypeParaVirtualSCSI SCSIControllerType = "ParaVirtual"
	SCSIControllerTypeBusLogic        SCSIControllerType = "BusLogic"
	SCSIControllerTypeLsiLogic        SCSIControllerType = "LsiLogic"
	SCSIControllerTypeLsiLogicSAS     SCSIControllerType = "LsiLogicSAS"
)

// +kubebuilder:validation:Enum=CDROM;Disk

type VirtualDeviceType string

const (
	VirtualDeviceTypeCDROM VirtualDeviceType = "CDROM"
	VirtualDeviceTypeDisk  VirtualDeviceType = "Disk"
)

type IDEControllerSpec struct {
	// +required
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1

	// BusNumber describes the desired bus number of the controller.
	BusNumber int32 `json:"busNumber"`
}

// MaxSlots returns the maximum number of slots per IDE controller.
func (c IDEControllerSpec) MaxSlots() int32 {
	return 2
}

// MaxCount returns the maximum number of IDE controllers per VM.
func (c IDEControllerSpec) MaxCount() int32 {
	return VirtualControllerTypeIDE.MaxCount()
}

// ReservedUnitNumber returns any reserved unit numbers or negative one.
func (c IDEControllerSpec) ReservedUnitNumber() int32 {
	return -1
}

type NVMEControllerSpec struct {
	// +required
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=3

	// BusNumber describes the desired bus number of the controller.
	BusNumber int32 `json:"busNumber"`

	// +optional
	// +kubebuilder:default=None
	// +kubebuilder:validation:Enum=None;Physical

	// SharingMode describes the sharing mode for the controller.
	//
	// Defaults to None.
	SharingMode VirtualControllerSharingMode `json:"sharingMode,omitempty"`
}

// MaxSlots returns the maximum number of slots per NVME controller.
func (c NVMEControllerSpec) MaxSlots() int32 {
	return 64
}

// MaxCount returns the maximum number of NVME controllers per VM.
func (c NVMEControllerSpec) MaxCount() int32 {
	return VirtualControllerTypeNVME.MaxCount()
}

// ReservedUnitNumber returns any reserved unit numbers or negative one.
func (c NVMEControllerSpec) ReservedUnitNumber() int32 {
	return -1
}

type SATAControllerSpec struct {
	// +required
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=3

	// BusNumber describes the desired bus number of the controller.
	BusNumber int32 `json:"busNumber"`
}

// MaxSlots returns the maximum number of slots per SATA controller.
func (c SATAControllerSpec) MaxSlots() int32 {
	return 30
}

// MaxCount returns the maximum number of SATA controllers per VM.
func (c SATAControllerSpec) MaxCount() int32 {
	return VirtualControllerTypeSATA.MaxCount()
}

// ReservedUnitNumber returns any reserved unit numbers or negative one.
func (c SATAControllerSpec) ReservedUnitNumber() int32 {
	return -1
}

type SCSIControllerSpec struct {
	// +required
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=3

	// BusNumber describes the desired bus number of the controller.
	BusNumber int32 `json:"busNumber"`

	// +optional
	// +kubebuilder:default=None
	// +kubebuilder:validation:Enum=None;Physical;Virtual

	// SharingMode describes the sharing mode for the controller.
	//
	// Defaults to None.
	SharingMode VirtualControllerSharingMode `json:"sharingMode,omitempty"`

	// +optional
	// +kubebuilder:default=ParaVirtual

	// Type describes the desired type of SCSI controller.
	//
	// Defaults to ParaVirtual.
	Type SCSIControllerType `json:"type,omitempty"`
}

// MaxSlots returns the maximum number of devices per SCSI controller type.
// The controller itself occupies one slot (unit number seven), which you should
// not use when assigning unit numbers to devices.
func (c SCSIControllerSpec) MaxSlots() int32 {
	switch c.Type {
	case SCSIControllerTypeParaVirtualSCSI:
		return 65 // There are 64 available slots and 1 reserved.
	case SCSIControllerTypeBusLogic,
		SCSIControllerTypeLsiLogic,
		SCSIControllerTypeLsiLogicSAS:
		return 16 // There are 15 available slots and 1 reserved.
	}
	return 0
}

// MaxCount returns the maximum number of SCSI controllers per VM.
func (c SCSIControllerSpec) MaxCount() int32 {
	return VirtualControllerTypeSCSI.MaxCount()
}

// ReservedUnitNumber returns any reserved unit numbers or negative one.
func (c SCSIControllerSpec) ReservedUnitNumber() int32 {
	return 7
}

type VirtualDeviceStatus struct {
	// +required

	// Type describes the type of the virtual device.
	Type VirtualDeviceType `json:"type"`

	// +required

	// UnitNumber describes the observed unit number of the device.
	UnitNumber int32 `json:"unitNumber"`
}

type VirtualControllerStatus struct {
	// +required

	// BusNumber describes the observed bus number of the controller.
	BusNumber int32 `json:"busNumber"`

	// +required

	// Type describes the observed type of the controller.
	Type VirtualControllerType `json:"type"`

	// +required

	// DeviceKey describes the observed device key of the controller.
	DeviceKey int32 `json:"deviceKey"`

	// +optional
	// +listType=map
	// +listMapKey=unitNumber

	// Devices describes the observed devices connected to the controller.
	Devices []VirtualDeviceStatus `json:"devices,omitempty"`
}

// VirtualMachineCdromSpec describes the desired state of a CD-ROM device.
type VirtualMachineCdromSpec struct {
	// +required
	// +kubebuilder:validation:Pattern="^[a-z0-9]{2,}$"

	// Name consists of at least two lowercase letters or digits of this CD-ROM.
	// It must be unique among all CD-ROM devices attached to the VM.
	//
	// This field is immutable when the VM is powered on.
	Name string `json:"name"`

	// Image describes the reference to an ISO type VirtualMachineImage or
	// ClusterVirtualMachineImage resource used as the backing for the CD-ROM.
	// If the image kind is omitted, it defaults to VirtualMachineImage.
	//
	// This field is immutable when the VM is powered on.
	//
	// Please note, unlike the spec.imageName field, the value of this
	// spec.cdrom.image.name MUST be a Kubernetes object name.
	Image VirtualMachineImageRef `json:"image"`

	// +optional

	// ControllerBusNumber describes the bus number of the controller to which
	// this CD-ROM should be attached.
	//
	// The bus number specifies a controller based on the value of the
	// controllerType field:
	//
	//   - IDE  -- spec.hardware.ideControllers
	//   - SATA -- spec.hardware.sataControllers
	//
	// If this and controllerType are both omitted, the CD-ROM will be attached
	// to the  first available IDE controller. If there is no IDE controller
	// with an available slot, a new SATA controller will be added as long as
	// there are currently three or fewer SATA controllers.
	//
	// If the specified controller has no available slots, the request will be
	// denied.
	ControllerBusNumber *int32 `json:"controllerBusNumber,omitempty"`

	// +optional
	// +kubebuilder:validation:Enum=IDE;SATA

	// ControllerType describes the type of the controller to which this CD-ROM
	// should be attached.
	//
	// Please keep in mind the number of devices supported by the different
	// types of controllers:
	//
	//   - IDE                -- 4 total (2 per controller)
	//   - SATA               -- 120 total (30 per controller)
	//
	// Defaults to IDE when controllerBusNumber is also omitted; otherwise the
	// default value is determined by the logic outlined in the description of
	// the controllerBusNumber field.
	ControllerType VirtualControllerType `json:"controllerType,omitempty"`

	// UnitNumber describes the desired unit number for attaching the CD-ROM to
	// a storage controller.
	//
	// When omitted, the next available unit number of the selected controller
	// is used.
	//
	// This value must be unique for the controller referenced by the
	// controllerBusNumber and controllerType properties. If the value is
	// already used by another device, this CD-ROM will not be attached.
	UnitNumber *int32 `json:"unitNumber,omitempty"`

	// +optional
	// +kubebuilder:default=true

	// Connected describes the desired connection state of the CD-ROM device.
	//
	// When true, the CD-ROM device is added and connected to the VM.
	// If the device already exists, it is updated to a connected state.
	//
	// When explicitly set to false, the CD-ROM device is added but remains
	// disconnected from the VM. If the CD-ROM device already exists, it is
	// updated to a disconnected state.
	//
	// Note: Before disconnecting a CD-ROM, the device may need to be unmounted
	// in the guest OS. Refer to the following KB article for more details:
	// https://knowledge.broadcom.com/external/article?legacyId=2144053
	//
	// Defaults to true if omitted.
	Connected *bool `json:"connected,omitempty"`

	// +optional
	// +kubebuilder:default=true

	// AllowGuestControl describes whether or not a web console connection
	// may be used to connect/disconnect the CD-ROM device.
	//
	// Defaults to true if omitted.
	AllowGuestControl *bool `json:"allowGuestControl,omitempty"`
}

type VirtualMachineHardwareSpec struct {
	// +optional
	// +listType=map
	// +listMapKey=name

	// Cdrom describes the desired state of the VM's CD-ROM devices.
	//
	// Each CD-ROM device requires a reference to an ISO-type
	// VirtualMachineImage or ClusterVirtualMachineImage resource as backing.
	//
	// Multiple CD-ROM devices using the same backing image, regardless of image
	// kinds (namespace or cluster scope), are not allowed.
	//
	// CD-ROM devices can be added, updated, or removed when the VM is powered
	// off. When the VM is powered on, only the connection state of existing
	// CD-ROM devices can be changed.
	// CD-ROM devices are attached to the VM in the specified list-order.
	Cdrom []VirtualMachineCdromSpec `json:"cdrom,omitempty"`

	// +optional
	// +listType=map
	// +listMapKey=busNumber
	// +kubebuilder:validation:MaxItems=2

	// IDEControllers describes the desired list of IDE controllers for the VM.
	//
	// Defaults to two IDE controllers, with bus 0 and bus 1.
	IDEControllers []IDEControllerSpec `json:"ideControllers,omitempty"`

	// +optional
	// +listType=map
	// +listMapKey=busNumber
	// +kubebuilder:validation:MaxItems=4

	// NVMEControllers describes the desired list of NVME controllers for the
	// VM.
	NVMEControllers []NVMEControllerSpec `json:"nvmeControllers,omitempty"`

	// +optional
	// +listType=map
	// +listMapKey=busNumber
	// +kubebuilder:validation:MaxItems=4

	// SATAControllers describes the desired list of SATA controllers for the
	// VM.
	//
	// Please note, all SATA controllers are VirtualAHCI.
	SATAControllers []SATAControllerSpec `json:"sataControllers,omitempty"`

	// +optional
	// +listType=map
	// +listMapKey=busNumber
	// +kubebuilder:validation:MaxItems=4

	// SCSIControllers describes the desired list of SCSI controllers for the
	// VM.
	SCSIControllers []SCSIControllerSpec `json:"scsiControllers,omitempty"`
}

type VirtualMachineCPUAllocationStatus struct {
	// +optional

	// Total describes the observed number of processors.
	Total int32 `json:"total,omitempty"`

	// +optional

	// Reservation describes the observed CPU reservation in MHz.
	Reservation int64 `json:"reservation,omitempty"`
}

type VirtualMachineMemoryAllocationStatus struct {
	// +optional

	// Total describes the observed amount of configured memory.
	Total *resource.Quantity `json:"total,omitempty"`

	// +optional

	// Reservation describes the observed memory reservation.
	Reservation *resource.Quantity `json:"reservation,omitempty"`
}

type VirtualMachineVGPUType string

const (
	VirtualMachineVGPUTypeNVIDIA VirtualMachineVGPUType = "Nvidia"
)

type VirtualMachineVGPUMigrationType string

const (
	VirtualMachineVGPUMigrationTypeNone     VirtualMachineVGPUMigrationType = "None"
	VirtualMachineVGPUMigrationTypeNormal   VirtualMachineVGPUMigrationType = "Normal"
	VirtualMachineVGPUMigrationTypeEnhanced VirtualMachineVGPUMigrationType = "Enhanced"
)

type VirtualMachineHardwareVGPUStatus struct {
	// +optional

	// Type describes the observed type of the vGPU.
	Type VirtualMachineVGPUType `json:"type,omitempty"`

	// +optional

	// Profile describes the observed profile used by the vGPU.
	//
	// Please note, this is only applicable to Nvidia vGPUs.
	Profile string `json:"profile,omitempty"`

	// +optional

	// MigrationType describes the vGPU's observed vMotion support.
	MigrationType VirtualMachineVGPUMigrationType `json:"migrationType,omitempty"`
}

type VirtualMachineHardwareStatus struct {
	// +optional
	// +listType=map
	// +listMapKey=busNumber
	// +listMapKey=type

	// Controllers describes the observed list of virtual controllers for the
	// VM.
	Controllers []VirtualControllerStatus `json:"controllers,omitempty"`

	// +optional

	// CPU describes the observed CPU allocation of the VM.
	CPU *VirtualMachineCPUAllocationStatus `json:"cpu,omitempty"`

	// +optional

	// Memory describes the observed memory allocation of the VM.
	Memory *VirtualMachineMemoryAllocationStatus `json:"memory,omitempty"`

	// +optional

	// VGPUs describes the observed vGPUs used by this VM.
	VGPUs []VirtualMachineHardwareVGPUStatus `json:"vGPUs,omitempty"`
}
