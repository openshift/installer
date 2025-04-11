/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta2

import "github.com/IBM/vpc-go-sdk/vpcv1"

const (
	// CIDRBlockAny is the CIDRBlock representing any allowable destination/source IP.
	CIDRBlockAny string = "0.0.0.0/0"

	// DefaultAPIServerPort is defuault API server port number.
	DefaultAPIServerPort int32 = 6443

	// UpdateMachineError indicates an error while trying to update a machine.
	UpdateMachineError string = "UpdateError"
)

// PowerVSInstanceState describes the state of an IBM Power VS instance.
type PowerVSInstanceState string

var (
	// PowerVSInstanceStateACTIVE is the string representing an instance in a ACTIVE state.
	PowerVSInstanceStateACTIVE = PowerVSInstanceState("ACTIVE")

	// PowerVSInstanceStateBUILD is the string representing an instance in a BUILD state.
	PowerVSInstanceStateBUILD = PowerVSInstanceState("BUILD")

	// PowerVSInstanceStateSHUTOFF is the string representing an instance in a SHUTOFF state.
	PowerVSInstanceStateSHUTOFF = PowerVSInstanceState("SHUTOFF")

	// PowerVSInstanceStateREBOOT is the string representing an instance in a REBOOT state.
	PowerVSInstanceStateREBOOT = PowerVSInstanceState("REBOOT")

	// PowerVSInstanceStateERROR is the string representing an instance in a ERROR state.
	PowerVSInstanceStateERROR = PowerVSInstanceState("ERROR")
)

// PowerVSImageState describes the state of an IBM Power VS image.
type PowerVSImageState string

var (
	// PowerVSImageStateACTIVE is the string representing an image in a active state.
	PowerVSImageStateACTIVE = PowerVSImageState("active")

	// PowerVSImageStateQue is the string representing an image in a queued state.
	PowerVSImageStateQue = PowerVSImageState("queued")

	// PowerVSImageStateFailed is the string representing an image in a failed state.
	PowerVSImageStateFailed = PowerVSImageState("failed")

	// PowerVSImageStateImporting is the string representing an image in a failed state.
	PowerVSImageStateImporting = PowerVSImageState("importing")
)

// ServiceInstanceState describes the state of a service instance.
type ServiceInstanceState string

var (
	// ServiceInstanceStateActive is the string representing a service instance in an active state.
	ServiceInstanceStateActive = ServiceInstanceState("active")

	// ServiceInstanceStateProvisioning is the string representing a service instance in a provisioning state.
	ServiceInstanceStateProvisioning = ServiceInstanceState("provisioning")

	// ServiceInstanceStateFailed is the string representing a service instance in a failed state.
	ServiceInstanceStateFailed = ServiceInstanceState("failed")

	// ServiceInstanceStateRemoved is the string representing a service instance in a removed state.
	ServiceInstanceStateRemoved = ServiceInstanceState("removed")
)

// TransitGatewayState describes the state of an IBM Transit Gateway.
type TransitGatewayState string

var (
	// TransitGatewayStateAvailable is the string representing a transit gateway in available state.
	TransitGatewayStateAvailable = TransitGatewayState("available")

	// TransitGatewayStatePending is the string representing a transit gateway in pending state.
	TransitGatewayStatePending = TransitGatewayState("pending")

	// TransitGatewayStateFailed is the string representing a transit gateway in failed state.
	TransitGatewayStateFailed = TransitGatewayState("failed")

	// TransitGatewayStateDeletePending is the string representing a transit gateway in deleting state.
	TransitGatewayStateDeletePending = TransitGatewayState("deleting")
)

// TransitGatewayConnectionState describes the state of an IBM Transit Gateway connection.
type TransitGatewayConnectionState string

var (
	// TransitGatewayConnectionStateAttached is the string representing a transit gateway connection in attached state.
	TransitGatewayConnectionStateAttached = TransitGatewayConnectionState("attached")

	// TransitGatewayConnectionStateFailed is the string representing a transit gateway connection in failed state.
	TransitGatewayConnectionStateFailed = TransitGatewayConnectionState("failed")

	// TransitGatewayConnectionStatePending is the string representing a transit gateway connection in pending state.
	TransitGatewayConnectionStatePending = TransitGatewayConnectionState("pending")

	// TransitGatewayConnectionStateDeleting is the string representing a transit gateway connection in deleting state.
	TransitGatewayConnectionStateDeleting = TransitGatewayConnectionState("deleting")
)

