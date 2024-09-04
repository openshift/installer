// Copyright (c) 2024 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// FoundationLoadBalancerConditionReady is added when all load balancer settings have been updated
	// and the load balancer is ready to be used.
	FoundationLoadBalancerConditionReady FoundationLoadBalancerConditionType = "Ready"
	// FoundationLoadBalancerConditionFailure is added when load balancer provider returns an error.
	FoundationLoadBalancerConditionFailure FoundationLoadBalancerConditionType = "Failure"

	FoundationLoadBalancerSizeSmall  FoundationLoadBalancerSize = "small"
	FoundationLoadBalancerSizeMedium FoundationLoadBalancerSize = "medium"
	FoundationLoadBalancerSizeLarge  FoundationLoadBalancerSize = "large"
	FoundationLoadBalancerSizeXL     FoundationLoadBalancerSize = "xlarge"

	FoundationAvailabilityModeActivePassive FoundationLoadBalancerAvailabilityMode = "active-passive"
	FoundationAvailabilityModeSingleNode    FoundationLoadBalancerAvailabilityMode = "single-node"
)

type FoundationLoadBalancerConditionType string
type FoundationLoadBalancerTopologyType string
type FoundationLoadBalancerSize string
type FoundationLoadBalancerAvailabilityMode string

// Spec objects. Input for FLB deployment.

// FoundationLoadBalancerDeploymentSpec describes how to deploy the load balancer.
type FoundationLoadBalancerDeploymentSpec struct {

	// Size describes the node form factor.
	//
	// +kubebuilder:validation:Enum=small;medium;large;xlarge
	// +kubebuilder:default=small
	Size FoundationLoadBalancerSize `json:"size"`

	// StoragePolicy is a vSphere Storage Policy ID which defines node storage placement.
	StoragePolicy string `json:"storagePolicy"`

	// Version number desired by the operator.
	//
	// Defaults to the latest available.
	//
	// +optional
	Version string `json:"version,omitempty"`

	// Zones contains the names of zones eligible for placing nodes. Zones must be one of the
	// AvailabilityZones defined and eligible for placement on the cluster.
	//
	// If no zones are provided, you must provide a PlacementSpec.
	//
	// +optional
	Zones []string `json:"zones,omitempty"`

	// AvailabilityMode defines how the availability of the solution is deployed and configured.
	// +kubebuilder:validation:Enum=active-passive;single-node
	// +kubebuilder:default=active-passive
	AvailabilityMode FoundationLoadBalancerAvailabilityMode `json:"availabilityMode"`

	// ActivePassiveAvailabilityMode configures the load balancer in active-passive configuration.
	// Active-passive configuration consists of a two node deployment with one node configured to
	// actively service traffic with the second node in standby mode. When the service detects the
	// active node is unhealthy, traffic will be moved to the passive node after a short delay.
	// Connections may be dropped on fail-over.
	//
	// +optional
	ActivePassiveAvailabilityMode *ActivePassiveAvailabilityMode `json:"activePassiveSpec,omitempty"`

	// SingleNodeAvailabilityMode deploys a single node to serve load balancer traffic. If the node
	// fails, the service will attempt to redeploy it, but redeployment is best-effort and depends on
	// the health of the underlying infrastructure. You must select
	//
	// +optional
	SingleNodeAvailabilityMode *SingleModeAvailabilityMode `json:"singleNodeSpec,omitempty"`

	// PlacementSpec is optional configuration defining custom placement of load balancer nodes.
	//
	// If Zones are specified, this field is ignored.
	// If Zones are not specified, this field must be set.
	//
	// +optional
	PlacementSpec []CustomPlacementSpec `json:"placementSpec,omitempty"`
}

// ActivePassiveAvailabilityMode deploys two nodes in Active-Passive mode where one node is set into
// active state and is responsible for serving traffic, and one node is passive -
// awaiting a fail-over event. When a fail-over occurs, connections to and from the load balancer
// may be reset.
type ActivePassiveAvailabilityMode struct {
	// Replicas describes the total number of deployed nodes. Defaults to 2.
	//
	// +kubebuilder:validation:Maximum=2
	// +kubebuilder:default=2
	Replicas uint32 `json:"replicas"`
}

// SingleModeAvailabilityMode defines single node configuration. Single node configuration involves
// trading availability in return for reduced resource consumption. Upon node failure, redeployment will
// be attempted on a best-effort basis.
type SingleModeAvailabilityMode struct {
	// Replicas describes the total number of deployed nodes. Defaults to 1.
	//
	// +kubebuilder:validation:Maximum=1
	// +kubebuilder:default=1
	Replicas uint32 `json:"replicas"`
}

