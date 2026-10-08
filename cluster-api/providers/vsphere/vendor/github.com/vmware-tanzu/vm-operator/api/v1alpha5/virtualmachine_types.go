// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

//go:generate go -C ../.. run ./pkg/gen/guestosid/guestosid.go v1alpha5 ./api/v1alpha5/zz_virtualmachine_guestosid_generated.go

package v1alpha5

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	vmopv1a5common "github.com/vmware-tanzu/vm-operator/api/v1alpha5/common"
)

const (
	// VirtualMachineConditionClassReady indicates that a referenced
	// VirtualMachineClass is ready.
	//
	// For more information please see VirtualMachineClass.Status.Ready.
	VirtualMachineConditionClassReady = "VirtualMachineClassReady"

	// VirtualMachineConditionImageReady indicates that a referenced
	// VirtualMachineImage is ready.
	//
	// For more information please see VirtualMachineImage.Status.Ready.
	VirtualMachineConditionImageReady = "VirtualMachineImageReady"

	// VirtualMachineConditionVMSetResourcePolicyReady indicates that a referenced
	// VirtualMachineSetResourcePolicy is Ready.
	VirtualMachineConditionVMSetResourcePolicyReady = "VirtualMachineConditionVMSetResourcePolicyReady"

	// VirtualMachineConditionStorageReady indicates that the storage prerequisites for the VM are ready.
	VirtualMachineConditionStorageReady = "VirtualMachineStorageReady"

	// VirtualMachineConditionBootstrapReady indicates that the bootstrap prerequisites for the VM are ready.
	VirtualMachineConditionBootstrapReady = "VirtualMachineBootstrapReady"

	// VirtualMachineConditionNetworkReady indicates that the network prerequisites for the VM are ready.
	VirtualMachineConditionNetworkReady = "VirtualMachineNetworkReady"

	// VirtualMachineConditionImageCacheReady indicates that the VirtualMachineImageCache for the VM is ready.
	VirtualMachineConditionImageCacheReady = "VirtualMachineConditionImageCacheReady"

	// VirtualMachineConditionPlacementReady indicates that the placement decision for the VM is ready.
	VirtualMachineConditionPlacementReady = "VirtualMachineConditionPlacementReady"

	// VirtualMachineEncryptionSynced indicates that the VirtualMachine's
	// encryption state is synced to the desired encryption state.
	VirtualMachineEncryptionSynced = "VirtualMachineEncryptionSynced"

	// VirtualMachineDiskPromotionStarted indicates that the VirtualMachine's
	// disk promotion process has started.
	VirtualMachineDiskPromotionStarted = "VirtualMachineDiskPromotionStarted"

	// VirtualMachineDiskPromotionSynced indicates that the VirtualMachine's
	// disk promotion state is synced to the desired promotion state.
	VirtualMachineDiskPromotionSynced = "VirtualMachineDiskPromotionSynced"

	// VirtualMachineConditionCreated indicates that the VM has been created.
	VirtualMachineConditionCreated = "VirtualMachineCreated"

	// VirtualMachineClassConfigurationSynced indicates that the VM's current configuration is synced to the
	// current version of its VirtualMachineClass.
	VirtualMachineClassConfigurationSynced = "VirtualMachineClassConfigurationSynced"

	// VirtualMachineHardwareDeviceConfigVerified indicates that the VM's hardware
	// device configuration (controllers, volumes, CD-ROM devices) matches the
	// desired state specified in the spec.
	VirtualMachineHardwareDeviceConfigVerified = "VirtualMachineHardwareDeviceConfigVerified"

	// VirtualMachineHardwareControllersVerified indicates that the VM's hardware
	// controllers match the desired state specified in the spec.
	VirtualMachineHardwareControllersVerified = "VirtualMachineHardwareControllersVerified"

	// VirtualMachineHardwareVolumesVerified indicates that the VM's hardware
	// volumes match the desired state specified in the spec.
	VirtualMachineHardwareVolumesVerified = "VirtualMachineHardwareVolumesVerified"

	// VirtualMachineHardwareCDROMVerified indicates that the VM's hardware
	// CD-ROM devices match the desired state specified in the spec.
	VirtualMachineHardwareCDROMVerified = "VirtualMachineHardwareCDROMVerified"

	// VirtualMachinePowerStateSynced indicates the VM's power state is synced
	// to the desired state.
	VirtualMachinePowerStateSynced = "VirtualMachinePowerStateSynced"

	// VirtualMachineGuestNetworkConfigSynced indicates the VM's guest network
	// configuration is synced to the desired state, determined by whether or
	// not the guest is reporting the expected IP address(es).
	VirtualMachineGuestNetworkConfigSynced = "VirtualMachineGuestNetworkConfigSynced"
)

const (
	// VirtualMachineSnapshotRevertSucceeded indicates that the VM
	// has been reverted to a snapshot.
	VirtualMachineSnapshotRevertSucceeded = "VirtualMachineSnapshotRevertSucceeded"

	// VirtualMachineSnapshotRevertInProgressReason indicates that the
	// revert operation is in progress.
	VirtualMachineSnapshotRevertInProgressReason = "VirtualMachineSnapshotRevertInProgress"

	// VirtualMachineSnapshotRevertTaskFailedReason indicates that the
	// revert operation is invalid.
	VirtualMachineSnapshotRevertTaskFailedReason = "VirtualMachineSnapshotRevertTaskFailed"

	// VirtualMachineSnapshotRevertFailedInvalidVMManifestReason indicates
	// that the revert operation has failed due to invalid VM spec to revert to.
	VirtualMachineSnapshotRevertFailedInvalidVMManifestReason = "VirtualMachineSnapshotRevertFailedInvalidVMManifest"

	// VirtualMachineSnapshotRevertSkippedReason indicates that the
	// revert operation was skipped.
	VirtualMachineSnapshotRevertSkippedReason = "VirtualMachineSnapshotRevertSkipped"

	// VirtualMachineSnapshotRevertFailedReason indicates that the
	// revert operation failed for some reason.
	VirtualMachineSnapshotRevertFailedReason = "VirtualMachineSnapshotRevertFailed"

	// VirtualMachineHardwareControllersMismatchReason indicates that the VM's
	// controller configuration does not match the desired state specified in the spec.
	VirtualMachineHardwareControllersMismatchReason = "HardwareControllersMismatch"

	// VirtualMachineHardwareVolumesMismatchReason indicates that the VM's
	// volume configuration does not match the desired state specified in the spec.
	VirtualMachineHardwareVolumesMismatchReason = "HardwareVolumesMismatch"

	// VirtualMachineHardwareCDROMMismatchReason indicates that the VM's
	// CD-ROM device configuration does not match the desired state specified in the spec.
	VirtualMachineHardwareCDROMMismatchReason = "HardwareCDROMMismatch"

	// VirtualMachineHardwareDeviceConfigMismatchReason indicates that the VM's
	// hardware device configuration does not match the desired state specified
	// in the spec. This is used for the aggregated condition.
	VirtualMachineHardwareDeviceConfigMismatchReason = "HardwareDeviceConfigMismatch"
)

