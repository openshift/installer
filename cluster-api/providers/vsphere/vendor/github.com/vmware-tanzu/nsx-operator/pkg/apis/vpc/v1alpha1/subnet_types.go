/* Copyright © 2022-2023 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: Apache-2.0 */

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AccessMode string
type DHCPConfigMode string

const (
	AccessModePublic          string = "Public"
	AccessModePrivate         string = "Private"
	AccessModeProject         string = "PrivateTGW"
	DHCPConfigModeDeactivated string = "DHCPDeactivated"
	DHCPConfigModeServer      string = "DHCPServer"
	DHCPConfigModeRelay       string = "DHCPRelay"
)

// SubnetSpec defines the desired state of Subnet.
// +kubebuilder:validation:XValidation:rule="has(oldSelf.subnetDHCPConfig) || !has(self.subnetDHCPConfig) || !has(self.subnetDHCPConfig.mode) || self.subnetDHCPConfig.mode=='DHCPDeactivated'", message="subnetDHCPConfig cannot switch from DHCPDeactivated to other modes"
// +kubebuilder:validation:XValidation:rule="!has(oldSelf.ipv4SubnetSize) || has(self.ipv4SubnetSize)", message="ipv4SubnetSize is required once set"
// +kubebuilder:validation:XValidation:rule="!has(oldSelf.accessMode) || has(self.accessMode)", message="accessMode is required once set"
// +kubebuilder:validation:XValidation:rule="!has(oldSelf.ipAddresses) || has(self.ipAddresses)", message="ipAddresses is required once set"
type SubnetSpec struct {
	// Size of Subnet based upon estimated workload count.
	// +kubebuilder:validation:Maximum:=65536
	// +kubebuilder:validation:Minimum:=16
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	IPv4SubnetSize int `json:"ipv4SubnetSize,omitempty"`
	// Access mode of Subnet, accessible only from within VPC or from outside VPC.
	// +kubebuilder:validation:Enum=Private;Public;PrivateTGW
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	AccessMode AccessMode `json:"accessMode,omitempty"`
	// Subnet CIDRS.
	// +kubebuilder:validation:MinItems=0
	// +kubebuilder:validation:MaxItems=2
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	IPAddresses []string `json:"ipAddresses,omitempty"`

	// DHCP mode of a SubnetSet cannot switch from DHCPDeactivated to DHCPServer or DHCPRelay.
	// If subnetDHCPConfig is not set, the DHCP mode is DHCPDeactivated by default.
	// In order to enforce this rule, two XValidation rules are defined.
	// The rule in SubnetSetSpec prevents the condition that subnetDHCPConfig is not set in
	// old SubnetSetSpec while the new SubnetSetSpec specifies a field other than DHCPDeactivated.
	// The rule in SubnetDHCPConfig prevents the mode changing from empty or
	// DHCPDeactivated to DHCPServer or DHCPRelay.

	// DHCP configuration for Subnet.
	SubnetDHCPConfig SubnetDHCPConfig `json:"subnetDHCPConfig,omitempty"`
}

// SubnetStatus defines the observed state of Subnet.
type SubnetStatus struct {
	NetworkAddresses    []string    `json:"networkAddresses,omitempty"`
	GatewayAddresses    []string    `json:"gatewayAddresses,omitempty"`
	DHCPServerAddresses []string    `json:"DHCPServerAddresses,omitempty"`
	Conditions          []Condition `json:"conditions,omitempty"`
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion

// Subnet is the Schema for the subnets API.
// +kubebuilder:printcolumn:name="AccessMode",type=string,JSONPath=`.spec.accessMode`,description="Access mode of Subnet"
// +kubebuilder:printcolumn:name="IPv4SubnetSize",type=string,JSONPath=`.spec.ipv4SubnetSize`,description="Size of Subnet"
// +kubebuilder:printcolumn:name="NetworkAddresses",type=string,JSONPath=`.status.networkAddresses[*]`,description="CIDRs for the Subnet"
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

// SubnetDHCPConfig is DHCP configuration for Subnet.
type SubnetDHCPConfig struct {
	// DHCP Mode. DHCPDeactivated will be used if it is not defined.
	// It cannot switch from DHCPDeactivated to DHCPServer or DHCPRelay.
	// +kubebuilder:validation:Enum=DHCPServer;DHCPRelay;DHCPDeactivated
	// +kubebuilder:validation:XValidation:rule="oldSelf!='DHCPDeactivated' || oldSelf==self", message="subnetDHCPConfig cannot switch from DHCPDeactivated to other modes"
	Mode DHCPConfigMode `json:"mode,omitempty"`
}

func init() {
	SchemeBuilder.Register(&Subnet{}, &SubnetList{})
}
