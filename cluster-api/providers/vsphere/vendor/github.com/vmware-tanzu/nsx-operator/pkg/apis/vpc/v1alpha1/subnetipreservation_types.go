/* Copyright © 2025 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: Apache-2.0 */

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SubnetIPReservationSpec defines the desired state of SubnetIPReservation
// +kubebuilder:validation:XValidation:rule="!has(self.numberOfIPs) || !has(self.reservedIPs)",message="Only one of numberOfIPs or reservedIPs can be specified"
// +kubebuilder:validation:XValidation:rule="has(self.numberOfIPs) || has(self.reservedIPs)",message="One of numberOfIPs or reservedIPs must be specified"
// +kubebuilder:validation:XValidation:rule="!has(oldSelf.reservedIPs) || has(self.reservedIPs)",message="reservedIPs cannot be unset once set"
// +kubebuilder:validation:XValidation:rule="!has(oldSelf.numberOfIPs) || has(self.numberOfIPs)",message="numberOfIPs cannot be unset once set"
// +kubebuilder:validation:XValidation:rule="!has(oldSelf.ipAddressType) || has(self.ipAddressType)",message="ipAddressType cannot be unset once set"
// +kubebuilder:validation:XValidation:rule="!has(self.ipAddressType) || has(self.numberOfIPs)",message="ipAddressType can only be set when numberOfIPs is specified"
type SubnetIPReservationSpec struct {
	// Subnet specifies the Subnet to reserve IPs from.
	// The Subnet needs to have static IP allocation activated.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="subnet is immutable"
	Subnet string `json:"subnet"`

	// NumberOfIPs defines number of IPs requested to be reserved.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="numberOfIPs is immutable"
	// +kubebuilder:validation:Maximum:=100
	// +kubebuilder:validation:Minimum:=1
	NumberOfIPs int `json:"numberOfIPs,omitempty"`

	// ReservedIPs represents array of Reserved IPs. It can contain IP addresses,
	// IP Address range and CIDRs.
	// Supported formats include: ["192.168.1.1", "192.168.1.3-192.168.1.100", "192.168.2.0/28",
	// "2001:db8::1", "2001:db8::1-2001:db8::ff", "2001:db8::1/64"]
	// +kubebuilder:validation:MinItems=1
	ReservedIPs []string `json:"reservedIPs,omitempty"`

	// IPAddressType defines the IP address type of the SubnetIPReservation.
	// +kubebuilder:validation:Enum=IPv4;IPv6;IPv4IPv6
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	IPAddressType IPAddressType `json:"ipAddressType,omitempty"`
}

// SubnetIPReservationStatus defines the observed state of SubnetIPReservation
type SubnetIPReservationStatus struct {
	// Conditions described if the SubnetIPReservation is configured on NSX or not.
	// Condition type ""
	Conditions []Condition `json:"conditions,omitempty"`
	// List of reserved IPs.
	// Supported formats include: ["192.168.1.1", "192.168.1.3-192.168.1.100", "192.168.2.0/28"]
	IPs []string `json:"ips,omitempty"`
}

//+genclient
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:storageversion
//+kubebuilder:selectablefield:JSONPath=.spec.subnet

// SubnetIPReservation is the Schema for the subnetipreservations API
// +kubebuilder:printcolumn:name="Subnet",type=string,JSONPath=`.spec.subnet`,description="The parent Subnet name of the SubnetIPReservation."
// +kubebuilder:printcolumn:name="NumberOfIPs",type=string,JSONPath=`.spec.numberOfIPs`,description="Number of IPs requested to be reserved."
// +kubebuilder:printcolumn:name="IPs",type=string,JSONPath=`.status.ips[:]`,description="List of reserved IPs."
// +kubebuilder:printcolumn:name="IPAddressType",type=string,JSONPath=`.spec.ipAddressType`,description="IP address type of the SubnetIPReservation."
type SubnetIPReservation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +kubebuilder:validation:Required
	Spec   SubnetIPReservationSpec   `json:"spec"`
	Status SubnetIPReservationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SubnetIPReservationList contains a list of SubnetIPReservation
type SubnetIPReservationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SubnetIPReservation `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SubnetIPReservation{}, &SubnetIPReservationList{})
}