const (
	// GuestBootstrapCondition exposes the status of guest bootstrap from within
	// the guest OS, when available.
	GuestBootstrapCondition = "GuestBootstrap"

	// GuestIDReconfiguredCondition exposes the status of guest ID
	// reconfiguration after a VM has been created, when available.
	GuestIDReconfiguredCondition = "GuestIDReconfigured"
)

const (
	// GuestCustomizationCondition exposes the status of guest customization
	// from within the guest OS, when available.
	GuestCustomizationCondition = "GuestCustomization"

	// GuestCustomizationIdleReason documents that guest
	// customizations were not applied for the VirtualMachine.
	GuestCustomizationIdleReason = "GuestCustomizationIdle"

	// GuestCustomizationPendingReason documents that guest
	// customization is still pending within the guest OS.
	GuestCustomizationPendingReason = "GuestCustomizationPending"

	// GuestCustomizationRunningReason documents that the guest
	// customization is now running on the guest OS.
	GuestCustomizationRunningReason = "GuestCustomizationRunning"

	// GuestCustomizationSucceededReason documents that the
	// guest customization succeeded within the guest OS.
	GuestCustomizationSucceededReason = "GuestCustomizationSucceeded"

	// GuestCustomizationFailedReason documents that the guest
	// customization failed within the guest OS.
	GuestCustomizationFailedReason = "GuestCustomizationFailed"
)

const (
	// VirtualMachineToolsCondition exposes the status of VMware Tools running
	// in the guest OS, when available.
	VirtualMachineToolsCondition = "VirtualMachineTools"

	// VirtualMachineToolsNotRunningReason documents that
	// VMware Tools is not running.
	VirtualMachineToolsNotRunningReason = "VirtualMachineToolsNotRunning"

	// VirtualMachineToolsRunningReason documents that VMware
	// Tools is running.
	VirtualMachineToolsRunningReason = "VirtualMachineToolsRunning"
)

const (
	// VirtualMachineReconcileReady exposes the status of VirtualMachine reconciliation.
	VirtualMachineReconcileReady = "VirtualMachineReconcileReady"

	// VirtualMachineReconcileRunningReason indicates that VirtualMachine
	// reconciliation is running.
	VirtualMachineReconcileRunningReason = "VirtualMachineReconcileRunning"

	// VirtualMachineReconcilePausedReason indicates that VirtualMachine
	// reconciliation is being paused.
	VirtualMachineReconcilePausedReason = "VirtualMachineReconcilePaused"
)

const (
	// PauseAnnotation is an annotation that prevents a VM from being
	// reconciled.
	//
	// This can be used when a VM needs to be modified directly on the
	// underlying infrastructure without VM Service attempting to direct the
	// VM back to its intended state.
	//
	// The VM will not be reconciled again until this annotation is removed.
	PauseAnnotation = GroupName + "/paused"

	// InstanceIDAnnotation is an annotation that can be applied to set Cloud-Init metadata Instance ID.
	//
	// This cannot be set by users. It is for VM Operator to handle corner cases.
	//
	// In a corner case where a VM first boot failed to bootstrap with Cloud-Init, VM Operator sets Instance ID
	// the same with the first boot Instance ID to prevent Cloud-Init from treating this VM as first boot
	// due to different Instance ID. This annotation is used in upgrade script.
	InstanceIDAnnotation = GroupName + "/cloud-init-instance-id"

	// FirstBootDoneAnnotation is an annotation that indicates the VM has been
	// booted at least once. This annotation cannot be set by users and will not
	// be removed once set until the VM is deleted.
	FirstBootDoneAnnotation = "virtualmachine." + GroupName + "/first-boot-done"

	// V1alpha1ConfigMapTransportAnnotation is an annotation that indicates that the VM
	// was created with the v1alpha1 API and specifies a configMap as the metadata transport resource type.
	V1alpha1ConfigMapTransportAnnotation = GroupName + "/v1a1-configmap-md-transport"

	// VirtualMachineSameVMClassResizeAnnotation is an annotation that indicates the VM
	// should be resized as the class it points to changes.
	VirtualMachineSameVMClassResizeAnnotation = GroupName + "/same-vm-class-resize"
)

const (
	// checkAnnotationSubDomain is the sub-domain to be used for all check-style
	// annotations that enable external components to participate in a VM's
	// lifecycle events.
	checkAnnotationSubDomain = "check.vmoperator.vmware.com"

	// CheckAnnotationPowerOn is an annotation that may be used to prevent a
	// VM from being powered on. A user can still set a VM's spec.powerState to
	// PoweredOn, but the VM will not be powered on until the check annotation
	// is removed.
	//
	// Please note, there may be multiple check annotations, ex.:
	//
	// - poweron.check.vmoperator.vmware.com/component1: "reason"
	// - poweron.check.vmoperator.vmware.com/component2: "reason"
	// - poweron.check.vmoperator.vmware.com/component3: "reason"
	//
	// All check annotations must be removed before a VM can be powered on.
	//
	// This annotation may only be applied when creating a new VM by any user.
	//
	// Only privileged users may apply this annotation to existing VMs.
	//
	// Only privileged users may remove this annotation from a VM. If a
	// non-privileged user accidentally adds this annotation when creating a VM,
	// the recourse is to delete the VM and recreate it without the annotation.
	CheckAnnotationPowerOn = "poweron." + checkAnnotationSubDomain

	// CheckAnnotationDelete is an annotation that may be used to prevent a
	// VM from being deleted. A user can still delete the VM, but the VM will
	// not *actually* be removed until the annotation is removed.
	//
	// Unlike a finalizer, this annotation *also* prevents the underlying
	// vSphere VM from being deleted as well.
	//
	// Please note, there may be multiple check annotations, ex.:
	//
	// - delete.check.vmoperator.vmware.com/component1: "reason"
	// - delete.check.vmoperator.vmware.com/component2: "reason"
	// - delete.check.vmoperator.vmware.com/component3: "reason"
	//
	// All check annotations must be removed before a VM can be deleted.
	//
	// Only privileged users may add or remove this annotation.
	CheckAnnotationDelete = "delete." + checkAnnotationSubDomain
)

const (
	// ManagedByExtensionKey and ManagedByExtensionType represent the ManagedBy
	// field on the VM. They are used to differentiate VM Service managed VMs
	// from traditional vSphere VMs.
	ManagedByExtensionKey  = "com.vmware.vcenter.wcp"
	ManagedByExtensionType = "VirtualMachine"
)

const (
	// VirtualMachineBackupUpToDateCondition exposes the status of the latest VirtualMachine Backup, when available.
	VirtualMachineBackupUpToDateCondition = "VirtualMachineBackupUpToDate"

	// VirtualMachineBackupPausedReason documents that VirtualMachine backup is paused.
	// This can happen after the backing virtual machine is restored by a backup/restore vendor, or a failover operation
	// by a data protection solution. In either of these operations, VM operator does not persist backup information and
	// waits for the virtual machine to be (re)-registered with VM Service.
	VirtualMachineBackupPausedReason = "VirtualMachineBackupPaused"

	// VirtualMachineBackupFailedReason documents that the VirtualMachine backup failed due to an error.
	VirtualMachineBackupFailedReason = "VirtualMachineBackupFailed"
)

