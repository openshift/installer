// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type VSphereDistributedNetworkConditionType string

const (
	// VSphereDistributedNetworkFailure is added when PortGroupID specified either doesn't exist, or
	// there was an error in communicating with vCenter Server.
	VSphereDistributedNetworkPortGroupFailure VSphereDistributedNetworkConditionType = "PortGroupFailure"
	// VSphereDistributedNetworkIPPoolInvalid is added when no valid IPPool references exists.
	VSphereDistributedNetworkIPPoolInvalid VSphereDistributedNetworkConditionType = "IPPoolInvalid"
)

type IPAssignmentModeType string

const (
	// IPAssignmentModeDHCP indicates IP address is assigned dynamically using DHCP.
	IPAssignmentModeDHCP IPAssignmentModeType = "dhcp"
	// IPAssignmentModeStaticPool indicates IP address is assigned from a static pool of IP addresses.
	IPAssignmentModeStaticPool IPAssignmentModeType = "staticpool"
)

// VSphereDistributedNetworkCondition describes the state of a VSphereDistributedNetwork at a certain point.
type VSphereDistributedNetworkCondition struct {
	// Type is the type of VSphereDistributedNetwork condition.
	Type VSphereDistributedNetworkConditionType `json:"type"`
	// Status is the status of the condition.
	// Can be True, False, Unknown.
	Status corev1.ConditionStatus `json:"status"`
	// Machine understandable string that gives the reason for condition's last transition.
	Reason string `json:"reason,omitempty"`
	// Human-readable message indicating details about last transition.
	Message string `json:"message,omitempty"`
}

type IPPoolReference struct {
	// Name of the IPPool resource being referenced.
	Name string `json:"name"`
	// API version of the referent.
	APIVersion string `json:"apiVersion,omitempty"`
}

// VSphereDistributedNetworkSpec defines the desired state of VSphereDistributedNetwork.
type VSphereDistributedNetworkSpec struct {
	// PortGroupID is an existing vSphere Distributed PortGroup identifier.
	PortGroupID string `json:"portGroupID"`

	// IPAssignmentMode to use for network interfaces. If unset, defaults to IPAssignmentModeStaticPool.
	// In case of IPAssignmentModeDHCP, IPPools, Gateway and SubnetMask fields are ignored.
	// +optional
	IPAssignmentMode IPAssignmentModeType `json:"ipAssignmentMode,omitempty"`

	// IPPools references list of IPPool objects. This field should be set to empty list for
	// IPAssignmentModeDHCP IPAssignmentMode.
	IPPools []IPPoolReference `json:"ipPools"`

	// Gateway setting to use for network interfaces. This field should be set to empty string
	// for IPAssignmentModeDHCP IPAssignmentMode.
	Gateway string `json:"gateway"`

	// SubnetMask setting to use for network interfaces. This field should be set to empty string
	// for IPAssignmentModeDHCP IPAssignmentMode.
	SubnetMask string `json:"subnetMask"`
}

// VSphereDistributedNetworkStatus defines the observed state of VSphereDistributedNetwork.
type VSphereDistributedNetworkStatus struct {
	// Conditions is an array of current observed vSphere Distributed network conditions.
	Conditions []VSphereDistributedNetworkCondition `json:"conditions,omitempty"`
}

// +genclient
// +genclient:nonNamespaced
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster

// VSphereDistributedNetwork represents schema for a network backed by a vSphere Distributed PortGroup on vSphere
// Distributed switch.
type VSphereDistributedNetwork struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VSphereDistributedNetworkSpec   `json:"spec,omitempty"`
	Status VSphereDistributedNetworkStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// VSphereDistributedNetworkList contains a list of VSphereDistributedNetwork
type VSphereDistributedNetworkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VSphereDistributedNetwork `json:"items"`
}

func init() {
	RegisterTypeWithScheme(&VSphereDistributedNetwork{}, &VSphereDistributedNetworkList{})
}
