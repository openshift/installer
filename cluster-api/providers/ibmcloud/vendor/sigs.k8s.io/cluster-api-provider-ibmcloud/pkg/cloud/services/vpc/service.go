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

package vpc

import (
	"fmt"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"

	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/authenticator"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/utils"
)

// SecurityGroupByNameNotFound represents an error when security group is not found by name.
type SecurityGroupByNameNotFound struct {
	Name string
}

func (s *SecurityGroupByNameNotFound) Error() string {
	return fmt.Sprintf("failed to find security group by name: %s", s.Name)
}

// Service holds the VPC Service specific information.
type Service struct {
	vpcService *vpcv1.VpcV1
}

// CreateInstance created an virtal server instance.
func (s *Service) CreateInstance(options *vpcv1.CreateInstanceOptions) (*vpcv1.Instance, *core.DetailedResponse, error) {
	return s.vpcService.CreateInstance(options)
}

// DeleteInstance deleted a virtal server instance.
func (s *Service) DeleteInstance(options *vpcv1.DeleteInstanceOptions) (*core.DetailedResponse, error) {
	return s.vpcService.DeleteInstance(options)
}

// GetInstance returns the virtal server instance.
func (s *Service) GetInstance(options *vpcv1.GetInstanceOptions) (*vpcv1.Instance, *core.DetailedResponse, error) {
	return s.vpcService.GetInstance(options)
}

// ListInstances returns list of virtual server instances.
func (s *Service) ListInstances(options *vpcv1.ListInstancesOptions) (*vpcv1.InstanceCollection, *core.DetailedResponse, error) {
	return s.vpcService.ListInstances(options)
}

// GetDedicatedHostByName returns Dedicated Host with given name. If not found, returns nil.
func (s *Service) GetDedicatedHostByName(dHostName string) (*vpcv1.DedicatedHost, error) {
	var dHost *vpcv1.DedicatedHost
	f := func(start string) (bool, string, error) {
		// check for existing Dedicated Hosts
		listDedicatedHostsOptions := &vpcv1.ListDedicatedHostsOptions{}
		if start != "" {
			listDedicatedHostsOptions.Start = &start
		}

		dHostsList, _, err := s.vpcService.ListDedicatedHosts(listDedicatedHostsOptions)
		if err != nil {
			return false, "", err
		}

		if dHostsList == nil {
			return false, "", fmt.Errorf("dedicated hosts list returned is nil")
		}

		for index, dH := range dHostsList.DedicatedHosts {
			if *dH.Name == dHostName {
				dHost = &dHostsList.DedicatedHosts[index]
				return true, "", nil
			}
		}

		if dHostsList.Next != nil && *dHostsList.Next.Href != "" {
			return false, *dHostsList.Next.Href, nil
		}
		return true, "", nil
	}

	if err := utils.PagingHelper(f); err != nil {
		return nil, err
	}

	return dHost, nil
}

// CreateVPC creates a new VPC.
func (s *Service) CreateVPC(options *vpcv1.CreateVPCOptions) (*vpcv1.VPC, *core.DetailedResponse, error) {
	return s.vpcService.CreateVPC(options)
}

// DeleteVPC deletes a VPC.
func (s *Service) DeleteVPC(options *vpcv1.DeleteVPCOptions) (*core.DetailedResponse, error) {
	return s.vpcService.DeleteVPC(options)
}

// ListVpcs returns list of VPCs in a region.
func (s *Service) ListVpcs(options *vpcv1.ListVpcsOptions) (*vpcv1.VPCCollection, *core.DetailedResponse, error) {
	return s.vpcService.ListVpcs(options)
}

// CreateSubnet creates a subnet.
func (s *Service) CreateSubnet(options *vpcv1.CreateSubnetOptions) (*vpcv1.Subnet, *core.DetailedResponse, error) {
	return s.vpcService.CreateSubnet(options)
}

// DeleteSubnet deletes a subnet.
func (s *Service) DeleteSubnet(options *vpcv1.DeleteSubnetOptions) (*core.DetailedResponse, error) {
	return s.vpcService.DeleteSubnet(options)
}