const (
	// ForceEnableBackupAnnotation is an annotation that instructs VM operator to
	// ignore all exclusion rules and persist the configuration of the resource in
	// virtual machine in relevant ExtraConfig fields.
	//
	// This is an experimental flag which only guarantees that the configuration
	// of the VirtualMachine resource will be persisted in the virtual machine on
	// vSphere.  There is no guarantee that the registration of the VM will be
	// successful post a restore or failover operation.
	ForceEnableBackupAnnotation = GroupName + "/force-enable-backup"

	// VirtualMachineBackupVersionAnnotation is an annotation that indicates the VM's
	// last backup version. It is a monotonically increasing counter and
	// is only supposed to be used by IaaS control plane and vCenter for virtual machine registration
	// post a restore operation.
	//
	// The VirtualMachineBackupVersionAnnotation on the VM resource in Supervisor and the BackupVersionExtraConfigKey on the vSphere VM
	// indicate whether the backups are in sync.
	VirtualMachineBackupVersionAnnotation = GroupName + "/backup-version"
)

const (
	// ManagerID on a VirtualMachine contains the UUID of the
	// VMware vCenter (VC) that is managing this virtual machine.
	ManagerID = GroupName + "/manager-id"

	// RestoredVMAnnotation on a VirtualMachine represents that a virtual
	// machine has been restored using the RegisterVM API, typically by a
	// VADP based data protection vendor. The presence of this annotation is
	// used to bypass some validation checks that are otherwise
	// applicable to all VirtualMachine create/update requests.
	RestoredVMAnnotation = GroupName + "/restored-vm"

	// ImportedVMAnnotation on a VirtualMachine represents that a traditional virtual
	// machine has been imported into Supervisor using the ImportVM API. The presence of this
	// annotation is used to bypass some validation checks that are otherwise applicable
	// to all VirtualMachine create/update requests.
	ImportedVMAnnotation = GroupName + "/imported-vm"

	// FailedOverVMAnnotation on a VirtualMachine resource represents that a virtual
	// machine has been failed over from one site to the other, typically as part of a
	// disaster recovery workflow.  The presence of this annotation is used to bypass
	// some validation checks that are otherwise applicable to all VirtualMachine
	// create/update requests.
	FailedOverVMAnnotation = GroupName + "/failed-over-vm"
)

const (
	// PauseVMExtraConfigKey is the ExtraConfig key to allow override
	// operations for admins to pause reconciliation of VM Service VM.
	//
	// Please note, the value that takes effect is the string "True"(case-insensitive).
	PauseVMExtraConfigKey = "vmservice.virtualmachine.pause"

	// PausedVMLabelKey is the label key to identify VMs that reconciliation
	// are paused. Value will specify whose operation is responsible for
	// the pause. It can be admins or devops or both.
	//
	// Only privileged user can edit this label.
	PausedVMLabelKey = GroupName + "/paused"
)

// +kubebuilder:validation:Enum=PoweredOff;PoweredOn;Suspended

// VirtualMachinePowerState defines a VM's desired and observed power states.
type VirtualMachinePowerState string

const (
	// VirtualMachinePowerStateOff indicates to shut down a VM and/or it is
	// shut down.
	VirtualMachinePowerStateOff VirtualMachinePowerState = "PoweredOff"

	// VirtualMachinePowerStateOn indicates to power on a VM and/or it is
	// powered on.
	VirtualMachinePowerStateOn VirtualMachinePowerState = "PoweredOn"

	// VirtualMachinePowerStateSuspended indicates to suspend a VM and/or it is
	// suspended.
	VirtualMachinePowerStateSuspended VirtualMachinePowerState = "Suspended"
)

// +kubebuilder:validation:Enum=Hard;Soft;TrySoft

// VirtualMachinePowerOpMode represents the various power operation modes when
// powering off or suspending a VM.
type VirtualMachinePowerOpMode string

const (
	// VirtualMachinePowerOpModeHard indicates to halt a VM when powering it
	// off or when suspending a VM to not involve the guest.
	VirtualMachinePowerOpModeHard VirtualMachinePowerOpMode = "Hard"

	// VirtualMachinePowerOpModeSoft indicates to ask VM Tools running
	// inside of a VM's guest to shutdown the guest gracefully when powering
	// off a VM or when suspending a VM to allow the guest to participate.
	//
	// If this mode is set on a VM whose guest does not have VM Tools or if
	// VM Tools is present but the operation fails, the VM may never realize
	// the desired power state. This can prevent a VM from being deleted as well
	// as many other unexpected issues. It is recommended to use trySoft
	// instead.
	VirtualMachinePowerOpModeSoft VirtualMachinePowerOpMode = "Soft"

	// VirtualMachinePowerOpModeTrySoft indicates to first attempt a Soft
	// operation and fall back to Hard if VM Tools is not present in the guest,
	// if the Soft operation fails, or if the VM is not in the desired power
	// state within five minutes.
	VirtualMachinePowerOpModeTrySoft VirtualMachinePowerOpMode = "TrySoft"
)

type VirtualMachineImageRef struct {
	// Kind describes the type of image, either a namespace-scoped
	// VirtualMachineImage or cluster-scoped ClusterVirtualMachineImage.
	Kind string `json:"kind"`

	// Name refers to the name of a VirtualMachineImage resource in the same
	// namespace as this VM or a cluster-scoped ClusterVirtualMachineImage.
	Name string `json:"name"`
}

// VirtualMachineCryptoSpec defines the desired state of a VirtualMachine's
// encryption state.
type VirtualMachineCryptoSpec struct {
	// +optional

	// EncryptionClassName describes the name of the EncryptionClass resource
	// used to encrypt this VM.
	//
	// Please note, this field is not required to encrypt the VM. If the
	// underlying platform has a default key provider, the VM may still be fully
	// or partially encrypted depending on the specified storage and VM classes.
	//
	// If there is a default key provider and an encryption storage class is
	// selected, the files in the VM's home directory and non-PVC virtual disks
	// will be encrypted
	//
	// If there is a default key provider and a VM Class with a virtual, trusted
	// platform module (vTPM) is selected, the files in the VM's home directory,
	// minus any virtual disks, will be encrypted.
	//
	// If the underlying vSphere platform does not have a default key provider,
	// then this field is required when specifying an encryption storage class
	// and/or a VM Class with a vTPM.
	//
	// If this field is set, spec.storageClass must use an encryption-enabled
	// storage class.
	EncryptionClassName string `json:"encryptionClassName,omitempty"`

	// +optional
	// +kubebuilder:default=true

	// UseDefaultKeyProvider describes the desired behavior for when an explicit
	// EncryptionClass is not provided.
	//
	// When an explicit EncryptionClass is not provided and this value is true:
	//
	// - Deploying a VirtualMachine with an encryption storage policy or vTPM
	//   will be encrypted using the default key provider.
	//
	// - If a VirtualMachine is not encrypted, uses an encryption storage
	//   policy or has a virtual, trusted platform module (vTPM), there is a
	//   default key provider, the VM will be encrypted using the default key
	//   provider.
	//
	// - If a VirtualMachine is encrypted with a provider other than the default
	//   key provider, the VM will be rekeyed using the default key provider.
	//
	// When an explicit EncryptionClass is not provided and this value is false:
	//
	// - Deploying a VirtualMachine with an encryption storage policy or vTPM
	//   will fail.
	//
	// - If a VirtualMachine is encrypted with a provider other than the default
	//   key provider, the VM will be not be rekeyed.
	//
	//   Please note, this could result in a VirtualMachine that cannot be
	//   powered on since it is encrypted using a provider or key that may have
	//   been removed. Without the key, the VM cannot be decrypted and thus
	//   cannot be powered on.
	//
	// Defaults to true if omitted.
	UseDefaultKeyProvider *bool `json:"useDefaultKeyProvider,omitempty"`

	// +kubebuilder:default=New

	// VTPMMode describes the desired behavior when deploying a VirtualMachine
	// using a VirtualMachine-backed image which created from an encrypted
	// VirtualMachine with a vTPM.
	//
	// The possible values for this field are:
	//
	// - Clone - The vTPM will be preserved from the VirtualMachineImage.
	// - New - The vTPM will not be preserved.
	//
	// The default value of this field is New.
	VTPMMode VirtualMachineCryptoVTPMMode `json:"vTPMMode,omitempty"`
}

