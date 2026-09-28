/*
Copyright 2026 The Kubernetes Authors.

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

package v1beta3

import "github.com/IBM/vpc-go-sdk/vpcv1"

const (
	// DefaultAPIServerPort is default API server port number.
	DefaultAPIServerPort int32 = 6443
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

	// PowerVSImageStateQueued is the string representing an image in a queued state.
	PowerVSImageStateQueued = PowerVSImageState("queued")

	// PowerVSImageStateFailed is the string representing an image in a failed state.
	PowerVSImageStateFailed = PowerVSImageState("failed")

	// PowerVSImageStateImporting is the string representing an image in a failed state.
	PowerVSImageStateImporting = PowerVSImageState("importing")

	// PowerVSImageStateCompleted is the string representing an image in a completed state.
	PowerVSImageStateCompleted = PowerVSImageState("completed")
)

// WorkspaceState describes the state of a PowerVS workspace.
type WorkspaceState string

var (
	// WorkspaceStateActive is the string representing a workspace in an active state.
	WorkspaceStateActive = WorkspaceState("active")

	// WorkspaceStateProvisioning is the string representing a workspace in a provisioning state.
	WorkspaceStateProvisioning = WorkspaceState("provisioning")

	// WorkspaceStateFailed is the string representing a workspace in a failed state.
	WorkspaceStateFailed = WorkspaceState("failed")

	// WorkspaceStateRemoved is the string representing a workspace in a removed state.
	WorkspaceStateRemoved = WorkspaceState("removed")
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

// LoadBalancerBackendPoolAlgorithm describes the backend pool's load balancing algorithm.
// +kubebuilder:validation:Enum=least_connections;round_robin;weighted_round_robin
type LoadBalancerBackendPoolAlgorithm string

var (
	// LoadBalancerBackendPoolAlgorithmLeastConnections is the string representing the least_connections load balancing algorithm.
	LoadBalancerBackendPoolAlgorithmLeastConnections LoadBalancerBackendPoolAlgorithm = vpcv1.CreateLoadBalancerPoolOptionsAlgorithmLeastConnectionsConst

	// LoadBalancerBackendPoolAlgorithmRoundRobin is the string representing the round_robin load balancing algorithm.
	LoadBalancerBackendPoolAlgorithmRoundRobin LoadBalancerBackendPoolAlgorithm = vpcv1.CreateLoadBalancerPoolOptionsAlgorithmRoundRobinConst

	// LoadBalancerBackendPoolAlgorithmWeightedRoundRobin is the string representing the weighted_round_robin load balancing algorithm.
	LoadBalancerBackendPoolAlgorithmWeightedRoundRobin LoadBalancerBackendPoolAlgorithm = vpcv1.CreateLoadBalancerPoolOptionsAlgorithmWeightedRoundRobinConst
)

// LoadBalancerBackendPoolProtocol describes the protocol for load balancer backend pools.
// We have unique types in case IBM Cloud Load Balancer Listener and Backend Pool supported algorithms ever diverage.
// +kubebuilder:validation:Enum=http;https;tcp;udp
type LoadBalancerBackendPoolProtocol string

var (
	// LoadBalancerBackendPoolProtocolHTTP is the string representing the http protocol for load balancer backend pools.
	LoadBalancerBackendPoolProtocolHTTP LoadBalancerBackendPoolProtocol = vpcv1.LoadBalancerPoolPrototypeLoadBalancerContextProtocolHTTPConst

	// LoadBalancerBackendPoolProtocolHTTPS is the string representing the https protocol for load balancer backend pools.
	LoadBalancerBackendPoolProtocolHTTPS LoadBalancerBackendPoolProtocol = vpcv1.LoadBalancerPoolPrototypeLoadBalancerContextProtocolHTTPSConst

	// LoadBalancerBackendPoolProtocolTCP is the string representing the tcp protocol for load balancer backend pools.
	LoadBalancerBackendPoolProtocolTCP LoadBalancerBackendPoolProtocol = vpcv1.LoadBalancerPoolPrototypeLoadBalancerContextProtocolTCPConst

	// LoadBalancerBackendPoolProtocolUDP is the string representing the tudp protocol for load balancer backend pools.
	LoadBalancerBackendPoolProtocolUDP LoadBalancerBackendPoolProtocol = vpcv1.LoadBalancerPoolPrototypeLoadBalancerContextProtocolUDPConst
)

// LoadBalancerListenerProtocol describes the protocol for load balancer listeners.
// We have unique types in case IBM Cloud Load Balancer Listener and Backend Pool supported algorithms ever diverage.
// +kubebuilder:validation:Enum=http;https;tcp;udp
type LoadBalancerListenerProtocol string

var (
	// LoadBalancerListenerProtocolHTTP is the string representing the http protocol for load balancer listeners.
	LoadBalancerListenerProtocolHTTP LoadBalancerListenerProtocol = vpcv1.LoadBalancerListenerProtocolHTTPConst

	// LoadBalancerListenerProtocolHTTPS is the string representing the https protocol for load balancer listeners.
	LoadBalancerListenerProtocolHTTPS LoadBalancerListenerProtocol = vpcv1.LoadBalancerListenerProtocolHTTPSConst

	// LoadBalancerListenerProtocolTCP is the string representing the tcp protocol for load balancer listeners.
	LoadBalancerListenerProtocolTCP LoadBalancerListenerProtocol = vpcv1.LoadBalancerListenerProtocolTCPConst

	// LoadBalancerListenerProtocolUDP is the string representing the tudp protocol for load balancer listeners.
	LoadBalancerListenerProtocolUDP LoadBalancerListenerProtocol = vpcv1.LoadBalancerListenerProtocolUDPConst
)

// LoadBalancerBackendPoolHealthMonitorType describes the backend pool's health check protocol type.
// +kubebuilder:validation:Enum=http;https;tcp
type LoadBalancerBackendPoolHealthMonitorType string

var (
	// LoadBalancerBackendPoolHealthMonitorTypeHTTP is the string representing the http health pool protocol type.
	LoadBalancerBackendPoolHealthMonitorTypeHTTP LoadBalancerBackendPoolHealthMonitorType = vpcv1.LoadBalancerPoolHealthMonitorTypeHTTPConst

	// LoadBalancerBackendPoolHealthMonitorTypeHTTPS is the string representing the https health pool protocol type.
	LoadBalancerBackendPoolHealthMonitorTypeHTTPS LoadBalancerBackendPoolHealthMonitorType = vpcv1.LoadBalancerPoolHealthMonitorTypeHTTPSConst

	// LoadBalancerBackendPoolHealthMonitorTypeTCP is the string representing the tcp health pool protocol type.
	LoadBalancerBackendPoolHealthMonitorTypeTCP LoadBalancerBackendPoolHealthMonitorType = vpcv1.LoadBalancerPoolHealthMonitorTypeTCPConst
)

// LoadBalancerState describes the state of the load balancer.
type LoadBalancerState string

var (
	// LoadBalancerStateActive is the string representing the load balancer in a active state.
	LoadBalancerStateActive = LoadBalancerState("active")

	// LoadBalancerStateCreatePending is the string representing the load balancer in a queued state.
	LoadBalancerStateCreatePending = LoadBalancerState("create_pending")

	// LoadBalancerStateUpdatePending is the string representing the load balancer in updating state.
	LoadBalancerStateUpdatePending = LoadBalancerState("update_pending")

	// LoadBalancerStateDeletePending is the string representing the load balancer in deleting state.
	LoadBalancerStateDeletePending = LoadBalancerState("delete_pending")
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
