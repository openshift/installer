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

package v1beta1

import (
	"github.com/pkg/errors"
	"k8s.io/utils/net"
)

// AzureManagedControlPlaneTemplateResourceSpec specifies an Azure managed control plane template resource.
type AzureManagedControlPlaneTemplateResourceSpec struct {
	AzureManagedControlPlaneClassSpec `json:",inline"`
}

// AzureManagedControlPlaneTemplateMachineTemplate is only used to fulfill the CAPI contract which expects a
// MachineTemplate field for any controlplane ref in a topology.
type AzureManagedControlPlaneTemplateMachineTemplate struct{}

// AzureManagedMachinePoolTemplateResourceSpec specifies an Azure managed control plane template resource.
type AzureManagedMachinePoolTemplateResourceSpec struct {
	AzureManagedMachinePoolClassSpec `json:",inline"`
}

// AzureManagedClusterTemplateResourceSpec specifies an Azure managed cluster template resource.
type AzureManagedClusterTemplateResourceSpec struct{}

// AzureClusterTemplateResourceSpec specifies an Azure cluster template resource.
type AzureClusterTemplateResourceSpec struct {
	AzureClusterClassSpec `json:",inline"`

	// NetworkSpec encapsulates all things related to Azure network.
	// +optional
	NetworkSpec NetworkTemplateSpec `json:"networkSpec,omitempty"`

	// BastionSpec encapsulates all things related to the Bastions in the cluster.
	// +optional
	BastionSpec BastionTemplateSpec `json:"bastionSpec,omitempty"`
}

// NetworkTemplateSpec specifies a network template.
type NetworkTemplateSpec struct {
	NetworkClassSpec `json:",inline"`

	// Vnet is the configuration for the Azure virtual network.
	// +optional
	Vnet VnetTemplateSpec `json:"vnet,omitempty"`

	// Subnets is the configuration for the control-plane subnet and the node subnet.
	// +optional
	Subnets SubnetTemplatesSpec `json:"subnets,omitempty"`

	// APIServerLB is the configuration for the control-plane load balancer.
	// +optional
	APIServerLB LoadBalancerClassSpec `json:"apiServerLB,omitempty"`

	// NodeOutboundLB is the configuration for the node outbound load balancer.
	// +optional
	NodeOutboundLB *LoadBalancerClassSpec `json:"nodeOutboundLB,omitempty"`

	// ControlPlaneOutboundLB is the configuration for the control-plane outbound load balancer.
	// This is different from APIServerLB, and is used only in private clusters (optionally) for enabling outbound traffic.
	// +optional
	ControlPlaneOutboundLB *LoadBalancerClassSpec `json:"controlPlaneOutboundLB,omitempty"`

	// AdditionalAPIServerLBPorts is the configuration for the additional inbound control-plane load balancer ports
	// Each port specified (e.g., 9345) creates an inbound rule where the frontend port and the backend port are the same.
	// +optional
	AdditionalAPIServerLBPorts []LoadBalancerPort `json:"additionalAPIServerLBPorts,omitempty"`
}

// GetSubnetTemplate returns the subnet template based on the subnet role.
func (n *NetworkTemplateSpec) GetSubnetTemplate(role SubnetRole) (SubnetTemplateSpec, error) {
	for _, sn := range n.Subnets {
		if sn.Role == role {
			return sn, nil
		}
	}
	return SubnetTemplateSpec{}, errors.Errorf("no subnet template found with role %s", role)
}

// UpdateSubnetTemplate updates the subnet template based on subnet role.
func (n *NetworkTemplateSpec) UpdateSubnetTemplate(subnet SubnetTemplateSpec, role SubnetRole) {
	for i, sn := range n.Subnets {
		if sn.Role == role {
			n.Subnets[i] = subnet
		}
	}
}

// VnetTemplateSpec defines the desired state of a virtual network.
type VnetTemplateSpec struct {
	VnetClassSpec `json:",inline"`

	// Peerings defines a list of peerings of the newly created virtual network with existing virtual networks.
	// +optional
	Peerings VnetPeeringsTemplateSpec `json:"peerings,omitempty"`
}

// VnetPeeringsTemplateSpec defines a list of peerings of the newly created virtual network with existing virtual networks.
type VnetPeeringsTemplateSpec []VnetPeeringClassSpec

// SubnetTemplateSpec specifies a template for a subnet.
type SubnetTemplateSpec struct {
	SubnetClassSpec `json:",inline"`

	// SecurityGroup defines the NSG (network security group) that should be attached to this subnet.
	// +optional
	SecurityGroup SecurityGroupClass `json:"securityGroup,omitempty"`

	// NatGateway associated with this subnet.
	// +optional
	NatGateway NatGatewayClassSpec `json:"natGateway,omitempty"`
}

// IsNatGatewayEnabled returns true if the NAT gateway is enabled.
func (s SubnetTemplateSpec) IsNatGatewayEnabled() bool {
	return s.NatGateway.Name != ""
}

// IsIPv6Enabled returns whether or not IPv6 is enabled on the subnet.
func (s SubnetTemplateSpec) IsIPv6Enabled() bool {
	for _, cidr := range s.CIDRBlocks {
		if net.IsIPv6CIDRString(cidr) {
			return true
		}
	}
	return false
}

// SubnetTemplatesSpec specifies a list of subnet templates.
// +listType=map
// +listMapKey=name
type SubnetTemplatesSpec []SubnetTemplateSpec

// BastionTemplateSpec specifies a template for a bastion host.
type BastionTemplateSpec struct {
	// +optional
	AzureBastion *AzureBastionTemplateSpec `json:"azureBastion,omitempty"`
}

// AzureBastionTemplateSpec specifies a template for an Azure Bastion host.
type AzureBastionTemplateSpec struct {
	// +optional
	Subnet SubnetTemplateSpec `json:"subnet,omitempty"`
}