// +kubebuilder:validation:Enum=Clone;New

// VirtualMachineCryptoVTPMMode represents whether to preserve the vTPM
// from an encrypted VirtualMachine-backed image when deploying a VirtualMachine.
type VirtualMachineCryptoVTPMMode string

const (
	VirtualMachineCryptoVTPMModeClone VirtualMachineCryptoVTPMMode = "Clone"
	VirtualMachineCryptoVTPMModeNew   VirtualMachineCryptoVTPMMode = "New"
)

// VirtualMachineBootOptionsBootableDevice represents a bootable device
// that a VM may be booted from.
type VirtualMachineBootOptionsBootableDevice struct {
	// +required

	// Type represents the type of bootable device.
	//
	// The available device types are:
	//
	// - Disk
	// - Network
	// - CDRom
	Type VirtualMachineBootOptionsBootableDeviceType `json:"type"`

	// +optional

	// Name represents the name of the bootable device. It is
	// required for Disk and Network device types, while ignored
	// for CDRom device types.
	Name string `json:"name,omitempty"`
}

// +kubebuilder:validation:Enum=Disk;Network;CDRom

// VirtualMachineBootOptionsBootableDeviceType represents the type of bootable device.
type VirtualMachineBootOptionsBootableDeviceType string

const (
	VirtualMachineBootOptionsBootableDiskDevice    VirtualMachineBootOptionsBootableDeviceType = "Disk"
	VirtualMachineBootOptionsBootableNetworkDevice VirtualMachineBootOptionsBootableDeviceType = "Network"
	VirtualMachineBootOptionsBootableCDRomDevice   VirtualMachineBootOptionsBootableDeviceType = "CDRom"
)

// +kubebuilder:validation:Enum=IP4;IP6

// VirtualMachineBootOptionsNetworkBootProtocol represents the protocol to
// use during PXE network boot or NetBoot.
type VirtualMachineBootOptionsNetworkBootProtocol string

const (
	VirtualMachineBootOptionsNetworkBootProtocolIP4 VirtualMachineBootOptionsNetworkBootProtocol = "IP4"
	VirtualMachineBootOptionsNetworkBootProtocolIP6 VirtualMachineBootOptionsNetworkBootProtocol = "IP6"
)

// +kubebuilder:validation:Enum=bios;efi

// VirtualMachineBootOptionsFirmwareType represents the firmware to use.
type VirtualMachineBootOptionsFirmwareType string

const (
	VirtualMachineBootOptionsFirmwareTypeBIOS VirtualMachineBootOptionsFirmwareType = "bios"
	VirtualMachineBootOptionsFirmwareTypeEFI  VirtualMachineBootOptionsFirmwareType = "efi"
)

// +kubebuilder:validation:Enum=Enabled;Disabled

// VirtualMachineBootOptionsForceBootEntry represents whether to force the virtual machine
// to enter bios/efi setup the next time the virtual machine boots.
type VirtualMachineBootOptionsForceBootEntry string

const (
	VirtualMachineBootOptionsForceBootEntryEnabled  VirtualMachineBootOptionsForceBootEntry = "Enabled"
	VirtualMachineBootOptionsForceBootEntryDisabled VirtualMachineBootOptionsForceBootEntry = "Disabled"
)

// +kubebuilder:validation:Enum=Enabled;Disabled

// VirtualMachineBootOptionsEFISecureBoot represents whether the virtual machine will
// perform EFI Secure Boot.
type VirtualMachineBootOptionsEFISecureBoot string

const (
	VirtualMachineBootOptionsEFISecureBootEnabled  VirtualMachineBootOptionsEFISecureBoot = "Enabled"
	VirtualMachineBootOptionsEFISecureBootDisabled VirtualMachineBootOptionsEFISecureBoot = "Disabled"
)

// VirtualMachineBootOptionsBootRetry represents whether a virtual machine that fails to boot
// will automatically try again.
type VirtualMachineBootOptionsBootRetry string

const (
	VirtualMachineBootOptionsBootRetryEnabled  VirtualMachineBootOptionsBootRetry = "Enabled"
	VirtualMachineBootOptionsBootRetryDisabled VirtualMachineBootOptionsBootRetry = "Disabled"
)

