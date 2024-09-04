// Copyright (c) 2020-2024 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// VirtualMachinePowerState represents the power state of a VirtualMachine.
// +kubebuilder:validation:Enum=poweredOn;poweredOff;suspended
type VirtualMachinePowerState string

const (
	// VirtualMachinePoweredOff indicates to shut down a VM and/or it is
	// shut down.
	VirtualMachinePoweredOff VirtualMachinePowerState = "poweredOff"

	// VirtualMachinePoweredOn indicates to power on a VM and/or it is
	// powered on.
	VirtualMachinePoweredOn VirtualMachinePowerState = "poweredOn"

	// VirtualMachineSuspended indicates to suspend a VM and/or it is
	// suspended.
	VirtualMachineSuspended VirtualMachinePowerState = "suspended"
)

// VirtualMachinePowerOpMode represents the various power operation modes when
// powering off or suspending a VM.
// +kubebuilder:validation:Enum=hard;soft;trySoft
type VirtualMachinePowerOpMode string

const (
	// VirtualMachinePowerOpModeHard indicates to halt a VM when powering it
	// off or when suspending a VM to not involve the guest.
	VirtualMachinePowerOpModeHard VirtualMachinePowerOpMode = "hard"

	// VirtualMachinePowerOpModeSoft indicates to ask VM Tools running
	// inside of a VM's guest to shutdown the guest gracefully when powering
	// off a VM or when suspending a VM to allow the guest to participate.
	//
	// If this mode is set on a VM whose guest does not have VM Tools or if
	// VM Tools is present but the operation fails, the VM may never realize
	// the desired power state. This can prevent a VM from being deleted as well
	// as many other unexpected issues. It is recommended to use trySoft
	// instead.
	VirtualMachinePowerOpModeSoft VirtualMachinePowerOpMode = "soft"

	// VirtualMachinePowerOpModeTrySoft indicates to first attempt a Soft
	// operation and fall back to hard if VM Tools is not present in the guest,
	// if the soft operation fails, or if the VM is not in the desired power
	// state within five minutes.
	VirtualMachinePowerOpModeTrySoft VirtualMachinePowerOpMode = "trySoft"
)

// VMStatusPhase is used to indicate the phase of a VirtualMachine's lifecycle.
type VMStatusPhase string

const (
	// Creating phase indicates that the VirtualMachine is being created by the backing infrastructure provider.
	Creating VMStatusPhase = "Creating"

	// Created phase indicates that the VirtualMachine has been already been created by the backing infrastructure
	// provider.
	Created VMStatusPhase = "Created"

	// Deleting phase indicates that the VirtualMachine is being deleted by the backing infrastructure provider.
	Deleting VMStatusPhase = "Deleting"

	// Deleted phase indicates that the VirtualMachine has been deleted by the backing infrastructure provider.
	Deleted VMStatusPhase = "Deleted"

	// Unknown phase indicates that the VirtualMachine status cannot be determined from the backing infrastructure
	// provider.
	Unknown VMStatusPhase = "Unknown"
)

const (
	// PauseAnnotation is an annotation that can be applied to any VirtualMachine object to prevent VM Operator from
	// reconciling the object with the vSphere infrastructure. VM Operator checks the presence of this annotation to
	// skip the reconcile of a VirtualMachine.
	//
	// This can be used when a Virtual Machine needs to be modified out-of-band of VM Operator on the infrastructure
	// directly (e.g., during a VADP based Restore operation).
	PauseAnnotation = GroupName + "/pause-reconcile"

	// NoDefaultNicAnnotation is an annotation that can be applied to prevent VM Operator from creating a default nic for
	// a VirtualMachine object with empty VirtualMachineNetworkInterfaces list.
	//
	// This can be used when users want to opt out a default network device when creating new VirtualMachines.
	//
	// When a VM without any VirtualMachineNetworkInterfaces is being created, VM Operator webhook checks the presence of
	// this annotation to skip adding a default nic. VM Operator won't add default NIC to any existing VMs or new VMs
	// with VirtualMachineNetworkInterfaces specified. This annotation is not required for such VMs.
	NoDefaultNicAnnotation = GroupName + "/no-default-nic"

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
)

