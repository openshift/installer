/* Copyright © 2022-2023 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: Apache-2.0 */

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SubnetPortSpec defines the desired state of SubnetPort.
type SubnetPortSpec struct {
	// Subnet defines the parent Subnet name of the SubnetPort.
	Subnet string `json:"subnet,omitempty"`
	// SubnetSet defines the parent SubnetSet name of the SubnetPort.
	SubnetSet string `json:"subnetSet,omitempty"`
}

// SubnetPortStatus defines the observed state of SubnetPort.
type SubnetPortStatus struct {
	// Conditions describes current state of SubnetPort.
	Conditions []Condition `json:"conditions,omitempty"`
	// Subnet port attachment state.
	Attachment             PortAttachment         `json:"attachment,omitempty"`
	NetworkInterfaceConfig NetworkInterfaceConfig `json:"networkInterfaceConfig,omitempty"`
}

// VIF attachment state of a subnet port.
type PortAttachment struct {
	ID string `json:"id,omitempty"`
}

type NetworkInterfaceConfig struct {
	LogicalSwitchUUID string                      `json:"logicalSwitchUUID,omitempty"`
	IPAddresses       []NetworkInterfaceIPAddress `json:"ipAddresses,omitempty"`
	MACAddress        string                      `json:"macAddress,omitempty"`
}

type NetworkInterfaceIPAddress struct {
	// IP address string with the prefix.
	IPAddress string `json:"ipAddress,omitempty"`
	Gateway   string `json:"gateway,omitempty"`
}

// +genclient
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:storageversion

// SubnetPort is the Schema for the subnetports API.
// +kubebuilder:printcolumn:name="VIFID",type=string,JSONPath=`.status.attachment.id`,description="Attachment VIF ID owned by the SubnetPort."
// +kubebuilder:printcolumn:name="IPAddress",type=string,JSONPath=`.status.networkInterfaceConfig.ipAddresses[0].ipAddress`,description="IP address string with the prefix."
// +kubebuilder:printcolumn:name="MACAddress",type=string,JSONPath=`.status.networkInterfaceConfig.macAddress`,description="MAC Address of the SubnetPort."
type SubnetPort struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SubnetPortSpec   `json:"spec,omitempty"`
	Status SubnetPortStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SubnetPortList contains a list of SubnetPort.
type SubnetPortList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SubnetPort `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SubnetPort{}, &SubnetPortList{})
}
