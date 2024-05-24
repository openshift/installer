/*
Copyright 2024 The Kubernetes Authors.

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

package scope

// ResourceNotFound is the string representing an error when a resource is not found in IBM Cloud.
type ResourceNotFound string

var (
	// VPCLoadBalancerNotFound is the error returned when a VPC load balancer is not found.
	VPCLoadBalancerNotFound = ResourceNotFound("cannot be found")

	// VPCSubnetNotFound is the error returned when a VPC subnet is not found.
	VPCSubnetNotFound = ResourceNotFound("Subnet not found")

	// VPCNotFound is the error returned when a VPC is not found.
	VPCNotFound = ResourceNotFound("VPC not found")

	// TransitGatewayNotFound is the error returned when a transit gateway is not found.
	TransitGatewayNotFound = ResourceNotFound("gateway was not found")

	// DHCPServerNotFound is the error returned when a DHCP server is not found.
	DHCPServerNotFound = ResourceNotFound("dhcp server does not exist")

	// COSInstanceNotFound is the error returned when a COS service instance is not found.
	COSInstanceNotFound = ResourceNotFound("COS instance unavailable")

	// VPCSecurityGroupNotFound is the error returned when a VPC security group is not found.
	VPCSecurityGroupNotFound = ResourceNotFound("Security group not found")
)
