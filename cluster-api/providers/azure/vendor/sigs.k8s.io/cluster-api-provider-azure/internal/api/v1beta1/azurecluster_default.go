/*
Copyright The Kubernetes Authors.

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
	"fmt"

	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/feature"
)

const (
	// DefaultVnetCIDR is the default Vnet CIDR.
	DefaultVnetCIDR = "10.0.0.0/8"
	// DefaultControlPlaneSubnetCIDR is the default Control Plane Subnet CIDR.
	DefaultControlPlaneSubnetCIDR = "10.0.0.0/16"
	// DefaultNodeSubnetCIDR is the default Node Subnet CIDR.
	DefaultNodeSubnetCIDR = "10.1.0.0/16"
	// DefaultClusterSubnetCIDR is the default Cluster Subnet CIDR.
	DefaultClusterSubnetCIDR = "10.0.0.0/16"
	// DefaultNodeSubnetCIDRPattern is the pattern that will be used to generate the default subnets CIDRs.
	DefaultNodeSubnetCIDRPattern = "10.%d.0.0/16"
	// DefaultAzureBastionSubnetCIDR is the default Subnet CIDR for AzureBastion.
	DefaultAzureBastionSubnetCIDR = "10.255.255.224/27"
	// DefaultAzureBastionSubnetName is the default Subnet Name for AzureBastion.
	DefaultAzureBastionSubnetName = "AzureBastionSubnet"
	// DefaultAzureBastionSubnetRole is the default Subnet role for AzureBastion.
	DefaultAzureBastionSubnetRole = infrav1.SubnetBastion
	// DefaultInternalLBIPAddress is the default internal load balancer ip address.
	DefaultInternalLBIPAddress = "10.0.0.100"
	// DefaultOutboundRuleIdleTimeoutInMinutes is the default for IdleTimeoutInMinutes for the load balancer.
	DefaultOutboundRuleIdleTimeoutInMinutes = 4
	// DefaultAzureCloud is the public cloud that will be used by most users.
	DefaultAzureCloud = "AzurePublicCloud"
)

// SetDefaultsAzureCluster sets default values for an AzureCluster.
func SetDefaultsAzureCluster(c *infrav1.AzureCluster) {
	AzureClusterClassSpecSetDefaults(&c.Spec.AzureClusterClassSpec)
	setDefaultAzureClusterResourceGroup(c)
	setDefaultAzureClusterNetworkSpec(c)
}

// setDefaultAzureClusterNetworkSpec sets default values for an AzureCluster's NetworkSpec.
func setDefaultAzureClusterNetworkSpec(c *infrav1.AzureCluster) {
	setDefaultAzureClusterVnet(c)
	setDefaultAzureClusterBastion(c)
	setDefaultAzureClusterSubnets(c)
	setDefaultAzureClusterVnetPeering(c)
	if c.Spec.ControlPlaneEnabled {
		setDefaultAzureClusterAPIServerLB(c)
	}
	setDefaultAzureClusterNodeOutboundLB(c)
	if c.Spec.ControlPlaneEnabled {
		setDefaultAzureClusterControlPlaneOutboundLB(c)
	}
	if !c.Spec.ControlPlaneEnabled {
		c.Spec.NetworkSpec.APIServerLB = nil
	}
}

// setDefaultAzureClusterResourceGroup sets the default resource group for an AzureCluster.
func setDefaultAzureClusterResourceGroup(c *infrav1.AzureCluster) {
	if c.Spec.ResourceGroup == "" {
		c.Spec.ResourceGroup = c.Name
	}
}

// SetDefaultAzureClusterAzureEnvironment sets the default Azure environment for an AzureCluster.
func SetDefaultAzureClusterAzureEnvironment(c *infrav1.AzureCluster) {
	if c.Spec.AzureEnvironment == "" {
		c.Spec.AzureEnvironment = DefaultAzureCloud
	}
}

// setDefaultAzureClusterVnet sets default values for an AzureCluster's VNet.
func setDefaultAzureClusterVnet(c *infrav1.AzureCluster) {
	if c.Spec.NetworkSpec.Vnet.ResourceGroup == "" {
		c.Spec.NetworkSpec.Vnet.ResourceGroup = c.Spec.ResourceGroup
	}
	if c.Spec.NetworkSpec.Vnet.Name == "" {
		c.Spec.NetworkSpec.Vnet.Name = generateVnetName(c.ObjectMeta.Name)
	}
	VnetClassSpecSetDefaults(&c.Spec.NetworkSpec.Vnet.VnetClassSpec)
}

// setDefaultAzureClusterSubnets ensures a fully populated, default subnet configuration
// and in certain scenarios creates new, default subnet configurations.
func setDefaultAzureClusterSubnets(c *infrav1.AzureCluster) {
	clusterSubnet, err := c.Spec.NetworkSpec.GetSubnet(infrav1.SubnetCluster)
	clusterSubnetExists := err == nil
	// If we already have a cluster subnet defined, ensure it has sensible defaults
	// for all properties.
	if clusterSubnetExists {
		setDefaultSubnetSpecClusterSubnet(&clusterSubnet, c.ObjectMeta.Name)
		c.Spec.NetworkSpec.UpdateSubnet(clusterSubnet, infrav1.SubnetCluster)
	}

	if c.Spec.ControlPlaneEnabled {
		cpSubnet, errcp := c.Spec.NetworkSpec.GetSubnet(infrav1.SubnetControlPlane)
		// If we already have a control plane subnet defined, ensure it has sensible defaults
		// for all properties.
		if errcp == nil {
			setDefaultSubnetSpecControlPlaneSubnet(&cpSubnet, c.ObjectMeta.Name)
			c.Spec.NetworkSpec.UpdateSubnet(cpSubnet, infrav1.SubnetControlPlane)
			// If we don't have either a control plane subnet or a cluster subnet,
			// create a new control plane subnet from scratch and populate with sensible defaults.
		} else if !clusterSubnetExists {
			cpSubnet = infrav1.SubnetSpec{SubnetClassSpec: infrav1.SubnetClassSpec{Role: infrav1.SubnetControlPlane}}
			setDefaultSubnetSpecControlPlaneSubnet(&cpSubnet, c.ObjectMeta.Name)
			c.Spec.NetworkSpec.Subnets = append(c.Spec.NetworkSpec.Subnets, cpSubnet)
		}
	}

	// anyNodeSubnetFound tracks whether or not we have one or more node subnets defined.
	var anyNodeSubnetFound bool
	// nodeSubnetCounter tracks all node subnets to aid automatic CIDR configuration.
	var nodeSubnetCounter int
	for i, subnet := range c.Spec.NetworkSpec.Subnets {
		// Skip all non-node subnets
		if subnet.Role != infrav1.SubnetNode {
			continue
		}
		nodeSubnetCounter++
		anyNodeSubnetFound = true
		// Set has sensible defaults for this existing node subnet.
		setDefaultSubnetSpecNodeSubnet(&subnet, c.ObjectMeta.Name, nodeSubnetCounter)
		// Because there can be multiple node subnets, we have to update any changes
		// after applying defaults to the explicit item at the current index.
		c.Spec.NetworkSpec.Subnets[i] = subnet
	}

	// We need at least one subnet for nodes.
	// If no node subnets are defined, and there is no cluster subnet defined,
	// create a default 10.1.0.0/16 node subnet.
	if !anyNodeSubnetFound && !clusterSubnetExists {
		nodeSubnet := infrav1.SubnetSpec{
			SubnetClassSpec: infrav1.SubnetClassSpec{
				Role:       infrav1.SubnetNode,
				CIDRBlocks: []string{DefaultNodeSubnetCIDR},
				Name:       generateNodeSubnetName(c.ObjectMeta.Name),
			},
			SecurityGroup: infrav1.SecurityGroup{
				Name: generateNodeSecurityGroupName(c.ObjectMeta.Name),
			},
			RouteTable: infrav1.RouteTable{
				Name: generateNodeRouteTableName(c.ObjectMeta.Name),
			},
			NatGateway: infrav1.NatGateway{
				NatGatewayClassSpec: infrav1.NatGatewayClassSpec{
					Name: generateNatGatewayName(c.ObjectMeta.Name),
				},
			},
		}
		c.Spec.NetworkSpec.Subnets = append(c.Spec.NetworkSpec.Subnets, nodeSubnet)
	}

	// Default the control plane subnet to share the same route table as the
	// node subnet. The cloud-controller-manager route controller programs
	// pod-CIDR routes into a single route table (the one on the node/cluster
	// subnet, see getOneNodeSubnet and the cloud config routeTableName). When
	// the control plane and node subnets differ and the CNI relies on Azure
	// routing (e.g. dual-stack/IPv6 with encapsulation disabled), the control
	// plane node must share that route table; otherwise the API server
	// aggregator cannot reach aggregated-API/webhook pods scheduled on worker
	// nodes. For overlay CNIs (e.g. single-stack VXLAN) the route table carries
	// no pod routes, so sharing it is a no-op. We copy the node subnet's actual
	// route table name (rather than the generated default) so this also works
	// when a custom node route table name is configured.
	if c.Spec.ControlPlaneEnabled {
		if cpSubnet, err := c.Spec.NetworkSpec.GetSubnet(infrav1.SubnetControlPlane); err == nil && cpSubnet.RouteTable.Name == "" {
			if rtName := nodeRouteTableName(c); rtName != "" {
				cpSubnet.RouteTable.Name = rtName
				c.Spec.NetworkSpec.UpdateSubnet(cpSubnet, infrav1.SubnetControlPlane)
			}
		}
	}
}

// nodeRouteTableName returns the route table name of the subnet that the
// cloud-controller-manager programs pod routes into. This mirrors the subnet
// selection in getOneNodeSubnet (the first node- or cluster-role subnet).
func nodeRouteTableName(c *infrav1.AzureCluster) string {
	for _, subnet := range c.Spec.NetworkSpec.Subnets {
		if subnet.Role == infrav1.SubnetNode || subnet.Role == infrav1.SubnetCluster {
			return subnet.RouteTable.Name
		}
	}
	return ""
}

// setDefaultSubnetSpecNodeSubnet sets default values for a node SubnetSpec.
func setDefaultSubnetSpecNodeSubnet(s *infrav1.SubnetSpec, clusterName string, index int) {
	if s.Name == "" {
		s.Name = withIndex(generateNodeSubnetName(clusterName), index)
	}
	SubnetClassSpecSetDefaults(&s.SubnetClassSpec, fmt.Sprintf(DefaultNodeSubnetCIDRPattern, index))

	if s.SecurityGroup.Name == "" {
		s.SecurityGroup.Name = generateNodeSecurityGroupName(clusterName)
	}
	SecurityGroupClassSetDefaults(&s.SecurityGroup.SecurityGroupClass)

	if s.RouteTable.Name == "" {
		s.RouteTable.Name = generateNodeRouteTableName(clusterName)
	}

	// NAT gateway only supports the use of IPv4 public IP addresses for outbound connectivity.
	// So default use the NAT gateway for outbound traffic in IPv4 cluster instead of loadbalancer.
	// We assume that if the ID is set, the subnet already exists so we shouldn't add a NAT gateway.
	if !s.IsIPv6Enabled() && s.ID == "" {
		if s.NatGateway.Name == "" {
			s.NatGateway.Name = withIndex(generateNatGatewayName(clusterName), index)
		}
		if s.NatGateway.NatGatewayIP.Name == "" {
			s.NatGateway.NatGatewayIP.Name = generateNatGatewayIPName(s.NatGateway.Name)
		}
	}
}

// setDefaultSubnetSpecControlPlaneSubnet sets default values for a control plane SubnetSpec.
func setDefaultSubnetSpecControlPlaneSubnet(s *infrav1.SubnetSpec, clusterName string) {
	if s.Name == "" {
		s.Name = generateControlPlaneSubnetName(clusterName)
	}

	SubnetClassSpecSetDefaults(&s.SubnetClassSpec, DefaultControlPlaneSubnetCIDR)

	if s.SecurityGroup.Name == "" {
		s.SecurityGroup.Name = generateControlPlaneSecurityGroupName(clusterName)
	}
	SecurityGroupClassSetDefaults(&s.SecurityGroup.SecurityGroupClass)
}

// setDefaultSubnetSpecClusterSubnet sets default values for a cluster SubnetSpec.
func setDefaultSubnetSpecClusterSubnet(s *infrav1.SubnetSpec, clusterName string) {
	if s.Name == "" {
		s.Name = generateClusterSubnetSubnetName(clusterName)
	}
	if s.SecurityGroup.Name == "" {
		s.SecurityGroup.Name = generateClusterSecurityGroupName(clusterName)
	}
	if s.RouteTable.Name == "" {
		s.RouteTable.Name = generateClusterRouteTableName(clusterName)
	}
	if s.ID == "" {
		if s.NatGateway.Name == "" {
			s.NatGateway.Name = generateClusterNatGatewayName(clusterName)
		}
		if !s.IsIPv6Enabled() && s.NatGateway.NatGatewayIP.Name == "" {
			s.NatGateway.NatGatewayIP.Name = generateNatGatewayIPName(s.NatGateway.Name)
		}
	}
	SubnetClassSpecSetDefaults(&s.SubnetClassSpec, DefaultClusterSubnetCIDR)
	SecurityGroupClassSetDefaults(&s.SecurityGroup.SecurityGroupClass)
}

// setDefaultAzureClusterVnetPeering sets default values for an AzureCluster's VNet peerings.
func setDefaultAzureClusterVnetPeering(c *infrav1.AzureCluster) {
	for i, peering := range c.Spec.NetworkSpec.Vnet.Peerings {
		if peering.ResourceGroup == "" {
			c.Spec.NetworkSpec.Vnet.Peerings[i].ResourceGroup = c.Spec.ResourceGroup
		}
	}
}

// setDefaultAzureClusterAPIServerLB sets default values for an AzureCluster's API server load balancer.
func setDefaultAzureClusterAPIServerLB(c *infrav1.AzureCluster) {
	if c.Spec.NetworkSpec.APIServerLB == nil {
		lbSpec := infrav1.LoadBalancerSpec{
			LoadBalancerClassSpec: infrav1.LoadBalancerClassSpec{
				Type: "Public",
			},
		}
		c.Spec.NetworkSpec.APIServerLB = &lbSpec
	}
	lb := c.Spec.NetworkSpec.APIServerLB

	setDefaultLoadBalancerClassSpecAPIServerLB(&lb.LoadBalancerClassSpec)

	if lb.Type == infrav1.Public {
		if lb.Name == "" {
			lb.Name = generatePublicLBName(c.ObjectMeta.Name)
		}
		if len(lb.FrontendIPs) == 0 {
			lb.FrontendIPs = []infrav1.FrontendIP{
				{
					Name: generateFrontendIPConfigName(lb.Name),
					PublicIP: &infrav1.PublicIPSpec{
						Name: generatePublicIPName(c.ObjectMeta.Name),
					},
				},
			}
		}
		// If the API Server ILB feature is enabled, create a default internal LB IP or use the specified one
		if feature.Gates.Enabled(feature.APIServerILB) {
			privateIPFound := false
			for i := range lb.FrontendIPs {
				if lb.FrontendIPs[i].FrontendIPClass.PrivateIPAddress != "" {
					if lb.FrontendIPs[i].Name == "" {
						lb.FrontendIPs[i].Name = generatePrivateIPConfigName(lb.Name)
					}
					privateIPFound = true
					break
				}
			}
			// if no private IP is found, we should create a default internal LB IP
			if !privateIPFound {
				privateIP := infrav1.FrontendIP{
					Name: generatePrivateIPConfigName(lb.Name),
					FrontendIPClass: infrav1.FrontendIPClass{
						PrivateIPAddress: DefaultInternalLBIPAddress,
					},
				}
				lb.FrontendIPs = append(lb.FrontendIPs, privateIP)
			}
		}
	} else if lb.Type == infrav1.Internal {
		if lb.Name == "" {
			lb.Name = generateInternalLBName(c.ObjectMeta.Name)
		}
		if len(lb.FrontendIPs) == 0 {
			lb.FrontendIPs = []infrav1.FrontendIP{
				{
					Name: generateFrontendIPConfigName(lb.Name),
					FrontendIPClass: infrav1.FrontendIPClass{
						PrivateIPAddress: DefaultInternalLBIPAddress,
					},
				},
			}
		}
	}
	setDefaultAzureClusterAPIServerLBBackendPoolName(c)
}

// setDefaultAzureClusterNodeOutboundLB sets the default values for the NodeOutboundLB.
func setDefaultAzureClusterNodeOutboundLB(c *infrav1.AzureCluster) {
	if c.Spec.NetworkSpec.NodeOutboundLB == nil {
		if !c.Spec.ControlPlaneEnabled || c.Spec.NetworkSpec.APIServerLB.Type == infrav1.Internal {
			return
		}

		var needsOutboundLB bool
		for _, subnet := range c.Spec.NetworkSpec.Subnets {
			if (subnet.Role == infrav1.SubnetNode || subnet.Role == infrav1.SubnetCluster) && subnet.IsIPv6Enabled() {
				needsOutboundLB = true
				break
			}
		}

		// If we don't default the outbound LB when there are some subnets with NAT gateway,
		// and some without, those without wouldn't have outbound traffic. So taking the
		// safer route, we configure the outbound LB in that scenario.
		if !needsOutboundLB {
			return
		}

		c.Spec.NetworkSpec.NodeOutboundLB = &infrav1.LoadBalancerSpec{}
	}

	lb := c.Spec.NetworkSpec.NodeOutboundLB
	setDefaultLoadBalancerClassSpecNodeOutboundLB(&lb.LoadBalancerClassSpec)

	if lb.Name == "" {
		lb.Name = c.ObjectMeta.Name
	}

	if lb.FrontendIPsCount == nil {
		lb.FrontendIPsCount = ptr.To[int32](1)
	}

	setDefaultAzureClusterOutboundLBFrontendIPs(c, lb, generateNodeOutboundIPName)
	setDefaultAzureClusterNodeOutboundLBBackendPoolName(c)
}

// setDefaultAzureClusterControlPlaneOutboundLB sets the default values for the control plane's outbound LB.
func setDefaultAzureClusterControlPlaneOutboundLB(c *infrav1.AzureCluster) {
	lb := c.Spec.NetworkSpec.ControlPlaneOutboundLB

	if lb == nil {
		return
	}

	setDefaultLoadBalancerClassSpecControlPlaneOutboundLB(&lb.LoadBalancerClassSpec)
	if lb.Name == "" {
		lb.Name = generateControlPlaneOutboundLBName(c.ObjectMeta.Name)
	}
	if lb.FrontendIPsCount == nil {
		lb.FrontendIPsCount = ptr.To[int32](1)
	}
	setDefaultAzureClusterOutboundLBFrontendIPs(c, lb, generateControlPlaneOutboundIPName)
	setDefaultAzureClusterControlPlaneOutboundLBBackendPoolName(c)
}

// SetDefaultAzureClusterBackendPoolName defaults the backend pool name of the LBs.
func SetDefaultAzureClusterBackendPoolName(c *infrav1.AzureCluster) {
	setDefaultAzureClusterAPIServerLBBackendPoolName(c)
	setDefaultAzureClusterNodeOutboundLBBackendPoolName(c)
	setDefaultAzureClusterControlPlaneOutboundLBBackendPoolName(c)
}

// setDefaultAzureClusterAPIServerLBBackendPoolName defaults the name of the backend pool for apiserver LB.
func setDefaultAzureClusterAPIServerLBBackendPoolName(c *infrav1.AzureCluster) {
	apiServerLB := c.Spec.NetworkSpec.APIServerLB
	if apiServerLB.BackendPool.Name == "" {
		apiServerLB.BackendPool.Name = generateBackendAddressPoolName(apiServerLB.Name)
	}
}

// setDefaultAzureClusterNodeOutboundLBBackendPoolName defaults the name of the backend pool for node outbound LB.
func setDefaultAzureClusterNodeOutboundLBBackendPoolName(c *infrav1.AzureCluster) {
	nodeOutboundLB := c.Spec.NetworkSpec.NodeOutboundLB
	if nodeOutboundLB != nil && nodeOutboundLB.BackendPool.Name == "" {
		nodeOutboundLB.BackendPool.Name = generateOutboundBackendAddressPoolName(nodeOutboundLB.Name)
	}
}

// setDefaultAzureClusterControlPlaneOutboundLBBackendPoolName defaults the name of the backend pool for control plane outbound LB.
func setDefaultAzureClusterControlPlaneOutboundLBBackendPoolName(c *infrav1.AzureCluster) {
	controlPlaneOutboundLB := c.Spec.NetworkSpec.ControlPlaneOutboundLB
	if controlPlaneOutboundLB != nil && controlPlaneOutboundLB.BackendPool.Name == "" {
		controlPlaneOutboundLB.BackendPool.Name = generateOutboundBackendAddressPoolName(generateControlPlaneOutboundLBName(c.ObjectMeta.Name))
	}
}

// setDefaultAzureClusterOutboundLBFrontendIPs sets the frontend IPs for the given load balancer.
// The name of the frontend IP is generated using generatePublicIPName function.
func setDefaultAzureClusterOutboundLBFrontendIPs(c *infrav1.AzureCluster, lb *infrav1.LoadBalancerSpec, generatePublicIPName func(string) string) {
	switch *lb.FrontendIPsCount {
	case 0:
		lb.FrontendIPs = []infrav1.FrontendIP{}
	case 1:
		lb.FrontendIPs = []infrav1.FrontendIP{
			{
				Name: generateFrontendIPConfigName(lb.Name),
				PublicIP: &infrav1.PublicIPSpec{
					Name: generatePublicIPName(c.ObjectMeta.Name),
				},
			},
		}
	default:
		lb.FrontendIPs = make([]infrav1.FrontendIP, *lb.FrontendIPsCount)
		for i := 0; i < int(*lb.FrontendIPsCount); i++ {
			lb.FrontendIPs[i] = infrav1.FrontendIP{
				Name: withIndex(generateFrontendIPConfigName(lb.Name), i+1),
				PublicIP: &infrav1.PublicIPSpec{
					Name: withIndex(generatePublicIPName(c.ObjectMeta.Name), i+1),
				},
			}
		}
	}
}

// setDefaultAzureClusterBastion sets default values for an AzureCluster's bastion configuration.
func setDefaultAzureClusterBastion(c *infrav1.AzureCluster) {
	if c.Spec.BastionSpec.AzureBastion != nil {
		if c.Spec.BastionSpec.AzureBastion.Name == "" {
			c.Spec.BastionSpec.AzureBastion.Name = generateAzureBastionName(c.ObjectMeta.Name)
		}
		// Ensure defaults for the Subnet settings.
		if c.Spec.BastionSpec.AzureBastion.Subnet.Name == "" {
			c.Spec.BastionSpec.AzureBastion.Subnet.Name = DefaultAzureBastionSubnetName
		}
		if len(c.Spec.BastionSpec.AzureBastion.Subnet.CIDRBlocks) == 0 {
			c.Spec.BastionSpec.AzureBastion.Subnet.CIDRBlocks = []string{DefaultAzureBastionSubnetCIDR}
		}
		if c.Spec.BastionSpec.AzureBastion.Subnet.Role == "" {
			c.Spec.BastionSpec.AzureBastion.Subnet.Role = DefaultAzureBastionSubnetRole
		}
		// Ensure defaults for the PublicIP settings.
		if c.Spec.BastionSpec.AzureBastion.PublicIP.Name == "" {
			c.Spec.BastionSpec.AzureBastion.PublicIP.Name = generateAzureBastionPublicIPName(c.ObjectMeta.Name)
		}
	}
}

// setDefaultLoadBalancerClassSpecAPIServerLB sets default values for an API server LoadBalancerClassSpec.
func setDefaultLoadBalancerClassSpecAPIServerLB(lb *infrav1.LoadBalancerClassSpec) {
	if lb.Type == "" {
		lb.Type = infrav1.Public
	}
	if lb.SKU == "" {
		lb.SKU = infrav1.SKUStandard
	}
	if lb.IdleTimeoutInMinutes == nil {
		lb.IdleTimeoutInMinutes = ptr.To[int32](DefaultOutboundRuleIdleTimeoutInMinutes)
	}
}

// setDefaultLoadBalancerClassSpecNodeOutboundLB sets default values for a node outbound LoadBalancerClassSpec.
func setDefaultLoadBalancerClassSpecNodeOutboundLB(lb *infrav1.LoadBalancerClassSpec) {
	setDefaultLoadBalancerClassSpecOutboundLB(lb)
}

// setDefaultLoadBalancerClassSpecControlPlaneOutboundLB sets default values for a control plane outbound LoadBalancerClassSpec.
func setDefaultLoadBalancerClassSpecControlPlaneOutboundLB(lb *infrav1.LoadBalancerClassSpec) {
	setDefaultLoadBalancerClassSpecOutboundLB(lb)
}

// setDefaultLoadBalancerClassSpecOutboundLB sets default values for an outbound LoadBalancerClassSpec.
func setDefaultLoadBalancerClassSpecOutboundLB(lb *infrav1.LoadBalancerClassSpec) {
	lb.Type = infrav1.Public
	lb.SKU = infrav1.SKUStandard
	if lb.IdleTimeoutInMinutes == nil {
		lb.IdleTimeoutInMinutes = ptr.To[int32](DefaultOutboundRuleIdleTimeoutInMinutes)
	}
}

// generateVnetName generates a virtual network name, based on the cluster name.
func generateVnetName(clusterName string) string {
	return fmt.Sprintf("%s-%s", clusterName, "vnet")
}

// generateClusterSubnetSubnetName generates a subnet name, based on the cluster name.
func generateClusterSubnetSubnetName(clusterName string) string {
	return fmt.Sprintf("%s-%s", clusterName, "subnet")
}

// generateControlPlaneSubnetName generates a node subnet name, based on the cluster name.
func generateControlPlaneSubnetName(clusterName string) string {
	return fmt.Sprintf("%s-%s", clusterName, "controlplane-subnet")
}

// generateNodeSubnetName generates a node subnet name, based on the cluster name.
func generateNodeSubnetName(clusterName string) string {
	return fmt.Sprintf("%s-%s", clusterName, "node-subnet")
}

// generateAzureBastionName generates an azure bastion name.
func generateAzureBastionName(clusterName string) string {
	return fmt.Sprintf("%s-azure-bastion", clusterName)
}

// generateAzureBastionPublicIPName generates an azure bastion public ip name.
func generateAzureBastionPublicIPName(clusterName string) string {
	return fmt.Sprintf("%s-azure-bastion-pip", clusterName)
}

// generateClusterSecurityGroupName generates a security group name, based on the cluster name.
func generateClusterSecurityGroupName(clusterName string) string {
	return fmt.Sprintf("%s-%s", clusterName, "nsg")
}

// generateControlPlaneSecurityGroupName generates a control plane security group name, based on the cluster name.
func generateControlPlaneSecurityGroupName(clusterName string) string {
	return fmt.Sprintf("%s-%s", clusterName, "controlplane-nsg")
}

// generateNodeSecurityGroupName generates a node security group name, based on the cluster name.
func generateNodeSecurityGroupName(clusterName string) string {
	return fmt.Sprintf("%s-%s", clusterName, "node-nsg")
}

// generateClusterRouteTableName generates a route table name, based on the cluster name.
func generateClusterRouteTableName(clusterName string) string {
	return fmt.Sprintf("%s-%s", clusterName, "routetable")
}

// generateNodeRouteTableName generates a node route table name, based on the cluster name.
func generateNodeRouteTableName(clusterName string) string {
	return fmt.Sprintf("%s-%s", clusterName, "node-routetable")
}

// generateInternalLBName generates a internal load balancer name, based on the cluster name.
func generateInternalLBName(clusterName string) string {
	return fmt.Sprintf("%s-%s", clusterName, "internal-lb")
}

// generatePublicLBName generates a public load balancer name, based on the cluster name.
func generatePublicLBName(clusterName string) string {
	return fmt.Sprintf("%s-%s", clusterName, "public-lb")
}

// generateControlPlaneOutboundLBName generates the name of the control plane outbound LB.
func generateControlPlaneOutboundLBName(clusterName string) string {
	return fmt.Sprintf("%s-outbound-lb", clusterName)
}

// generatePublicIPName generates a public IP name, based on the cluster name and a hash.
func generatePublicIPName(clusterName string) string {
	return fmt.Sprintf("pip-%s-apiserver", clusterName)
}

// generateFrontendIPConfigName generates a load balancer frontend IP config name.
func generateFrontendIPConfigName(lbName string) string {
	return fmt.Sprintf("%s-%s", lbName, "frontEnd")
}

// generatePrivateIPConfigName generates a load balancer frontend private IP config name.
func generatePrivateIPConfigName(lbName string) string {
	return fmt.Sprintf("%s-%s", lbName, "frontEnd-internal-ip")
}

// generateNodeOutboundIPName generates a public IP name, based on the cluster name.
func generateNodeOutboundIPName(clusterName string) string {
	return fmt.Sprintf("pip-%s-node-outbound", clusterName)
}

// generateControlPlaneOutboundIPName generates a public IP name, based on the cluster name.
func generateControlPlaneOutboundIPName(clusterName string) string {
	return fmt.Sprintf("pip-%s-controlplane-outbound", clusterName)
}

// generateClusterNatGatewayName generates a NAT gateway name.
func generateClusterNatGatewayName(clusterName string) string {
	return fmt.Sprintf("%s-%s", clusterName, "natgw")
}

// generateNatGatewayName generates a NAT gateway name.
func generateNatGatewayName(clusterName string) string {
	return fmt.Sprintf("%s-%s", clusterName, "node-natgw")
}

// generateNatGatewayIPName generates a NAT gateway IP name.
func generateNatGatewayIPName(natGatewayName string) string {
	return fmt.Sprintf("pip-%s", natGatewayName)
}

// withIndex appends the index as suffix to a generated name.
func withIndex(name string, n int) string {
	return fmt.Sprintf("%s-%d", name, n)
}

// generateBackendAddressPoolName generates a load balancer backend address pool name.
func generateBackendAddressPoolName(lbName string) string {
	return fmt.Sprintf("%s-%s", lbName, "backendPool")
}

// generateOutboundBackendAddressPoolName generates a load balancer outbound backend address pool name.
func generateOutboundBackendAddressPoolName(lbName string) string {
	return fmt.Sprintf("%s-%s", lbName, "outboundBackendPool")
}