// VPCLoadBalancerBackendPoolAlgorithm describes the backend pool's load balancing algorithm.
// +kubebuilder:validation:Enum=least_connections;round_robin;weighted_round_robin
type VPCLoadBalancerBackendPoolAlgorithm string

var (
	// VPCLoadBalancerBackendPoolAlgorithmLeastConnections is the string representing the least_connections load balancing algorithm.
	VPCLoadBalancerBackendPoolAlgorithmLeastConnections VPCLoadBalancerBackendPoolAlgorithm = vpcv1.CreateLoadBalancerPoolOptionsAlgorithmLeastConnectionsConst

	// VPCLoadBalancerBackendPoolAlgorithmRoundRobin is the string representing the round_robin load balancing algorithm.
	VPCLoadBalancerBackendPoolAlgorithmRoundRobin VPCLoadBalancerBackendPoolAlgorithm = vpcv1.CreateLoadBalancerPoolOptionsAlgorithmRoundRobinConst

	// VPCLoadBalancerBackendPoolAlgorithmWeightedRoundRobin is the string representing the weighted_round_robin load balancing algorithm.
	VPCLoadBalancerBackendPoolAlgorithmWeightedRoundRobin VPCLoadBalancerBackendPoolAlgorithm = vpcv1.CreateLoadBalancerPoolOptionsAlgorithmWeightedRoundRobinConst
)

// VPCLoadBalancerBackendPoolProtocol describes the protocol for load balancer backend pools.
// We have unique types in case IBM Cloud Load Balancer Listener and Backend Pool supported algorithms ever diverage.
// +kubebuilder:validation:Enum=http;https;tcp;udp
type VPCLoadBalancerBackendPoolProtocol string

var (
	// VPCLoadBalancerBackendPoolProtocolHTTP is the string representing the http protocol for load balancer backend pools.
	VPCLoadBalancerBackendPoolProtocolHTTP VPCLoadBalancerBackendPoolProtocol = vpcv1.LoadBalancerPoolPrototypeLoadBalancerContextProtocolHTTPConst

	// VPCLoadBalancerBackendPoolProtocolHTTPS is the string representing the https protocol for load balancer backend pools.
	VPCLoadBalancerBackendPoolProtocolHTTPS VPCLoadBalancerBackendPoolProtocol = vpcv1.LoadBalancerPoolPrototypeLoadBalancerContextProtocolHTTPSConst

	// VPCLoadBalancerBackendPoolProtocolTCP is the string representing the tcp protocol for load balancer backend pools.
	VPCLoadBalancerBackendPoolProtocolTCP VPCLoadBalancerBackendPoolProtocol = vpcv1.LoadBalancerPoolPrototypeLoadBalancerContextProtocolTCPConst

	// VPCLoadBalancerBackendPoolProtocolUDP is the string representing the tudp protocol for load balancer backend pools.
	VPCLoadBalancerBackendPoolProtocolUDP VPCLoadBalancerBackendPoolProtocol = vpcv1.LoadBalancerPoolPrototypeLoadBalancerContextProtocolUDPConst
)

// VPCLoadBalancerListenerProtocol describes the protocol for load balancer listeners.
// We have unique types in case IBM Cloud Load Balancer Listener and Backend Pool supported algorithms ever diverage.
// +kubebuilder:validation:Enum=http;https;tcp;udp
type VPCLoadBalancerListenerProtocol string