// VirtualMachineBootOptions defines the boot-time behavior of a virtual machine.
type VirtualMachineBootOptions struct {
	// +optional

	// Firmware represents the firmware for the virtual machine to use. Any update
	// to this value after the virtual machine has already been created will be
	// ignored. Setting will happen in the following manner:
	//
	// 1. If this value is specified, it will be used to set the VM firmware. If this
	//    value is unset, then
	// 2. the virtual machine image will be checked. If that value is set, then it will
	//    be used to set the VM firmware. If that value is unset, then
	// 3. the virtual machine class will be checked. If that value is set, then it will
	//    be used to set the VM firmware. If that value is unset, then
	// 4. the VM firmware will be set, by default, to BIOS.
	//
	// The available values of this field are:
	//
	// - bios
	// - efi
	Firmware VirtualMachineBootOptionsFirmwareType `json:"firmware,omitempty"`

	// +optional

	// BootDelay is the delay before starting the boot sequence. The boot delay
	// specifies a time interval between virtual machine power on or restart and
	// the beginning of the boot sequence.
	BootDelay *metav1.Duration `json:"bootDelay,omitempty"`

	// +optional

	// BootOrder represents the boot order of the virtual machine. After list is exhausted,
	// default BIOS boot device algorithm is used for booting. Note that order of the entries
	// in the list is important: device listed first is used for boot first, if that one
	// fails second entry is used, and so on. Platform may have some internal limit on the
	// number of devices it supports. If bootable device is not reached before platform's limit
	// is hit, boot will fail. At least single entry is supported by all products supporting
	// boot order settings.
	BootOrder []VirtualMachineBootOptionsBootableDevice `json:"bootOrder,omitempty"`

	// +optional
	// +kubebuilder:default=Disabled

	// BootRetry specifies whether a virtual machine that fails to boot
	// will try again. The available values are:
	//
	// - Enabled -- A virtual machine that fails to boot will try again
	//              after BootRetryDelay time period has expired.
	// - Disabled -- The virtual machine waits indefinitely for you to
	//               initiate boot retry.
	BootRetry VirtualMachineBootOptionsBootRetry `json:"bootRetry,omitempty"`

	// +optional

	// BootRetryDelay specifies a time interval between virtual machine boot failure
	// and the subsequent attempt to boot again. The virtual machine uses this value
	// only if BootRetry is Enabled.
	BootRetryDelay *metav1.Duration `json:"bootRetryDelay,omitempty"`

	// +optional
	// +kubebuilder:default=Disabled

	// EnterBootSetup specifies whether to automatically enter BIOS/EFI setup the next
	// time the virtual machine boots. The virtual machine resets this flag to false
	// so that subsequent boots proceed normally. The available values are:
	//
	// - Enabled -- The virtual machine will automatically enter BIOS/EFI setup the next
	//              time the virtual machine boots.
	// - Disabled -- The virtual machine will boot normaally.
	// TODO: (abaruni) Revisit the need for this option
	// EnterBootSetup VirtualMachineBootOptionsForceBootEntry `json:"enterBootSetup,omitempty"`

	// +optional
	// +kubebuilder:default=Disabled

	// EFISecureBoot specifies whether the virtual machine's firmware will
	// perform signature checks of any EFI images loaded during startup. If set to
	// true, signature checks will be performed and the virtual machine's firmware
	// will refuse to start any images which do not pass those signature checks.
	//
	// Please note, this field will not be honored unless the value of
	// spec.bootOptions.firmware is "EFI". The available values are:
	//
	// - Enabled -- Signature checks will be performed and the virtual machine's firmware
	//              will refuse to start any images which do not pass those signature checks.
	// - Disabled -- No signature checks will be performed.
	EFISecureBoot VirtualMachineBootOptionsEFISecureBoot `json:"efiSecureBoot,omitempty"`

	// +optional
	// +kubebuilder:default=IP4

	// NetworkBootProtocol is the protocol to attempt during PXE network boot or NetBoot.
	// The available protocols are:
	//
	// - IP4 -- PXE (or Apple NetBoot) over IPv4. The default.
	// - IP6 -- PXE over IPv6. Only meaningful for EFI virtual machines.
	NetworkBootProtocol VirtualMachineBootOptionsNetworkBootProtocol `json:"networkBootProtocol,omitempty"`
}

// +kubebuilder:validation:Enum=Direct;Linked

// VirtualMachineDeployMode represents the available modes in which a VM may be
// deployed.
type VirtualMachineDeployMode string

const (
	VirtualMachineDeployModeDirect VirtualMachineDeployMode = "Direct"
	VirtualMachineDeployModeLinked VirtualMachineDeployMode = "Linked"
)

// +kubebuilder:validation:Enum=Online;Offline;Disabled

// VirtualMachinePromoteDisksMode represents the available modes for promoting
// child disks to full clones.
type VirtualMachinePromoteDisksMode string

const (
	VirtualMachinePromoteDisksModeDisabled VirtualMachinePromoteDisksMode = "Disabled"
	VirtualMachinePromoteDisksModeOnline   VirtualMachinePromoteDisksMode = "Online"
	VirtualMachinePromoteDisksModeOffline  VirtualMachinePromoteDisksMode = "Offline"
)