// ListSubnets returns list of subnets in a region.
func (s *Service) ListSubnets(options *vpcv1.ListSubnetsOptions) (*vpcv1.SubnetCollection, *core.DetailedResponse, error) {
	return s.vpcService.ListSubnets(options)
}

// GetSubnetPublicGateway returns a public gateway attached to the subnet.
func (s *Service) GetSubnetPublicGateway(options *vpcv1.GetSubnetPublicGatewayOptions) (*vpcv1.PublicGateway, *core.DetailedResponse, error) {
	return s.vpcService.GetSubnetPublicGateway(options)
}

// CreatePublicGateway creates a public gateway for the VPC.
func (s *Service) CreatePublicGateway(options *vpcv1.CreatePublicGatewayOptions) (*vpcv1.PublicGateway, *core.DetailedResponse, error) {
	return s.vpcService.CreatePublicGateway(options)
}

// DeletePublicGateway deletes a public gateway.
func (s *Service) DeletePublicGateway(options *vpcv1.DeletePublicGatewayOptions) (*core.DetailedResponse, error) {
	return s.vpcService.DeletePublicGateway(options)
}

// UnsetSubnetPublicGateway detaches a public gateway from the subnet.
func (s *Service) UnsetSubnetPublicGateway(options *vpcv1.UnsetSubnetPublicGatewayOptions) (*core.DetailedResponse, error) {
	return s.vpcService.UnsetSubnetPublicGateway(options)
}

// SetSubnetPublicGateway attaches a public gateway to the subnet.
func (s *Service) SetSubnetPublicGateway(options *vpcv1.SetSubnetPublicGatewayOptions) (*vpcv1.PublicGateway, *core.DetailedResponse, error) {
	return s.vpcService.SetSubnetPublicGateway(options)
}

// ListVPCAddressPrefixes returns list of all address prefixes for a VPC.
func (s *Service) ListVPCAddressPrefixes(options *vpcv1.ListVPCAddressPrefixesOptions) (*vpcv1.AddressPrefixCollection, *core.DetailedResponse, error) {
	return s.vpcService.ListVPCAddressPrefixes(options)
}

// CreateSecurityGroupRule creates a rule for a security group.
func (s *Service) CreateSecurityGroupRule(options *vpcv1.CreateSecurityGroupRuleOptions) (vpcv1.SecurityGroupRuleIntf, *core.DetailedResponse, error) {
	return s.vpcService.CreateSecurityGroupRule(options)
}

// CreateLoadBalancer creates a new load balancer.
func (s *Service) CreateLoadBalancer(options *vpcv1.CreateLoadBalancerOptions) (*vpcv1.LoadBalancer, *core.DetailedResponse, error) {
	return s.vpcService.CreateLoadBalancer(options)
}

// DeleteLoadBalancer deletes a load balancer.
func (s *Service) DeleteLoadBalancer(options *vpcv1.DeleteLoadBalancerOptions) (*core.DetailedResponse, error) {
	return s.vpcService.DeleteLoadBalancer(options)
}

// ListLoadBalancers returns list of load balancers in a region.
func (s *Service) ListLoadBalancers(options *vpcv1.ListLoadBalancersOptions) (*vpcv1.LoadBalancerCollection, *core.DetailedResponse, error) {
	return s.vpcService.ListLoadBalancers(options)
}

// GetLoadBalancer returns a load balancer.
func (s *Service) GetLoadBalancer(options *vpcv1.GetLoadBalancerOptions) (*vpcv1.LoadBalancer, *core.DetailedResponse, error) {
	return s.vpcService.GetLoadBalancer(options)
}

// CreateLoadBalancerPoolMember creates a new member and adds the member to the pool.
func (s *Service) CreateLoadBalancerPoolMember(options *vpcv1.CreateLoadBalancerPoolMemberOptions) (*vpcv1.LoadBalancerPoolMember, *core.DetailedResponse, error) {
	return s.vpcService.CreateLoadBalancerPoolMember(options)
}

// DeleteLoadBalancerPoolMember deletes a member from the load balancer pool.
func (s *Service) DeleteLoadBalancerPoolMember(options *vpcv1.DeleteLoadBalancerPoolMemberOptions) (*core.DetailedResponse, error) {
	return s.vpcService.DeleteLoadBalancerPoolMember(options)
}

