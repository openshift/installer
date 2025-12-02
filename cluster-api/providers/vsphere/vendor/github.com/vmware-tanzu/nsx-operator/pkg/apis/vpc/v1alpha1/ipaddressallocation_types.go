/* Copyright © 2024 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: Apache-2.0 */

package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type IPAddressVisibility string

var (
	IPAddressVisibilityExternal   IPAddressVisibility = "External"
	IPAddressVisibilityPrivate    IPAddressVisibility = "Private"
	IPAddressVisibilityPrivateTGW IPAddressVisibility = "PrivateTGW"
)

// +genclient
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:storageversion

// IPAddressAllocation is the Schema for the IP allocation API.
// +kubebuilder:printcolumn:name="IPAddressBlockVisibility",type=string,JSONPath=`.spec.ipAddressBlockVisibility`,description="IPAddressBlockVisibility of IPAddressAllocation"
// +kubebuilder:printcolumn:name="AllocationIPs",type=string,JSONPath=`.status.allocationIPs`, description="AllocationIPs for the IPAddressAllocation"
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
type IPAddressAllocationSpec struct {
	// IPAddressBlockVisibility specifies the visibility of the IPBlocks to allocate IP addresses. Can be External, Private or PrivateTGW.
	// +kubebuilder:validation:Enum=External;Private;PrivateTGW
	// +kubebuilder:default=Private
	// +optional
	IPAddressBlockVisibility IPAddressVisibility `json:"ipAddressBlockVisibility,omitempty"`
	// AllocationSize specifies the size of allocationIPs to be allocated.
	AllocationSize int `json:"allocationSize,omitempty"`
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