var (
	// VPCLoadBalancerListenerProtocolHTTP is the string representing the http protocol for load balancer listeners.
	VPCLoadBalancerListenerProtocolHTTP VPCLoadBalancerListenerProtocol = vpcv1.LoadBalancerListenerProtocolHTTPConst

	// VPCLoadBalancerListenerProtocolHTTPS is the string representing the https protocol for load balancer listeners.
	VPCLoadBalancerListenerProtocolHTTPS VPCLoadBalancerListenerProtocol = vpcv1.LoadBalancerListenerProtocolHTTPSConst

	// VPCLoadBalancerListenerProtocolTCP is the string representing the tcp protocol for load balancer listeners.
	VPCLoadBalancerListenerProtocolTCP VPCLoadBalancerListenerProtocol = vpcv1.LoadBalancerListenerProtocolTCPConst

	// VPCLoadBalancerListenerProtocolUDP is the string representing the tudp protocol for load balancer listeners.
	VPCLoadBalancerListenerProtocolUDP VPCLoadBalancerListenerProtocol = vpcv1.LoadBalancerListenerProtocolUDPConst
)

// VPCLoadBalancerBackendPoolHealthMonitorType describes the backend pool's health check protocol type.
// +kubebuilder:validation:Enum=http;https;tcp
type VPCLoadBalancerBackendPoolHealthMonitorType string

var (
	// VPCLoadBalancerBackendPoolHealthMonitorTypeHTTP is the string representing the http health pool protocol type.
	VPCLoadBalancerBackendPoolHealthMonitorTypeHTTP VPCLoadBalancerBackendPoolHealthMonitorType = vpcv1.LoadBalancerPoolHealthMonitorTypeHTTPConst

	// VPCLoadBalancerBackendPoolHealthMonitorTypeHTTPS is the string representing the https health pool protocol type.
	VPCLoadBalancerBackendPoolHealthMonitorTypeHTTPS VPCLoadBalancerBackendPoolHealthMonitorType = vpcv1.LoadBalancerPoolHealthMonitorTypeHTTPSConst

	// VPCLoadBalancerBackendPoolHealthMonitorTypeTCP is the string representing the tcp health pool protocol type.
	VPCLoadBalancerBackendPoolHealthMonitorTypeTCP VPCLoadBalancerBackendPoolHealthMonitorType = vpcv1.LoadBalancerPoolHealthMonitorTypeTCPConst
)

// VPCLoadBalancerState describes the state of the load balancer.
type VPCLoadBalancerState string

var (
	// VPCLoadBalancerStateActive is the string representing the load balancer in a active state.
	VPCLoadBalancerStateActive = VPCLoadBalancerState("active")

	// VPCLoadBalancerStateCreatePending is the string representing the load balancer in a queued state.
	VPCLoadBalancerStateCreatePending = VPCLoadBalancerState("create_pending")

	// VPCLoadBalancerStateUpdatePending is the string representing the load balancer in updating state.
	VPCLoadBalancerStateUpdatePending = VPCLoadBalancerState("update_pending")

	// VPCLoadBalancerStateDeletePending is the string representing the load balancer in deleting state.
	VPCLoadBalancerStateDeletePending = VPCLoadBalancerState("delete_pending")
)

// VPCSubnetState describes the state of a VPC Subnet.
type VPCSubnetState string

var (
	// VPCSubnetStateDeleting is the string representing a VPC subnet in deleting state.
	VPCSubnetStateDeleting = VPCSubnetState("deleting")
)

// VPCState describes the state of a VPC.
type VPCState string

var (
	// VPCStatePending is the string representing a VPC in pending state.
	VPCStatePending = VPCState("pending")

	// VPCStateDeleting is the string representing a VPC in deleting state.
	VPCStateDeleting = VPCState("deleting")
)

// DHCPServerState describes the state of the DHCP Server.
type DHCPServerState string

var (
	// DHCPServerStateActive indicates the active state of DHCP server.
	DHCPServerStateActive = DHCPServerState("ACTIVE")

	// DHCPServerStateBuild indicates the build state of DHCP server.
	DHCPServerStateBuild = DHCPServerState("BUILD")

	// DHCPServerStateError indicates the error state of DHCP server.
	DHCPServerStateError = DHCPServerState("ERROR")
)

// DeletePolicy defines the policy used to identify images to be preserved.
type DeletePolicy string

var (
	// DeletePolicyRetain is the string representing an image to be retained.
	DeletePolicyRetain = DeletePolicy("retain")
)

// ResourceType describes IBM Cloud resource name.
type ResourceType string