// VirtualMachineSpec defines the desired state of a VirtualMachine.
type VirtualMachineSpec struct {
	// +optional

	// Image describes the reference to the VirtualMachineImage or
	// ClusterVirtualMachineImage resource used to deploy this VM.
	//
	// Please note, unlike the field spec.imageName, the value of
	// spec.image.name MUST be a Kubernetes object name.
	//
	// Please also note, when creating a new VirtualMachine, if this field and
	// spec.imageName are both non-empty, then they must refer to the same
	// resource or an error is returned.
	//
	// Please note, this field *may* be empty if the VM was imported instead of
	// deployed by VM Operator. An imported VirtualMachine resource references
	// an existing VM on the underlying platform that was not deployed from a
	// VM image.
	Image *VirtualMachineImageRef `json:"image,omitempty"`

	// +optional

	// ImageName describes the name of the image resource used to deploy this
	// VM.
	//
	// This field may be used to specify the name of a VirtualMachineImage
	// or ClusterVirtualMachineImage resource. The resolver first checks to see
	// if there is a VirtualMachineImage with the specified name in the
	// same namespace as the VM being deployed. If no such resource exists, the
	// resolver then checks to see if there is a ClusterVirtualMachineImage
	// resource with the specified name.
	//
	// This field may also be used to specify the display name (vSphere name) of
	// a VirtualMachineImage or ClusterVirtualMachineImage resource. If the
	// display name unambiguously resolves to a distinct VM image (among all
	// existing VirtualMachineImages in the VM's namespace and all existing
	// ClusterVirtualMachineImages), then a mutation webhook updates the
	// spec.image field with the reference to the resolved VM image. If the
	// display name resolves to multiple or no VM images, then the mutation
	// webhook denies the request and returns an error.
	//
	// Please also note, when creating a new VirtualMachine, if this field and
	// spec.image are both non-empty, then they must refer to the same
	// resource or an error is returned.
	//
	// Please note, this field *may* be empty if the VM was imported instead of
	// deployed by VM Operator. An imported VirtualMachine resource references
	// an existing VM on the underlying platform that was not deployed from a
	// VM image.
	ImageName string `json:"imageName,omitempty"`

	// +optional

	// ClassName describes the name of the VirtualMachineClass resource used to
	// deploy this VM.
	//
	// When creating a virtual machine, if this field is empty and a
	// VirtualMachineClassInstance is specified in spec.class, then
	// this field is populated with the VirtualMachineClass object's
	// name.
	//
	// Please also note, when creating a new VirtualMachine, if this field and
	// spec.class are both non-empty, then they must refer to the same
	// VirtualMachineClass or an error is returned.
	//
	// Please note, this field *may* be empty if the VM was imported instead of
	// deployed by VM Operator. An imported VirtualMachine resource references
	// an existing VM on the underlying platform that was not deployed from a
	// VM class.
	//
	// If a VM is using a class, a different value in spec.className
	// leads to the VM being resized.
	ClassName string `json:"className,omitempty"`

	// +optional

	// Class describes the VirtualMachineClassInstance resource that is
	// referenced by this virtual machine. This can be the
	// VirtualMachineClassInstance that the virtual machine was
	// created, or later resized with.
	//
	// The value of spec.class.Name must be the Kubernetes object name
	// of a valid VirtualMachineClassInstance resource.
	//
	// Please also note, if this field and spec.className are both
	// non-empty, then they must refer to the same VirtualMachineClass
	// or an error is returned.
	//
	// If a className is specified, but this field is omitted, VM operator
	// picks the latest instance for the VM class to create the VM.
	//
	// If a VM class has been modified and thus, the newly available
	// VirtualMachineClassInstance can be specified in spec.class to
	// trigger a resize operation.
	Class *vmopv1a5common.LocalObjectRef `json:"class,omitempty"`

	// +optional

	// Affinity describes the VM's scheduling constraints.
	Affinity *AffinitySpec `json:"affinity,omitempty"`

	// +optional

	// Crypto describes the desired encryption state of the VirtualMachine.
	Crypto *VirtualMachineCryptoSpec `json:"crypto,omitempty"`

	// +optional

	// StorageClass describes the name of a Kubernetes StorageClass resource
	// used to configure this VM's storage-related attributes.
	//
	// Please see https://kubernetes.io/docs/concepts/storage/storage-classes/
	// for more information on Kubernetes storage classes.
	StorageClass string `json:"storageClass,omitempty"`

	// +optional

	// Bootstrap describes the desired state of the guest's bootstrap
	// configuration.
	//
	// If omitted, a default bootstrap method may be selected based on the
	// guest OS identifier. If Linux, then the LinuxPrep method is used.
	Bootstrap *VirtualMachineBootstrapSpec `json:"bootstrap,omitempty"`

	// +optional

	// Network describes the desired network configuration for the VM.
	//
	// Please note this value may be omitted entirely and the VM will be
	// assigned a single, virtual network interface that is connected to the
	// Namespace's default network.
	Network *VirtualMachineNetworkSpec `json:"network,omitempty"`

	// +optional

	// PowerState describes the desired power state of a VirtualMachine.
	//
	// Please note this field may be omitted when creating a new VM and will
	// default to "PoweredOn." However, once the field is set to a non-empty
	// value, it may no longer be set to an empty value.
	//
	// Additionally, setting this value to "Suspended" is not supported when
	// creating a new VM. The valid values when creating a new VM are
	// "PoweredOn" and "PoweredOff." An empty value is also allowed on create
	// since this value defaults to "PoweredOn" for new VMs.
	PowerState VirtualMachinePowerState `json:"powerState,omitempty"`

	// +optional
	// +kubebuilder:default=TrySoft

	// PowerOffMode describes the desired behavior when powering off a VM.
	//
	// There are three, supported power off modes: Hard, Soft, and
	// TrySoft. The first mode, Hard, is the equivalent of a physical
	// system's power cord being ripped from the wall. The Soft mode
	// requires the VM's guest to have VM Tools installed and attempts to
	// gracefully shutdown the VM. Its variant, TrySoft, first attempts
	// a graceful shutdown, and if that fails or the VM is not in a powered off
	// state after five minutes, the VM is halted.
	//
	// If omitted, the mode defaults to TrySoft.
	PowerOffMode VirtualMachinePowerOpMode `json:"powerOffMode,omitempty"`

	// +optional
	// +kubebuilder:default=TrySoft

	// SuspendMode describes the desired behavior when suspending a VM.
	//
	// There are three, supported suspend modes: Hard, Soft, and
	// TrySoft. The first mode, Hard, is where vSphere suspends the VM to
	// disk without any interaction inside of the guest. The Soft mode
	// requires the VM's guest to have VM Tools installed and attempts to
	// gracefully suspend the VM. Its variant, TrySoft, first attempts
	// a graceful suspend, and if that fails or the VM is not in a put into
	// standby by the guest after five minutes, the VM is suspended.
	//
	// If omitted, the mode defaults to TrySoft.
	SuspendMode VirtualMachinePowerOpMode `json:"suspendMode,omitempty"`

	// +optional

	// NextRestartTime may be used to restart the VM, in accordance with
	// RestartMode, by setting the value of this field to "now"
	// (case-insensitive).
	//
	// A mutating webhook changes this value to the current time (UTC), which
	// the VM controller then uses to determine the VM should be restarted by
	// comparing the value to the timestamp of the last time the VM was
	// restarted.
	//
	// Please note it is not possible to schedule future restarts using this
	// field. The only value that users may set is the string "now"
	// (case-insensitive).
	NextRestartTime string `json:"nextRestartTime,omitempty"`

	// +optional
	// +kubebuilder:default=TrySoft

	// RestartMode describes the desired behavior for restarting a VM when
	// spec.nextRestartTime is set to "now" (case-insensitive).
	//
	// There are three, supported suspend modes: Hard, Soft, and
	// TrySoft. The first mode, Hard, is where vSphere resets the VM without any
	// interaction inside of the guest. The Soft mode requires the VM's guest to
	// have VM Tools installed and asks the guest to restart the VM. Its
	// variant, TrySoft, first attempts a soft restart, and if that fails or
	// does not complete within five minutes, the VM is hard reset.
	//
	// If omitted, the mode defaults to TrySoft.
	RestartMode VirtualMachinePowerOpMode `json:"restartMode,omitempty"`

	// +optional
	// +listType=map
	// +listMapKey=name

	// Volumes describes a list of volumes that can be mounted to the VM.
	Volumes []VirtualMachineVolume `json:"volumes,omitempty"`

	// +optional

	// ReadinessProbe describes a probe used to determine the VM's ready state.
	ReadinessProbe *VirtualMachineReadinessProbeSpec `json:"readinessProbe,omitempty"`

	// +optional

	// Advanced describes a set of optional, advanced VM configuration options.
	Advanced *VirtualMachineAdvancedSpec `json:"advanced,omitempty"`

	// +optional

	// Reserved describes a set of VM configuration options reserved for system
	// use.
	//
	// Please note attempts to modify the value of this field by a DevOps user
	// will result in a validation error.
	Reserved *VirtualMachineReservedSpec `json:"reserved,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=13

	// MinHardwareVersion describes the desired, minimum hardware version.
	//
	// The logic that determines the hardware version is as follows:
	//
	// 1. If this field is set, then its value is used.
	// 2. Otherwise, if the VirtualMachineClass used to deploy the VM contains a
	//    non-empty hardware version, then it is used.
	// 3. Finally, if the hardware version is still undetermined, the value is
	//    set to the default hardware version for the Datacenter/Cluster/Host
	//    where the VM is provisioned.
	//
	// This field is never updated to reflect the derived hardware version.
	// Instead, VirtualMachineStatus.HardwareVersion surfaces
	// the observed hardware version.
	//
	// Please note, setting this field's value to N ensures a VM's hardware
	// version is equal to or greater than N. For example, if a VM's observed
	// hardware version is 10 and this field's value is 13, then the VM will be
	// upgraded to hardware version 13. However, if the observed hardware
	// version is 17 and this field's value is 13, no change will occur.
	//
	// Several features are hardware version dependent, for example:
	//
	// * NVMe Controllers                >= 14
	// * Dynamic Direct Path I/O devices >= 17
	//
	// Please refer to https://kb.vmware.com/s/article/1003746 for a list of VM
	// hardware versions.
	//
	// It is important to remember that a VM's hardware version may not be
	// downgraded and upgrading a VM deployed from an image based on an older
	// hardware version to a more recent one may result in unpredictable
	// behavior. In other words, please be careful when choosing to upgrade a
	// VM to a newer hardware version.
	MinHardwareVersion int32 `json:"minHardwareVersion,omitempty"`

	// +optional
	// +kubebuilder:validation:Format:=uuid

	// InstanceUUID describes the desired Instance UUID for a VM.
	// If omitted, this field defaults to a random UUID.
	// This value is only used for the VM Instance UUID,
	// it is not used within cloudInit.
	// This identifier is used by VirtualCenter to uniquely identify all
	// virtual machine instances, including those that may share the same BIOS UUID.
	InstanceUUID string `json:"instanceUUID,omitempty"`

	// +optional
	// +kubebuilder:validation:Format:=uuid

	// BiosUUID describes the desired BIOS UUID for a VM.
	// If omitted, this field defaults to a random UUID.
	// When the bootstrap provider is Cloud-Init, this value is used as the
	// default value for spec.bootstrap.cloudInit.instanceID if it is omitted.
	BiosUUID string `json:"biosUUID,omitempty"`

	// +optional

	// GuestID describes the desired guest operating system identifier for a VM.
	//
	// The logic that determines the guest ID is as follows:
	//
	// If this field is set, then its value is used.
	// Otherwise, if the VM is deployed from an OVF template that defines a
	// guest ID, then that value is used.
	// The guest ID from VirtualMachineClass used to deploy the VM is ignored.
	//
	// For a complete list of supported values, please refer to
	// https://developer.broadcom.com/xapis/vsphere-web-services-api/latest/vim.vm.GuestOsDescriptor.GuestOsIdentifier.html.
	//
	// Please note that some guest ID values may require a minimal hardware
	// version, which can be set using the `spec.minHardwareVersion` field.
	// To see the mapping between virtual hardware versions and the product
	// versions that support a specific guest ID, please refer to
	// https://knowledge.broadcom.com/external/article/315655/virtual-machine-hardware-versions.html.
	//
	// Please note that this field is immutable after the VM is powered on.
	// To change the guest ID after the VM is powered on, the VM must be powered
	// off and then powered on again with the updated guest ID spec.
	//
	// This field is required when the VM has any CD-ROM devices attached.
	GuestID string `json:"guestID,omitempty"`

	// +optional
	// +kubebuilder:default=Online

	// PromoteDisksMode describes the mode used to promote a VM's delta disks to
	// full disks. The available modes are:
	//
	// - Disabled -- Do not promote disks.
	// - Online   -- Promote disks while the VM is powered on. VMs with
	//               snapshots do not support online promotion.
	// - Offline  -- Promote disks while the VM is powered off.
	//
	// Please note, this field is ignored for encrypted VMs since they do not
	// use delta disks.
	//
	// Defaults to Online.
	PromoteDisksMode VirtualMachinePromoteDisksMode `json:"promoteDisksMode,omitempty"`

	// +optional

	// BootOptions describes the settings that control the boot behavior of the
	// virtual machine. These settings take effect during the next power-on of the
	// virtual machine.
	BootOptions *VirtualMachineBootOptions `json:"bootOptions,omitempty"`

	// +optional

	// CurrentSnapshotName represents the desired snapshot that the VM
	// should point to. This field can be specified to revert the VM
	// to a given snapshot. Once the virtual machine has been
	// successfully reverted to the desired snapshot, the value of
	// this field is cleared.
	//
	// The value of this field must be the name of an existing
	// VirtualMachineSnapshot resource in the same namespace.
	//
	// Reverting a virtual machine to a snapshot rolls back the data
	// and the configuration of the virtual machine to that of the
	// specified snapshot. The VirtualMachineSpec of the
	// VirtualMachine resource is replaced from the one stored with
	// the snapshot.
	//
	// If the virtual machine is currently powered off, but you revert to
	// a snapshot that was taken while the VM was powered on, then the
	// VM will be automatically powered on during the revert.
	// Additionally, the VirtualMachineSpec will be updated to match
	// the power state from the snapshot (i.e., powered on). This can
	// be overridden by specifying the PowerState to PoweredOff in the
	// VirtualMachineSpec.
	CurrentSnapshotName string `json:"currentSnapshotName,omitempty"`

	// +optional

	// GroupName indicates the name of the VirtualMachineGroup to which this
	// VM belongs.
	//
	// VMs that belong to a group do not drive their own placement, rather that
	// is handled by the group.
	//
	// When this field is set to a valid group that contains this VM as a
	// member, an owner reference to that group is added to this VM.
	//
	// When this field is deleted or changed, any existing owner reference to
	// the previous group will be removed from this VM.
	GroupName string `json:"groupName,omitempty"`

	// +optional

	// Hardware describes the VM's desired hardware.
	Hardware *VirtualMachineHardwareSpec `json:"hardware,omitempty"`

	// +optional

	// Policies describes a list of policies that should be explicitly applied
	// to this VM.
	//
	// Please note, not all policies may be applied explicitly to a VM. Please
	// consult a policy to determine if it may be applied directly. For example,
	// the ComputePolicy object from the vsphere.policy.vmware.com API group
	// has a field named spec.type. Only ComputePolicy objects with
	// type=Optional may be applied explicitly to a VM.
	//
	// Valid policy types are: ComputePolicy.
	Policies []PolicySpec `json:"policies,omitempty"`
}

// VirtualMachineReservedSpec describes a set of VM configuration options
// reserved for system use. Modification attempts by DevOps users will result
// in a validation error.
type VirtualMachineReservedSpec struct {
	// +optional

	// ResourcePolicyName describes the name of a
	// VirtualMachineSetResourcePolicy resource used to configure the VM's

	ResourcePolicyName string `json:"resourcePolicyName,omitempty"`
}

// VirtualMachineAdvancedSpec describes a set of optional, advanced VM
// configuration options.
type VirtualMachineAdvancedSpec struct {
	// +optional

	// BootDiskCapacity is the capacity of the VM's boot disk -- the first disk
	// from the VirtualMachineImage from which the VM was deployed.
	//
	// Please note it is not advised to change this value while the VM is
	// running. Also, resizing the VM's boot disk may require actions inside of
	// the guest to take advantage of the additional capacity. Finally, changing
	// the size of the VM's boot disk, even increasing it, could adversely
	// affect the VM.
	//
	// Please note this field is ignored if the VM is deployed from an ISO with
	// CD-ROM devices attached.
	BootDiskCapacity *resource.Quantity `json:"bootDiskCapacity,omitempty"`

	// +optional

	// DefaultVolumeProvisioningMode specifies the default provisioning mode for
	// persistent volumes managed by this VM.
	DefaultVolumeProvisioningMode VolumeProvisioningMode `json:"defaultVolumeProvisioningMode,omitempty"`

	// +optional

	// ChangeBlockTracking is a flag that enables incremental backup support
	// for this VM, a feature utilized by external backup systems such as
	// VMware Data Recovery.
	ChangeBlockTracking *bool `json:"changeBlockTracking,omitempty"`
}

type VirtualMachineEncryptionType string

const (
	VirtualMachineEncryptionTypeConfig VirtualMachineEncryptionType = "Config"
	VirtualMachineEncryptionTypeDisks  VirtualMachineEncryptionType = "Disks"
)

type VirtualMachineCryptoStatus struct {
	// +optional

	// Encrypted describes the observed state of the VirtualMachine's
	// encryption. There may be two values in this list:
	//
	// - Config -- This refers to all of the files related to a VM except any
	//             virtual disks.
	// - Disks  -- This refers to at least one of the VM's attached disks. To
	//             determine the encryption state of the individual disks,
	//             please refer to status.volumes[].crypto.
	Encrypted []VirtualMachineEncryptionType `json:"encrypted,omitempty"`

	// +optional

	// ProviderID describes the provider ID used to encrypt the VirtualMachine.
	// Please note, this field will be empty if the VirtualMachine is not
	// encrypted.
	ProviderID string `json:"providerID,omitempty"`

	// +optional

	// KeyID describes the key ID used to encrypt the VirtualMachine.
	// Please note, this field will be empty if the VirtualMachine is not
	// encrypted.
	KeyID string `json:"keyID,omitempty"`

	// +optional

	// HasVTPM indicates whether or not the VM has a vTPM.
	HasVTPM bool `json:"hasVTPM,omitempty"`
}

type VirtualMachineGuestStatus struct {
	// +optional

	// GuestID describes the ID of the observed operating system.
	GuestID string `json:"guestID,omitempty"`

	// +optional

	// GuestFullName describes the full name of the observed operating system.
	GuestFullName string `json:"guestFullName,omitempty"`
}

// VirtualMachineProviderStatus describes the observed state of the
// VirtualMachine from the underlying infrastructure provider (vSphere/vCenter).
type VirtualMachineProviderStatus struct {
	// +optional

	// CreationTimestamp describes the timestamp the underlying
	// hypervisor reports when the VM is successfully created.
	//
	// This value should be immutable, but ultimately it is up to the
	// provider as to whether or not the value ever changes.
	CreationTimestamp *metav1.Time `json:"creationTimestamp,omitempty"`
}

// VirtualMachineStatus defines the observed state of a VirtualMachine instance.
type VirtualMachineStatus struct {
	// +optional

	// Class is a reference to the VirtualMachineClass resource used to deploy
	// this VM.
	Class *vmopv1a5common.LocalObjectRef `json:"class,omitempty"`

	// +optional

	// NodeName describes the observed name of the node where the VirtualMachine
	// is scheduled.
	NodeName string `json:"nodeName,omitempty"`

	// +optional

	// PowerState describes the observed power state of the VirtualMachine.
	PowerState VirtualMachinePowerState `json:"powerState,omitempty"`

	// +optional

	// Conditions describes the observed conditions of the VirtualMachine.
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// +optional

	// Crypto describes the observed state of the VirtualMachine's encryption
	// configuration.
	Crypto *VirtualMachineCryptoStatus `json:"crypto,omitempty"`

	// +optional

	// Network describes the observed state of the VM's network configuration.
	// Please note much of the network status information is only available if
	// the guest has VM Tools installed.
	Network *VirtualMachineNetworkStatus `json:"network,omitempty"`

	// +optional

	// UniqueID describes a unique identifier that is provided by the underlying
	// infrastructure provider, such as vSphere.
	UniqueID string `json:"uniqueID,omitempty"`

	// +optional

	// BiosUUID describes a unique identifier provided by the underlying
	// infrastructure provider that is exposed to the Guest OS BIOS as a unique
	// hardware identifier.
	BiosUUID string `json:"biosUUID,omitempty"`

	// +optional

	// InstanceUUID describes the unique instance UUID provided by the
	// underlying infrastructure provider, such as vSphere.
	InstanceUUID string `json:"instanceUUID,omitempty"`

	// +optional
	// +listType=map
	// +listMapKey=name
	// +listMapKey=type

	// Volumes describes the observed state of the volumes that are intended to
	// be attached to the VirtualMachine.
	Volumes []VirtualMachineVolumeStatus `json:"volumes,omitempty"`

	// +optional

	// ChangeBlockTracking describes whether or not change block tracking is
	// enabled for the VirtualMachine.
	ChangeBlockTracking *bool `json:"changeBlockTracking,omitempty"`

	// +optional

	// Zone describes the availability zone where the VirtualMachine has been
	// scheduled.
	//
	// Please note this field may be empty when the cluster is not zone-aware.
	Zone string `json:"zone,omitempty"`

	// +optional

	// LastRestartTime describes the last time the VM was restarted.
	LastRestartTime *metav1.Time `json:"lastRestartTime,omitempty"`

	// +optional

	// HardwareVersion describes the VirtualMachine resource's observed
	// hardware version.
	//
	// Please refer to VirtualMachineSpec.MinHardwareVersion for more
	// information on the topic of a VM's hardware version.
	HardwareVersion int32 `json:"hardwareVersion,omitempty"`

	// +optional

	// Storage describes the observed state of the VirtualMachine's storage.
	Storage *VirtualMachineStorageStatus `json:"storage,omitempty"`

	// +optional

	// Provider describes the observed state of the VirtualMachine from
	// the underlying infrastructure provider (vSphere/vCenter).
	Provider *VirtualMachineProviderStatus `json:"provider,omitempty"`

	// +optional

	// CurrentSnapshot describes the observed working snapshot of the VirtualMachine.
	// This field contains the name of the current snapshot.
	CurrentSnapshot *VirtualMachineSnapshotReference `json:"currentSnapshot,omitempty"`

	// +optional

	// RootSnapshots represents the observed list of root snapshots of
	// a VM. Since each snapshot includes the list of its child
	// snapshots, these root snapshot references can effectively be
	// used to construct the entire snapshot chain of a virtual
	// machine.
	RootSnapshots []VirtualMachineSnapshotReference `json:"rootSnapshots,omitempty"`

	// Guest describes the observed state of the VM's guest.
	Guest *VirtualMachineGuestStatus `json:"guest,omitempty"`

	// Hardware describes the observed state of the VM's hardware.
	Hardware *VirtualMachineHardwareStatus `json:"hardware,omitempty"`

	// +optional

	// Policies describes the observed policies applied to this VM.
	Policies []PolicyStatus `json:"policies,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=vm
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Power-State",type="string",JSONPath=".status.powerState"
// +kubebuilder:printcolumn:name="Class",type="string",priority=1,JSONPath=".spec.className"
// +kubebuilder:printcolumn:name="Image",type="string",priority=1,JSONPath=".spec.image.name"
// +kubebuilder:printcolumn:name="Primary-IP4",type="string",priority=1,JSONPath=".status.network.primaryIP4"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// VirtualMachine is the schema for the virtualmachines API and represents the
// desired state and observed status of a virtualmachines resource.
type VirtualMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualMachineSpec   `json:"spec,omitempty"`
	Status VirtualMachineStatus `json:"status,omitempty"`
}

