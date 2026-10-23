/* Copyright © 2024 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: Apache-2.0 */

package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type IPAddressVisibility string
type IPAllocationAddressType string

var (
	IPAddressVisibilityExternal   IPAddressVisibility     = "External"
	IPAddressVisibilityPrivate    IPAddressVisibility     = "Private"
	IPAddressVisibilityPrivateTGW IPAddressVisibility     = "PrivateTGW"
	IPAllocationIPAddressTypeIPv4 IPAllocationAddressType = "IPv4"
	IPAllocationIPAddressTypeIPv6 IPAllocationAddressType = "IPv6"
)

// +genclient
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:storageversion

// IPAddressAllocation is the Schema for the IP allocation API.
// +kubebuilder:printcolumn:name="IPAddressBlockVisibility",type=string,JSONPath=`.spec.ipAddressBlockVisibility`,description="IPAddressBlockVisibility of IPAddressAllocation"
// +kubebuilder:printcolumn:name="AllocationIPs",type=string,JSONPath=`.status.allocationIPs`, description="AllocationIPs for the IPAddressAllocation"
// +kubebuilder:printcolumn:name="IPv6AllocationPrefixLength",type=string,JSONPath=`.spec.ipv6AllocationPrefixLength`,description="IPv6 allocation prefix length of the IPAddressAllocation."
type IPAddressAllocation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   IPAddressAllocationSpec   `json:"spec"`
	Status IPAddressAllocationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// IPAddressAllocationList contains a list of IPAddressAllocation.
type IPAddressAllocationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IPAddressAllocation `json:"items"`
}

// IPAddressAllocationSpec defines the desired state of IPAddressAllocation.
// +kubebuilder:validation:XValidation:rule="!has(self.allocationSize) || !has(self.allocationIPs)", message="Only one of allocationSize or allocationIPs can be specified"
// +kubebuilder:validation:XValidation:rule="!has(self.ipv6AllocationPrefixLength) || !has(self.allocationIPs)", message="Only one of ipv6AllocationPrefixLength or allocationIPs can be specified"
// +kubebuilder:validation:XValidation:rule="!has(oldSelf.allocationSize) || has(self.allocationSize)", message="allocationSize is required once set"
// +kubebuilder:validation:XValidation:rule="!has(oldSelf.allocationIPs) || has(self.allocationIPs)", message="allocationIPs is required once set"
// +kubebuilder:validation:XValidation:rule="!has(oldSelf.ipv6AllocationPrefixLength) || has(self.ipv6AllocationPrefixLength)", message="ipv6AllocationPrefixLength is required once set"
// +kubebuilder:validation:XValidation:rule="!has(self.allocationSize) || !has(self.ipAddressType) || self.ipAddressType == 'IPv4'", message="allocationSize can only be set when ipAddressType is IPv4"
// +kubebuilder:validation:XValidation:rule="!has(self.ipv6AllocationPrefixLength) || self.ipAddressType == 'IPv6'", message="ipv6AllocationPrefixLength can only be set when ipAddressType is IPv6"
// +kubebuilder:validation:XValidation:rule="!has(self.ipAddressBlockVisibility) || !has(self.ipAddressType) || self.ipAddressType != 'IPv6'", message="ipAddressBlockVisibility cannot be set when ipAddressType is IPv6"
type IPAddressAllocationSpec struct {
	// IPAddressBlockVisibility specifies the visibility of the IPBlocks to allocate IP addresses. Can be External, Private or PrivateTGW.
	// This field is not applicable if ipAddressType is IPv6.
	// +kubebuilder:validation:Enum=External;Private;PrivateTGW
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	IPAddressBlockVisibility IPAddressVisibility `json:"ipAddressBlockVisibility,omitempty"`
	// AllocationSize specifies the size of IPv4 allocationIPs to be allocated.
	// It should be a power of 2.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	// +kubebuilder:validation:Minimum:=1
	AllocationSize int `json:"allocationSize,omitempty"`
	// AllocationIPs specifies the Allocated IP addresses in CIDR or single IP Address format.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	AllocationIPs string `json:"allocationIPs,omitempty"`
	// IPv6AllocationPrefixLength specifies the prefix length of IPv6 addresses.
	// Defaults to 64 when ipAddressType is IPv6 and this field is not specified.
	// +kubebuilder:validation:Minimum:=64
	// +kubebuilder:validation:Maximum:=128
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	IPv6AllocationPrefixLength int `json:"ipv6AllocationPrefixLength,omitempty"`
	// IPAddressType specifies the IP address type of the IPAddressAllocation.
	// +kubebuilder:validation:Enum=IPv4;IPv6
	// +kubebuilder:default=IPv4
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	IPAddressType IPAllocationAddressType `json:"ipAddressType,omitempty"`
}

// IPAddressAllocationStatus defines the observed state of IPAddressAllocation.
type IPAddressAllocationStatus struct {
	// AllocationIPs is the allocated IP addresses
	AllocationIPs string      `json:"allocationIPs"`
	Conditions    []Condition `json:"conditions,omitempty"`
}

func init() {
	SchemeBuilder.Register(&IPAddressAllocation{}, &IPAddressAllocationList{})
}