const (
	// ManagedByExtensionKey and ManagedByExtensionType represent the ManagedBy
	// field on the VM. They are used to differentiate VM Service managed VMs
	// from traditional vSphere VMs.
	ManagedByExtensionKey  = "com.vmware.vcenter.wcp"
	ManagedByExtensionType = "VirtualMachine"
)

// VirtualMachinePort is unused and can be considered deprecated.
type VirtualMachinePort struct {
	Port     int             `json:"port"`
	Ip       string          `json:"ip"` //nolint:revive,stylecheck
	Name     string          `json:"name"`
	Protocol corev1.Protocol `json:"protocol"`
}

// NetworkInterfaceProviderReference contains info to locate a network interface provider object.
type NetworkInterfaceProviderReference struct {
	// APIGroup is the group for the resource being referenced.
	APIGroup string `json:"apiGroup"`
	// Kind is the type of resource being referenced
	Kind string `json:"kind"`
	// Name is the name of resource being referenced
	Name string `json:"name"`
	// API version of the referent.
	APIVersion string `json:"apiVersion,omitempty"`
}

// VirtualMachineNetworkInterface defines the properties of a network interface to attach to a VirtualMachine
// instance.  A VirtualMachineNetworkInterface describes network interface configuration that is used by the
// VirtualMachine controller when integrating the VirtualMachine into a VirtualNetwork. Currently, only NSX-T
// and vSphere Distributed Switch (VDS) type network integrations are supported using this VirtualMachineNetworkInterface
// structure.
type VirtualMachineNetworkInterface struct {
	// NetworkType describes the type of VirtualNetwork that is referenced by the NetworkName. Currently, the supported
	// NetworkTypes are "nsx-t", "nsx-t-subnet", "nsx-t-subnetset" and "vsphere-distributed".
	// +optional
	NetworkType string `json:"networkType,omitempty"`

	// NetworkName describes the name of an existing virtual network that this interface should be added to.
	// For "nsx-t" NetworkType, this is the name of a pre-existing NSX-T VirtualNetwork. If unspecified,
	// the default network for the namespace will be used. For "vsphere-distributed" NetworkType, the
	// NetworkName must be specified.
	// +optional
	NetworkName string `json:"networkName,omitempty"`

	// ProviderRef is reference to a network interface provider object that specifies the network interface configuration.
	// If unset, default configuration is assumed.
	// +optional
	ProviderRef *NetworkInterfaceProviderReference `json:"providerRef,omitempty"`

	// EthernetCardType describes an optional ethernet card that should be used by the VirtualNetworkInterface (vNIC)
	// associated with this network integration.  The default is "vmxnet3".
	// +optional
	EthernetCardType string `json:"ethernetCardType,omitempty"`
}

// VirtualMachineMetadataTransport is used to indicate the transport used by VirtualMachineMetadata
// Valid values are "ExtraConfig", "OvfEnv", "vAppConfig", "CloudInit", and "Sysprep".
// +kubebuilder:validation:Enum=ExtraConfig;OvfEnv;vAppConfig;CloudInit;Sysprep
type VirtualMachineMetadataTransport string

