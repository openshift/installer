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

import (
	"fmt"

	"k8s.io/utils/ptr"

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
	DefaultAzureBastionSubnetRole = SubnetBastion
	// DefaultInternalLBIPAddress is the default internal load balancer ip address.
	DefaultInternalLBIPAddress = "10.0.0.100"
	// DefaultOutboundRuleIdleTimeoutInMinutes is the default for IdleTimeoutInMinutes for the load balancer.
	DefaultOutboundRuleIdleTimeoutInMinutes = 4
	// DefaultAzureCloud is the public cloud that will be used by most users.
	DefaultAzureCloud = "AzurePublicCloud"
)

func (c *AzureCluster) setDefaults() {
	c.Spec.AzureClusterClassSpec.setDefaults()
	c.setResourceGroupDefault()
	c.setNetworkSpecDefaults()
}

func (c *AzureCluster) setNetworkSpecDefaults() {
	c.setVnetDefaults()
	c.setBastionDefaults()
	c.setSubnetDefaults()
	c.setVnetPeeringDefaults()
	if c.Spec.ControlPlaneEnabled {
		c.setAPIServerLBDefaults()
	}
	c.SetNodeOutboundLBDefaults()
	if c.Spec.ControlPlaneEnabled {
		c.SetControlPlaneOutboundLBDefaults()
	}
	if !c.Spec.ControlPlaneEnabled {
		c.Spec.NetworkSpec.APIServerLB = nil
	}
}

func (c *AzureCluster) setResourceGroupDefault() {
	if c.Spec.ResourceGroup == "" {
		c.Spec.ResourceGroup = c.Name
	}
}

func (c *AzureCluster) setAzureEnvironmentDefault() {
	if c.Spec.AzureEnvironment == "" {
		c.Spec.AzureEnvironment = DefaultAzureCloud
	}
}

func (c *AzureCluster) setVnetDefaults() {
	if c.Spec.NetworkSpec.Vnet.ResourceGroup == "" {
		c.Spec.NetworkSpec.Vnet.ResourceGroup = c.Spec.ResourceGroup
	}
	if c.Spec.NetworkSpec.Vnet.Name == "" {
		c.Spec.NetworkSpec.Vnet.Name = generateVnetName(c.ObjectMeta.Name)
	}
	c.Spec.NetworkSpec.Vnet.VnetClassSpec.setDefaults()
}

func (c *AzureCluster) setSubnetDefaults() {
	clusterSubnet, err := c.Spec.NetworkSpec.GetSubnet(SubnetCluster)
	clusterSubnetExists := err == nil
	if clusterSubnetExists {
		clusterSubnet.setClusterSubnetDefaults(c.ObjectMeta.Name)
		c.Spec.NetworkSpec.UpdateSubnet(clusterSubnet, SubnetCluster)
	}

	if c.Spec.ControlPlaneEnabled {
		/* if there is a cp subnet set defaults
		   if no cp subnet and cluster subnet create a default cp subnet */
		cpSubnet, errcp := c.Spec.NetworkSpec.GetSubnet(SubnetControlPlane)
		if errcp == nil {
			cpSubnet.setControlPlaneSubnetDefaults(c.ObjectMeta.Name)
			c.Spec.NetworkSpec.UpdateSubnet(cpSubnet, SubnetControlPlane)
		} else if !clusterSubnetExists {
			cpSubnet = SubnetSpec{SubnetClassSpec: SubnetClassSpec{Role: SubnetControlPlane}}
			cpSubnet.setControlPlaneSubnetDefaults(c.ObjectMeta.Name)
			c.Spec.NetworkSpec.Subnets = append(c.Spec.NetworkSpec.Subnets, cpSubnet)
		}
	}

	var nodeSubnetFound bool
	var nodeSubnetCounter int
	for i, subnet := range c.Spec.NetworkSpec.Subnets {
		if subnet.Role != SubnetNode {
			continue
		}
		nodeSubnetCounter++
		nodeSubnetFound = true
		subnet.setNodeSubnetDefaults(c.ObjectMeta.Name, nodeSubnetCounter)
		c.Spec.NetworkSpec.Subnets[i] = subnet
	}

	if !nodeSubnetFound && !clusterSubnetExists {
		nodeSubnet := SubnetSpec{
			SubnetClassSpec: SubnetClassSpec{
				Role:       SubnetNode,
				CIDRBlocks: []string{DefaultNodeSubnetCIDR},
				Name:       generateNodeSubnetName(c.ObjectMeta.Name),
			},
			SecurityGroup: SecurityGroup{
				Name: generateNodeSecurityGroupName(c.ObjectMeta.Name),
			},
			RouteTable: RouteTable{
				Name: generateNodeRouteTableName(c.ObjectMeta.Name),
			},
			NatGateway: NatGateway{
				NatGatewayClassSpec: NatGatewayClassSpec{
					Name: generateNatGatewayName(c.ObjectMeta.Name),
				},
			},
		}
		c.Spec.NetworkSpec.Subnets = append(c.Spec.NetworkSpec.Subnets, nodeSubnet)
	}
}

