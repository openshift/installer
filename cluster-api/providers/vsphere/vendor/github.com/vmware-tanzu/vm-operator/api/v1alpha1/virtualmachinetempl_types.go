// Copyright (c) 2022 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

// NetworkDeviceStatus defines the network interface IP configuration including
// gateway, subnetmask and IP address as seen by OVF properties.
type NetworkDeviceStatus struct {
	// Gateway4 is the gateway for the IPv4 address family for this device.
	// +optional
	Gateway4 string

	// MacAddress is the MAC address of the network device.
	// +optional
	MacAddress string

	// IpAddresses represents one or more IP addresses assigned to the network
	// device in CIDR notation, ex. "192.0.2.1/16".
	// +optional
	IPAddresses []string
}

// NetworkStatus describes the observed state of the VM's network configuration.
type NetworkStatus struct {
	// Devices describe a list of current status information for each
	// network interface that is desired to be attached to the
	// VirtualMachineTemplate.
	// +optional
	Devices []NetworkDeviceStatus

	// Nameservers describe a list of the DNS servers accessible by one of the
	// VM's configured network devices.
	// +optional
	Nameservers []string
}

// VirtualMachineTemplate defines the specification for configuring
// VirtualMachine Template. A Virtual Machine Template is created during VM
// customization to populate OVF properties. Then by utilizing Golang-based
// templating, Virtual Machine Template provides access to dynamic configuration
// data.
type VirtualMachineTemplate struct {
	// Net describes the observed state of the VM's network configuration.
	// +optional
	Net NetworkStatus

	// VM represents a pointer to a VirtualMachine instance that consist of the
	// desired specification and the observed status
	VM *VirtualMachine
}
