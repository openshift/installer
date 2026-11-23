/* Copyright © 2022-2026 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: Apache-2.0 */

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AccessMode string
type DHCPConfigMode string
type DHCPv6ConfigMode string
type ConnectivityState string
type IPAddressType string

const (
	AccessModePublic                string            = "Public"
	AccessModePrivate               string            = "Private"
	AccessModeProject               string            = "PrivateTGW"
	AccessModeL2Only                string            = "L2Only"
	DHCPConfigModeDeactivated       string            = "DHCPDeactivated"
	DHCPConfigModeServer            string            = "DHCPServer"
	DHCPConfigModeRelay             string            = "DHCPRelay"
	DHCPv6ConfigModeDeactivated     DHCPv6ConfigMode  = "DHCPDeactivated"
	DHCPv6ConfigModeServer          DHCPv6ConfigMode  = "DHCPServer"
	DHCPv6ConfigModeRelay           DHCPv6ConfigMode  = "DHCPRelay"
	DHCPv6ConfigModeServerStateless DHCPv6ConfigMode  = "DHCPServerStateless"
	ConnectivityStateConnected      ConnectivityState = "Connected"
	ConnectivityStateDisconnected   ConnectivityState = "Disconnected"
	IPAddressTypeIPv4               IPAddressType     = "IPv4"
	IPAddressTypeIPv6               IPAddressType     = "IPv6"
	IPAddressTypeIPv4IPv6           IPAddressType     = "IPv4IPv6"
)

// SubnetSpec defines the desired state of Subnet.
// +kubebuilder:validation:XValidation:rule="!has(oldSelf.vpcName) || self.vpcName == oldSelf.vpcName",message="vpcName is immutable after set"
// +kubebuilder:validation:XValidation:rule="!has(oldSelf.ipv4SubnetSize) || has(self.ipv4SubnetSize)", message="ipv4SubnetSize is required once set"
// +kubebuilder:validation:XValidation:rule="!has(oldSelf.ipv6PrefixLength) || has(self.ipv6PrefixLength)", message="ipv6PrefixLength is required once set"
// +kubebuilder:validation:XValidation:rule="!has(oldSelf.accessMode) || has(self.accessMode)", message="accessMode is required once set"
// +kubebuilder:validation:XValidation:rule="!(has(oldSelf.advancedConfig) && has(oldSelf.advancedConfig.staticIPAllocation) && has(oldSelf.advancedConfig.staticIPAllocation.enabled) && (!has(self.advancedConfig.staticIPAllocation.enabled) || oldSelf.advancedConfig.staticIPAllocation.enabled != self.advancedConfig.staticIPAllocation.enabled))", message="staticIPAllocation enabled cannot be changed once set"
// +kubebuilder:validation:XValidation:rule="!(has(self.advancedConfig) && has(self.advancedConfig.staticIPAllocation) && has(self.advancedConfig.staticIPAllocation.enabled) && self.advancedConfig.staticIPAllocation.enabled==true && has(self.subnetDHCPConfig) && has(self.subnetDHCPConfig.mode) && (self.subnetDHCPConfig.mode=='DHCPServer' || self.subnetDHCPConfig.mode=='DHCPRelay'))", message="Static IP allocation and Subnet DHCP configuration cannot be enabled simultaneously on a Subnet"
// +kubebuilder:validation:XValidation:rule="!(has(self.advancedConfig) && has(self.advancedConfig.dhcpServerAddresses) && size(self.advancedConfig.dhcpServerAddresses)>0) || (has(self.subnetDHCPConfig) && has(self.subnetDHCPConfig.mode) && self.subnetDHCPConfig.mode=='DHCPServer') || (has(self.subnetDHCPv6Config) && has(self.subnetDHCPv6Config.mode) && self.subnetDHCPv6Config.mode=='DHCPServer')", message="DHCPServerAddresses can only be set when DHCP mode or DHCPv6 mode is DHCPServer"
// +kubebuilder:validation:XValidation:rule="!has(self.ipAddresses) && !(has(self.subnetDHCPConfig) && has(self.subnetDHCPConfig.dhcpServerAdditionalConfig) && has(self.subnetDHCPConfig.dhcpServerAdditionalConfig.reservedIPRanges)) || has(self.ipAddresses)", message="ipAddresses is required to configure subnet reserved ip ranges."
// +kubebuilder:validation:XValidation:rule="!has(self.ipAddresses) && !(has(self.subnetDHCPv6Config) && has(self.subnetDHCPv6Config.dhcpv6ServerAdditionalConfig) && has(self.subnetDHCPv6Config.dhcpv6ServerAdditionalConfig.reservedIPRanges)) || has(self.ipAddresses)", message="ipAddresses is required to configure subnet DHCPv6 reserved ip ranges."
// +kubebuilder:validation:XValidation:rule="!(has(self.advancedConfig) && has(self.advancedConfig.gatewayAddresses) && size(self.advancedConfig.gatewayAddresses)>0) || has(self.ipAddresses)", message="ipAddresses is required when custom gatewayAddresses are specified"
// +kubebuilder:validation:XValidation:rule="!(has(self.advancedConfig) && has(self.advancedConfig.dhcpServerAddresses) && size(self.advancedConfig.dhcpServerAddresses)>0) || has(self.ipAddresses)", message="ipAddresses is required when custom dhcpServerAddresses are specified"
type SubnetSpec struct {
	// VPC name of the Subnet.
	VPCName string `json:"vpcName,omitempty"`
	// IPAddressType defines the IP address type that will be allocated for the Subnet.
	// +kubebuilder:validation:Enum=IPv4;IPv6;IPv4IPv6
	// +kubebuilder:default=IPv4
	IPAddressType IPAddressType `json:"ipAddressType,omitempty"`
	// Size of IPv4 Subnet based upon estimated workload count.
	// +kubebuilder:validation:Maximum:=65536
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	IPv4SubnetSize int `json:"ipv4SubnetSize,omitempty"`
	// IPv6 prefix length for the Subnet (e.g. 64 means /64).
	// +kubebuilder:validation:Minimum:=2
	// +kubebuilder:validation:Maximum:=127
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	IPv6PrefixLength int `json:"ipv6PrefixLength,omitempty"`
	// Access mode of IPv4 Subnet, accessible only from within VPC or from outside VPC.
	// +kubebuilder:validation:Enum=Private;Public;PrivateTGW;L2Only
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	AccessMode AccessMode `json:"accessMode,omitempty"`
	// Subnet CIDRS.
	// +kubebuilder:validation:MinItems=0
	// +kubebuilder:validation:MaxItems=2
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	IPAddresses []string `json:"ipAddresses,omitempty"`
	// DHCP configuration for Subnet.
	SubnetDHCPConfig SubnetDHCPConfig `json:"subnetDHCPConfig,omitempty"`
	// DHCPv6 configuration for Subnet.
	SubnetDHCPv6Config SubnetDHCPv6Config `json:"subnetDHCPv6Config,omitempty"`
	// VPC Subnet advanced configuration.
	AdvancedConfig SubnetAdvancedConfig `json:"advancedConfig,omitempty"`
	// Distributed VLAN Connection name.
	VLANConnectionName string `json:"vlanConnectionName,omitempty"`
}

