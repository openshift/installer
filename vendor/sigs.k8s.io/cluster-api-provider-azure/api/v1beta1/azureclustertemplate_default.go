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
	"fmt"
)

func (c *AzureClusterTemplate) setDefaults() {
	c.Spec.Template.Spec.AzureClusterClassSpec.setDefaults()
	c.setNetworkTemplateSpecDefaults()
}

func (c *AzureClusterTemplate) setNetworkTemplateSpecDefaults() {
	c.setVnetTemplateDefaults()
	c.setBastionTemplateDefaults()
	c.setSubnetsTemplateDefaults()

	apiServerLB := &c.Spec.Template.Spec.NetworkSpec.APIServerLB
	apiServerLB.setAPIServerLBDefaults()
	c.setNodeOutboundLBDefaults()
	c.setControlPlaneOutboundLBDefaults()
}

func (c *AzureClusterTemplate) setVnetTemplateDefaults() {
	c.Spec.Template.Spec.NetworkSpec.Vnet.VnetClassSpec.setDefaults()
}

func (c *AzureClusterTemplate) setBastionTemplateDefaults() {
	if c.Spec.Template.Spec.BastionSpec.AzureBastion != nil {
		// Ensure defaults for Subnet settings.
		if len(c.Spec.Template.Spec.BastionSpec.AzureBastion.Subnet.CIDRBlocks) == 0 {
			c.Spec.Template.Spec.BastionSpec.AzureBastion.Subnet.CIDRBlocks = []string{DefaultAzureBastionSubnetCIDR}
		}
		if c.Spec.Template.Spec.BastionSpec.AzureBastion.Subnet.Role == "" {
			c.Spec.Template.Spec.BastionSpec.AzureBastion.Subnet.Role = DefaultAzureBastionSubnetRole
		}
	}
}

func (c *AzureClusterTemplate) setSubnetsTemplateDefaults() {
	cpSubnet, err := c.Spec.Template.Spec.NetworkSpec.GetControlPlaneSubnetTemplate()
	if err != nil {
		cpSubnet = SubnetTemplateSpec{SubnetClassSpec: SubnetClassSpec{Role: SubnetControlPlane}}
		c.Spec.Template.Spec.NetworkSpec.Subnets = append(c.Spec.Template.Spec.NetworkSpec.Subnets, cpSubnet)
	}
	cpSubnet.SubnetClassSpec.setDefaults(DefaultControlPlaneSubnetCIDR)
	cpSubnet.SecurityGroup.setDefaults()
	c.Spec.Template.Spec.NetworkSpec.UpdateControlPlaneSubnetTemplate(cpSubnet)

	var nodeSubnetFound bool
	var nodeSubnetCounter int
	for i, subnet := range c.Spec.Template.Spec.NetworkSpec.Subnets {
		if subnet.Role != SubnetNode {
			continue
		}
		nodeSubnetCounter++
		nodeSubnetFound = true
		subnet.SubnetClassSpec.setDefaults(fmt.Sprintf(DefaultNodeSubnetCIDRPattern, nodeSubnetCounter))
		subnet.SecurityGroup.setDefaults()
		c.Spec.Template.Spec.NetworkSpec.Subnets[i] = subnet
	}

	if !nodeSubnetFound {
		nodeSubnet := SubnetTemplateSpec{
			SubnetClassSpec: SubnetClassSpec{
				Role:       SubnetNode,
				CIDRBlocks: []string{DefaultNodeSubnetCIDR},
			},
		}
		c.Spec.Template.Spec.NetworkSpec.Subnets = append(c.Spec.Template.Spec.NetworkSpec.Subnets, nodeSubnet)
	}
}

func (c *AzureClusterTemplate) setNodeOutboundLBDefaults() {
	if c.Spec.Template.Spec.NetworkSpec.NodeOutboundLB == nil {
		if c.Spec.Template.Spec.NetworkSpec.APIServerLB.Type == Internal {
			return
		}

		var needsOutboundLB bool
		for _, subnet := range c.Spec.Template.Spec.NetworkSpec.Subnets {
			if subnet.Role == SubnetNode && subnet.IsIPv6Enabled() {
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

		c.Spec.Template.Spec.NetworkSpec.NodeOutboundLB = &LoadBalancerClassSpec{}
	}

	c.Spec.Template.Spec.NetworkSpec.NodeOutboundLB.setNodeOutboundLBDefaults()
}

func (c *AzureClusterTemplate) setControlPlaneOutboundLBDefaults() {
	lb := c.Spec.Template.Spec.NetworkSpec.ControlPlaneOutboundLB
	if lb == nil {
		return
	}
	lb.setControlPlaneOutboundLBDefaults()
}
