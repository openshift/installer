// Copyright (c) 2020-2025 Broadcom. All Rights Reserved.
// Broadcom Confidential. The term "Broadcom" refers to Broadcom Inc.
// and/or its subsidiaries.

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type VSphereDistributedNetworkConditionType string

const (
	// VSphereDistributedNetworkPortGroupFailure is added when PortGroupID specified either doesn't exist, or
	// there was an error in communicating with vCenter Server.
	VSphereDistributedNetworkPortGroupFailure VSphereDistributedNetworkConditionType = "PortGroupFailure"
	// VSphereDistributedNetworkIPPoolInvalid is added when no valid IPPool references exists.
	VSphereDistributedNetworkIPPoolInvalid VSphereDistributedNetworkConditionType = "IPPoolInvalid"
	// VsphereDistributedNetworkIPPoolPressure condition status is set to True when IPPool is low on free IPs.
	VsphereDistributedNetworkIPPoolPressure VSphereDistributedNetworkConditionType = "IPPoolPressure"
)

type IPAssignmentModeType string

const (
	// IPAssignmentModeDHCP indicates IP address is assigned dynamically using DHCP.
	IPAssignmentModeDHCP IPAssignmentModeType = "dhcp"
	// IPAssignmentModeStaticPool indicates IP address is assigned from a static pool of IP addresses.
	IPAssignmentModeStaticPool IPAssignmentModeType = "staticpool"
	// IPAssignmentModeNone indicates that no IP assignment will be performed.
	// The operator will not assign an IP and no DHCP client will be configured.
	IPAssignmentModeNone IPAssignmentModeType = "none"
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
	// Provides a timestamp for when the VSphereDistributedNetwork object last transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" patchStrategy:"replace"`
}

// VSphereDistributedNetworkSpec defines the desired state of VSphereDistributedNetwork.
type VSphereDistributedNetworkSpec struct {
	// PortGroupID is an existing vSphere Distributed PortGroup identifier.
	PortGroupID string `json:"portGroupID"`

	// IPAssignmentMode to use for network interfaces. If unset, defaults to IPAssignmentModeStaticPool.
	// For IPAssignmentModeDHCP and IPAssignmentModeNone, the IPPools, Gateway and SubnetMask
	// fields should be empty/unset. When using IPAssignmentModeNone, no IP will be assigned
	// and no DHCP client will be configured.
	// +optional
	IPAssignmentMode IPAssignmentModeType `json:"ipAssignmentMode,omitempty"`

	// IPPools references list of IPPool objects. This field should only be set when using
	// IPAssignmentModeStaticPool. For all other modes (IPAssignmentModeDHCP, IPAssignmentModeNone), this should be set
	// to an empty list.
	IPPools []IPPoolReference `json:"ipPools"`

	// Gateway setting to use for network interfaces. This field should only be set when using
	// IPAssignmentModeStaticPool. For all other modes (IPAssignmentModeDHCP, IPAssignmentModeNone), this should be set
	// to an empty string.
	Gateway string `json:"gateway"`

	// SubnetMask setting to use for network interfaces. This field should only be set when using
	// IPAssignmentModeStaticPool. For all other modes (IPAssignmentModeDHCP, IPAssignmentModeNone), this should be set
	// to an empty string.
	SubnetMask string `json:"subnetMask"`
}

// VLANType represents the type of VLAN configuration
type VLANType string

const (
	// VLANTypeStandard represents a standard VLAN configuration with a single VLAN ID
	VLANTypeStandard VLANType = "standard"
	// VLANTypeTrunk represents a VLAN trunk configuration that allows multiple VLANs
	VLANTypeTrunk VLANType = "trunk"
	// VLANTypePrivate represents a private VLAN configuration
	VLANTypePrivate VLANType = "private"
)

// VLANTrunkRange represents a range of VLAN IDs for trunk configuration
type VLANTrunkRange struct {
	// Start represents the beginning of the VLAN ID range (inclusive).
	Start int32 `json:"start"`

	// End represents the end of the VLAN ID range (inclusive).
	End int32 `json:"end"`
}

// VlanSpec represents the VLAN configuration.
type VlanSpec struct {
	// Type indicates the type of VLAN configuration (standard, trunk, or private).
	Type VLANType `json:"type"`

	// VlanID specifies the VLAN ID when Type is VLANTypeStandard.
	// This field is ignored for other VLAN types.
	// Possible values:
	// - A value of 0 indicates there is no VLAN configuration for the port.
	// - A value from 1 to 4094 specifies a VLAN ID for the port.
	// +optional
	VlanID *int32 `json:"vlanID,omitempty"`

	// TrunkRange specifies the ranges of allowed VLANs when Type is VLANTypeTrunk.
	// This field is ignored for other VLAN types.
	// Each range's Start and End values must be between 0 and 4094 inclusive.
	// Overlapping ranges are allowed.
	// +optional
	TrunkRange []VLANTrunkRange `json:"trunkRange,omitempty"`

	// PrivateVlanID specifies the private VLAN ID when Type is VLANTypePrivate.
	// This field is ignored for other VLAN types.
	// +optional
	PrivateVlanID *int32 `json:"privateVlanID,omitempty"`
}

// VSphereDistributedPortConfig represents the port-level configuration for a vSphere Distributed Network
type VSphereDistributedPortConfig struct {
	// Vlan represents the VLAN configuration for this port.
	// If unset, indicates that no VLAN configuration has been retrieved yet for this port.
	// +optional
	Vlan *VlanSpec `json:"vlan,omitempty"`
}

// VSphereDistributedNetworkStatus defines the observed state of VSphereDistributedNetwork.
type VSphereDistributedNetworkStatus struct {
	// Conditions is an array of current observed vSphere Distributed network conditions.
	Conditions []VSphereDistributedNetworkCondition `json:"conditions,omitempty"`

	// DefaultPortConfig represents the default port-level configuration that applies to all ports
	// unless overridden at the individual port level.
	// +optional
	DefaultPortConfig *VSphereDistributedPortConfig `json:"defaultPortConfig,omitempty"`
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