var (
	// ResourceTypeServiceInstance is Power VS service instance resource.
	ResourceTypeServiceInstance = ResourceType("serviceInstance")
	// ResourceTypeNetwork is Power VS network resource.
	ResourceTypeNetwork = ResourceType("network")
	// ResourceTypeDHCPServer is Power VS DHCP server.
	ResourceTypeDHCPServer = ResourceType("dhcpServer")
	// ResourceTypeLoadBalancer VPC loadBalancer resource.
	ResourceTypeLoadBalancer = ResourceType("loadBalancer")
	// ResourceTypeLoadBalancerPool is a Load Balancer Pool resource.
	ResourceTypeLoadBalancerPool = ResourceType("loadBalancerPool")
	// ResourceTypeTransitGateway is transit gateway resource.
	ResourceTypeTransitGateway = ResourceType("transitGateway")
	// ResourceTypeVPC is Power VS network resource.
	ResourceTypeVPC = ResourceType("vpc")
	// ResourceTypeSubnet is VPC subnet resource.
	ResourceTypeSubnet = ResourceType("subnet")
	// ResourceTypeControlPlaneSubnet is a VPC subnet resource designated for the Control Plane.
	ResourceTypeControlPlaneSubnet = ResourceType("controlPlaneSubnet")
	// ResourceTypeWorkerSubnet is a VPC subnet resource designated for the Worker (Data) Plane.
	ResourceTypeWorkerSubnet = ResourceType("workerSubnet")
	// ResourceTypeSecurityGroup is a VPC Security Group resource.
	ResourceTypeSecurityGroup = ResourceType("securityGroup")
	// ResourceTypeCOSInstance is IBM COS instance resource.
	ResourceTypeCOSInstance = ResourceType("cosInstance")
	// ResourceTypeCOSBucket is IBM COS bucket resource.
	ResourceTypeCOSBucket = ResourceType("cosBucket")
	// ResourceTypeResourceGroup is IBM Resource Group.
	ResourceTypeResourceGroup = ResourceType("resourceGroup")
	// ResourceTypePublicGateway is a VPC Public Gatway.
	ResourceTypePublicGateway = ResourceType("publicGateway")
	// ResourceTypeCustomImage is a VPC Custom Image.
	ResourceTypeCustomImage = ResourceType("customImage")
)

const (
	// VPCSecurityGroupRuleProtocolAllType is a string representation of the 'SecurityGroupRuleSecurityGroupRuleProtocolAll' type.
	VPCSecurityGroupRuleProtocolAllType = "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolAll"

	// VPCSecurityGroupRuleProtocolIcmpType is a string representation of the 'SecurityGroupRuleSecurityGroupRuleProtocolIcmp' type.
	VPCSecurityGroupRuleProtocolIcmpType = "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp"

	// VPCSecurityGroupRuleProtocolTcpudpType is a string representation of the 'SecurityGroupRuleSecurityGroupRuleProtocolTcpudp' type.
	VPCSecurityGroupRuleProtocolTcpudpType = "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp"
)

// VPCSecurityGroupRuleAction represents the actions for a Security Group Rule.
// +kubebuilder:validation:Enum=allow;deny
type VPCSecurityGroupRuleAction string

const (
	// VPCSecurityGroupRuleActionAllow defines that the Rule should allow traffic.
	VPCSecurityGroupRuleActionAllow VPCSecurityGroupRuleAction = vpcv1.NetworkACLRuleActionAllowConst
	// VPCSecurityGroupRuleActionDeny defines that the Rule should deny traffic.
	VPCSecurityGroupRuleActionDeny VPCSecurityGroupRuleAction = vpcv1.NetworkACLRuleActionDenyConst
)

// VPCSecurityGroupRuleDirection represents the directions for a Security Group Rule.
// +kubebuilder:validation:Enum=inbound;outbound
type VPCSecurityGroupRuleDirection string

const (
	// VPCSecurityGroupRuleDirectionInbound defines the Rule is for inbound traffic.
	VPCSecurityGroupRuleDirectionInbound VPCSecurityGroupRuleDirection = vpcv1.NetworkACLRuleDirectionInboundConst
	// VPCSecurityGroupRuleDirectionOutbound defines the Rule is for outbound traffic.
	VPCSecurityGroupRuleDirectionOutbound VPCSecurityGroupRuleDirection = vpcv1.NetworkACLRuleDirectionOutboundConst
)

