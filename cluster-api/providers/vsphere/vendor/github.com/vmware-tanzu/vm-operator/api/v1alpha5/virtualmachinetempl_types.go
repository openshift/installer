// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package v1alpha5

// NetworkDeviceStatus defines the network interface IP configuration including
// gateway, subnet mask and IP address as seen by OVF properties.
type NetworkDeviceStatus struct {
	// +optional

	// Gateway4 is the gateway for the IPv4 address family for this device.
	Gateway4 string

	// +optional

	// MacAddress is the MAC address of the network device.
	MacAddress string

	// +optional

	// IpAddresses represents one or more IP addresses assigned to the network
	// device in CIDR notation, ex. "192.0.2.1/16".
	IPAddresses []string
}

// NetworkStatus describes the observed state of the VM's network configuration.
type NetworkStatus struct {
	// +optional

	// Devices describe a list of current status information for each
	// network interface that is desired to be attached to the
	// VirtualMachineTemplate.
	Devices []NetworkDeviceStatus

	// +optional

	// Nameservers describe a list of the DNS servers accessible by one of the
	// VM's configured network devices.
	Nameservers []string
}

// VirtualMachineTemplate defines the specification for configuring
// VirtualMachine Template. A Virtual Machine Template is created during VM
// customization to populate OVF properties. Then by utilizing Golang-based
// templating, Virtual Machine Template provides access to dynamic configuration
// data.
type VirtualMachineTemplate struct {
	// +optional

	// Net describes the observed state of the VM's network configuration.
	Net NetworkStatus

	// VM represents a pointer to a VirtualMachine instance that consist of the
	// desired specification and the observed status
	VM *VirtualMachine
}
