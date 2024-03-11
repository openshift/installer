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
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"

	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/authenticator"
)

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

// CreateFloatingIP reserves a floating IP.
func (s *Service) CreateFloatingIP(options *vpcv1.CreateFloatingIPOptions) (*vpcv1.FloatingIP, *core.DetailedResponse, error) {
	return s.vpcService.CreateFloatingIP(options)
}

// DeleteFloatingIP releases a floating IP.
func (s *Service) DeleteFloatingIP(options *vpcv1.DeleteFloatingIPOptions) (*core.DetailedResponse, error) {
	return s.vpcService.DeleteFloatingIP(options)
}

// ListFloatingIps returns list of the floating IPs in a region.
func (s *Service) ListFloatingIps(options *vpcv1.ListFloatingIpsOptions) (*vpcv1.FloatingIPCollection, *core.DetailedResponse, error) {
	return s.vpcService.ListFloatingIps(options)
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

// AddInstanceNetworkInterfaceFloatingIP associates a floating IP with a network interface.
func (s *Service) AddInstanceNetworkInterfaceFloatingIP(options *vpcv1.AddInstanceNetworkInterfaceFloatingIPOptions) (*vpcv1.FloatingIP, *core.DetailedResponse, error) {
	return s.vpcService.AddInstanceNetworkInterfaceFloatingIP(options)
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

// ListImages returns list of images in a region.
func (s *Service) ListImages(options *vpcv1.ListImagesOptions) (*vpcv1.ImageCollection, *core.DetailedResponse, error) {
	return s.vpcService.ListImages(options)
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