func (s *SubnetSpec) setNodeSubnetDefaults(clusterName string, index int) {
	if s.Name == "" {
		s.Name = withIndex(generateNodeSubnetName(clusterName), index)
	}
	s.SubnetClassSpec.setDefaults(fmt.Sprintf(DefaultNodeSubnetCIDRPattern, index))

	if s.SecurityGroup.Name == "" {
		s.SecurityGroup.Name = generateNodeSecurityGroupName(clusterName)
	}
	s.SecurityGroup.SecurityGroupClass.setDefaults()

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

func (s *SubnetSpec) setControlPlaneSubnetDefaults(clusterName string) {
	if s.Name == "" {
		s.Name = generateControlPlaneSubnetName(clusterName)
	}

	s.SubnetClassSpec.setDefaults(DefaultControlPlaneSubnetCIDR)

	if s.SecurityGroup.Name == "" {
		s.SecurityGroup.Name = generateControlPlaneSecurityGroupName(clusterName)
	}
	s.SecurityGroup.SecurityGroupClass.setDefaults()
}

func (s *SubnetSpec) setClusterSubnetDefaults(clusterName string) {
	if s.Name == "" {
		s.Name = generateClusterSubnetSubnetName(clusterName)
	}
	if s.SecurityGroup.Name == "" {
		s.SecurityGroup.Name = generateClusterSecurityGroupName(clusterName)
	}
	if s.RouteTable.Name == "" {
		s.RouteTable.Name = generateClustereRouteTableName(clusterName)
	}
	if s.NatGateway.Name == "" {
		s.NatGateway.Name = generateClusterNatGatewayName(clusterName)
	}
	if !s.IsIPv6Enabled() && s.ID == "" && s.NatGateway.NatGatewayIP.Name == "" {
		s.NatGateway.NatGatewayIP.Name = generateNatGatewayIPName(s.NatGateway.Name)
	}
	s.setDefaults(DefaultClusterSubnetCIDR)
	s.SecurityGroup.SecurityGroupClass.setDefaults()
}

func (c *AzureCluster) setVnetPeeringDefaults() {
	for i, peering := range c.Spec.NetworkSpec.Vnet.Peerings {
		if peering.ResourceGroup == "" {
			c.Spec.NetworkSpec.Vnet.Peerings[i].ResourceGroup = c.Spec.ResourceGroup
		}
	}
}

func (c *AzureCluster) setAPIServerLBDefaults() {
	if c.Spec.NetworkSpec.APIServerLB == nil {
		lbSpec := LoadBalancerSpec{
			LoadBalancerClassSpec: LoadBalancerClassSpec{
				Type: "Public",
			},
		}
		c.Spec.NetworkSpec.APIServerLB = &lbSpec
	}
	lb := c.Spec.NetworkSpec.APIServerLB

	lb.LoadBalancerClassSpec.setAPIServerLBDefaults()

	if lb.Type == Public {
		if lb.Name == "" {
			lb.Name = generatePublicLBName(c.ObjectMeta.Name)
		}
		if len(lb.FrontendIPs) == 0 {
			lb.FrontendIPs = []FrontendIP{
				{
					Name: generateFrontendIPConfigName(lb.Name),
					PublicIP: &PublicIPSpec{
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
				privateIP := FrontendIP{
					Name: generatePrivateIPConfigName(lb.Name),
					FrontendIPClass: FrontendIPClass{
						PrivateIPAddress: DefaultInternalLBIPAddress,
					},
				}
				lb.FrontendIPs = append(lb.FrontendIPs, privateIP)
			}
		}
	} else if lb.Type == Internal {
		if lb.Name == "" {
			lb.Name = generateInternalLBName(c.ObjectMeta.Name)
		}
		if len(lb.FrontendIPs) == 0 {
			lb.FrontendIPs = []FrontendIP{
				{
					Name: generateFrontendIPConfigName(lb.Name),
					FrontendIPClass: FrontendIPClass{
						PrivateIPAddress: DefaultInternalLBIPAddress,
					},
				},
			}
		}
	}
	c.SetAPIServerLBBackendPoolNameDefault()
}

// SetNodeOutboundLBDefaults sets the default values for the NodeOutboundLB.
func (c *AzureCluster) SetNodeOutboundLBDefaults() {
	if c.Spec.NetworkSpec.NodeOutboundLB == nil {
		if !c.Spec.ControlPlaneEnabled || c.Spec.NetworkSpec.APIServerLB.Type == Internal {
			return
		}

		var needsOutboundLB bool
		for _, subnet := range c.Spec.NetworkSpec.Subnets {
			if (subnet.Role == SubnetNode || subnet.Role == SubnetCluster) && subnet.IsIPv6Enabled() {
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

		c.Spec.NetworkSpec.NodeOutboundLB = &LoadBalancerSpec{}
	}

	lb := c.Spec.NetworkSpec.NodeOutboundLB
	lb.LoadBalancerClassSpec.setNodeOutboundLBDefaults()

	if lb.Name == "" {
		lb.Name = c.ObjectMeta.Name
	}

	if lb.FrontendIPsCount == nil {
		lb.FrontendIPsCount = ptr.To[int32](1)
	}

	c.setOutboundLBFrontendIPs(lb, generateNodeOutboundIPName)
	c.SetNodeOutboundLBBackendPoolNameDefault()
}

// SetControlPlaneOutboundLBDefaults sets the default values for the control plane's outbound LB.
func (c *AzureCluster) SetControlPlaneOutboundLBDefaults() {
	lb := c.Spec.NetworkSpec.ControlPlaneOutboundLB

	if lb == nil {
		return
	}

	lb.LoadBalancerClassSpec.setControlPlaneOutboundLBDefaults()
	if lb.Name == "" {
		lb.Name = generateControlPlaneOutboundLBName(c.ObjectMeta.Name)
	}
	if lb.FrontendIPsCount == nil {
		lb.FrontendIPsCount = ptr.To[int32](1)
	}
	c.setOutboundLBFrontendIPs(lb, generateControlPlaneOutboundIPName)
	c.SetControlPlaneOutboundLBBackendPoolNameDefault()
}

// SetBackendPoolNameDefault defaults the backend pool name of the LBs.
func (c *AzureCluster) SetBackendPoolNameDefault() {
	c.SetAPIServerLBBackendPoolNameDefault()
	c.SetNodeOutboundLBBackendPoolNameDefault()
	c.SetControlPlaneOutboundLBBackendPoolNameDefault()
}

// SetAPIServerLBBackendPoolNameDefault defaults the name of the backend pool for apiserver LB.
func (c *AzureCluster) SetAPIServerLBBackendPoolNameDefault() {
	apiServerLB := c.Spec.NetworkSpec.APIServerLB
	if apiServerLB.BackendPool.Name == "" {
		apiServerLB.BackendPool.Name = generateBackendAddressPoolName(apiServerLB.Name)
	}
}

// SetNodeOutboundLBBackendPoolNameDefault defaults the name of the backend pool for node outbound LB.
func (c *AzureCluster) SetNodeOutboundLBBackendPoolNameDefault() {
	nodeOutboundLB := c.Spec.NetworkSpec.NodeOutboundLB
	if nodeOutboundLB != nil && nodeOutboundLB.BackendPool.Name == "" {
		nodeOutboundLB.BackendPool.Name = generateOutboundBackendAddressPoolName(nodeOutboundLB.Name)
	}
}

// SetControlPlaneOutboundLBBackendPoolNameDefault defaults the name of the backend pool for control plane outbound LB.
func (c *AzureCluster) SetControlPlaneOutboundLBBackendPoolNameDefault() {
	controlPlaneOutboundLB := c.Spec.NetworkSpec.ControlPlaneOutboundLB
	if controlPlaneOutboundLB != nil && controlPlaneOutboundLB.BackendPool.Name == "" {
		controlPlaneOutboundLB.BackendPool.Name = generateOutboundBackendAddressPoolName(generateControlPlaneOutboundLBName(c.ObjectMeta.Name))
	}
}

// setOutboundLBFrontendIPs sets the frontend ips for the given load balancer.
// The name of the frontend ip is generated using generatePublicIPName function.
func (c *AzureCluster) setOutboundLBFrontendIPs(lb *LoadBalancerSpec, generatePublicIPName func(string) string) {
	switch *lb.FrontendIPsCount {
	case 0:
		lb.FrontendIPs = []FrontendIP{}
	case 1:
		lb.FrontendIPs = []FrontendIP{
			{
				Name: generateFrontendIPConfigName(lb.Name),
				PublicIP: &PublicIPSpec{
					Name: generatePublicIPName(c.ObjectMeta.Name),
				},
			},
		}
	default:
		lb.FrontendIPs = make([]FrontendIP, *lb.FrontendIPsCount)
		for i := 0; i < int(*lb.FrontendIPsCount); i++ {
			lb.FrontendIPs[i] = FrontendIP{
				Name: withIndex(generateFrontendIPConfigName(lb.Name), i+1),
				PublicIP: &PublicIPSpec{
					Name: withIndex(generatePublicIPName(c.ObjectMeta.Name), i+1),
				},
			}
		}
	}
}

func (c *AzureCluster) setBastionDefaults() {
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

func (lb *LoadBalancerClassSpec) setAPIServerLBDefaults() {
	if lb.Type == "" {
		lb.Type = Public
	}
	if lb.SKU == "" {
		lb.SKU = SKUStandard
	}
	if lb.IdleTimeoutInMinutes == nil {
		lb.IdleTimeoutInMinutes = ptr.To[int32](DefaultOutboundRuleIdleTimeoutInMinutes)
	}
}

func (lb *LoadBalancerClassSpec) setNodeOutboundLBDefaults() {
	lb.setOutboundLBDefaults()
}

func (lb *LoadBalancerClassSpec) setControlPlaneOutboundLBDefaults() {
	lb.setOutboundLBDefaults()
}

func (lb *LoadBalancerClassSpec) setOutboundLBDefaults() {
	lb.Type = Public
	lb.SKU = SKUStandard
	if lb.IdleTimeoutInMinutes == nil {
		lb.IdleTimeoutInMinutes = ptr.To[int32](DefaultOutboundRuleIdleTimeoutInMinutes)
	}
}

func setControlPlaneOutboundLBDefaults(lb *LoadBalancerClassSpec, apiserverLBType LBType) {
	// public clusters don't need control plane outbound lb
	if apiserverLBType == Public {
		return
	}

	// private clusters can disable control plane outbound lb by setting it to nil.
	if lb == nil {
		return
	}

	lb.Type = Public
	lb.SKU = SKUStandard

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

// generateClustereRouteTableName generates a route table name, based on the cluster name.
func generateClustereRouteTableName(clusterName string) string {
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

// generateFrontendIPConfigName generates a load balancer frontend IP config name.
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
