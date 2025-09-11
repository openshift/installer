/* Copyright © 2022-2023 VMware, Inc. All Rights Reserved.
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
	// If VPC is set, only defaultIPv4SubnetSize and defaultSubnetAccessMode
	// take effect, other fields are ignored.
	// +optional
	VPC string `json:"vpc,omitempty"`

	// NSX Project the Namespace is associated with.
	NSXProject string `json:"nsxProject,omitempty"`

	// VPCConnectivityProfile ID. This profile has configuration related to creating VPC transit gateway attachment.
	VPCConnectivityProfile string `json:"vpcConnectivityProfile,omitempty"`

	// Private IPs.
	PrivateIPs []string `json:"privateIPs,omitempty"`

	// Default size of Subnets.
	// Defaults to 32.
	// +kubebuilder:default=32
	// +kubebuilder:validation:Maximum:=65536
	// +kubebuilder:validation:Minimum:=16
	DefaultSubnetSize int `json:"defaultSubnetSize,omitempty"`
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
