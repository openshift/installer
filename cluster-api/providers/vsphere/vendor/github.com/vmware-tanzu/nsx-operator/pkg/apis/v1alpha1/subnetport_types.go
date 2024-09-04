/* Copyright © 2022-2023 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: Apache-2.0 */

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SubnetPortSpec defines the desired state of SubnetPort.
type SubnetPortSpec struct {
	// Subnet defines the parent Subnet name of the SubnetPort.
	Subnet string `json:"subnet,omitempty"`
	// SubnetSet defines the parent SubnetSet name of the SubnetPort.
	SubnetSet string `json:"subnetSet,omitempty"`
	// AttachmentRef refers to the virtual machine which the SubnetPort is attached.
	AttachmentRef corev1.ObjectReference `json:"attachmentRef,omitempty"`
}

// SubnetPortStatus defines the observed state of SubnetPort.
type SubnetPortStatus struct {
	// Conditions describes current state of SubnetPort.
	Conditions []Condition `json:"conditions,omitempty"`
	// VIFID describes the attachment VIF ID owned by the SubnetPort in NSX-T.
	VIFID string `json:"vifID,omitempty"`
	// IPAddresses describes the IP addresses of the SubnetPort.
	IPAddresses []SubnetPortIPAddress `json:"ipAddresses,omitempty"`
	// MACAddress describes the MAC address of the SubnetPort.
	MACAddress string `json:"macAddress,omitempty"`
	// LogicalSwitchID defines the logical switch ID in NSX-T.
	LogicalSwitchID string `json:"logicalSwitchID,omitempty"`
}

type SubnetPortIPAddress struct {
	Gateway string `json:"gateway,omitempty"`
	IP      string `json:"ip,omitempty"`
	Netmask string `json:"netmask,omitempty"`
}

// +genclient
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:storageversion

// SubnetPort is the Schema for the subnetports API.
// +kubebuilder:printcolumn:name="VIFID",type=string,JSONPath=`.status.vifID`,description="Attachment VIF ID owned by the SubnetPort"
// +kubebuilder:printcolumn:name="IPAddress",type=string,JSONPath=`.status.ipAddresses[0].ip`,description="IP Address of the SubnetPort"
// +kubebuilder:printcolumn:name="MACAddress",type=string,JSONPath=`.status.macAddress`,description="MAC Address of the SubnetPort"
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