// VPCSecurityGroupRuleProtocol represents the protocols for a Security Group Rule.
// +kubebuilder:validation:Enum=all;icmp;tcp;udp
type VPCSecurityGroupRuleProtocol string

const (
	// VPCSecurityGroupRuleProtocolAll defines the Rule is for all network protocols.
	VPCSecurityGroupRuleProtocolAll VPCSecurityGroupRuleProtocol = vpcv1.NetworkACLRuleProtocolAllConst
	// VPCSecurityGroupRuleProtocolIcmp defiens the Rule is for ICMP network protocol.
	VPCSecurityGroupRuleProtocolIcmp VPCSecurityGroupRuleProtocol = vpcv1.NetworkACLRuleProtocolIcmpConst
	// VPCSecurityGroupRuleProtocolTCP defines the Rule is for TCP network protocol.
	VPCSecurityGroupRuleProtocolTCP VPCSecurityGroupRuleProtocol = vpcv1.NetworkACLRuleProtocolTCPConst
	// VPCSecurityGroupRuleProtocolUDP defines the Rule is for UDP network protocol.
	VPCSecurityGroupRuleProtocolUDP VPCSecurityGroupRuleProtocol = vpcv1.NetworkACLRuleProtocolUDPConst
)

// VPCSecurityGroupRuleRemoteType represents the type of Security Group Rule's destination or source is
// intended. This is intended to define the VPCSecurityGroupRulePrototype subtype.
// For example:
// - any - Any source or destination (0.0.0.0/0)
// - cidr - A CIDR representing a set of IP's (10.0.0.0/28)
// - address - A specific address (192.168.0.1)
// - sg - A Security Group.
// +kubebuilder:validation:Enum=any;cidr;address;sg
type VPCSecurityGroupRuleRemoteType string

const (
	// VPCSecurityGroupRuleRemoteTypeAny defines the destination or source for the Rule is anything/anywhere.
	VPCSecurityGroupRuleRemoteTypeAny VPCSecurityGroupRuleRemoteType = VPCSecurityGroupRuleRemoteType("any")
	// VPCSecurityGroupRuleRemoteTypeCIDR defines the destination or source for the Rule is a CIDR block.
	VPCSecurityGroupRuleRemoteTypeCIDR VPCSecurityGroupRuleRemoteType = VPCSecurityGroupRuleRemoteType("cidr")
	// VPCSecurityGroupRuleRemoteTypeAddress defines the destination or source for the Rule is an address.
	VPCSecurityGroupRuleRemoteTypeAddress VPCSecurityGroupRuleRemoteType = VPCSecurityGroupRuleRemoteType("address")
	// VPCSecurityGroupRuleRemoteTypeSG defines the destination or source for the Rule is a VPC Security Group.
	VPCSecurityGroupRuleRemoteTypeSG VPCSecurityGroupRuleRemoteType = VPCSecurityGroupRuleRemoteType("sg")
)

// IBMCloudResourceReference represents an IBM Cloud resource.
type IBMCloudResourceReference struct {
	// id defines the IBM Cloud Resource ID.
	// +required
	ID string `json:"id"`

	// name defines the IBM Cloud Resource Name.
	// +optional
	Name *string `json:"name,omitempty"`
}

// IBMCloudCatalogOffering represents an IBM Cloud Catalog Offering resource.
// +kubebuilder:validation:XValidation:rule="(has(self.offeringCRN) && !has(self.versionCRN)) || (!has(self.offeringCRN) && has(self.versionCRN))",message="either offeringCRN or version CRN must be provided, not both"
type IBMCloudCatalogOffering struct {
	// OfferingCRN defines the IBM Cloud Catalog Offering CRN. Using the OfferingCRN expects that the latest version of the Offering will be used.
	// If a specific version should be used instead, rely on VersionCRN.
	// +optional
	OfferingCRN *string `json:"offeringCRN,omitempty"`

	// PlanCRN defines the IBM Cloud Catalog Offering Plan CRN to use for the Offering.
	// +optional
	PlanCRN *string `json:"planCRN,omitempty"`

	// VersionCRN defines the IBM Cloud Catalog Offering Version CRN. A specific version of the Catalog Offering will be used, as defined by this CRN.
	// +optional
	VersionCRN *string `json:"versionCRN,omitempty"`
}