// ListLoadBalancerPoolMembers returns members of a load balancer pool.
func (s *Service) ListLoadBalancerPoolMembers(options *vpcv1.ListLoadBalancerPoolMembersOptions) (*vpcv1.LoadBalancerPoolMemberCollection, *core.DetailedResponse, error) {
	return s.vpcService.ListLoadBalancerPoolMembers(options)
}

// ListKeys returns list of keys in a region.
func (s *Service) ListKeys(options *vpcv1.ListKeysOptions) (*vpcv1.KeyCollection, *core.DetailedResponse, error) {
	return s.vpcService.ListKeys(options)
}

// CreateImage creates a new VPC Custom Image.
func (s *Service) CreateImage(options *vpcv1.CreateImageOptions) (*vpcv1.Image, *core.DetailedResponse, error) {
	return s.vpcService.CreateImage(options)
}

// ListImages returns list of images in a region.
func (s *Service) ListImages(options *vpcv1.ListImagesOptions) (*vpcv1.ImageCollection, *core.DetailedResponse, error) {
	return s.vpcService.ListImages(options)
}

// GetImage returns a VPC Custom image.
func (s *Service) GetImage(options *vpcv1.GetImageOptions) (*vpcv1.Image, *core.DetailedResponse, error) {
	return s.vpcService.GetImage(options)
}

// GetInstanceProfile returns instance profile.
func (s *Service) GetInstanceProfile(options *vpcv1.GetInstanceProfileOptions) (*vpcv1.InstanceProfile, *core.DetailedResponse, error) {
	return s.vpcService.GetInstanceProfile(options)
}

// GetVPC returns VPC details.
func (s *Service) GetVPC(options *vpcv1.GetVPCOptions) (*vpcv1.VPC, *core.DetailedResponse, error) {
	return s.vpcService.GetVPC(options)
}

// GetVPCByName returns VPC with given name. If not found, returns nil.
func (s *Service) GetVPCByName(vpcName string) (*vpcv1.VPC, error) {
	var vpc *vpcv1.VPC
	f := func(start string) (bool, string, error) {
		// check for existing vpcs
		listVpcsOptions := &vpcv1.ListVpcsOptions{}
		if start != "" {
			listVpcsOptions.Start = &start
		}

		vpcsList, _, err := s.ListVpcs(listVpcsOptions)
		if err != nil {
			return false, "", err
		}

		if vpcsList == nil {
			return false, "", fmt.Errorf("vpc list returned is nil")
		}

		for i, v := range vpcsList.Vpcs {
			if (*v.Name) == vpcName {
				vpc = &vpcsList.Vpcs[i]
				return true, "", nil
			}
		}

		if vpcsList.Next != nil && *vpcsList.Next.Href != "" {
			return false, *vpcsList.Next.Href, nil
		}
		return true, "", nil
	}

	if err := utils.PagingHelper(f); err != nil {
		return nil, err
	}

	return vpc, nil
}

// GetImageByName returns the VPC Custom Image with given name. If not found, returns nil.
func (s *Service) GetImageByName(imageName string) (*vpcv1.Image, error) {
	var image *vpcv1.Image
	f := func(start string) (bool, string, error) {
		// check for existing images
		listImagesOptions := &vpcv1.ListImagesOptions{}
		if start != "" {
			listImagesOptions.Start = &start
		}

		imagesList, _, err := s.ListImages(listImagesOptions)
		if err != nil {
			return false, "", err
		}

		if imagesList == nil {
			return false, "", fmt.Errorf("image list returned is nil")
		}

		for i, v := range imagesList.Images {
			if *v.Name == imageName {
				image = &imagesList.Images[i]
				return true, "", nil
			}
		}

		if imagesList.Next != nil && *imagesList.Next.Href != "" {
			return false, *imagesList.Next.Href, nil
		}
		return true, "", nil
	}

	if err := utils.PagingHelper(f); err != nil {
		return nil, err
	}

	return image, nil
}

