/*
Copyright 2025 The Kubernetes Authors.

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

package v1beta2

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/resource"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
)

const (
	// AnnotationClusterInfrastructureReady indicates the cluster's
	// infrastructure sources are ready and machines may be created.
	AnnotationClusterInfrastructureReady = "vsphere.infrastructure.cluster.x-k8s.io/infrastructure-ready"

	// AnnotationControlPlaneReady indicates the cluster's control plane is
	// ready.
	AnnotationControlPlaneReady = "vsphere.infrastructure.cluster.x-k8s.io/control-plane-ready"

	// ValueReady is the ready value for *Ready annotations.
	ValueReady = "true"
)

// CloneMode is the type of clone operation used to clone a VM from a template.
// +kubebuilder:validation:Enum=fullClone;linkedClone
type CloneMode string

const (
	// FullClone indicates a VM will have no relationship to the source of the
	// clone operation once the operation is complete. This is the safest clone
	// mode, but it is not the fastest.
	FullClone CloneMode = "fullClone"

	// LinkedClone means resulting VMs will be dependent upon the snapshot of
	// the source VM/template from which the VM was cloned. This is the fastest
	// clone mode, but it also prevents expanding a VMs disk beyond the size of
	// the source VM/template.
	LinkedClone CloneMode = "linkedClone"
)

// FtEncryptionMode represents the encrypted fault tolerance mode.
// +kubebuilder:validation:Enum=ftEncryptionDisabled;ftEncryptionOpportunistic;ftEncryptionRequired
type FtEncryptionMode string

const (
	// FtEncryptionDisabled indicates a VM will not use encrypted Fault Tolerance,
	// even if available.
	FtEncryptionDisabled FtEncryptionMode = "ftEncryptionDisabled"

	// FtEncryptionOpportunistic indicates a VM will use encrypted Fault Tolerance
	// if source and destination hosts support it, fall back to unencrypted Fault Tolerance otherwise.
	FtEncryptionOpportunistic FtEncryptionMode = "ftEncryptionOpportunistic"

	// FtEncryptionRequired indicates a VM will allow only encrypted Fault Tolerance.
	// If either the source or destination host does not support encrypted Fault Tolerance,
	// do not allow the Fault Tolerance to occur.
	FtEncryptionRequired FtEncryptionMode = "ftEncryptionRequired"
)

// MigrateEncryption represents the encrypted vMotion mode.
// +kubebuilder:validation:Enum=disabled;opportunistic;required
type MigrateEncryption string

const (
	// DisabledMigrateEncryption indicates a VM will not use encrypted vMotion, even if available.
	DisabledMigrateEncryption MigrateEncryption = "disabled"

	// OpportunisticMigrateEncryption indicates a VM Use encrypted vMotion
	// if source and destination hosts support it, fall back to unencrypted vMotion otherwise.
	OpportunisticMigrateEncryption MigrateEncryption = "opportunistic"

	// RequiredMigrateEncryption indicates a VM will Allow only encrypted vMotion.
	// If the source or destination host does not support vMotion encryption,
	// do not allow the vMotion to occur.
	RequiredMigrateEncryption MigrateEncryption = "required"
)

// OS is the type of Operating System the virtual machine uses.
type OS string

const (
	// Linux indicates the VM uses a Linux Operating System.
	Linux OS = "Linux"

	// Windows indicates the VM uses Windows Server 2019 as the OS.
	Windows OS = "Windows"
)

// VirtualMachinePowerOpMode represents the various power operation modes
// when powering off or suspending a VM.
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
	// state within the configured timeout (default 5m).
	VirtualMachinePowerOpModeTrySoft VirtualMachinePowerOpMode = "trySoft"
)

// VirtualMachineCloneSpec is information used to clone a virtual machine.
type VirtualMachineCloneSpec struct {
	// template is the name, inventory path, managed object reference or the managed
	// object ID of the template used to clone the virtual machine.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	Template string `json:"template,omitempty"`

	// cloneMode specifies the type of clone operation.
	// The linkedClone mode is only support for templates that have at least
	// one snapshot. If the template has no snapshots, then CloneMode defaults
	// to fullClone.
	// When linkedClone mode is enabled the DiskGiB field is ignored as it is
	// not possible to expand disks of linked clones.
	// Defaults to linkedClone, but fails gracefully to fullClone if the source
	// of the clone operation has no snapshots.
	// +optional
	CloneMode CloneMode `json:"cloneMode,omitempty"`

	// snapshot is the name of the snapshot from which to create a linked clone.
	// This field is ignored if linkedClone is not enabled.
	// Defaults to the source's current snapshot.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	Snapshot string `json:"snapshot,omitempty"`

	// server is the IP address or FQDN of the vSphere server on which
	// the virtual machine is created/located.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	Server string `json:"server,omitempty"`

	// thumbprint is the colon-separated SHA-1 checksum of the given vCenter server's host certificate
	// When this is set to empty, this VirtualMachine would be created
	// without TLS certificate validation of the communication between Cluster API Provider vSphere
	// and the VMware vCenter server.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	Thumbprint string `json:"thumbprint,omitempty"`

	// datacenter is the name, inventory path, managed object reference or the managed
	// object ID of the datacenter in which the virtual machine is created/located.
	// Defaults to * which selects the default datacenter.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	Datacenter string `json:"datacenter,omitempty"`

	// folder is the name, inventory path, managed object reference or the managed
	// object ID of the folder in which the virtual machine is created/located.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	Folder string `json:"folder,omitempty"`

	// datastore is the name, inventory path, managed object reference or the managed
	// object ID of the datastore in which the virtual machine is created/located.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	Datastore string `json:"datastore,omitempty"`

	// storagePolicyName of the storage policy to use with this
	// Virtual Machine
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	StoragePolicyName string `json:"storagePolicyName,omitempty"`

	// resourcePool is the name, inventory path, managed object reference or the managed
	// object ID in which the virtual machine is created/located.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	ResourcePool string `json:"resourcePool,omitempty"`

	// network is the network configuration for this machine's VM.
	// +required
	Network NetworkSpec `json:"network,omitzero"`

	// numCPUs is the number of virtual processors in a virtual machine.
	// Defaults to the eponymous property value in the template from which the
	// virtual machine is cloned.
	// +optional
	// +kubebuilder:validation:Minimum=2
	NumCPUs int32 `json:"numCPUs,omitempty"`

	// numCoresPerSocket is the number of cores among which to distribute CPUs in this
	// virtual machine.
	// Defaults to the eponymous property value in the template from which the
	// virtual machine is cloned.
	// Note: Starting with vSphere 8 numCoresPerSocket can be set to 0 to enable "Assigned at power on".
	// +optional
	// +kubebuilder:validation:Minimum=0
	NumCoresPerSocket *int32 `json:"numCoresPerSocket,omitempty"`

	// resources is the definition of the VM's cpu and memory
	// reservations, limits and shares.
	// +optional
	Resources VirtualMachineResources `json:"resources,omitempty,omitzero"`

	// memoryMiB is the size of a virtual machine's memory, in MiB.
	// Defaults to the eponymous property value in the template from which the
	// virtual machine is cloned.
	// +optional
	// +kubebuilder:validation:Minimum=1
	MemoryMiB int64 `json:"memoryMiB,omitempty"`

	// diskGiB is the size of a virtual machine's disk, in GiB.
	// Defaults to the eponymous property value in the template from which the
	// virtual machine is cloned.
	// +optional
	// +kubebuilder:validation:Minimum=1
	DiskGiB int32 `json:"diskGiB,omitempty"`

	// additionalDisksGiB holds the sizes of additional disks of the virtual machine, in GiB
	// Defaults to the eponymous property value in the template from which the
	// virtual machine is cloned.
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=128
	AdditionalDisksGiB []int32 `json:"additionalDisksGiB,omitempty"`

	// customVMXKeys is a dictionary of advanced VMX options that can be set on VM
	// Defaults to empty map
	// +optional
	CustomVMXKeys map[string]string `json:"customVMXKeys,omitempty"`

	// tagIDs is an optional set of tags to add to an instance. Specified tagIDs
	// must use URN-notation instead of display names.
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=128
	// +kubebuilder:validation:items:MinLength=1
	// +kubebuilder:validation:items:MaxLength=1024
	TagIDs []string `json:"tagIDs,omitempty"`

	// pciDevices is the list of pci devices used by the virtual machine.
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=128
	PciDevices []PCIDeviceSpec `json:"pciDevices,omitempty"`

	// os is the Operating System of the virtual machine
	// Defaults to Linux
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=128
	OS OS `json:"os,omitempty"`

	// hardwareVersion is the hardware version of the virtual machine.
	// Defaults to the eponymous property value in the template from which the
	// virtual machine is cloned.
	// Check the compatibility with the ESXi version before setting the value.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=128
	HardwareVersion string `json:"hardwareVersion,omitempty"`

	// dataDisks are additional disks to add to the VM that are not part of the VM's OVA template.
	// +optional
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MaxItems=29
	DataDisks []VSphereDisk `json:"dataDisks,omitempty"`

	// nestedHV controls nested hardware-assisted virtualization.
	// Defaults to the eponymous property value in the template from which the
	// virtual machine is cloned.
	// Check the compatibility with the ESXi version before setting the value.
	// +optional
	NestedHV *bool `json:"nestedHV,omitempty"`

	// ftEncryptionMode is the encrypted fault tolerance mode.
	// Defaults to the eponymous property value in the template from which the
	// virtual machine is cloned.
	// Check the compatibility with the ESXi version before setting the value.
	// +optional
	FtEncryptionMode FtEncryptionMode `json:"ftEncryptionMode,omitempty"`

	// migrateEncryption is the encrypted vMotion mode.
	// Defaults to the eponymous property value in the template from which the
	// virtual machine is cloned.
	// Check the compatibility with the ESXi version before setting the value.
	// +optional
	MigrateEncryption MigrateEncryption `json:"migrateEncryption,omitempty"`

	// cryptoKeyID is the crypto key id.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=128
	CryptoKeyID string `json:"cryptoKeyID,omitempty"`

	// cryptoProfile of the storage encryption policy to use with this
	// Virtual Machine.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=128
	CryptoProfile string `json:"cryptoProfile,omitempty"`
}

// VirtualMachineResources is the definition of the VM's cpu and memory
// reservations, limits and shares.
// +kubebuilder:validation:MinProperties=1
type VirtualMachineResources struct {
	// requests is the definition of the VM's cpu (in hertz, rounded up to the nearest MHz)
	// and memory (in bytes, rounded up to the nearest MiB) reservations
	// +optional
	Requests VirtualMachineResourceSpec `json:"requests,omitempty,omitzero"`

	// limits is the definition of the VM's cpu (in hertz, rounded up to the nearest MHz)
	// and memory (in bytes, rounded up to the nearest MiB) limits
	// +optional
	Limits VirtualMachineResourceSpec `json:"limits,omitempty,omitzero"`

	// shares is the definition of the VM's cpu and memory shares
	// +optional
	Shares VirtualMachineResourceShares `json:"shares,omitempty,omitzero"`
}

// VirtualMachineResourceSpec is the numerical definition of memory and cpu quantity for the
// given VM hardware policy.
// +kubebuilder:validation:MinProperties=1
type VirtualMachineResourceSpec struct {
	// cpu is the definition of the cpu quantity for the given VM hardware policy
	// +optional
	CPU resource.Quantity `json:"cpu,omitempty"`

	// memory is the definition of the memory quantity for the given VM hardware policy
	// +optional
	Memory resource.Quantity `json:"memory,omitempty"`
}

// VirtualMachineResourceShares is the numerical definition of memory and cpu shares for the
// given VM
// +kubebuilder:validation:MinProperties=1
type VirtualMachineResourceShares struct {
	// cpu is the number of spu shares to assign to the VM
	// +kubebuilder:validation:Minimum=1
	// +optional
	CPU int32 `json:"cpu,omitempty"`

	// memory is the number of memory shares to assign to the VM
	// +kubebuilder:validation:Minimum=1
	// +optional
	Memory int32 `json:"memory,omitempty"`
}

// VSphereDisk is an additional disk to add to the VM that is not part of the VM OVA template.
type VSphereDisk struct {
	// name is used to identify the disk definition. Name is required and needs to be unique so that it can be used to
	// clearly identify purpose of the disk.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	Name string `json:"name,omitempty"`

	// sizeGiB is the size of the disk in GiB.
	// +required
	// +kubebuilder:validation:Minimum=1
	SizeGiB int32 `json:"sizeGiB,omitempty"`

	// provisioningMode specifies the provisioning type to be used by this vSphere data disk.
	// If not set, the setting will be provided by the default storage policy.
	// +optional
	ProvisioningMode ProvisioningMode `json:"provisioningMode,omitempty"`
}

// ProvisioningMode represents the various provisioning types available to a VMs disk.
// +kubebuilder:validation:Enum=Thin;Thick;EagerlyZeroed
type ProvisioningMode string

var (
	// ThinProvisioningMode creates the disk using thin provisioning. This means a sparse (allocate on demand)
	// format with additional space optimizations.
	ThinProvisioningMode ProvisioningMode = "Thin"

	// ThickProvisioningMode creates the disk with all space allocated.
	ThickProvisioningMode ProvisioningMode = "Thick"

	// EagerlyZeroedProvisioningMode creates the disk using eager zero provisioning. An eager zeroed thick disk
	// has all space allocated and wiped clean of any previous contents on the physical media at
	// creation time. Such disks may take longer time during creation compared to other disk formats.
	EagerlyZeroedProvisioningMode ProvisioningMode = "EagerlyZeroed"
)

// VSphereMachineTemplateResource describes the data needed to create a VSphereMachine from a template.
// +kubebuilder:validation:MinProperties=1
type VSphereMachineTemplateResource struct {
	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	ObjectMeta clusterv1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec is the specification of the desired behavior of the machine.
	// +optional
	Spec VSphereMachineSpec `json:"spec,omitzero"`
}

// APIEndpoint represents a reachable Kubernetes API endpoint.
// +kubebuilder:validation:MinProperties=1
type APIEndpoint struct {
	// host is the hostname on which the API server is serving.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=512
	Host string `json:"host,omitempty"`

	// port is the port on which the API server is serving.
	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port int32 `json:"port,omitempty"`
}

// String returns a formatted version HOST:PORT of this APIEndpoint.
func (v APIEndpoint) String() string {
	return fmt.Sprintf("%s:%d", v.Host, v.Port)
}

// PCIDeviceSpec defines virtual machine's PCI configuration.
type PCIDeviceSpec struct {
	// deviceId is the device ID of a virtual machine's PCI, in integer.
	// Defaults to the eponymous property value in the template from which the
	// virtual machine is cloned.
	// Mutually exclusive with VGPUProfile as VGPUProfile and DeviceID + VendorID
	// are two independent ways to define PCI devices.
	// +optional
	DeviceID *int32 `json:"deviceId,omitempty"`

	// vendorId is the vendor ID of a virtual machine's PCI, in integer.
	// Defaults to the eponymous property value in the template from which the
	// virtual machine is cloned.
	// Mutually exclusive with VGPUProfile as VGPUProfile and DeviceID + VendorID
	// are two independent ways to define PCI devices.
	// +optional
	VendorID *int32 `json:"vendorId,omitempty"`

	// vGPUProfile is the profile name of a virtual machine's vGPU, in string.
	// Defaults to the eponymous property value in the template from which the
	// virtual machine is cloned.
	// Mutually exclusive with DeviceID and VendorID as VGPUProfile and DeviceID + VendorID
	// are two independent ways to define PCI devices.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	VGPUProfile string `json:"vGPUProfile,omitempty"`

	// customLabel is the hardware label of a virtual machine's PCI device.
	// Defaults to the eponymous property value in the template from which the
	// virtual machine is cloned.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	CustomLabel string `json:"customLabel,omitempty"`
}

// NetworkSpec defines the virtual machine's network configuration.
type NetworkSpec struct {
	// devices is the list of network devices used by the virtual machine.
	//
	// TODO(akutz) Make sure at least one network matches the ClusterSpec.CloudProviderConfiguration.Network.Name
	// +required
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=128
	Devices []NetworkDeviceSpec `json:"devices,omitempty"`

	// routes is a list of optional, static routes applied to the virtual
	// machine.
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=512
	Routes []NetworkRouteSpec `json:"routes,omitempty"`
}

// NetworkDeviceSpec defines the network configuration for a virtual machine's
// network device.
type NetworkDeviceSpec struct {
	// networkName is the name, managed object reference or the managed
	// object ID of the vSphere network to which the device will be connected.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	NetworkName string `json:"networkName,omitempty"`

	// deviceName may be used to explicitly assign a name to the network device
	// as it exists in the guest operating system.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	DeviceName string `json:"deviceName,omitempty"`

	// dhcp4 is a flag that indicates whether or not to use DHCP for IPv4
	// on this device.
	// If true then IPAddrs should not contain any IPv4 addresses.
	// +optional
	DHCP4 *bool `json:"dhcp4,omitempty"`

	// dhcp6 is a flag that indicates whether or not to use DHCP for IPv6
	// on this device.
	// If true then IPAddrs should not contain any IPv6 addresses.
	// +optional
	DHCP6 *bool `json:"dhcp6,omitempty"`

	// gateway4 is the IPv4 gateway used by this device.
	// Required when DHCP4 is false.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=64
	Gateway4 string `json:"gateway4,omitempty"`

	// gateway6 is the IPv6 gateway used by this device.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=64
	Gateway6 string `json:"gateway6,omitempty"`

	// ipAddrs is a list of one or more IPv4 and/or IPv6 addresses to assign
	// to this device. IP addresses must also specify the segment length in
	// CIDR notation.
	// Required when DHCP4, DHCP6 and SkipIPAllocation are false.
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=128
	// +kubebuilder:validation:items:MinLength=1
	// +kubebuilder:validation:items:MaxLength=39
	IPAddrs []string `json:"ipAddrs,omitempty"`

	// mtu is the device’s Maximum Transmission Unit size in bytes.
	// +optional
	MTU *int64 `json:"mtu,omitempty"`

	// macAddr is the MAC address used by this device.
	// It is generally a good idea to omit this field and allow a MAC address
	// to be generated.
	// Please note that this value must use the VMware OUI to work with the
	// in-tree vSphere cloud provider.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=23
	MACAddr string `json:"macAddr,omitempty"`

	// nameservers is a list of IPv4 and/or IPv6 addresses used as DNS
	// nameservers.
	// Please note that Linux allows only three nameservers (https://linux.die.net/man/5/resolv.conf).
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=128
	// +kubebuilder:validation:items:MinLength=1
	// +kubebuilder:validation:items:MaxLength=64
	Nameservers []string `json:"nameservers,omitempty"`

	// routes is a list of optional, static routes applied to the device.
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=512
	Routes []NetworkRouteSpec `json:"routes,omitempty"`

	// searchDomains is a list of search domains used when resolving IP
	// addresses with DNS.
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=128
	// +kubebuilder:validation:items:MinLength=1
	// +kubebuilder:validation:items:MaxLength=1024
	SearchDomains []string `json:"searchDomains,omitempty"`

	// addressesFromPools is a list of IPAddressPools that should be assigned
	// to IPAddressClaims. The machine's cloud-init metadata will be populated
	// with IPAddresses fulfilled by an IPAM provider.
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=128
	AddressesFromPools []IPPoolReference `json:"addressesFromPools,omitempty"`

	// dhcp4Overrides allows for the control over several DHCP behaviors.
	// Overrides will only be applied when the corresponding DHCP flag is set.
	// Only configured values will be sent, omitted values will default to
	// distribution defaults.
	// Dependent on support in the network stack for your distribution.
	// For more information see the netplan reference (https://netplan.io/reference#dhcp-overrides)
	// +optional
	DHCP4Overrides *DHCPOverrides `json:"dhcp4Overrides,omitempty"`

	// dhcp6Overrides allows for the control over several DHCP behaviors.
	// Overrides will only be applied when the corresponding DHCP flag is set.
	// Only configured values will be sent, omitted values will default to
	// distribution defaults.
	// Dependent on support in the network stack for your distribution.
	// For more information see the netplan reference (https://netplan.io/reference#dhcp-overrides)
	// +optional
	DHCP6Overrides *DHCPOverrides `json:"dhcp6Overrides,omitempty"`

	// skipIPAllocation allows the device to not have IP address or DHCP configured.
	// This is suitable for devices for which IP allocation is handled externally, eg. using Multus CNI.
	// If true, CAPV will not verify IP address allocation.
	// +optional
	SkipIPAllocation *bool `json:"skipIPAllocation,omitempty"`
}

// IPPoolReference is a reference to an IPPool.
type IPPoolReference struct {
	// name of the IPPool.
	// name must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`
	Name string `json:"name,omitempty"`

	// kind of the IPPool.
	// kind must consist of alphanumeric characters or '-', start with an alphabetic character, and end with an alphanumeric character.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern=`^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$`
	Kind string `json:"kind,omitempty"`

	// apiGroup of the IPPool.
	// apiGroup must be fully qualified domain name.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`
	APIGroup string `json:"apiGroup,omitempty"`
}

// DHCPOverrides allows for the control over several DHCP behaviors.
// Overrides will only be applied when the corresponding DHCP flag is set.
// Only configured values will be sent, omitted values will default to
// distribution defaults.
// Dependent on support in the network stack for your distribution.
// For more information see the netplan reference (https://netplan.io/reference#dhcp-overrides)
type DHCPOverrides struct {
	// hostname is the name which will be sent to the DHCP server instead of
	// the machine's hostname.
	// +optional
	// +kubebuilder:validation:MaxLength=1024
	Hostname *string `json:"hostname,omitempty"`

	// routeMetric is used to prioritize routes for devices. A lower metric for
	// an interface will have a higher priority.
	// +optional
	// +kubebuilder:validation:Minimum=0
	RouteMetric *int32 `json:"routeMetric,omitempty"`

	// sendHostname when `true`, the hostname of the machine will be sent to the
	// DHCP server.
	// +optional
	SendHostname *bool `json:"sendHostname,omitempty"`

	// useDNS when `true`, the DNS servers in the DHCP server will be used and
	// take precedence.
	// +optional
	UseDNS *bool `json:"useDNS,omitempty"`

	// useDomains can take the values `true`, `false`, or `route`. When `true`,
	// the domain name from the DHCP server will be used as the DNS search
	// domain for this device. When `route`, the domain name from the DHCP
	// response will be used for routing DNS only, not for searching.
	// +optional
	// +kubebuilder:validation:MaxLength=128
	UseDomains *string `json:"useDomains,omitempty"`

	// useHostname when `true`, the hostname from the DHCP server will be set
	// as the transient hostname of the machine.
	// +optional
	UseHostname *bool `json:"useHostname,omitempty"`

	// useMTU when `true`, the MTU from the DHCP server will be set as the
	// MTU of the device.
	// +optional
	UseMTU *bool `json:"useMTU,omitempty"`

	// useNTP when `true`, the NTP servers from the DHCP server will be used
	// by systemd-timesyncd and take precedence.
	// +optional
	UseNTP *bool `json:"useNTP,omitempty"`

	// useRoutes when `true`, the routes from the DHCP server will be installed
	// in the routing table.
	// +optional
	// +kubebuilder:validation:MaxLength=128
	UseRoutes *string `json:"useRoutes,omitempty"`
}

// NetworkRouteSpec defines a static network route.
type NetworkRouteSpec struct {
	// to is an IPv4 or IPv6 address.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=39
	To string `json:"to,omitempty"`

	// via is an IPv4 or IPv6 address.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=39
	Via string `json:"via,omitempty"`

	// metric is the weight/priority of the route.
	// +required
	// +kubebuilder:validation:Minimum=0
	Metric *int32 `json:"metric,omitempty"`
}

// NetworkStatus provides information about one of a VM's networks.
type NetworkStatus struct {
	// connected is a flag that indicates whether this network is currently
	// connected to the VM.
	// +optional
	Connected *bool `json:"connected,omitempty"`

	// ipAddrs is one or more IP addresses reported by vm-tools.
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=128
	// +kubebuilder:validation:items:MinLength=1
	// +kubebuilder:validation:items:MaxLength=39
	IPAddrs []string `json:"ipAddrs,omitempty"`

	// macAddr is the MAC address of the network device.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=17
	MACAddr string `json:"macAddr,omitempty"`

	// networkName is the name of the network.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	NetworkName string `json:"networkName,omitempty"`
}

// VirtualMachinePowerState describe the power state of a VM.
type VirtualMachinePowerState string

const (
	// VirtualMachinePowerStatePoweredOn is the string representing a VM in powered on state.
	VirtualMachinePowerStatePoweredOn VirtualMachinePowerState = "poweredOn"

	// VirtualMachinePowerStatePoweredOff is the string representing a VM in powered off state.
	VirtualMachinePowerStatePoweredOff = "poweredOff"

	// VirtualMachinePowerStateSuspended is the string representing a VM in suspended state.
	VirtualMachinePowerStateSuspended = "suspended"
)

// SSHUser is granted remote access to a system.
type SSHUser struct {
	// name is the name of the SSH user.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	Name string `json:"name,omitempty"`

	// authorizedKeys is one or more public SSH keys that grant remote access.
	// +required
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=128
	// +kubebuilder:validation:items:MinLength=1
	// +kubebuilder:validation:items:MaxLength=10240
	AuthorizedKeys []string `json:"authorizedKeys,omitempty"`
}