// NetworkInterface holds the network interface information like subnet id.
type NetworkInterface struct {
	// SecurityGroups defines a set of IBM Cloud VPC Security Groups to attach to the network interface.
	// +optional
	SecurityGroups []VPCResource `json:"securityGroups,omitempty"`

	// Subnet ID of the network interface.
	Subnet string `json:"subnet,omitempty"`
}

// VPCLoadBalancerBackendPoolMember represents a VPC Load Balancer Backend Pool Member.
type VPCLoadBalancerBackendPoolMember struct {
	// LoadBalancer defines the Load Balancer the Pool Member is for.
	// +required
	LoadBalancer VPCResource `json:"loadBalancer"`

	// Pool defines the Load Balancer Pool the Pool Member should be in.
	// +required
	Pool VPCResource `json:"pool"`

	// Port defines the Port the Load Balancer Pool Member listens for traffic.
	// +required
	Port int64 `json:"port"`

	// Weight of the service member. Only applicable if the pool algorithm is "weighted_round_robin".
	// +optional
	Weight *int64 `json:"weight,omitempty"`
}

// VPCMachinePlacementTarget represents a VPC Machine's placement restrictions.
// +kubebuilder:validation:XValidation:rule="(has(self.dedicatedHost) && !has(self.dedicatedHostGroup) && !has(self.placementGroup)) || (!has(self.dedicatedHost) && has(self.dedicatedHostGroup) && !has(self.placementGroup)) || (!has(self.dedicatedHost) && !has(self.dedicatedHostGroup) && has(self.placementGroup))",message="only one of dedicatedHost, dedicatedHostGroup, or placementGroup must be defined for machine placement"
type VPCMachinePlacementTarget struct {
	// DedicatedHost defines the Dedicated Host to place a VPC Machine (Instance) on.
	// +optional
	DedicatedHost *VPCResource `json:"dedicatedHost,omitempty"`

	// DedicatedHostGroup defines the Dedicated Host Group to use when placing a VPC Machine (Instance).
	// +optional
	DedicatedHostGroup *VPCResource `json:"dedicatedHostGroup,omitempty"`

	// PlacementGroup defines the Placement Group to use when placing a VPC Machine (Instance).
	// +optional
	PlacementGroup *VPCResource `json:"placementGroup,omitempty"`
}

// VPCSecurityGroupPortRange represents a range of ports, minimum to maximum.
// +kubebuilder:validation:XValidation:rule="self.maximumPort >= self.minimumPort",message="maximum port must be greater than or equal to minimum port"
type VPCSecurityGroupPortRange struct {
	// maximumPort is the inclusive upper range of ports.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	MaximumPort int64 `json:"maximumPort,omitempty"`

	// minimumPort is the inclusive lower range of ports.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	MinimumPort int64 `json:"minimumPort,omitempty"`
}

// VPCSecurityGroup defines a VPC Security Group that should exist or be created within the specified VPC, with the specified Security Group Rules.
// +kubebuilder:validation:XValidation:rule="has(self.id) || has(self.name)",message="either an id or name must be specified"
type VPCSecurityGroup struct {
	// id of the Security Group.
	// +optional
	ID *string `json:"id,omitempty"`

	// name of the Security Group.
	// +optional
	Name *string `json:"name,omitempty"`

	// rules are the Security Group Rules for the Security Group.
	// +optional
	Rules []*VPCSecurityGroupRule `json:"rules,omitempty"`

	// tags are tags to add to the Security Group.
	// +optional
	Tags []*string `json:"tags,omitempty"`
}