// GetVPCPublicGatewayByName returns the VPC Public Gateway with given name. If not found, returns nil.
func (s *Service) GetVPCPublicGatewayByName(publicGatewayName string, resourceGroupID string) (*vpcv1.PublicGateway, error) {
	var publicGateway *vpcv1.PublicGateway
	f := func(start string) (bool, string, error) {
		// check for existing public gateways
		listPublicGatewaysOptions := s.vpcService.NewListPublicGatewaysOptions().SetResourceGroupID(resourceGroupID)
		if start != "" {
			listPublicGatewaysOptions.Start = &start
		}

		publicGatewaysList, _, err := s.vpcService.ListPublicGateways(listPublicGatewaysOptions)
		if err != nil {
			return false, "", err
		}

		if publicGatewaysList == nil {
			return false, "", fmt.Errorf("public gateways list returned is nil")
		}

		for index, pg := range publicGatewaysList.PublicGateways {
			if *pg.Name == publicGatewayName {
				publicGateway = &publicGatewaysList.PublicGateways[index]
				return true, "", nil
			}
		}

		if publicGatewaysList.Next != nil && *publicGatewaysList.Next.Href != "" {
			return false, *publicGatewaysList.Next.Href, nil
		}
		return true, "", nil
	}

	if err := utils.PagingHelper(f); err != nil {
		return nil, err
	}

	return publicGateway, nil
}

// GetSubnet return subnet.
func (s *Service) GetSubnet(options *vpcv1.GetSubnetOptions) (*vpcv1.Subnet, *core.DetailedResponse, error) {
	return s.vpcService.GetSubnet(options)
}

// GetVPCSubnetByName returns subnet with given name. If not found, returns nil.
func (s *Service) GetVPCSubnetByName(subnetName string) (*vpcv1.Subnet, error) {
	var subnet *vpcv1.Subnet
	f := func(start string) (bool, string, error) {
		// check for existing subnets
		listSubnetsOptions := &vpcv1.ListSubnetsOptions{}
		if start != "" {
			listSubnetsOptions.Start = &start
		}

		subnetsList, _, err := s.ListSubnets(listSubnetsOptions)
		if err != nil {
			return false, "", err
		}

		if subnetsList == nil {
			return false, "", fmt.Errorf("subnet list returned is nil")
		}

		for i, s := range subnetsList.Subnets {
			if (*s.Name) == subnetName {
				subnet = &subnetsList.Subnets[i]
				return true, "", nil
			}
		}

		if subnetsList.Next != nil && *subnetsList.Next.Href != "" {
			return false, *subnetsList.Next.Href, nil
		}
		return true, "", nil
	}

	if err := utils.PagingHelper(f); err != nil {
		return nil, err
	}

	return subnet, nil
}

// GetLoadBalancerPoolByName returns a Load Balancer Pool with the given name, in the provided Load Balancer. If not found, returns nil.
func (s *Service) GetLoadBalancerPoolByName(loadBalancerID string, poolName string) (*vpcv1.LoadBalancerPool, error) {
	listLoadBalancerPoolsOptions := &vpcv1.ListLoadBalancerPoolsOptions{}
	listLoadBalancerPoolsOptions.SetLoadBalancerID(loadBalancerID)

	pools, _, err := s.vpcService.ListLoadBalancerPools(listLoadBalancerPoolsOptions)
	if err != nil {
		return nil, fmt.Errorf("error listing pools for load balancer %s: %w", loadBalancerID, err)
	} else if pools == nil {
		return nil, fmt.Errorf("error no pools for load balancer: %s", loadBalancerID)
	}

	for _, pool := range pools.Pools {
		if pool.Name != nil && *pool.Name == poolName {
			return &pool, nil
		}
	}

	return nil, nil
}

