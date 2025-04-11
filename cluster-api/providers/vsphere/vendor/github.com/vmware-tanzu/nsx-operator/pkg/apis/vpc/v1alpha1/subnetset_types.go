/* Copyright © 2022-2023 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: Apache-2.0 */

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SubnetSetSpec defines the desired state of SubnetSet.
type SubnetSetSpec struct {
	// Size of Subnet based upon estimated workload count.
	// +kubebuilder:validation:Maximum:=65536
	// +kubebuilder:validation:Minimum:=16
	IPv4SubnetSize int `json:"ipv4SubnetSize,omitempty"`
	// Access mode of Subnet, accessible only from within VPC or from outside VPC.
	// +kubebuilder:validation:Enum=Private;Public;PrivateTGW
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	AccessMode AccessMode `json:"accessMode,omitempty"`
	// DHCPConfig DHCP configuration.
	DHCPConfig DHCPConfig `json:"DHCPConfig,omitempty"`
}

// SubnetInfo defines the observed state of a single Subnet of a SubnetSet.
type SubnetInfo struct {
	NetworkAddresses    []string `json:"networkAddresses,omitempty"`
	GatewayAddresses    []string `json:"gatewayAddresses,omitempty"`
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
// +kubebuilder:printcolumn:name="AccessMode",type=string,JSONPath=`.spec.accessMode`,description="Access mode of Subnet"
// +kubebuilder:printcolumn:name="IPv4SubnetSize",type=string,JSONPath=`.spec.ipv4SubnetSize`,description="Size of Subnet"
// +kubebuilder:printcolumn:name="NetworkAddresses",type=string,JSONPath=`.status.subnets[*].networkAddresses[*]`,description="CIDRs for the SubnetSet"
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