// CustomPlacementSpec defines specific configurations for placing load balancer nodes.
type CustomPlacementSpec struct {

	// Cluster is the Managed Object ID of a vSphere ClusterComputeResource for placement outside a Supervisor.
	Cluster string `json:"cluster"`

	// ResourcePool is the Managed Object ID of a vSphere ResourcePool for placement outside a Supervisor.
	ResourcePool string `json:"resourcePool"`

	// Folder is the Managed Object ID of a vSphere Folder for placement outside a Supervisor.
	// Defaults to the Namespaces folder created on the cluster.
	//
	// +optional
	Folder string `json:"folder,omitempty"`
}

// Status objects. Specs are realized into Statuses.

// FoundationLoadBalancerNodeStatus describes the per-node status of the load balancer.
type FoundationLoadBalancerNodeStatus struct {
	// NodeID is a node's unique identifier.
	NodeID string `json:"nodeID"`
	// ManagementNetworkInterface defines the management NetworkInterface if it exists.
	//
	// +optional
	ManagementNetworkInterface NetworkInterfaceReference `json:"managementNetworkInterface,omitempty"`
	// WorkloadNetworkInterface defines the workload NetworkInterfaces if they exist.
	//
	// +optional
	WorkloadNetworkInterfaces []NetworkInterfaceReference `json:"workloadNetworkInterfaces,omitempty"`
	// VIPNetworkInterface is the interface bound to the Virtual IP Network.
	VIPNetworkInterface NetworkInterfaceReference `json:"vipNetworkInterface"`
}

// FoundationLoadBalancerConfigStatus describes the observed state of the Foundation Load Balancer.
type FoundationLoadBalancerConfigStatus struct {
	// Nodes list specific information about each deployed node.
	//
	// +optional
	Nodes []FoundationLoadBalancerNodeStatus `json:"nodes,omitempty"`
	// Conditions describes states of the load balancer at specific points in time.
	//
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// FoundationLoadBalancerConfigSpec defines the configuration for a vSphere Foundation Load Balancer.
// This specification is used to configure the resources for the load balancer on vCenter Server.
type FoundationLoadBalancerConfigSpec struct {
	DeploymentSpec FoundationLoadBalancerDeploymentSpec `json:"deploymentSpec"`
	// ManagementNetwork points to the Network used to program node management network interfaces.
	ManagementNetwork NetworkReference `json:"managementNetwork"`
	// WorkloadNetwork points to the Network used to program node workload network interfaces.
	//
	// If unset, workload data traffic will be routed out of the same NIF bound to VirtualIPNetwork.
	//
	// +kubebuilder:validation:MaxItems:=1
	// +optional
	WorkloadNetworks []NetworkReference `json:"workloadNetworks,omitempty"`
	// VirtualIPNetwork points to the Network used to program node VIP network interfaces.
	VirtualIPNetwork NetworkReference `json:"virtualIPNetwork"`

	// NetworkSpec contains values for configuring networks on the load balancer.
	// If unset, default settings will be applied.
	//
	// +optional
	NetworkSpec FoundationLoadBalancerNetworkConfigSpec `json:"networkSpec,omitempty"`
}

// FoundationLoadBalancerNetworkConfigSpec contains values for configuring networks on the load balancer.
type FoundationLoadBalancerNetworkConfigSpec struct {
	// StrictVIPPools indicates NetworkInterfaces provisioned by this load balancer should
	// only use vip-labeled IPPools for provisioning load balancer IP addresses.
	//
	// If this option is toggled while IP addresses are already provisioned, then
	// the IP addresses of the Services may be changed if VIPs are re-provisioned.
	//
	// By default, load balancer IP addresses may be consumed from general IP pools with
	// preference given to vip-labeled pools if they are defined.
	//
	// +optional
	StrictVIPPools bool `json:"strictVIPPools,omitempty"`
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=flb

// FoundationLoadBalancerConfig is the Schema for the FoundationLoadBalancerConfig API
type FoundationLoadBalancerConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FoundationLoadBalancerConfigSpec   `json:"spec,omitempty"`
	Status FoundationLoadBalancerConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// FoundationLoadBalancerConfigList contains a list of FoundationLoadBalancerConfig.
type FoundationLoadBalancerConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []FoundationLoadBalancerConfig `json:"items"`
}

func init() {
	RegisterTypeWithScheme(&FoundationLoadBalancerConfig{}, &FoundationLoadBalancerConfigList{})
}
