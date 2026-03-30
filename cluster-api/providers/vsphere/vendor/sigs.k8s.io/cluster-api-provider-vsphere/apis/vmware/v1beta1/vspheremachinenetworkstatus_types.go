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

package v1beta1

import (
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VSphereMachineNetworkDNSStatus describes the observed state of the guest's
// RFC 1034 client-side DNS settings.
// +kubebuilder:validation:MinProperties=1
type VSphereMachineNetworkDNSStatus struct {
	// dhcp indicates whether or not dynamic host control protocol (DHCP) was
	// used to configure DNS configuration.
	//
	// +optional
	DHCP *bool `json:"dhcp,omitempty"`

	// domainName is the domain name portion of the DNS name. For example,
	// the "domain.local" part of "my-vm.domain.local".
	//
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=512
	DomainName string `json:"domainName,omitempty"`

	// hostName is the host name portion of the DNS name. For example,
	// the "my-vm" part of "my-vm.domain.local".
	//
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=512
	HostName string `json:"hostName,omitempty"`

	// nameservers is a list of the IP addresses for the DNS servers to use.
	//
	// IP4 addresses are specified using dotted decimal notation. For example,
	// "192.0.2.1".
	//
	// IP6 addresses are 128-bit addresses represented as eight fields of up to
	// four hexadecimal digits. A colon separates each field (:). For example,
	// 2001:DB8:101::230:6eff:fe04:d9ff. The address can also consist of the
	// symbol '::' to represent multiple 16-bit groups of contiguous 0's only
	// once in an address as described in RFC 2373.
	//
	// +optional
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=100
	// +kubebuilder:validation:items:MinLength=1
	// +kubebuilder:validation:items:MaxLength=512
	// +listType=atomic
	Nameservers []string `json:"nameservers,omitempty"`

	// searchDomains is a list of domains in which to search for hosts, in the
	// order of preference.
	//
	// +optional
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=100
	// +kubebuilder:validation:items:MinLength=1
	// +kubebuilder:validation:items:MaxLength=512
	// +listType=atomic
	SearchDomains []string `json:"searchDomains,omitempty"`
}

// IsDefined returns true if the VSphereMachineNetworkDNSStatus is defined.
func (r *VSphereMachineNetworkDNSStatus) IsDefined() bool {
	return !reflect.DeepEqual(r, &VSphereMachineNetworkDNSStatus{})
}

// KeyValuePair is useful when wanting to realize a map as a list of key/value
// pairs.
type KeyValuePair struct {
	// key is the key part of the key/value pair.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=100
	Key string `json:"key,omitempty"`
	// value is the optional value part of the key/value pair.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=10000
	Value string `json:"value,omitempty"`
}

// VSphereMachineNetworkDHCPOptionsStatus describes the observed state of
// DHCP options.
// +kubebuilder:validation:MinProperties=1
type VSphereMachineNetworkDHCPOptionsStatus struct {
	// config describes platform-dependent settings for the DHCP client.
	//
	// The key part is a unique number while the value part is the platform
	// specific configuration command. For example on Linux and BSD systems
	// using the file dhclient.conf output would be reported at system scope:
	// key='1', value='timeout 60;' key='2', value='reboot 10;'. The output
	// reported per interface would be:
	// key='1', value='prepend domain-name-servers 192.0.2.1;'
	// key='2', value='require subnet-mask, domain-name-servers;'.
	//
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=100
	Config []KeyValuePair `json:"config,omitempty"`

	// enabled reports the status of the DHCP client services.
	// +optional
	Enabled *bool `json:"enabled,omitempty"`
}

// IsDefined returns true if the VSphereMachineNetworkDHCPOptionsStatus is defined.
func (r *VSphereMachineNetworkDHCPOptionsStatus) IsDefined() bool {
	return !reflect.DeepEqual(r, &VSphereMachineNetworkDHCPOptionsStatus{})
}

// VSphereMachineNetworkDHCPStatus describes the observed state of the
// client-side, system-wide DHCP settings for IP4 and IP6.
// +kubebuilder:validation:MinProperties=1
type VSphereMachineNetworkDHCPStatus struct {

	// ip4 describes the observed state of the IP4 DHCP client settings.
	//
	// +optional
	IP4 VSphereMachineNetworkDHCPOptionsStatus `json:"ip4,omitempty,omitzero"`

	// ip6 describes the observed state of the IP6 DHCP client settings.
	//
	// +optional
	IP6 VSphereMachineNetworkDHCPOptionsStatus `json:"ip6,omitempty,omitzero"`
}

// IsDefined returns true if the VSphereMachineNetworkDHCPStatus is defined.
func (r *VSphereMachineNetworkDHCPStatus) IsDefined() bool {
	return !reflect.DeepEqual(r, &VSphereMachineNetworkDHCPStatus{})
}

// VSphereMachineNetworkInterfaceIPAddrStatus describes information about a
// specific IP address.
type VSphereMachineNetworkInterfaceIPAddrStatus struct {
	// address is an IP4 or IP6 address and their network prefix length.
	//
	// An IP4 address is specified using dotted decimal notation. For example,
	// "192.0.2.1".
	//
	// IP6 addresses are 128-bit addresses represented as eight fields of up to
	// four hexadecimal digits. A colon separates each field (:). For example,
	// 2001:DB8:101::230:6eff:fe04:d9ff. The address can also consist of the
	// symbol '::' to represent multiple 16-bit groups of contiguous 0's only
	// once in an address as described in RFC 2373.
	//
	// +required
	// +kubebuilder:validation:MinLength=7
	// +kubebuilder:validation:MaxLength=512
	Address string `json:"address,omitempty"`

	// lifetime describes when this address will expire.
	//
	// +optional
	Lifetime metav1.Time `json:"lifetime,omitempty,omitzero"`

	// origin describes how this address was configured.
	//
	// +optional
	// +kubebuilder:validation:Enum=dhcp;linklayer;manual;other;random
	Origin string `json:"origin,omitempty"`

	// state describes the state of this IP address.
	//
	// +optional
	// +kubebuilder:validation:Enum=deprecated;duplicate;inaccessible;invalid;preferred;tentative;unknown
	State string `json:"state,omitempty"`
}

// VSphereMachineNetworkInterfaceIPStatus describes the observed state of a
// VM's network interface's IP configuration.
// +kubebuilder:validation:MinProperties=1
type VSphereMachineNetworkInterfaceIPStatus struct {
	// autoConfigurationEnabled describes whether or not ICMPv6 router
	// solicitation requests are enabled or disabled from a given interface.
	//
	// These requests acquire an IP6 address and default gateway route from
	// zero-to-many routers on the connected network.
	//
	// If not set then ICMPv6 is not available on this VM.
	//
	// +optional
	AutoConfigurationEnabled *bool `json:"autoConfigurationEnabled,omitempty"`

	// dhcp describes the VM's observed, client-side, interface-specific DHCP
	// options.
	//
	// +optional
	DHCP VSphereMachineNetworkDHCPStatus `json:"dhcp,omitempty,omitzero"`

	// addresses describes observed IP addresses for this interface.
	//
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=10
	// +listType=atomic
	// +optional
	Addresses []VSphereMachineNetworkInterfaceIPAddrStatus `json:"addresses,omitempty"`

	// macAddr describes the observed MAC address for this interface.
	//
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=512
	// +optional
	MACAddr string `json:"macAddr,omitempty"`
}

// IsDefined returns true if the VSphereMachineNetworkInterfaceIPStatus is defined.
func (r *VSphereMachineNetworkInterfaceIPStatus) IsDefined() bool {
	return !reflect.DeepEqual(r, &VSphereMachineNetworkInterfaceIPStatus{})
}

// VSphereMachineNetworkInterfaceStatus describes the observed state of a
// VM's network interface.
// +kubebuilder:validation:MinProperties=1
type VSphereMachineNetworkInterfaceStatus struct {
	// name describes the corresponding network interface with the same name
	// in the VM's desired network interface list. If unset, then there is no
	// corresponding entry for this interface.
	//
	// Please note this name is not necessarily related to the name of the
	// device as it is surfaced inside of the guest.
	//
	// +kubebuilder:validation:MinLength=2
	// +kubebuilder:validation:MaxLength=512
	// +optional
	Name string `json:"name,omitempty"`

	// deviceKey describes the unique hardware device key of this network
	// interface.
	//
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=1000000
	// +optional
	DeviceKey int32 `json:"deviceKey,omitempty"`

	// ip describes the observed state of the interface's IP configuration.
	//
	// +optional
	IP VSphereMachineNetworkInterfaceIPStatus `json:"ip,omitempty,omitzero"`

	// dns describes the observed state of the interface's DNS configuration.
	//
	// +optional
	DNS VSphereMachineNetworkDNSStatus `json:"dns,omitempty,omitzero"`
}

// VSphereMachineNetworkStatus defines the observed state of a VM's
// network configuration.
//
// This a mirror of the v1alpha2 VirtualMachineNetworkStatus. See
// https://github.com/vmware-tanzu/vm-operator/blob/main/api/v1alpha2/virtualmachine_network_types.go
// for more information.
// +kubebuilder:validation:MinProperties=1
type VSphereMachineNetworkStatus struct {
	// interfaces describes the status of the VM's network interfaces.
	//
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=100
	// +listType=atomic
	// +optional
	Interfaces []VSphereMachineNetworkInterfaceStatus `json:"interfaces,omitempty"`
}

// IsDefined returns true if the VSphereMachineNetworkStatus is defined.
func (r *VSphereMachineNetworkStatus) IsDefined() bool {
	return !reflect.DeepEqual(r, &VSphereMachineNetworkStatus{})
}
