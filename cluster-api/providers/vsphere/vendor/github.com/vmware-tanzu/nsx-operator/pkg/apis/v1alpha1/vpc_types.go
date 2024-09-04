/* Copyright © 2022-2023 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: Apache-2.0 */

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:storageversion

// VPC is the Schema for the VPC API
// +kubebuilder:printcolumn:name="PrivateIPv4CIDRs",type=string,JSONPath=`.status.privateIPv4CIDRs`,description="Private IPv4 CIDRs"
// +kubebuilder:printcolumn:name="SNATIP",type=string,JSONPath=`.status.defaultSNATIP`,description="Default SNAT IP for Private Subnets"
// +kubebuilder:printcolumn:name="LBSubnetCIDR",type=string,JSONPath=`.status.lbSubnetCIDR`,description="CIDR for the load balancer Subnet"
type VPC struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VPCSpec   `json:"spec,omitempty"`
	Status VPCStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// VPCList contains a list of VPC
type VPCList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VPC `json:"items"`
}

// VPCSpec defines VPC configuration
type VPCSpec struct {
}

// VPCStatus defines the observed state of VPC
type VPCStatus struct {
	Conditions []Condition `json:"conditions"`
	// NSX VPC Policy API resource path.
	NSXResourcePath string `json:"nsxResourcePath"`
	// Default SNAT IP for Private Subnets.
	DefaultSNATIP string `json:"defaultSNATIP"`
	// NSX PolicyPath for the load balancer Subnet.
	LBSubnetPath string `json:"lbSubnetPath"`
	// CIDR for the load balancer Subnet.
	LBSubnetCIDR string `json:"lbSubnetCIDR"`
        // Private CIDRs used for the VPC.
        PrivateIPv4CIDRs []string `json:"privateIPv4CIDRs"`
}

func init() {
	SchemeBuilder.Register(&VPC{}, &VPCList{})
}
