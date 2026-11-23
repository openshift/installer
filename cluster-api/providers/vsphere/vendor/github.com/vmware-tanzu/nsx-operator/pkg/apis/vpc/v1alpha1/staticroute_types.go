/* Copyright © 2022-2023 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: Apache-2.0 */

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type StaticRouteStatusCondition string

// StaticRouteCondition defines condition of StaticRoute.
type StaticRouteCondition Condition

// StaticRouteSpec defines static routes configuration on VPC.
// +kubebuilder:validation:XValidation:rule="!has(self.network) || !has(self.networkIpAllocation)",message="spec.network and spec.networkIpAllocation are mutually exclusive"
type StaticRouteSpec struct {
	// Specify network address in CIDR format.
	// Mutually exclusive with networkIpAllocation.
	// +kubebuilder:validation:Format=cidr
	// +optional
	Network string `json:"network,omitempty"`
	// Specify the name of an IPAddressAllocation CR whose allocated CIDR is used as
	// the static route network. Mutually exclusive with network.
	// +optional
	NetworkIPAllocation string `json:"networkIpAllocation,omitempty"`
	// Next hop gateway
	// +kubebuilder:validation:MinItems=1
	NextHops []NextHop `json:"nextHops"`
}

// NextHop defines next hop configuration for network.
type NextHop struct {
	// Next hop gateway IP address.
	// +kubebuilder:validation:Format=ip
	IPAddress string `json:"ipAddress"`
}

// StaticRouteStatus defines the observed state of StaticRoute.
type StaticRouteStatus struct {
	Conditions []StaticRouteCondition `json:"conditions"`
}

// +genclient
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:storageversion

// StaticRoute is the Schema for the staticroutes API.
// +kubebuilder:printcolumn:name="Network",type=string,JSONPath=`.spec.network`,description="Network in CIDR format"
// +kubebuilder:printcolumn:name="NetworkIPAllocation",type=string,JSONPath=`.spec.networkIpAllocation`,description="IPAddressAllocation CR name"
// +kubebuilder:printcolumn:name="NextHops",type=string,JSONPath=`.spec.nextHops[*].ipAddress`,description="Next Hops"
type StaticRoute struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StaticRouteSpec   `json:"spec,omitempty"`
	Status StaticRouteStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// StaticRouteList contains a list of StaticRoute.
type StaticRouteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StaticRoute `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StaticRoute{}, &StaticRouteList{})
}
