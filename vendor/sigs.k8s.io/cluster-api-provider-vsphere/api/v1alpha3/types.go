/*
Copyright 2018 The Kubernetes Authors.

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

package v1alpha3

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
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

// VSphereMachineTemplateResource describes the data needed to create a VSphereMachine from a template
type VSphereMachineTemplateResource struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	clusterv1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the specification of the desired behavior of the machine.
	Spec VSphereMachineSpec `json:"spec"`
}

// VSphereMachineProviderConditionType is a valid value for VSphereMachineProviderCondition.Type
type VSphereMachineProviderConditionType string

// Valid conditions for an VSphere machine instance
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
	Port int `json:"port"`
}

// IsZero returns true if either the host or the port are zero values.
func (v APIEndpoint) IsZero() bool {
	return v.Host == "" || v.Port == 0
}

// String returns a formatted version HOST:PORT of this APIEndpoint.
func (v APIEndpoint) String() string {
	return fmt.Sprintf("%s:%d", v.Host, v.Port)
}

// NetworkSpec defines the virtual machine's network configuration.
type NetworkSpec struct {
	// Devices is the list of network devices used by the virtual machine.
	// TODO(akutz) Make sure at least one network matches the
	//             ClusterSpec.CloudProviderConfiguration.Network.Name
	Devices []NetworkDeviceSpec `json:"devices"`

	// Routes is a list of optional, static routes applied to the virtual
	// machine.
	// +optional
	Routes []NetworkRouteSpec `json:"routes,omitempty"`

	// PreferredAPIServeCIDR is the preferred CIDR for the Kubernetes API
	// server endpoint on this machine
	// +optional
	PreferredAPIServerCIDR string `json:"preferredAPIServerCidr,omitempty"`
}

// NetworkDeviceSpec defines the network configuration for a virtual machine's
// network device.
type NetworkDeviceSpec struct {
	// NetworkName is the name of the vSphere network to which the device
	// will be connected.
	NetworkName string `json:"networkName"`

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
	// Required when DHCP6 is false.
	// +optional
	Gateway6 string `json:"gateway6,omitempty"`

	// IPAddrs is a list of one or more IPv4 and/or IPv6 addresses to assign
	// to this device.
	// Required when DHCP4 and DHCP6 are both false.
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

// VirtualMachinePowerState describe the power state of a VM
type VirtualMachinePowerState string

const (
	// VirtualMachinePowerStatePoweredOn is the string representing a VM in powered on state
	VirtualMachinePowerStatePoweredOn VirtualMachinePowerState = "poweredOn"

	// VirtualMachinePowerStatePoweredOff is the string representing a VM in powered off state
	VirtualMachinePowerStatePoweredOff = "poweredOff"

	// VirtualMachinePowerStateSuspended is the string representing a VM in suspended state
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
}