func (vm *VirtualMachine) SetAnnotation(k, v string) {
	if vm.Annotations == nil {
		vm.Annotations = map[string]string{}
	}
	vm.Annotations[k] = v
}

func (vm VirtualMachine) NamespacedName() string {
	return vm.Namespace + "/" + vm.Name
}

func (vm VirtualMachine) GetConditions() []metav1.Condition {
	return vm.Status.Conditions
}

func (vm *VirtualMachine) SetConditions(conditions []metav1.Condition) {
	vm.Status.Conditions = conditions
}

func (vm VirtualMachine) GetMemberKind() string {
	return "VirtualMachine"
}

func (vm VirtualMachine) GetGroupName() string {
	return vm.Spec.GroupName
}

func (vm *VirtualMachine) SetGroupName(value string) {
	vm.Spec.GroupName = value
}

func (vm VirtualMachine) GetPowerState() VirtualMachinePowerState {
	return vm.Status.PowerState
}

func (vm *VirtualMachine) SetPowerState(value VirtualMachinePowerState) {
	vm.Spec.PowerState = value
}

// +kubebuilder:object:root=true

// VirtualMachineList contains a list of VirtualMachine.
type VirtualMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualMachine `json:"items"`
}

func init() {
	objectTypes = append(objectTypes, &VirtualMachine{}, &VirtualMachineList{})
}
