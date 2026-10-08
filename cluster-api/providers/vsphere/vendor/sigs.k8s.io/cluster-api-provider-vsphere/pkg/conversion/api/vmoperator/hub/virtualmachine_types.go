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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	conversionmeta "sigs.k8s.io/cluster-api-provider-vsphere/pkg/conversion/api/meta"
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

// VirtualMachineSpec defines the desired state of a VirtualMachine.
type VirtualMachineSpec struct {
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

	// Affinity describes the VM's scheduling constraints.
	Affinity *AffinitySpec `json:"affinity,omitempty"`

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
	// +listType=map
	// +listMapKey=name

	// Volumes describes a list of volumes that can be mounted to the VM.
	Volumes []VirtualMachineVolume `json:"volumes,omitempty"`

	// +optional

	// ReadinessProbe describes a probe used to determine the VM's ready state.
	ReadinessProbe *VirtualMachineReadinessProbeSpec `json:"readinessProbe,omitempty"`

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

	// Policies describes a list of policies that should be explicitly applied
	// to this VM.
	//
	// Please note, not all policies may be applied explicitly to a VM. Please
	// consult a policy to determine if it may be applied directly. For example,
	// the ComputePolicy object from the vsphere.policy.vmware.com API group
	// has a field named spec.type. Only ComputePolicy objects with
	// type=Optional may be applied explicitly to a VM.
	// +optional
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

// VirtualMachineStatus defines the observed state of a VirtualMachine instance.
type VirtualMachineStatus struct {
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

	// Network describes the observed state of the VM's network configuration.
	// Please note much of the network status information is only available if
	// the guest has VM Tools installed.
	Network *VirtualMachineNetworkStatus `json:"network,omitempty"`

	// +optional

	// BiosUUID describes a unique identifier provided by the underlying
	// infrastructure provider that is exposed to the Guest OS BIOS as a unique
	// hardware identifier.
	BiosUUID string `json:"biosUUID,omitempty"`

	// +optional

	// Zone describes the availability zone where the VirtualMachine has been
	// scheduled.
	//
	// Please note this field may be empty when the cluster is not zone-aware.
	Zone string `json:"zone,omitempty"`

	// Policies describes the observed policies applied to this VM.
	// +optional
	Policies []PolicyStatus `json:"policies,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=vm
// +kubebuilder:storageversion
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

	Source conversionmeta.SourceTypeMeta `json:"source,omitempty,omitzero"`
}

// GetConditions returns the set of conditions for this object.
func (vm VirtualMachine) GetConditions() []metav1.Condition {
	return vm.Status.Conditions
}

// SetConditions sets conditions for an API object.
func (vm *VirtualMachine) SetConditions(conditions []metav1.Condition) {
	vm.Status.Conditions = conditions
}

// GetSource returns the Source for this object.
func (vm *VirtualMachine) GetSource() conversionmeta.SourceTypeMeta {
	return vm.Source
}

// SetSource sets Source for an API object.
func (vm *VirtualMachine) SetSource(source conversionmeta.SourceTypeMeta) {
	vm.Source = source
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