const (
	// VirtualMachineMetadataExtraConfigTransport indicates that the data set in
	// the VirtualMachineMetadata Transport Resource, i.e., a ConfigMap or Secret,
	// will be extraConfig key value fields on the VM.
	// Only keys prefixed with "guestinfo." will be set.
	VirtualMachineMetadataExtraConfigTransport VirtualMachineMetadataTransport = "ExtraConfig"

	// VirtualMachineMetadataOvfEnvTransport indicates that the data set in
	// the VirtualMachineMetadata Transport Resource, i.e., a ConfigMap or Secret,
	// will be vApp properties on the VM, which will be exposed as OvfEnv to the Guest VM.
	// Only properties marked userConfigurable and already present in either
	// OVF Properties of a VirtualMachineImage or as vApp properties on an existing VM
	// or VMTX will be set, all others will be ignored.
	//
	// This transport uses Guest OS customization for networking.
	VirtualMachineMetadataOvfEnvTransport VirtualMachineMetadataTransport = "OvfEnv"

	// VirtualMachineMetadataVAppConfigTransport indicates that the data set in
	// the VirtualMachineMetadata Transport Resource, i.e., a ConfigMap or Secret,
	// will be vApp properties on the VM, which will be exposed as vAppConfig to the Guest VM.
	// Only properties marked userConfigurable and already present in either
	// OVF Properties of a VirtualMachineImage or as vApp properties on an existing VM
	// or VMTX will be set, all others will be ignored.
	//
	// Selecting this transport means the guest's network is not automatically
	// configured by vm-tools. This transport should only be selected if the image
	// exposes OVF/vApp properties that are used by the guest to bootstrap
	// its networking configuration.
	VirtualMachineMetadataVAppConfigTransport VirtualMachineMetadataTransport = "vAppConfig"

	// VirtualMachineMetadataCloudInitTransport indicates the data set in
	// the VirtualMachineMetadata Transport Resource, i.e., a ConfigMap or Secret,
	// in the "user-data" key is cloud-init userdata.
	//
	// Please note that, despite the name, VirtualMachineMetadata has no
	// relationship to cloud-init instance metadata.
	//
	// For more information, please refer to cloud-init's official documentation.
	VirtualMachineMetadataCloudInitTransport VirtualMachineMetadataTransport = "CloudInit"

	// VirtualMachineMetadataSysprepTransport indicates the data set in
	// the VirtualMachineMetadata Transport Resource, i.e., a ConfigMap or Secret,
	// in the "unattend" key is an XML, Sysprep answers file.
	//
	// For more information, please refer to Microsoft's documentation on
	// "Answer files (unattend.xml)" and "Unattended Windows Setup Reference".
	VirtualMachineMetadataSysprepTransport VirtualMachineMetadataTransport = "Sysprep"
)

// VirtualMachineMetadata defines any metadata that should be passed to the VirtualMachine instance.  A typical use
// case is for this metadata to be used for Guest Customization, however the intended use of the metadata is
// agnostic to the VirtualMachine controller.  VirtualMachineMetadata is read from a configured ConfigMap or a Secret and then
// propagated to the VirtualMachine instance using a desired "Transport" mechanism.
type VirtualMachineMetadata struct {
	// ConfigMapName describes the name of the ConfigMap, in the same Namespace as the VirtualMachine, that should be
	// used for VirtualMachine metadata.  The contents of the Data field of the ConfigMap is used as the VM Metadata.
	// The format of the contents of the VM Metadata are not parsed or interpreted by the VirtualMachine controller.
	// Please note, this field and SecretName are mutually exclusive.
	// +optional
	ConfigMapName string `json:"configMapName,omitempty"`

	// SecretName describes the name of the Secret, in the same Namespace as the VirtualMachine, that should be used
	// for VirtualMachine metadata. The contents of the Data field of the Secret is used as the VM Metadata.
	// The format of the contents of the VM Metadata are not parsed or interpreted by the VirtualMachine controller.
	// Please note, this field and ConfigMapName are mutually exclusive.
	// +optional
	SecretName string `json:"secretName,omitempty"`

	// Transport describes the name of a supported VirtualMachineMetadata transport protocol.  Currently, the only supported
	// transport protocols are "ExtraConfig", "OvfEnv" and "CloudInit".
	Transport VirtualMachineMetadataTransport `json:"transport,omitempty"`
}

