// Copyright (c) 2024 Broadcom. All Rights Reserved.
// Broadcom Confidential. The term "Broadcom" refers to Broadcom Inc.
// and/or its subsidiaries.

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// FoundationLoadBalancerConditionHealthy reflects the health status of the load balancer data-plane's runtime.
	FoundationLoadBalancerConditionHealthy FoundationLoadBalancerConditionType = "Healthy"
	// FoundationLoadBalancerConditionDeploymentStatusReady reflects the deployment status of the load balancer node(s).
	FoundationLoadBalancerConditionDeploymentStatusReady FoundationLoadBalancerConditionType = "DeploymentStatusReady"
	// FoundationLoadBalancerConditionOperationStatusReady reflects the operation status of the load balancer instance.
	FoundationLoadBalancerConditionOperationStatusReady FoundationLoadBalancerConditionType = "OperationStatusReady"

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
	// +kubebuilder:default:=small
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
	Zones []string `json:"zones"`

	// AvailabilityMode defines how the availability of the solution is deployed and configured.
	// +kubebuilder:validation:Enum=active-passive;single-node
	// +kubebuilder:default:=active-passive
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
	SingleNodeAvailabilityMode *SingleNodeAvailabilityMode `json:"singleNodeSpec,omitempty"`
}

// ActivePassiveAvailabilityMode deploys two nodes in Active-Passive mode where one node is set into
// active state and is responsible for serving traffic, and one node is passive -
// awaiting a fail-over event. When a fail-over occurs, connections to and from the load balancer
// may be reset.
type ActivePassiveAvailabilityMode struct {
	// Replicas describes the total number of deployed nodes. Defaults to 2.
	//
	// +kubebuilder:validation:Maximum=2
	// +kubebuilder:default:=2
	Replicas uint32 `json:"replicas"`
}

// SingleNodeAvailabilityMode defines single node configuration. Single node configuration involves
// trading availability in return for reduced resource consumption. Upon node failure, redeployment will
// be attempted on a best-effort basis.
type SingleNodeAvailabilityMode struct {
	// Replicas describes the total number of deployed nodes. Defaults to 1.
	//
	// +kubebuilder:validation:Maximum=1
	// +kubebuilder:default:=1
	Replicas uint32 `json:"replicas"`
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
	// Version describes the current version of the Foundation Load Balancer.
	//
	// +optional
	Version string `json:"version,omitempty"`

	// Nodes list specific information about each deployed node.
	//
	// +optional
	Nodes []FoundationLoadBalancerNodeStatus `json:"nodes,omitempty"`

	// VirtualServerIPPoolsUtilization describes the current states of virtual server IP addresses utilization.
	//
	// +optional
	VirtualServerIPPoolsUtilization VirtualIPPoolsUtilization `json:"virtualServerIPPoolsUtilization,omitempty"`

	// Conditions describes states of the load balancer at specific points in time.
	//
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// VirtualIPPoolsUtilization defines the IP addresses utilization for virtual IPPools resource.
type VirtualIPPoolsUtilization struct {
	// IPsAllocated represents the total number of virtual IP addresses currently allocated to services.
	//
	// +optional
	IPsAllocated int64 `json:"ipsAllocated,omitempty"`

	// IPsAvailable represents the total number of virtual IP addresses eligible to be used for services.
	//
	// +optional
	IPsAvailable int64 `json:"ipsAvailable,omitempty"`
}

// FoundationLoadBalancerConfigSpec defines the configuration for a vSphere Foundation Load Balancer.
// This specification is used to configure the resources for the load balancer on vCenter Server.
type FoundationLoadBalancerConfigSpec struct {
	// DeploymentSpec describes sizing and placement constraints of the load balancer.
	DeploymentSpec FoundationLoadBalancerDeploymentSpec `json:"deploymentSpec"`

	// ManagementNetwork points to the Network used to program node management network interfaces.
	//
	// If unset, the VirtualIPNetwork will be used for management traffic.
	//
	// +optional
	ManagementNetwork *NetworkReference `json:"managementNetwork,omitempty"`

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
	// VirtualServerIPPools are the list of IPPools that are
	// used for load balancer IP addresses.
	VirtualServerIPPools []IPPoolReference `json:"virtualServerIPPools"`

	// VirtualServerSubnets are the list of subnets specified in CIDR notation
	// that are directly connected to the VirtualIPNetwork.
	//
	// The VirtualServerIPPools must fall within the subnet of the VirtualIPNetwork
	// or one of these subnets.
	//
	// +kubebuilder:default:={}
	// +optional
	VirtualServerSubnets []string `json:"virtualServerSubnets"`

	// DNSServers is the list of servers used for DNS traffic.
	// These servers must be reachable from the network configured
	// for management traffic.
	//
	// +kubebuilder:default:={}
	// +optional
	DNSServers []string `json:"dnsServers"`

	// DNSSearchDomains are the domains resolvable on the specified DNSServers.
	//
	// +kubebuilder:default:={}
	// +optional
	DNSSearchDomains []string `json:"dnsSearchDomains"`

	// NTPServers are the servers used to sync time across nodes.
	// These servers must be reachable from the network configured
	// for management traffic.
	//
	// +kubebuilder:default:={}
	// +optional
	NTPServers []string `json:"ntpServers"`

	// SyslogEndpoint configures the syslog server. It accepts a protocol, host and port.
	// If using TLS, you must configure a TLS CA that is capable of verifying the endpoint certificate.
	// E.g. [protocol://]host[:port]
	// This server must be reachable from the network configured for management traffic.
	//
	// If empty, data will be logged locally to load balancer nodes.
	// Defaults to port 514 if using UDP and 6514 if using TLS.
	//
	// +optional
	SyslogEndpoint string `json:"syslogEndpoint,omitempty"`

	// SyslogCertificateSecretName is the certificate required to verify
	// the TLS syslog endpoint in PEM format.
	//
	// +optional
	SyslogCertificate string `json:"syslogCertificate,omitempty"`
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

func (flb *FoundationLoadBalancerConfig) GetConditions() []metav1.Condition {
	return flb.Status.Conditions
}

func (flb *FoundationLoadBalancerConfig) SetConditions(conditions []metav1.Condition) {
	flb.Status.Conditions = conditions
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
