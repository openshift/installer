/* Copyright © 2024-2026 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: Apache-2.0 */

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NetworkStackType string

const (
	FullStackVPC  NetworkStackType = "FullStackVPC"
	VLANBackedVPC NetworkStackType = "VLANBackedVPC"
)

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:storageversion

// NetworkInfo is used to report the network information for a namespace.
// +kubebuilder:resource:path=networkinfos
type NetworkInfo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	VPCs []VPCState `json:"vpcs"`
	// AllowedDNSDomains specifies the permitted DNS domain names in the namespace.
	AllowedDNSDomains []string `json:"allowedDNSDomains,omitempty"`
}

// +kubebuilder:object:root=true

// NetworkInfoList contains a list of NetworkInfo.
type NetworkInfoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NetworkInfo `json:"items"`
}

// VPCState defines information for VPC.
type VPCState struct {
	// VPC name.
	Name string `json:"name"`
	// Default SNAT IP for Private Subnets.
	DefaultSNATIP string `json:"defaultSNATIP"`
	// LoadBalancerIPAddresses (AVI SE Subnet CIDR or NSX LB SNAT IPs).
	LoadBalancerIPAddresses string `json:"loadBalancerIPAddresses,omitempty"`
	// LoadBalancerBackendIPs is a list of IPv4 and/or IPv6 IP addresses
	// which the VPC load balancer uses to reach backend servers.
	LoadBalancerBackendIPs []string `json:"loadBalancerBackendIPs,omitempty"`
	// Private CIDRs used for the VPC.
	PrivateIPs []string `json:"privateIPs,omitempty"`
	// NetworkStack indicates the networking stack for the VPC.
	// Valid values: FullStackVPC, VLANBackedVPC
	// +kubebuilder:validation:Enum=FullStackVPC;VLANBackedVPC
	NetworkStack NetworkStackType `json:"networkStack,omitempty"`
}

func init() {
	SchemeBuilder.Register(&NetworkInfo{}, &NetworkInfoList{})
}