// GetLoadBalancerByName returns loadBalancer with given name. If not found, returns nil.
func (s *Service) GetLoadBalancerByName(loadBalancerName string) (*vpcv1.LoadBalancer, error) {
	var loadBalancer *vpcv1.LoadBalancer
	f := func(start string) (bool, string, error) {
		// check for existing loadBalancers
		listLoadBalancersOptions := &vpcv1.ListLoadBalancersOptions{}
		if start != "" {
			listLoadBalancersOptions.Start = &start
		}

		loadBalancersList, _, err := s.ListLoadBalancers(listLoadBalancersOptions)
		if err != nil {
			return false, "", err
		}

		if loadBalancersList == nil {
			return false, "", fmt.Errorf("loadBalancer list returned is nil")
		}

		for i, lb := range loadBalancersList.LoadBalancers {
			if (*lb.Name) == loadBalancerName {
				loadBalancer = &loadBalancersList.LoadBalancers[i]
				return true, "", nil
			}
		}

		if loadBalancersList.Next != nil && *loadBalancersList.Next.Href != "" {
			return false, *loadBalancersList.Next.Href, nil
		}
		return true, "", nil
	}

	if err := utils.PagingHelper(f); err != nil {
		return nil, err
	}

	return loadBalancer, nil
}

// CreateSecurityGroup creates a new security group.
func (s *Service) CreateSecurityGroup(options *vpcv1.CreateSecurityGroupOptions) (*vpcv1.SecurityGroup, *core.DetailedResponse, error) {
	return s.vpcService.CreateSecurityGroup(options)
}

// DeleteSecurityGroup deletes the security group passed.
func (s *Service) DeleteSecurityGroup(options *vpcv1.DeleteSecurityGroupOptions) (*core.DetailedResponse, error) {
	return s.vpcService.DeleteSecurityGroup(options)
}

// ListSecurityGroups lists security group.
func (s *Service) ListSecurityGroups(options *vpcv1.ListSecurityGroupsOptions) (*vpcv1.SecurityGroupCollection, *core.DetailedResponse, error) {
	return s.vpcService.ListSecurityGroups(options)
}

// GetSecurityGroup gets a specific security group by id.
func (s *Service) GetSecurityGroup(options *vpcv1.GetSecurityGroupOptions) (*vpcv1.SecurityGroup, *core.DetailedResponse, error) {
	return s.vpcService.GetSecurityGroup(options)
}

// GetSecurityGroupByName gets a specific security group by name.
func (s *Service) GetSecurityGroupByName(name string) (*vpcv1.SecurityGroup, error) {
	securityGroupPager, err := s.vpcService.NewSecurityGroupsPager(&vpcv1.ListSecurityGroupsOptions{})
	if err != nil {
		return nil, fmt.Errorf("error listing security group: %v", err)
	}

	for {
		if !securityGroupPager.HasNext() {
			break
		}

		securityGroups, err := securityGroupPager.GetNext()
		if err != nil {
			return nil, fmt.Errorf("error retrieving next page of security groups: %v", err)
		}

		for _, sg := range securityGroups {
			if *sg.Name == name {
				return &sg, nil
			}
		}
	}

	return nil, &SecurityGroupByNameNotFound{Name: name}
}

// GetSecurityGroupRule gets a specific security group rule.
func (s *Service) GetSecurityGroupRule(options *vpcv1.GetSecurityGroupRuleOptions) (vpcv1.SecurityGroupRuleIntf, *core.DetailedResponse, error) {
	return s.vpcService.GetSecurityGroupRule(options)
}

// ListSecurityGroupRules returns a list of security group rules.
func (s *Service) ListSecurityGroupRules(options *vpcv1.ListSecurityGroupRulesOptions) (*vpcv1.SecurityGroupRuleCollection, *core.DetailedResponse, error) {
	return s.vpcService.ListSecurityGroupRules(options)
}

// GetVPCZonesByRegion gets the VPC availability zones for a specific IBM Cloud region.
func (s *Service) GetVPCZonesByRegion(region string) ([]string, error) {
	zones := make([]string, 0)
	options := s.vpcService.NewListRegionZonesOptions(region)
	result, _, err := s.vpcService.ListRegionZones(options)
	if err != nil {
		return zones, err
	}
	for _, zone := range result.Zones {
		zones = append(zones, *zone.Name)
	}
	return zones, nil
}

// NewService returns a new VPC Service.
func NewService(svcEndpoint string) (Vpc, error) {
	service := &Service{}
	auth, err := authenticator.GetAuthenticator()
	if err != nil {
		return nil, err
	}

	service.vpcService, err = vpcv1.NewVpcV1(&vpcv1.VpcV1Options{
		Authenticator: auth,
		URL:           svcEndpoint,
	})

	return service, err
}
