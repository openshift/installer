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

// NetworkScope is the interface for the scope to be used with the network services.
type NetworkScope interface {
	cloud.ClusterScoper

	// Network returns the cluster network object.
	Network() *infrav1.NetworkStatus
	// VPC returns the cluster VPC.
	VPC() *infrav1.VPCSpec
	// Subnets returns the cluster subnets.
	Subnets() infrav1.Subnets
	// SetSubnets updates the clusters subnets.
	SetSubnets(subnets infrav1.Subnets)
	// CNIIngressRules returns the CNI spec ingress rules.
	CNIIngressRules() infrav1.CNIIngressRules
	// SecurityGroups returns the cluster security groups as a map, it creates the map if empty.
	SecurityGroups() map[infrav1.SecurityGroupRole]infrav1.SecurityGroup
	// SecondaryCidrBlock returns the optional secondary CIDR block to use for pod IPs. This may later be renamed since
	// it should not be confused with SecondaryCidrBlocks.
	SecondaryCidrBlock() *string
	// SecondaryCidrBlocks returns the additional CIDR blocks to be associated with the managed VPC.
	SecondaryCidrBlocks() []infrav1.VpcCidrBlock
	// AllSecondaryCidrBlocks returns a unique list of all secondary CIDR blocks (combining `SecondaryCidrBlock` and
	// `SecondaryCidrBlocks`).
	AllSecondaryCidrBlocks() []infrav1.VpcCidrBlock

	// Bastion returns the bastion details for the cluster.
	Bastion() *infrav1.Bastion

	// Bucket returns the cluster bucket.
	Bucket() *infrav1.S3Bucket

	// TagUnmanagedNetworkResources returns is tagging unmanaged network resources is set.
	TagUnmanagedNetworkResources() bool

	// SetNatGatewaysIPs sets the Nat Gateways Public IPs.
	SetNatGatewaysIPs(ips []string)
	// GetNatGatewaysIPs gets the Nat Gateways Public IPs.
	GetNatGatewaysIPs() []string
}