// VirtualMachineVolume describes a Volume that should be attached to a specific VirtualMachine.
// Only one of PersistentVolumeClaim, VsphereVolume should be specified.
type VirtualMachineVolume struct {
	// Name specifies the name of the VirtualMachineVolume.  Each volume within the scope of a VirtualMachine must
	// have a unique name.
	Name string `json:"name"`

	// PersistentVolumeClaim represents a reference to a PersistentVolumeClaim
	// in the same namespace. The PersistentVolumeClaim must match one of the
	// following:
	//
	//   * A volume provisioned (either statically or dynamically) by the
	//     cluster's CSI provider.
	//
	//   * An instance volume with a lifecycle coupled to the VM.
	// +optional
	PersistentVolumeClaim *PersistentVolumeClaimVolumeSource `json:"persistentVolumeClaim,omitempty"`

	// VsphereVolume represents a reference to a VsphereVolumeSource in the same namespace. Only one of PersistentVolumeClaim or
	// VsphereVolume can be specified. This is enforced via a webhook
	// +optional
	VsphereVolume *VsphereVolumeSource `json:"vSphereVolume,omitempty"`
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

// VsphereVolumeSource describes a volume source that represent static disks that belong to a VirtualMachine.
type VsphereVolumeSource struct {
	// A description of the virtual volume's resources and capacity
	// +optional
	Capacity corev1.ResourceList `json:"capacity,omitempty"`

	// Device key of vSphere disk.
	// +optional
	DeviceKey *int `json:"deviceKey,omitempty"`
}

// Probe describes a health check to be performed against a VirtualMachine to determine whether it is
// alive or ready to receive traffic. Only one probe action can be specified.
type Probe struct {
	// TCPSocket specifies an action involving a TCP port.
	//
	// Deprecated: The TCPSocket action requires network connectivity that is not supported in all environments.
	// This field will be removed in a later API version.
	// +optional
	TCPSocket *TCPSocketAction `json:"tcpSocket,omitempty"`

	// GuestHeartbeat specifies an action involving the guest heartbeat status.
	// +optional
	GuestHeartbeat *GuestHeartbeatAction `json:"guestHeartbeat,omitempty"`

	// TimeoutSeconds specifies a number of seconds after which the probe times out.
	// Defaults to 10 seconds. Minimum value is 1.
	// +optional
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=60
	TimeoutSeconds int32 `json:"timeoutSeconds,omitempty"`

	// PeriodSeconds specifics how often (in seconds) to perform the probe.
	// Defaults to 10 seconds. Minimum value is 1.
	// +optional
	// +kubebuilder:validation:Minimum:=1
	PeriodSeconds int32 `json:"periodSeconds,omitempty"`
}

// TCPSocketAction describes an action based on opening a socket.
type TCPSocketAction struct {
	// Port specifies a number or name of the port to access on the VirtualMachine.
	// If the format of port is a number, it must be in the range 1 to 65535.
	// If the format of name is a string, it must be an IANA_SVC_NAME.
	Port intstr.IntOrString `json:"port"`

	// Host is an optional host name to connect to.  Host defaults to the VirtualMachine IP.
	// +optional
	Host string `json:"host,omitempty"`
}

// GuestHeartbeatStatus is the status type for a GuestHeartbeat.
type GuestHeartbeatStatus string

// See govmomi.vim25.types.ManagedEntityStatus.
const (
	// VMware Tools are not installed or not running.
	GrayHeartbeatStatus GuestHeartbeatStatus = "gray"
	// No heartbeat. Guest operating system may have stopped responding.
	RedHeartbeatStatus GuestHeartbeatStatus = "red"
	// Intermittent heartbeat. May be due to guest load.
	YellowHeartbeatStatus GuestHeartbeatStatus = "yellow"
	// Guest operating system is responding normally.
	GreenHeartbeatStatus GuestHeartbeatStatus = "green"
)

// GuestHeartbeatAction describes an action based on the guest heartbeat.
type GuestHeartbeatAction struct {
	// ThresholdStatus is the value that the guest heartbeat status must be at or above to be
	// considered successful.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=green
	// +kubebuilder:validation:Enum=yellow;green
	ThresholdStatus GuestHeartbeatStatus `json:"thresholdStatus,omitempty"`
}

// VirtualMachineSpec defines the desired state of a VirtualMachine.
type VirtualMachineSpec struct {
	// ImageName describes the name of the image resource used to deploy this
	// VM.
	//
	// This field may be used to specify the name of a VirtualMachineImage
	// or ClusterVirtualMachineImage resource. The resolver first checks to see
	// if there is a VirtualMachineImage with the specified name. If no
	// such resource exists, the resolver then checks to see if there is a
	// ClusterVirtualMachineImage resource with the specified name in the same
	// Namespace as the VM being deployed.
	//
	// This field may also be used to specify the display name (vSphere name) of
	// a VirtualMachineImage or ClusterVirtualMachineImage resource. If the
	// display name unambiguously resolves to a distinct VM image (among all
	// existing VirtualMachineImages in the VM's namespace and all existing
	// ClusterVirtualMachineImages), then a mutation webhook updates this field
	// with the VM image resource name. If the display name resolves to multiple
	// or no VM images, then the mutation webhook denies the request and outputs
	// an error message accordingly.
	ImageName string `json:"imageName"`

	// ClassName describes the name of a VirtualMachineClass that is to be used as the overlaid resource configuration
	// of VirtualMachine.  A VirtualMachineClass is used to further customize the attributes of the VirtualMachine
	// instance.  See VirtualMachineClass for more description.
	ClassName string `json:"className"`

	// PowerState describes the desired power state of a VirtualMachine.
	//
	// Please note this field may be omitted when creating a new VM and will
	// default to "poweredOn." However, once the field is set to a non-empty
	// value, it may no longer be set to an empty value.
	//
	// Additionally, setting this value to "suspended" is not supported when
	// creating a new VM. The valid values when creating a new VM are
	// "poweredOn" and "poweredOff." An empty value is also allowed on create
	// since this value defaults to "poweredOn" for new VMs.
	//
	// +optional
	PowerState VirtualMachinePowerState `json:"powerState,omitempty"`

	// PowerOffMode describes the desired behavior when powering off a VM.
	//
	// There are three, supported power off modes: hard, soft, and
	// trySoft. The first mode, hard, is the equivalent of a physical
	// system's power cord being ripped from the wall. The soft mode
	// requires the VM's guest to have VM Tools installed and attempts to
	// gracefully shutdown the VM. Its variant, trySoft, first attempts
	// a graceful shutdown, and if that fails or the VM is not in a powered off
	// state after five minutes, the VM is halted.
	//
	// If omitted, the mode defaults to hard.
	//
	// +optional
	// +kubebuilder:default=hard
	PowerOffMode VirtualMachinePowerOpMode `json:"powerOffMode,omitempty"`

	// SuspendMode describes the desired behavior when suspending a VM.
	//
	// There are three, supported suspend modes: hard, soft, and
	// trySoft. The first mode, hard, is where vSphere suspends the VM to
	// disk without any interaction inside of the guest. The soft mode
	// requires the VM's guest to have VM Tools installed and attempts to
	// gracefully suspend the VM. Its variant, trySoft, first attempts
	// a graceful suspend, and if that fails or the VM is not in a put into
	// standby by the guest after five minutes, the VM is suspended.
	//
	// If omitted, the mode defaults to hard.
	//
	// +optional
	// +kubebuilder:default=hard
	SuspendMode VirtualMachinePowerOpMode `json:"suspendMode,omitempty"`

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
	//
	// +optional
	NextRestartTime string `json:"nextRestartTime,omitempty"`

	// RestartMode describes the desired behavior for restarting a VM when
	// spec.nextRestartTime is set to "now" (case-insensitive).
	//
	// There are three, supported suspend modes: hard, soft, and
	// trySoft. The first mode, hard, is where vSphere resets the VM without any
	// interaction inside of the guest. The soft mode requires the VM's guest to
	// have VM Tools installed and asks the guest to restart the VM. Its
	// variant, trySoft, first attempts a soft restart, and if that fails or
	// does not complete within five minutes, the VM is hard reset.
	//
	// If omitted, the mode defaults to hard.
	//
	// +optional
	// +kubebuilder:default=hard
	RestartMode VirtualMachinePowerOpMode `json:"restartMode,omitempty"`

	// Ports is currently unused and can be considered deprecated.
	// +optional
	Ports []VirtualMachinePort `json:"ports,omitempty"`

	// VmMetadata describes any optional metadata that should be passed to the Guest OS.
	// +optional
	VmMetadata *VirtualMachineMetadata `json:"vmMetadata,omitempty"` //nolint:revive,stylecheck

	// StorageClass describes the name of a StorageClass that should be used to configure storage-related attributes of the VirtualMachine
	// instance.
	// +optional
	StorageClass string `json:"storageClass,omitempty"`

	// NetworkInterfaces describes a list of VirtualMachineNetworkInterfaces to be configured on the VirtualMachine instance.
	// Each of these VirtualMachineNetworkInterfaces describes external network integration configurations that are to be
	// used by the VirtualMachine controller when integrating the VirtualMachine into one or more external networks.
	//
	// The maximum number of network interface allowed is 10 because of the limit built into vSphere.
	//
	// +optional
	// +kubebuilder:validation:MaxItems=10
	NetworkInterfaces []VirtualMachineNetworkInterface `json:"networkInterfaces,omitempty"`

	// ResourcePolicyName describes the name of a VirtualMachineSetResourcePolicy to be used when creating the
	// VirtualMachine instance.
	// +optional
	ResourcePolicyName string `json:"resourcePolicyName"`

	// Volumes describes the list of VirtualMachineVolumes that are desired to be attached to the VirtualMachine.  Each of
	// these volumes specifies a volume identity that the VirtualMachine controller will attempt to satisfy, potentially
	// with an external Volume Management service.
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	Volumes []VirtualMachineVolume `json:"volumes,omitempty" patchStrategy:"merge" patchMergeKey:"name"`

	// ReadinessProbe describes a network probe that can be used to determine if the VirtualMachine is available and
	// responding to the probe.
	// +optional
	ReadinessProbe *Probe `json:"readinessProbe,omitempty"`

	// AdvancedOptions describes a set of optional, advanced options for configuring a VirtualMachine
	AdvancedOptions *VirtualMachineAdvancedOptions `json:"advancedOptions,omitempty"`

	// MinHardwareVersion specifies the desired minimum hardware version
	// for this VM.
	//
	// Usually the VM's hardware version is derived from:
	// 1. the VirtualMachineClass used to deploy the VM provided by the ClassName field
	// 2. the datacenter/cluster/host default hardware version
	// Setting this field will ensure that the hardware version of the VM
	// is at least set to the specified value. To enforce this, it will override
	// the value from the VirtualMachineClass.
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
	//
	// +optional
	// +kubebuilder:validation:Minimum=13
	MinHardwareVersion int32 `json:"minHardwareVersion,omitempty"`
}

// VirtualMachineAdvancedOptions describes a set of optional, advanced options for configuring a VirtualMachine.
type VirtualMachineAdvancedOptions struct {
	// DefaultProvisioningOptions specifies the provisioning type to be used by default for VirtualMachine volumes exclusively
	// owned by this VirtualMachine. This does not apply to PersistentVolumeClaim volumes that are created and managed externally.
	DefaultVolumeProvisioningOptions *VirtualMachineVolumeProvisioningOptions `json:"defaultVolumeProvisioningOptions,omitempty"`

	// ChangeBlockTracking specifies the enablement of incremental backup support for this VirtualMachine, which can be utilized
	// by external backup systems such as VMware Data Recovery.
	ChangeBlockTracking *bool `json:"changeBlockTracking,omitempty"`
}

// VirtualMachineVolumeProvisioningOptions specifies the provisioning options for a VirtualMachineVolume.
type VirtualMachineVolumeProvisioningOptions struct {
	// ThinProvisioned specifies whether to use thin provisioning for the VirtualMachineVolume.
	// This means a sparse (allocate on demand) format with additional space optimizations.
	ThinProvisioned *bool `json:"thinProvisioned,omitempty"`

	// EagerZeroed specifies whether to use eager zero provisioning for the VirtualMachineVolume.
	// An eager zeroed thick disk has all space allocated and wiped clean of any previous contents
	// on the physical media at creation time. Such disks may take longer time during creation
	// compared to other disk formats.
	// EagerZeroed is only applicable if ThinProvisioned is false. This is validated by the webhook.
	EagerZeroed *bool `json:"eagerZeroed,omitempty"`
}

// VirtualMachineVolumeStatus defines the observed state of a VirtualMachineVolume instance.
type VirtualMachineVolumeStatus struct {
	// Name is the name of the volume in a VirtualMachine.
	Name string `json:"name"`

	// Attached represents whether a volume has been successfully attached to the VirtualMachine or not.
	Attached bool `json:"attached"`

	// DiskUuid represents the underlying virtual disk UUID and is present when attachment succeeds.
	DiskUuid string `json:"diskUUID"` //nolint:revive,stylecheck

	// Error represents the last error seen when attaching or detaching a volume.  Error will be empty if attachment succeeds.
	Error string `json:"error"`
}

// NetworkInterfaceStatus defines the observed state of network interfaces attached to the VirtualMachine
// as seen by the Guest OS and VMware tools.
type NetworkInterfaceStatus struct {
	// Connected represents whether the network interface is connected or not.
	Connected bool `json:"connected"`

	// MAC address of the network adapter
	MacAddress string `json:"macAddress,omitempty"`

	// IpAddresses represents zero, one or more IP addresses assigned to the network interface in CIDR notation.
	// For eg, "192.0.2.1/16".
	IpAddresses []string `json:"ipAddresses,omitempty"` //nolint:revive,stylecheck
}

// VirtualMachineStatus defines the observed state of a VirtualMachine instance.
type VirtualMachineStatus struct {
	// Host describes the hostname or IP address of the infrastructure host that the VirtualMachine is executing on.
	// +optional
	Host string `json:"host,omitempty"`

	// PowerState describes the current power state of the VirtualMachine.
	// +optional
	PowerState VirtualMachinePowerState `json:"powerState,omitempty"`

	// Phase describes the current phase information of the VirtualMachine.
	// +optional
	Phase VMStatusPhase `json:"phase,omitempty"`

	// Conditions describes the current condition information of the VirtualMachine.
	// +optional
	Conditions []Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// VmIp describes the Primary IP address assigned to the guest operating system, if known.
	// Multiple IPs can be available for the VirtualMachine. Refer to networkInterfaces in the VirtualMachine
	// status for additional IPs
	// +optional
	VmIp string `json:"vmIp,omitempty"` //nolint:revive,stylecheck

	// UniqueID describes a unique identifier that is provided by the underlying infrastructure provider, such as
	// vSphere.
	// +optional
	UniqueID string `json:"uniqueID,omitempty"`

	// BiosUUID describes a unique identifier provided by the underlying infrastructure provider that is exposed to the
	// Guest OS BIOS as a unique hardware identifier.
	// +optional
	BiosUUID string `json:"biosUUID,omitempty"`

	// InstanceUUID describes the unique instance UUID provided by the underlying infrastructure provider, such as vSphere.
	// +optional
	InstanceUUID string `json:"instanceUUID,omitempty"`

	// Volumes describes a list of current status information for each Volume that is desired to be attached to the
	// VirtualMachine.
	// +optional
	Volumes []VirtualMachineVolumeStatus `json:"volumes,omitempty"`

	// ChangeBlockTracking describes the CBT enablement status on the VirtualMachine.
	// +optional
	ChangeBlockTracking *bool `json:"changeBlockTracking,omitempty"`

	// NetworkInterfaces describes a list of current status information for each network interface that is desired to
	// be attached to the VirtualMachine.
	// +optional
	NetworkInterfaces []NetworkInterfaceStatus `json:"networkInterfaces,omitempty"`

	// Zone describes the availability zone where the VirtualMachine has been scheduled.
	// Please note this field may be empty when the cluster is not zone-aware.
	// +optional
	Zone string `json:"zone,omitempty"`

	// LastRestartTime describes the last time the VM was restarted.
	// +optional
	LastRestartTime *metav1.Time `json:"lastRestartTime,omitempty"`

	// HardwareVersion describes the VirtualMachine resource's observed
	// hardware version.
	//
	// Please refer to VirtualMachineSpec.MinHardwareVersion for more
	// information on the topic of a VM's hardware version.
	//
	// +optional
	HardwareVersion int32 `json:"hardwareVersion,omitempty"`
}

func (vm *VirtualMachine) GetConditions() Conditions {
	return vm.Status.Conditions
}

func (vm *VirtualMachine) SetConditions(conditions Conditions) {
	vm.Status.Conditions = conditions
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=vm
// +kubebuilder:storageversion:false
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Power-State",type="string",JSONPath=".status.powerState"
// +kubebuilder:printcolumn:name="Class",type="string",priority=1,JSONPath=".spec.className"
// +kubebuilder:printcolumn:name="Image",type="string",priority=1,JSONPath=".spec.imageName"
// +kubebuilder:printcolumn:name="Primary-IP",type="string",priority=1,JSONPath=".status.vmIp"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// VirtualMachine is the Schema for the virtualmachines API.
// A VirtualMachine represents the desired specification and the observed status of a VirtualMachine instance.  A
// VirtualMachine is realized by the VirtualMachine controller on a backing Virtual Infrastructure provider such as
// vSphere.
type VirtualMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualMachineSpec   `json:"spec,omitempty"`
	Status VirtualMachineStatus `json:"status,omitempty"`
}

func (vm VirtualMachine) NamespacedName() string {
	return vm.Namespace + "/" + vm.Name
}

// VirtualMachineList contains a list of VirtualMachine.
//
// +kubebuilder:object:root=true
type VirtualMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualMachine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VirtualMachine{}, &VirtualMachineList{})
}
