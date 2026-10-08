/* Copyright © 2022-2025 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: Apache-2.0 */

// +kubebuilder:object:generate=true
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VPCNetworkConfigurationSpec defines the desired state of VPCNetworkConfiguration.
// There is a default VPCNetworkConfiguration that applies to Namespaces
// do not have a VPCNetworkConfiguration assigned. When a field is not set
// in a Namespace's VPCNetworkConfiguration, the Namespace will use the value
// in the default VPCNetworkConfiguration.
type VPCNetworkConfigurationSpec struct {
	// NSX path of the VPC the Namespace is associated with.
	// If vpc is set, only defaultSubnetSize and defaultIPv6PrefixLength take effect, other fields are ignored.
	// +optional
	VPC string `json:"vpc,omitempty"`

	// Shared Subnets the Namespace is associated with.
	// +optional
	Subnets []SharedSubnet `json:"subnets,omitempty"`

	// NSX Project the Namespace is associated with.
	NSXProject string `json:"nsxProject,omitempty"`

	// VPCConnectivityProfile Path. This profile has configuration related to creating VPC transit gateway attachment.
	VPCConnectivityProfile string `json:"vpcConnectivityProfile,omitempty"`

	// Private IPs.
	PrivateIPs []string `json:"privateIPs,omitempty"`

	// Default size of IPv4 Subnets.
	// Defaults to 32.
	// +kubebuilder:default=32
	// +kubebuilder:validation:Maximum:=65536
	DefaultSubnetSize int `json:"defaultSubnetSize,omitempty"`

	// DNSZones specifies the list of permitted DNS zones, identified by their NSX paths.
	DNSZones []string `json:"dnsZones,omitempty"`

	// Default prefix length of IPv6 Subnets.
	// Defaults to 64.
	// +kubebuilder:default=64
	// +kubebuilder:validation:Minimum:=2
	// +kubebuilder:validation:Maximum:=127
	DefaultIPv6PrefixLength int `json:"defaultIPv6PrefixLength,omitempty"`
}

// SharedSubnet defines the information for a Subnet shared with vSphere Namespace.
type SharedSubnet struct {
	// NSX path of Subnets created outside of the Supervisor to be associated with this vSphere Namespace
	Path string `json:"path"`
	// Indicates if this Subnet is used for the Pod default network.
	PodDefault bool `json:"podDefault,omitempty"`
	// Indicates if this Subnet is used for the VM default network.
	VMDefault bool `json:"vmDefault,omitempty"`
	// Name of the Subnet. If the name is empty, it will be derived from the shared Subnet path.
	// This field is immutable.
	Name string `json:"name,omitempty"`
}

// VPCNetworkConfigurationStatus defines the observed state of VPCNetworkConfiguration
type VPCNetworkConfigurationStatus struct {
	// VPCs describes VPC info, now it includes Load Balancer Subnet info which are needed
	// for the Avi Kubernetes Operator (AKO).
	VPCs []VPCInfo `json:"vpcs,omitempty"`
	// Conditions describe current state of VPCNetworkConfiguration.
	Conditions []Condition `json:"conditions,omitempty"`
}

// VPCInfo defines VPC info needed by tenant admin.
type VPCInfo struct {
	// VPC name.
	Name string `json:"name"`
	// AVISESubnetPath is the NSX Policy Path for the AVI SE Subnet.
	AVISESubnetPath string `json:"lbSubnetPath,omitempty"`
	// NSXLoadBalancerPath is the NSX Policy path for the NSX Load Balancer.
	NSXLoadBalancerPath string `json:"nsxLoadBalancerPath,omitempty"`
	// NSX Policy path for VPC.
	VPCPath string `json:"vpcPath"`
}

// +genclient
// +genclient:nonNamespaced
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion

// VPCNetworkConfiguration is the Schema for the vpcnetworkconfigurations API.
// +kubebuilder:resource:scope="Cluster"
// +kubebuilder:printcolumn:name="VPCPath",type=string,JSONPath=`.status.vpcs[0].vpcPath`,description="NSX VPC path the Namespace is associated with"
type VPCNetworkConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VPCNetworkConfigurationSpec   `json:"spec,omitempty"`
	Status VPCNetworkConfigurationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// VPCNetworkConfigurationList contains a list of VPCNetworkConfiguration.
type VPCNetworkConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VPCNetworkConfiguration `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VPCNetworkConfiguration{}, &VPCNetworkConfigurationList{})
}
