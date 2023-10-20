/* **********************************************************
 * Copyright 2019 VMware, Inc.  All rights reserved. -- VMware Confidential
 * **********************************************************/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VirtualNetworkInterface describe a vnetif resource
// +k8s:openapi-gen=true
type VirtualNetworkInterface struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualNetworkInterfaceSpec   `json:"spec,omitempty"`
	Status VirtualNetworkInterfaceStatus `json:"status,omitempty"`
}

// VirtualNetworkInterfaceSpec defines the desired state of VirtualNetworkInterface
type VirtualNetworkInterfaceSpec struct {
	VirtualNetwork string `json:"virtualNetwork,omitempty"`
}

// VirtualNetworkInterfaceStatus defines the observed state of the VirtualNetworkInterface
type VirtualNetworkInterfaceStatus struct {
	Conditions     []VirtualNetworkCondition              `json:"conditions,omitempty"`
	InterfaceID    string                                 `json:"interfaceID,omitempty"`
	IPAddresses    []VirtualNetworkInterfaceIP            `json:"ipAddresses,omitempty"`
	MacAddress     string                                 `json:"macAddress,omitempty"`
	ProviderStatus *VirtualNetworkInterfaceProviderStatus `json:"providerStatus,omitempty"`
}

// VirtualNetworkInterfaceIP defines the interface status
type VirtualNetworkInterfaceIP struct {
	Gateway    string `json:"gateway,omitempty"`
	IP         string `json:"ip,omitempty"`
	SubnetMask string `json:"subnetMask,omitempty"`
}

// VirtualNetworkInterfaceProviderStatus defines the nsx-t resource provider status
type VirtualNetworkInterfaceProviderStatus struct {
	NsxLogicalPortID   string `json:"nsxLogicalPortID,omitempty"`
	NsxLogicalSwitchID string `json:"nsxLogicalSwitchID,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type VirtualNetworkInterfaceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualNetworkInterface `json:"items"`
}
