/* Copyright © 2022-2023 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: Apache-2.0 */

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AccessMode string

// SubnetSpec defines the desired state of Subnet.
type SubnetSpec struct {
	// Size of Subnet based upon estimated workload count.
	// +kubebuilder:validation:Maximum:=65536
	// +kubebuilder:validation:Minimum:=16
	IPv4SubnetSize int `json:"ipv4SubnetSize,omitempty"`
	// Access mode of Subnet, accessible only from within VPC or from outside VPC.
	// +kubebuilder:validation:Enum=Private;Public
	AccessMode AccessMode `json:"accessMode,omitempty"`
	// Subnet CIDRS.
	// +kubebuilder:validation:MinItems=0
	// +kubebuilder:validation:MaxItems=2
	IPAddresses []string `json:"ipAddresses,omitempty"`
	// Subnet advanced configuration.
	AdvancedConfig AdvancedConfig `json:"advancedConfig,omitempty"`
	// DHCPConfig DHCP configuration.
	DHCPConfig DHCPConfig `json:"DHCPConfig,omitempty"`
}

// SubnetStatus defines the observed state of Subnet.
type SubnetStatus struct {
	NSXResourcePath string      `json:"nsxResourcePath,omitempty"`
	IPAddresses     []string    `json:"ipAddresses,omitempty"`
	Conditions      []Condition `json:"conditions,omitempty"`
}

// +genclient
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:storageversion

// Subnet is the Schema for the subnets API.
// +kubebuilder:printcolumn:name="AccessMode",type=string,JSONPath=`.spec.accessMode`,description="Access mode of Subnet"
// +kubebuilder:printcolumn:name="IPv4SubnetSize",type=string,JSONPath=`.spec.ipv4SubnetSize`,description="Size of Subnet"
// +kubebuilder:printcolumn:name="IPAddresses",type=string,JSONPath=`.status.ipAddresses[*]`,description="CIDRs for the Subnet"
type Subnet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SubnetSpec   `json:"spec,omitempty"`
	Status SubnetStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SubnetList contains a list of Subnet.
type SubnetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Subnet `json:"items"`
}

// AdvancedConfig is Subnet advanced configuration.
type AdvancedConfig struct {
	// StaticIPAllocation configuration for subnet ports with VIF attachment.
	StaticIPAllocation StaticIPAllocation `json:"staticIPAllocation,omitempty"`
}

// StaticIPAllocation is static IP allocation for subnet ports with VIF attachment.
type StaticIPAllocation struct {
	// Enable or disable static IP allocation for subnet ports with VIF attachment.
	// +kubebuilder:default:=false
	Enable bool `json:"enable,omitempty"`
}

// DHCPConfig is DHCP configuration.
type DHCPConfig struct {
	// +kubebuilder:default:=false
	EnableDHCP bool `json:"enableDHCP,omitempty"`
	// DHCPRelayConfigPath is policy path of DHCP-relay-config.
	DHCPRelayConfigPath string `json:"dhcpRelayConfigPath,omitempty"`
	// DHCPV4PoolSize IPs in % to be reserved for DHCP ranges.
	// By default, 80% of IPv4 IPs will be reserved for DHCP.
	// Configure 0 if no pool is required.
	// +kubebuilder:default:=80
	// +kubebuilder:validation:Maximum:=100
	// +kubebuilder:validation:Minimum:=0
	DHCPV4PoolSize int `json:"dhcpV4PoolSize,omitempty"`
	// DHCPV6PoolSize number of IPs to be reserved for DHCP ranges.
	// By default, 2000 IPv6 IPs will be reserved for DHCP.
	// +kubebuilder:default:=2000
	DHCPV6PoolSize  int             `json:"dhcpV6PoolSize,omitempty"`
	DNSClientConfig DNSClientConfig `json:"dnsClientConfig,omitempty"`
}

// DNSClientConfig holds DNS configurations.
type DNSClientConfig struct {
	DNSServersIPs []string `json:"dnsServersIPs,omitempty"`
}

func init() {
	SchemeBuilder.Register(&Subnet{}, &SubnetList{})
}