// VPCSecurityGroupRule defines a VPC Security Group Rule for a specified Security Group.
// +kubebuilder:validation:XValidation:rule="(has(self.destination) && !has(self.source)) || (!has(self.destination) && has(self.source))",message="both destination and source cannot be provided"
// +kubebuilder:validation:XValidation:rule="self.direction == 'inbound' ? has(self.source) : true",message="source must be set for VPCSecurityGroupRuleDirectionInbound direction"
// +kubebuilder:validation:XValidation:rule="self.direction == 'inbound' ? !has(self.destination) : true",message="destination is not valid for VPCSecurityGroupRuleDirectionInbound direction"
// +kubebuilder:validation:XValidation:rule="self.direction == 'outbound' ? has(self.destination) : true",message="destination must be set for VPCSecurityGroupRuleDirectionOutbound direction"
// +kubebuilder:validation:XValidation:rule="self.direction == 'outbound' ? !has(self.source) : true",message="source is not valid for VPCSecurityGroupRuleDirectionOutbound direction"
type VPCSecurityGroupRule struct {
	// action defines whether to allow or deny traffic defined by the Security Group Rule.
	// +required
	Action VPCSecurityGroupRuleAction `json:"action"`

	// destination is a VPCSecurityGroupRulePrototype which defines the destination of outbound traffic for the Security Group Rule.
	// Only used when direction is VPCSecurityGroupRuleDirectionOutbound.
	// +optional
	Destination *VPCSecurityGroupRulePrototype `json:"destination,omitempty"`

	// direction defines whether the traffic is inbound or outbound for the Security Group Rule.
	// +required
	Direction VPCSecurityGroupRuleDirection `json:"direction"`

	// securityGroupID is the ID of the Security Group for the Security Group Rule.
	// +optional
	SecurityGroupID *string `json:"securityGroupID,omitempty"`

	// source is a VPCSecurityGroupRulePrototype which defines the source of inbound traffic for the Security Group Rule.
	// Only used when direction is VPCSecurityGroupRuleDirectionInbound.
	// +optional
	Source *VPCSecurityGroupRulePrototype `json:"source,omitempty"`
}

// VPCSecurityGroupRuleRemote defines a VPC Security Group Rule's remote details.
// The type of remote defines the additional remote details where are used for defining the remote.
// +kubebuilder:validation:XValidation:rule="self.remoteType == 'any' ? (!has(self.cidrSubnetName) && !has(self.address) && !has(self.securityGroupName)) : true",message="cidrSubnetName, addresss, and securityGroupName are not valid for VPCSecurityGroupRuleRemoteTypeAny remoteType"
// +kubebuilder:validation:XValidation:rule="self.remoteType == 'cidr' ? (has(self.cidrSubnetName) && !has(self.address) && !has(self.securityGroupName)) : true",message="only cidrSubnetName is valid for VPCSecurityGroupRuleRemoteTypeCIDR remoteType"
// +kubebuilder:validation:XValidation:rule="self.remoteType == 'address' ? (has(self.address) && !has(self.cidrSubnetName) && !has(self.securityGroupName)) : true",message="only address is valid for VPCSecurityGroupRuleRemoteTypeIP remoteType"
// +kubebuilder:validation:XValidation:rule="self.remoteType == 'sg' ? (has(self.securityGroupName) && !has(self.cidrSubnetName) && !has(self.address)) : true",message="only securityGroupName is valid for VPCSecurityGroupRuleRemoteTypeSG remoteType"
type VPCSecurityGroupRuleRemote struct {
	// cidrSubnetName is the name of the VPC Subnet to retrieve the CIDR from, to use for the remote's destination/source.
	// Only used when remoteType is VPCSecurityGroupRuleRemoteTypeCIDR.
	// +optional
	CIDRSubnetName *string `json:"cidrSubnetName,omitempty"`

	//  address is the address to use for the remote's destination/source.
	// Only used when remoteType is VPCSecurityGroupRuleRemoteTypeAddress.
	// +optional
	Address *string `json:"address,omitempty"`

	// remoteType defines the type of filter to define for the remote's destination/source.
	// +required
	RemoteType VPCSecurityGroupRuleRemoteType `json:"remoteType"`

	// securityGroupName is the name of the VPC Security Group to use for the remote's destination/source.
	// Only used when remoteType is VPCSecurityGroupRuleRemoteTypeSG
	// +optional
	SecurityGroupName *string `json:"securityGroupName,omitempty"`
}

