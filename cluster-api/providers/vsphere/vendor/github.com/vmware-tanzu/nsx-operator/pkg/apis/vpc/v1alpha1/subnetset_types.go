/* Copyright © 2022-2026 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: Apache-2.0 */

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SubnetSetSpec defines the desired state of SubnetSet.
// +kubebuilder:validation:XValidation:rule="!has(oldSelf.accessMode) || has(self.accessMode)", message="accessMode is required once set"
// +kubebuilder:validation:XValidation:rule="!has(oldSelf.ipv4SubnetSize) || has(self.ipv4SubnetSize)", message="ipv4SubnetSize is required once set"
// +kubebuilder:validation:XValidation:rule="!has(oldSelf.ipv6PrefixLength) || has(self.ipv6PrefixLength)", message="ipv6PrefixLength is required once set"
// +kubebuilder:validation:XValidation:rule="!has(self.subnetDHCPConfig) || has(self.subnetDHCPConfig) && !has(self.subnetDHCPConfig.dhcpServerAdditionalConfig) || has(self.subnetDHCPConfig) && has(self.subnetDHCPConfig.dhcpServerAdditionalConfig) && !has(self.subnetDHCPConfig.dhcpServerAdditionalConfig.reservedIPRanges)", message="reservedIPRanges is not supported in SubnetSet"
// +kubebuilder:validation:XValidation:rule="!has(self.subnetDHCPv6Config) || has(self.subnetDHCPv6Config) && !has(self.subnetDHCPv6Config.dhcpv6ServerAdditionalConfig) || has(self.subnetDHCPv6Config) && has(self.subnetDHCPv6Config.dhcpv6ServerAdditionalConfig) && !has(self.subnetDHCPv6Config.dhcpv6ServerAdditionalConfig.reservedIPRanges)", message="reservedIPRanges is not supported in SubnetSet"
// +kubebuilder:validation:XValidation:rule="!has(self.subnetDHCPConfig) || !has(self.subnetDHCPConfig.mode) || self.subnetDHCPConfig.mode!='DHCPRelay'", message="DHCPRelay is not supported in SubnetSet"
// +kubebuilder:validation:XValidation:rule="!has(self.subnetDHCPv6Config) || !has(self.subnetDHCPv6Config.mode) || self.subnetDHCPv6Config.mode!='DHCPRelay'", message="DHCPRelay is not supported in SubnetSet"
type SubnetSetSpec struct {
	// IPAddressType defines the IP address type that will be allocated for subnets in the SubnetSet.
	// +kubebuilder:validation:Enum=IPv4;IPv6;IPv4IPv6
	IPAddressType IPAddressType `json:"ipAddressType,omitempty"`
	// Size of IPv4 Subnet based upon estimated workload count.
	// +kubebuilder:validation:Maximum:=65536
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	IPv4SubnetSize int `json:"ipv4SubnetSize,omitempty"`
	// IPv6 prefix length for subnets in the SubnetSet (e.g. 64 means /64).
	// +kubebuilder:validation:Minimum:=2
	// +kubebuilder:validation:Maximum:=127
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	IPv6PrefixLength int `json:"ipv6PrefixLength,omitempty"`
	// Access mode of IPv4 Subnet, accessible only from within VPC or from outside VPC.
	// +kubebuilder:validation:Enum=Private;Public;PrivateTGW
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	AccessMode AccessMode `json:"accessMode,omitempty"`
	// Subnet DHCP configuration.
	SubnetDHCPConfig SubnetDHCPConfig `json:"subnetDHCPConfig,omitempty"`
	// DHCPv6 configuration for subnets in the SubnetSet.
	SubnetDHCPv6Config SubnetDHCPv6Config `json:"subnetDHCPv6Config,omitempty"`
	// The names of the Subnets that have been created in advance.
	// It is mutually exclusive with the other fields like IPv4SubnetSize, AccessMode, and SubnetDHCPConfig.
	// Once this field is set, the other fields cannot be set.
	SubnetNames *[]string `json:"subnetNames,omitempty"`
}

// SubnetInfo defines the observed state of a single Subnet of a SubnetSet.
type SubnetInfo struct {
	// Network address of the Subnet.
	NetworkAddresses []string `json:"networkAddresses,omitempty"`
	// Gateway address of the Subnet.
	GatewayAddresses []string `json:"gatewayAddresses,omitempty"`
	// Dhcp server IP address.
	DHCPServerAddresses []string `json:"DHCPServerAddresses,omitempty"`
}

// SubnetSetStatus defines the observed state of SubnetSet.
type SubnetSetStatus struct {
	Conditions []Condition  `json:"conditions,omitempty"`
	Subnets    []SubnetInfo `json:"subnets,omitempty"`
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion

// SubnetSet is the Schema for the subnetsets API.
// +kubebuilder:printcolumn:name="AccessMode",type=string,JSONPath=`.spec.accessMode`,description="Access mode of IPv4 Subnet"
// +kubebuilder:printcolumn:name="IPAddressType",type=string,JSONPath=`.spec.ipAddressType`,description="IP address type of Subnet"
// +kubebuilder:printcolumn:name="IPv4SubnetSize",type=string,JSONPath=`.spec.ipv4SubnetSize`,description="Size of IPv4 Subnet"
// +kubebuilder:printcolumn:name="IPv6PrefixLength",type=string,JSONPath=`.spec.ipv6PrefixLength`,description="Prefix length of IPv6 Subnet"
// +kubebuilder:printcolumn:name="NetworkAddresses",type=string,JSONPath=`.status.subnets[*].networkAddresses[*]`,description="CIDRs for the SubnetSet"
// +kubebuilder:validation:XValidation:rule="!has(oldSelf.spec) || has(self.spec)", message="spec is required once set"
type SubnetSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SubnetSetSpec   `json:"spec,omitempty"`
	Status SubnetSetStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SubnetSetList contains a list of SubnetSet.
type SubnetSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SubnetSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SubnetSet{}, &SubnetSetList{})
}