// SubnetStatus defines the observed state of Subnet.
type SubnetStatus struct {
	// Network address of the Subnet.
	NetworkAddresses []string `json:"networkAddresses,omitempty"`
	// Gateway address of the Subnet.
	GatewayAddresses []string `json:"gatewayAddresses,omitempty"`
	// DHCP server IP address.
	DHCPServerAddresses []string `json:"DHCPServerAddresses,omitempty"`
	// VLAN extension configured for VPC Subnet.
	VLANExtension VLANExtension `json:"vlanExtension,omitempty"`
	// Whether this is a pre-created Subnet shared with the Namespace.
	// +kubebuilder:default=false
	Shared     bool        `json:"shared,omitempty"`
	Conditions []Condition `json:"conditions,omitempty"`
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion

// Subnet is the Schema for the subnets API.
// +kubebuilder:printcolumn:name="AccessMode",type=string,JSONPath=`.spec.accessMode`,description="Access mode of IPv4 Subnet"
// +kubebuilder:printcolumn:name="IPAddressType",type=string,JSONPath=`.spec.ipAddressType`,description="IP address type of Subnet"
// +kubebuilder:printcolumn:name="IPv4SubnetSize",type=string,JSONPath=`.spec.ipv4SubnetSize`,description="Size of IPv4 Subnet"
// +kubebuilder:printcolumn:name="IPv6PrefixLength",type=string,JSONPath=`.spec.ipv6PrefixLength`,description="Prefix length of IPv6 Subnet"
// +kubebuilder:printcolumn:name="NetworkAddresses",type=string,JSONPath=`.status.networkAddresses[*]`,description="CIDRs for the Subnet"
// +kubebuilder:validation:XValidation:rule="!has(oldSelf.spec) || has(self.spec)", message="spec is required once set"
type Subnet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SubnetSpec   `json:"spec,omitempty"`
	Status SubnetStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SubnetList contains a list of Subnet.
type SubnetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Subnet `json:"items"`
}

