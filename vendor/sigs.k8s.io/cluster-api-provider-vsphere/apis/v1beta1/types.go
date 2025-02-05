/*
Copyright 2021 The Kubernetes Authors.

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

package v1beta1

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
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
	// Template is the name, inventory path, managed object reference or the managed
	// object ID of the template used to clone the virtual machine.
	// +kubebuilder:validation:MinLength=1
	Template string `json:"template"`

	// CloneMode specifies the type of clone operation.
	// The LinkedClone mode is only support for templates that have at least
	// one snapshot. If the template has no snapshots, then CloneMode defaults
	// to FullClone.
	// When LinkedClone mode is enabled the DiskGiB field is ignored as it is
	// not possible to expand disks of linked clones.
	// Defaults to LinkedClone, but fails gracefully to FullClone if the source
	// of the clone operation has no snapshots.
	// +optional
	CloneMode CloneMode `json:"cloneMode,omitempty"`

	// Snapshot is the name of the snapshot from which to create a linked clone.
	// This field is ignored if LinkedClone is not enabled.
	// Defaults to the source's current snapshot.
	// +optional
	Snapshot string `json:"snapshot,omitempty"`

	// Server is the IP address or FQDN of the vSphere server on which
	// the virtual machine is created/located.
	// +optional
	Server string `json:"server,omitempty"`

	// Thumbprint is the colon-separated SHA-1 checksum of the given vCenter server's host certificate
	// When this is set to empty, this VirtualMachine would be created
	// without TLS certificate validation of the communication between Cluster API Provider vSphere
	// and the VMware vCenter server.
	// +optional
	Thumbprint string `json:"thumbprint,omitempty"`

	// Datacenter is the name, inventory path, managed object reference or the managed
	// object ID of the datacenter in which the virtual machine is created/located.
	// Defaults to * which selects the default datacenter.
	// +optional
	Datacenter string `json:"datacenter,omitempty"`

	// Folder is the name, inventory path, managed object reference or the managed
	// object ID of the folder in which the virtual machine is created/located.
	// +optional
	Folder string `json:"folder,omitempty"`

	// Datastore is the name, inventory path, managed object reference or the managed
	// object ID of the datastore in which the virtual machine is created/located.
	// +optional
	Datastore string `json:"datastore,omitempty"`

	// StoragePolicyName of the storage policy to use with this
	// Virtual Machine
	// +optional
	StoragePolicyName string `json:"storagePolicyName,omitempty"`

	// ResourcePool is the name, inventory path, managed object reference or the managed
	// object ID in which the virtual machine is created/located.
	// +optional
	ResourcePool string `json:"resourcePool,omitempty"`

	// Network is the network configuration for this machine's VM.
	Network NetworkSpec `json:"network"`

	// NumCPUs is the number of virtual processors in a virtual machine.
	// Defaults to the eponymous property value in the template from which the
	// virtual machine is cloned.
	// +optional
	NumCPUs int32 `json:"numCPUs,omitempty"`
	// NumCPUs is the number of cores among which to distribute CPUs in this
	// virtual machine.
	// Defaults to the eponymous property value in the template from which the
	// virtual machine is cloned.
	// +optional
	NumCoresPerSocket int32 `json:"numCoresPerSocket,omitempty"`
	// MemoryMiB is the size of a virtual machine's memory, in MiB.
	// Defaults to the eponymous property value in the template from which the
	// virtual machine is cloned.
	// +optional
	MemoryMiB int64 `json:"memoryMiB,omitempty"`
	// DiskGiB is the size of a virtual machine's disk, in GiB.
	// Defaults to the eponymous property value in the template from which the
	// virtual machine is cloned.
	// +optional
	DiskGiB int32 `json:"diskGiB,omitempty"`
	// AdditionalDisksGiB holds the sizes of additional disks of the virtual machine, in GiB
	// Defaults to the eponymous property value in the template from which the
	// virtual machine is cloned.
	// +optional
	AdditionalDisksGiB []int32 `json:"additionalDisksGiB,omitempty"`
	// CustomVMXKeys is a dictionary of advanced VMX options that can be set on VM
	// Defaults to empty map
	// +optional
	CustomVMXKeys map[string]string `json:"customVMXKeys,omitempty"`
	// TagIDs is an optional set of tags to add to an instance. Specified tagIDs
	// must use URN-notation instead of display names.
	// +optional
	TagIDs []string `json:"tagIDs,omitempty"`
	// PciDevices is the list of pci devices used by the virtual machine.
	// +optional
	PciDevices []PCIDeviceSpec `json:"pciDevices,omitempty"`
	// OS is the Operating System of the virtual machine
	// Defaults to Linux
	// +optional
	OS OS `json:"os,omitempty"`
	// HardwareVersion is the hardware version of the virtual machine.
	// Defaults to the eponymous property value in the template from which the
	// virtual machine is cloned.
	// Check the compatibility with the ESXi version before setting the value.
	// +optional
	HardwareVersion string `json:"hardwareVersion,omitempty"`
	// DataDisks are additional disks to add to the VM that are not part of the VM's OVA template.
	// +optional
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MaxItems=29
	DataDisks []VSphereDisk `json:"dataDisks,omitempty"`
}

// VSphereDisk is an additional disk to add to the VM that is not part of the VM OVA template.
type VSphereDisk struct {
	// Name is used to identify the disk definition. Name is required and needs to be unique so that it can be used to
	// clearly identify purpose of the disk.
	// +kubebuilder:validation:Required
	Name string `json:"name,omitempty"`
	// SizeGiB is the size of the disk in GiB.
	// +kubebuilder:validation:Required
	SizeGiB int32 `json:"sizeGiB"`
	// ProvisioningMode specifies the provisioning type to be used by this vSphere data disk.
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
type VSphereMachineTemplateResource struct {

	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	ObjectMeta clusterv1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the specification of the desired behavior of the machine.
	Spec VSphereMachineSpec `json:"spec"`
}

// VSphereMachineProviderConditionType is a valid value for VSphereMachineProviderCondition.Type.
type VSphereMachineProviderConditionType string

// Valid conditions for an VSphere machine instance.
const (
	// MachineCreated indicates whether the machine has been created or not. If not,
	// it should include a reason and message for the failure.
	MachineCreated VSphereMachineProviderConditionType = "MachineCreated"
)

// APIEndpoint represents a reachable Kubernetes API endpoint.
type APIEndpoint struct {
	// The hostname on which the API server is serving.
	Host string `json:"host"`

	// The port on which the API server is serving.
	Port int32 `json:"port"`
}

// IsZero returns true if either the host or the port are zero values.
func (v APIEndpoint) IsZero() bool {
	return v.Host == "" || v.Port == 0
}

// String returns a formatted version HOST:PORT of this APIEndpoint.
func (v APIEndpoint) String() string {
	return fmt.Sprintf("%s:%d", v.Host, v.Port)
}

// PCIDeviceSpec defines virtual machine's PCI configuration.
type PCIDeviceSpec struct {
	// DeviceID is the device ID of a virtual machine's PCI, in integer.
	// Defaults to the eponymous property value in the template from which the
	// virtual machine is cloned.
	// Mutually exclusive with VGPUProfile as VGPUProfile and DeviceID + VendorID
	// are two independent ways to define PCI devices.
	// +optional
	DeviceID *int32 `json:"deviceId,omitempty"`
	// VendorId is the vendor ID of a virtual machine's PCI, in integer.
	// Defaults to the eponymous property value in the template from which the
	// virtual machine is cloned.
	// Mutually exclusive with VGPUProfile as VGPUProfile and DeviceID + VendorID
	// are two independent ways to define PCI devices.
	// +optional
	VendorID *int32 `json:"vendorId,omitempty"`
	// VGPUProfile is the profile name of a virtual machine's vGPU, in string.
	// Defaults to the eponymous property value in the template from which the
	// virtual machine is cloned.
	// Mutually exclusive with DeviceID and VendorID as VGPUProfile and DeviceID + VendorID
	// are two independent ways to define PCI devices.
	// +optional
	VGPUProfile string `json:"vGPUProfile,omitempty"`
	// CustomLabel is the hardware label of a virtual machine's PCI device.
	// Defaults to the eponymous property value in the template from which the
	// virtual machine is cloned.
	// +optional
	CustomLabel string `json:"customLabel,omitempty"`
}

// NetworkSpec defines the virtual machine's network configuration.
type NetworkSpec struct {
	// Devices is the list of network devices used by the virtual machine.
	//
	// TODO(akutz) Make sure at least one network matches the ClusterSpec.CloudProviderConfiguration.Network.Name
	Devices []NetworkDeviceSpec `json:"devices"`

	// Routes is a list of optional, static routes applied to the virtual
	// machine.
	// +optional
	Routes []NetworkRouteSpec `json:"routes,omitempty"`

	// PreferredAPIServeCIDR is the preferred CIDR for the Kubernetes API
	// server endpoint on this machine
	// +optional
	//
	// Deprecated: This field is going to be removed in a future release.
	PreferredAPIServerCIDR string `json:"preferredAPIServerCidr,omitempty"`
}

// NetworkDeviceSpec defines the network configuration for a virtual machine's
// network device.
type NetworkDeviceSpec struct {
	// NetworkName is the name, managed object reference or the managed
	// object ID of the vSphere network to which the device will be connected.
	NetworkName string `json:"networkName"`

	// DeviceName may be used to explicitly assign a name to the network device
	// as it exists in the guest operating system.
	// +optional
	DeviceName string `json:"deviceName,omitempty"`

	// DHCP4 is a flag that indicates whether or not to use DHCP for IPv4
	// on this device.
	// If true then IPAddrs should not contain any IPv4 addresses.
	// +optional
	DHCP4 bool `json:"dhcp4,omitempty"`

	// DHCP6 is a flag that indicates whether or not to use DHCP for IPv6
	// on this device.
	// If true then IPAddrs should not contain any IPv6 addresses.
	// +optional
	DHCP6 bool `json:"dhcp6,omitempty"`

	// Gateway4 is the IPv4 gateway used by this device.
	// Required when DHCP4 is false.
	// +optional
	Gateway4 string `json:"gateway4,omitempty"`

	// Gateway4 is the IPv4 gateway used by this device.
	// +optional
	Gateway6 string `json:"gateway6,omitempty"`

	// IPAddrs is a list of one or more IPv4 and/or IPv6 addresses to assign
	// to this device. IP addresses must also specify the segment length in
	// CIDR notation.
	// Required when DHCP4, DHCP6 and SkipIPAllocation are false.
	// +optional
	IPAddrs []string `json:"ipAddrs,omitempty"`

	// MTU is the device’s Maximum Transmission Unit size in bytes.
	// +optional
	MTU *int64 `json:"mtu,omitempty"`

	// MACAddr is the MAC address used by this device.
	// It is generally a good idea to omit this field and allow a MAC address
	// to be generated.
	// Please note that this value must use the VMware OUI to work with the
	// in-tree vSphere cloud provider.
	// +optional
	MACAddr string `json:"macAddr,omitempty"`

	// Nameservers is a list of IPv4 and/or IPv6 addresses used as DNS
	// nameservers.
	// Please note that Linux allows only three nameservers (https://linux.die.net/man/5/resolv.conf).
	// +optional
	Nameservers []string `json:"nameservers,omitempty"`

	// Routes is a list of optional, static routes applied to the device.
	// +optional
	Routes []NetworkRouteSpec `json:"routes,omitempty"`

	// SearchDomains is a list of search domains used when resolving IP
	// addresses with DNS.
	// +optional
	SearchDomains []string `json:"searchDomains,omitempty"`

	// AddressesFromPools is a list of IPAddressPools that should be assigned
	// to IPAddressClaims. The machine's cloud-init metadata will be populated
	// with IPAddresses fulfilled by an IPAM provider.
	// +optional
	AddressesFromPools []corev1.TypedLocalObjectReference `json:"addressesFromPools,omitempty"`

	// DHCP4Overrides allows for the control over several DHCP behaviors.
	// Overrides will only be applied when the corresponding DHCP flag is set.
	// Only configured values will be sent, omitted values will default to
	// distribution defaults.
	// Dependent on support in the network stack for your distribution.
	// For more information see the netplan reference (https://netplan.io/reference#dhcp-overrides)
	// +optional
	DHCP4Overrides *DHCPOverrides `json:"dhcp4Overrides,omitempty"`

	// DHCP6Overrides allows for the control over several DHCP behaviors.
	// Overrides will only be applied when the corresponding DHCP flag is set.
	// Only configured values will be sent, omitted values will default to
	// distribution defaults.
	// Dependent on support in the network stack for your distribution.
	// For more information see the netplan reference (https://netplan.io/reference#dhcp-overrides)
	// +optional
	DHCP6Overrides *DHCPOverrides `json:"dhcp6Overrides,omitempty"`

	// SkipIPAllocation allows the device to not have IP address or DHCP configured.
	// This is suitable for devices for which IP allocation is handled externally, eg. using Multus CNI.
	// If true, CAPV will not verify IP address allocation.
	// +optional
	SkipIPAllocation bool `json:"skipIPAllocation,omitempty"`
}

// DHCPOverrides allows for the control over several DHCP behaviors.
// Overrides will only be applied when the corresponding DHCP flag is set.
// Only configured values will be sent, omitted values will default to
// distribution defaults.
// Dependent on support in the network stack for your distribution.
// For more information see the netplan reference (https://netplan.io/reference#dhcp-overrides)
type DHCPOverrides struct {
	// Hostname is the name which will be sent to the DHCP server instead of
	// the machine's hostname.
	// +optional
	Hostname *string `json:"hostname,omitempty"`
	// RouteMetric is used to prioritize routes for devices. A lower metric for
	// an interface will have a higher priority.
	// +optional
	RouteMetric *int `json:"routeMetric,omitempty"`
	// SendHostname when `true`, the hostname of the machine will be sent to the
	// DHCP server.
	// +optional
	SendHostname *bool `json:"sendHostname,omitempty"`
	// UseDNS when `true`, the DNS servers in the DHCP server will be used and
	// take precedence.
	// +optional
	UseDNS *bool `json:"useDNS,omitempty"`
	// UseDomains can take the values `true`, `false`, or `route`. When `true`,
	// the domain name from the DHCP server will be used as the DNS search
	// domain for this device. When `route`, the domain name from the DHCP
	// response will be used for routing DNS only, not for searching.
	// +optional
	UseDomains *string `json:"useDomains,omitempty"`
	// UseHostname when `true`, the hostname from the DHCP server will be set
	// as the transient hostname of the machine.
	// +optional
	UseHostname *bool `json:"useHostname,omitempty"`
	// UseMTU when `true`, the MTU from the DHCP server will be set as the
	// MTU of the device.
	// +optional
	UseMTU *bool `json:"useMTU,omitempty"`
	// UseNTP when `true`, the NTP servers from the DHCP server will be used
	// by systemd-timesyncd and take precedence.
	// +optional
	UseNTP *bool `json:"useNTP,omitempty"`
	// UseRoutes when `true`, the routes from the DHCP server will be installed
	// in the routing table.
	// +optional
	UseRoutes *string `json:"useRoutes,omitempty"`
}

// NetworkRouteSpec defines a static network route.
type NetworkRouteSpec struct {
	// To is an IPv4 or IPv6 address.
	To string `json:"to"`
	// Via is an IPv4 or IPv6 address.
	Via string `json:"via"`
	// Metric is the weight/priority of the route.
	Metric int32 `json:"metric"`
}

// NetworkStatus provides information about one of a VM's networks.
type NetworkStatus struct {
	// Connected is a flag that indicates whether this network is currently
	// connected to the VM.
	Connected bool `json:"connected,omitempty"`

	// IPAddrs is one or more IP addresses reported by vm-tools.
	// +optional
	IPAddrs []string `json:"ipAddrs,omitempty"`

	// MACAddr is the MAC address of the network device.
	MACAddr string `json:"macAddr"`

	// NetworkName is the name of the network.
	// +optional
	NetworkName string `json:"networkName,omitempty"`
}

// VirtualMachineState describes the state of a VM.
type VirtualMachineState string

const (
	// VirtualMachineStateNotFound is the string representing a VM that
	// cannot be located.
	VirtualMachineStateNotFound VirtualMachineState = "notfound"

	// VirtualMachineStatePending is the string representing a VM with an in-flight task.
	VirtualMachineStatePending = "pending"

	// VirtualMachineStateReady is the string representing a powered-on VM with reported IP addresses.
	VirtualMachineStateReady = "ready"
)

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

// VirtualMachine represents data about a vSphere virtual machine object.
type VirtualMachine struct {
	// Name is the VM's name.
	Name string `json:"name"`

	// BiosUUID is the VM's BIOS UUID.
	BiosUUID string `json:"biosUUID"`

	// State is the VM's state.
	State VirtualMachineState `json:"state"`

	// Network is the status of the VM's network devices.
	Network []NetworkStatus `json:"network"`

	// VMRef is the VM's Managed Object Reference on vSphere.
	VMRef string `json:"vmRef"`
}

// SSHUser is granted remote access to a system.
type SSHUser struct {
	// Name is the name of the SSH user.
	Name string `json:"name"`
	// AuthorizedKeys is one or more public SSH keys that grant remote access.
	AuthorizedKeys []string `json:"authorizedKeys"`
}
