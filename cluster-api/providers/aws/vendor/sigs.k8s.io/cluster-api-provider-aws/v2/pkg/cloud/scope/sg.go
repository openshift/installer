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

package scope

import (
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
)

// SGScope is the interface for the scope to be used with the sg service.
type SGScope interface {
	cloud.ClusterScoper

	// Network returns the cluster network object.
	Network() *infrav1.NetworkStatus

	// SecurityGroups returns the cluster security groups as a map, it creates the map if empty.
	SecurityGroups() map[infrav1.SecurityGroupRole]infrav1.SecurityGroup

	// SecurityGroupOverrides returns the security groups that are used as overrides in the cluster spec
	SecurityGroupOverrides() map[infrav1.SecurityGroupRole]string

	// VPC returns the cluster VPC.
	VPC() *infrav1.VPCSpec

	// CNIIngressRules returns the CNI spec ingress rules.
	CNIIngressRules() infrav1.CNIIngressRules

	// Bastion returns the bastion details for the cluster.
	Bastion() *infrav1.Bastion

	// ControlPlaneLoadBalancer returns the load balancer settings that are requested.
	// Deprecated: Use ControlPlaneLoadBalancers()
	ControlPlaneLoadBalancer() *infrav1.AWSLoadBalancerSpec

	// SetNatGatewaysIPs sets the Nat Gateways Public IPs.
	SetNatGatewaysIPs(ips []string)

	// GetNatGatewaysIPs gets the Nat Gateways Public IPs.
	GetNatGatewaysIPs() []string

	// AdditionalControlPlaneIngressRules returns the additional ingress rules for the control plane security group.
	AdditionalControlPlaneIngressRules() []infrav1.IngressRule

	// ControlPlaneLoadBalancers returns both the ControlPlaneLoadBalancer and SecondaryControlPlaneLoadBalancer AWSLoadBalancerSpecs.
	// The control plane load balancers should always be returned in the above order.
	ControlPlaneLoadBalancers() []*infrav1.AWSLoadBalancerSpec

	// NodePortIngressRuleCidrBlocks returns the CIDR blocks for the node NodePort ingress rules.
	NodePortIngressRuleCidrBlocks() []string
}