// VPCSecurityGroupRulePrototype defines a VPC Security Group Rule's traffic specifics for a series of remotes (destinations or sources).
// +kubebuilder:validation:XValidation:rule="self.protocol != 'icmp' ? (!has(self.icmpCode) && !has(self.icmpType)) : true",message="icmpCode and icmpType are only supported for VPCSecurityGroupRuleProtocolIcmp protocol"
// +kubebuilder:validation:XValidation:rule="self.protocol == 'all' ? !has(self.portRange) : true",message="portRange is not valid for VPCSecurityGroupRuleProtocolAll protocol"
// +kubebuilder:validation:XValidation:rule="self.protocol == 'icmp' ? !has(self.portRange) : true",message="portRange is not valid for VPCSecurityGroupRuleProtocolIcmp protocol"
type VPCSecurityGroupRulePrototype struct {
	// icmpCode is the ICMP code for the Rule.
	// Only used when Protocol is VPCSecurityGroupRuleProtocolIcmp.
	// +optional
	ICMPCode *int64 `json:"icmpCode,omitempty"`

	// icmpType is the ICMP type for the Rule.
	// Only used when Protocol is VPCSecurityGroupRuleProtocolIcmp.
	// +optional
	ICMPType *int64 `json:"icmpType,omitempty"`

	// portRange is a range of ports allowed for the Rule's remote.
	// +optional
	PortRange *VPCSecurityGroupPortRange `json:"portRange,omitempty"`

	// protocol defines the traffic protocol used for the Security Group Rule.
	// +required
	Protocol VPCSecurityGroupRuleProtocol `json:"protocol"`

	// remotes is a set of VPCSecurityGroupRuleRemote's that define the traffic allowed by the Rule's remote.
	// Specifying multiple VPCSecurityGroupRuleRemote's creates a unique Security Group Rule with the shared Protocol, PortRange, etc.
	// This allows for easier management of Security Group Rule's for sets of CIDR's, IP's, etc.
	Remotes []VPCSecurityGroupRuleRemote `json:"remotes"`
}

// Subnet describes a subnet.
type Subnet struct {
	Ipv4CidrBlock *string `json:"cidr,omitempty"`
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength:=63
	// +kubebuilder:validation:Pattern=`^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`
	Name *string `json:"name,omitempty"`
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength:=64
	// +kubebuilder:validation:Pattern=`^[-0-9a-z_]+$`
	ID   *string `json:"id,omitempty"`
	Zone *string `json:"zone,omitempty"`
}

// VPCEndpoint describes a VPCEndpoint.
type VPCEndpoint struct {
	Address *string `json:"address"`
	// +optional
	// Deprecated: This field has no function and is going to be removed in the next release.
	FIPID *string `json:"floatingIPID,omitempty"`
	// +optional
	LBID *string `json:"loadBalancerIPID,omitempty"`
}

// ResourceStatus identifies a resource by id (and name) and whether it is ready.
type ResourceStatus struct {
	// id defines the Id of the IBM Cloud resource status.
	// +required
	ID string `json:"id"`

	// name defines the name of the IBM Cloud resource status.
	// +optional
	Name *string `json:"name,omitempty"`

	// ready defines whether the IBM Cloud resource is ready.
	// +required
	Ready bool `json:"ready"`
}

// Set sets the ResourceStatus fields.
func (s *ResourceStatus) Set(resource ResourceStatus) {
	s.ID = resource.ID
	// Set the name if it hasn't been, or the incoming name won't remove it (nil).
	if s.Name == nil && resource.Name != nil {
		s.Name = resource.Name
	}
	s.Ready = resource.Ready
}

// VPCResource represents a VPC resource.
// +kubebuilder:validation:XValidation:rule="has(self.id) || has(self.name)",message="an id or name must be provided"
type VPCResource struct {
	// id of the resource.
	// +kubebuilder:validation:MinLength=1
	// +optional
	ID *string `json:"id,omitempty"`

	// name of the resource.
	// +kubebuilder:validation:MinLength=1
	// +optional
	Name *string `json:"name,omitempty"`
}
