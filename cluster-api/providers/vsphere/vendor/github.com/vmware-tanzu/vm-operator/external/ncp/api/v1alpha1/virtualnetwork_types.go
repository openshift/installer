/* **********************************************************
 * Copyright 2018 VMware, Inc.  All rights reserved. -- VMware Confidential
 * **********************************************************/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VirtualNetwork describe a vnet resource
// +k8s:openapi-gen=true
type VirtualNetwork struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualNetworkSpec   `json:"spec,omitempty"`
	Status VirtualNetworkStatus `json:"status,omitempty"`
}

// VirtualNetworkSpec defines the desired state of VirtualMachineClass
type VirtualNetworkSpec struct {
	WhitelistSourceRanges string `json:"whitelist_source_ranges,omitempty"`
}

// VirtualNetworkStatus defines the observed state of VirtualMachineClass
type VirtualNetworkStatus struct {
	Conditions    []VirtualNetworkCondition `json:"conditions,omitempty"`
	DefaultSNATIP string                    `json:"defaultSNATIP,omitempty"`
}

// VirtualNetworkCondition defines the condition for the VirtualNetwork
type VirtualNetworkCondition struct {
	Status  string `json:"status,omitempty"`
	Type    string `json:"type,omitempty"`
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type VirtualNetworkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualNetwork `json:"items"`
}
