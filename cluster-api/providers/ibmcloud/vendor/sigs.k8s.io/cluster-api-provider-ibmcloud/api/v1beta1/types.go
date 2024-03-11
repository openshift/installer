/*
Copyright 2021 The Kubernetes Authors.

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

package v1beta1

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

// VPCLoadBalancerState describes the state of the load balancer.
type VPCLoadBalancerState string

var (
	// VPCLoadBalancerStateActive is the string representing the load balancer in a active state.
	VPCLoadBalancerStateActive = VPCLoadBalancerState("active")

	// VPCLoadBalancerStateCreatePending is the string representing the load balancer in a queued state.
	VPCLoadBalancerStateCreatePending = VPCLoadBalancerState("create_pending")

	// VPCLoadBalancerStateDeletePending is the string representing the load balancer in a failed state.
	VPCLoadBalancerStateDeletePending = VPCLoadBalancerState("delete_pending")
)

// DeletePolicy defines the policy used to identify images to be preserved.
type DeletePolicy string

var (
	// DeletePolicyRetain is the string representing an image to be retained.
	DeletePolicyRetain = DeletePolicy("retain")
)

// NetworkInterface holds the network interface information like subnet id.
type NetworkInterface struct {
	// Subnet ID of the network interface.
	Subnet string `json:"subnet,omitempty"`
}

// Subnet describes a subnet.
type Subnet struct {
	Ipv4CidrBlock *string `json:"cidr"`
	Name          *string `json:"name"`
	ID            *string `json:"id"`
	Zone          *string `json:"zone"`
}

// VPCEndpoint describes a VPCEndpoint.
type VPCEndpoint struct {
	Address *string `json:"address"`
	// +optional
	FIPID *string `json:"floatingIPID,omitempty"`
	// +optional
	LBID *string `json:"loadBalancerIPID,omitempty"`
}