type SubnetAdvancedConfig struct {
	// Connectivity status of the Subnet from other Subnets of the VPC.
	// The default value is "Connected".
	// +kubebuilder:validation:Enum=Connected;Disconnected
	// +kubebuilder:default=Connected
	ConnectivityState ConnectivityState `json:"connectivityState,omitempty"`
	// Static IP allocation for VPC Subnet Ports.
	StaticIPAllocation StaticIPAllocation `json:"staticIPAllocation,omitempty"`
	// GatewayAddresses specifies custom gateway IP addresses for the Subnet.
	// Supports up to 2 addresses for dual-stack Subnets (1 IPv4 + 1 IPv6).
	// +kubebuilder:validation:MaxItems=2
	GatewayAddresses []string `json:"gatewayAddresses,omitempty"`
	// DHCPServerAddresses specifies custom DHCP server IP addresses for the Subnet.
	// Supports up to 2 addresses for dual-stack Subnets (1 IPv4 + 1 IPv6).
	// +kubebuilder:validation:MaxItems=2
	DHCPServerAddresses []string `json:"dhcpServerAddresses,omitempty"`
}

type StaticIPAllocation struct {
	// Activate or deactivate static IP allocation for VPC Subnet Ports.
	// If the DHCP mode is DHCPDeactivated or not set, its default value is true.
	// If the DHCP mode is DHCPServer or DHCPRelay, its default value is false.
	// The value cannot be set to true when the DHCP mode is DHCPServer or DHCPRelay.
	Enabled *bool `json:"enabled,omitempty"`
}

// Additional DHCP server config for a VPC Subnet.
// The additional configuration must not be set when the Subnet has DHCP relay enabled or DHCP is deactivated.
type DHCPServerAdditionalConfig struct {
	// Reserved IP ranges.
	// Supported formats include: ["192.168.1.1", "192.168.1.3-192.168.1.100"]
	// +kubebuilder:validation::MaxItems=10
	ReservedIPRanges []string `json:"reservedIPRanges,omitempty"`
}

// SubnetDHCPConfig is a DHCP configuration for Subnet.
// +kubebuilder:validation:XValidation:rule="(!has(self.mode)|| self.mode=='DHCPDeactivated' || self.mode=='DHCPRelay' ) && (!has(self.dhcpServerAdditionalConfig) || !has(self.dhcpServerAdditionalConfig.reservedIPRanges) || size(self.dhcpServerAdditionalConfig.reservedIPRanges)==0) || has(self.mode) && self.mode=='DHCPServer'", message="DHCPServerAdditionalConfig must be cleared when Subnet has DHCP relay enabled or DHCP is deactivated."
type SubnetDHCPConfig struct {
	// DHCP Mode. DHCPDeactivated will be used if it is not defined.
	// +kubebuilder:validation:Enum=DHCPServer;DHCPRelay;DHCPDeactivated
	Mode DHCPConfigMode `json:"mode,omitempty"`
	// Additional DHCP server config for a VPC Subnet.
	DHCPServerAdditionalConfig DHCPServerAdditionalConfig `json:"dhcpServerAdditionalConfig,omitempty"`
}

// DHCPv6ServerAdditionalConfig is additional DHCPv6 server config for a VPC Subnet.
// The additional configuration must not be set when the Subnet has DHCP relay enabled or DHCP is deactivated.
type DHCPv6ServerAdditionalConfig struct {
	// Reserved IPv6 ranges.
	// Supported formats include: ["2001:db8::1", "2001:db8::1-2001:db8::ff"]
	// +kubebuilder:validation::MaxItems=10
	ReservedIPRanges []string `json:"reservedIPRanges,omitempty"`
}

// SubnetDHCPv6Config is a DHCPv6 configuration for Subnet.
// +kubebuilder:validation:XValidation:rule="(!has(self.mode) || self.mode=='DHCPDeactivated' || self.mode=='DHCPRelay') && (!has(self.dhcpv6ServerAdditionalConfig) || !has(self.dhcpv6ServerAdditionalConfig.reservedIPRanges) || size(self.dhcpv6ServerAdditionalConfig.reservedIPRanges)==0) || has(self.mode) && self.mode=='DHCPServer'", message="DHCPv6ServerAdditionalConfig must be cleared when Subnet has DHCP relay enabled or DHCP is deactivated."
type SubnetDHCPv6Config struct {
	// DHCPv6 Mode. DHCPDeactivated will be used if it is not defined.
	// +kubebuilder:validation:Enum=DHCPServer;DHCPRelay;DHCPDeactivated;DHCPServerStateless
	Mode DHCPv6ConfigMode `json:"mode,omitempty"`
	// Additional DHCPv6 server config for a VPC Subnet.
	DHCPv6ServerAdditionalConfig DHCPv6ServerAdditionalConfig `json:"dhcpv6ServerAdditionalConfig,omitempty"`
}

// VLANExtension describes VLAN extension configuration for the VPC Subnet.
type VLANExtension struct {
	// Flag to control whether the VLAN extension Subnet connects to the VPC gateway.
	VPCGatewayConnectionEnable bool `json:"vpcGatewayConnectionEnable,omitempty"`
	// VLAN ID of the VLAN extension Subnet.
	VLANID int `json:"vlanId,omitempty"`
}

func init() {
	SchemeBuilder.Register(&Subnet{}, &SubnetList{})
}
