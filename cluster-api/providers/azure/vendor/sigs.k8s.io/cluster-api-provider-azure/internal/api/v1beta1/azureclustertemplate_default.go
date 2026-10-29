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

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
)

// SetDefaultsAzureClusterTemplate sets default values for an AzureClusterTemplate.
func SetDefaultsAzureClusterTemplate(c *infrav1.AzureClusterTemplate) {
	AzureClusterClassSpecSetDefaults(&c.Spec.Template.Spec.AzureClusterClassSpec)
	setDefaultAzureClusterTemplateNetworkTemplateSpec(c)
}

// setDefaultAzureClusterTemplateNetworkTemplateSpec sets default values for an AzureClusterTemplate's NetworkTemplateSpec.
func setDefaultAzureClusterTemplateNetworkTemplateSpec(c *infrav1.AzureClusterTemplate) {
	setDefaultAzureClusterTemplateVnetTemplate(c)
	setDefaultAzureClusterTemplateBastionTemplate(c)
	setDefaultAzureClusterTemplateSubnetsTemplate(c)

	apiServerLB := &c.Spec.Template.Spec.NetworkSpec.APIServerLB
	setDefaultLoadBalancerClassSpecAPIServerLB(apiServerLB)
	setDefaultAzureClusterTemplateNodeOutboundLB(c)
	setDefaultAzureClusterTemplateControlPlaneOutboundLB(c)
}

// setDefaultAzureClusterTemplateVnetTemplate sets default values for an AzureClusterTemplate's VNet template.
func setDefaultAzureClusterTemplateVnetTemplate(c *infrav1.AzureClusterTemplate) {
	VnetClassSpecSetDefaults(&c.Spec.Template.Spec.NetworkSpec.Vnet.VnetClassSpec)
}

// setDefaultAzureClusterTemplateBastionTemplate sets default values for an AzureClusterTemplate's bastion template.
func setDefaultAzureClusterTemplateBastionTemplate(c *infrav1.AzureClusterTemplate) {
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

// setDefaultAzureClusterTemplateSubnetsTemplate sets default values for an AzureClusterTemplate's subnet templates.
func setDefaultAzureClusterTemplateSubnetsTemplate(c *infrav1.AzureClusterTemplate) {
	clusterSubnet, err := c.Spec.Template.Spec.NetworkSpec.GetSubnetTemplate(infrav1.SubnetCluster)
	clusterSubnetExists := err == nil
	if clusterSubnetExists {
		SubnetClassSpecSetDefaults(&clusterSubnet.SubnetClassSpec, DefaultClusterSubnetCIDR)
		SecurityGroupClassSetDefaults(&clusterSubnet.SecurityGroup)
		c.Spec.Template.Spec.NetworkSpec.UpdateSubnetTemplate(clusterSubnet, infrav1.SubnetCluster)
	}

	cpSubnet, errcp := c.Spec.Template.Spec.NetworkSpec.GetSubnetTemplate(infrav1.SubnetControlPlane)
	if errcp == nil {
		SubnetClassSpecSetDefaults(&cpSubnet.SubnetClassSpec, DefaultControlPlaneSubnetCIDR)
		SecurityGroupClassSetDefaults(&cpSubnet.SecurityGroup)
		c.Spec.Template.Spec.NetworkSpec.UpdateSubnetTemplate(cpSubnet, infrav1.SubnetControlPlane)
	} else if errcp != nil && !clusterSubnetExists {
		cpSubnet = infrav1.SubnetTemplateSpec{SubnetClassSpec: infrav1.SubnetClassSpec{Role: infrav1.SubnetControlPlane}}
		SubnetClassSpecSetDefaults(&cpSubnet.SubnetClassSpec, DefaultControlPlaneSubnetCIDR)
		SecurityGroupClassSetDefaults(&cpSubnet.SecurityGroup)
		c.Spec.Template.Spec.NetworkSpec.Subnets = append(c.Spec.Template.Spec.NetworkSpec.Subnets, cpSubnet)
	}

	var nodeSubnetFound bool
	var nodeSubnetCounter int
	for i, subnet := range c.Spec.Template.Spec.NetworkSpec.Subnets {
		if subnet.Role != infrav1.SubnetNode {
			continue
		}
		nodeSubnetCounter++
		nodeSubnetFound = true
		SubnetClassSpecSetDefaults(&subnet.SubnetClassSpec, fmt.Sprintf(DefaultNodeSubnetCIDRPattern, nodeSubnetCounter))
		SecurityGroupClassSetDefaults(&subnet.SecurityGroup)
		c.Spec.Template.Spec.NetworkSpec.Subnets[i] = subnet
	}

	if !nodeSubnetFound && !clusterSubnetExists {
		nodeSubnet := infrav1.SubnetTemplateSpec{
			SubnetClassSpec: infrav1.SubnetClassSpec{
				Role:       infrav1.SubnetNode,
				CIDRBlocks: []string{DefaultNodeSubnetCIDR},
			},
		}
		c.Spec.Template.Spec.NetworkSpec.Subnets = append(c.Spec.Template.Spec.NetworkSpec.Subnets, nodeSubnet)
	}
}

// setDefaultAzureClusterTemplateNodeOutboundLB sets default values for an AzureClusterTemplate's node outbound LB.
func setDefaultAzureClusterTemplateNodeOutboundLB(c *infrav1.AzureClusterTemplate) {
	if c.Spec.Template.Spec.NetworkSpec.NodeOutboundLB == nil {
		if c.Spec.Template.Spec.NetworkSpec.APIServerLB.Type == infrav1.Internal {
			return
		}

		var needsOutboundLB bool
		for _, subnet := range c.Spec.Template.Spec.NetworkSpec.Subnets {
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

		c.Spec.Template.Spec.NetworkSpec.NodeOutboundLB = &infrav1.LoadBalancerClassSpec{}
	}

	setDefaultLoadBalancerClassSpecNodeOutboundLB(c.Spec.Template.Spec.NetworkSpec.NodeOutboundLB)
}

// setDefaultAzureClusterTemplateControlPlaneOutboundLB sets default values for an AzureClusterTemplate's control plane outbound LB.
func setDefaultAzureClusterTemplateControlPlaneOutboundLB(c *infrav1.AzureClusterTemplate) {
	lb := c.Spec.Template.Spec.NetworkSpec.ControlPlaneOutboundLB
	if lb == nil {
		return
	}
	setDefaultLoadBalancerClassSpecControlPlaneOutboundLB(lb)
}
