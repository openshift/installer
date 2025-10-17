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
)

//go:generate ../../../../hack/tools/bin/mockgen -source=./vpc.go -destination=./mock/vpc_generated.go -package=mock
//go:generate /usr/bin/env bash -c "cat ../../../../hack/boilerplate/boilerplate.generatego.txt ./mock/vpc_generated.go > ./mock/_vpc_generated.go && mv ./mock/_vpc_generated.go ./mock/vpc_generated.go"

// Vpc interface defines methods that a Cluster API IBMCLOUD object should implement.
type Vpc interface {
	CreateInstance(options *vpcv1.CreateInstanceOptions) (*vpcv1.Instance, *core.DetailedResponse, error)
	DeleteInstance(options *vpcv1.DeleteInstanceOptions) (*core.DetailedResponse, error)
	GetInstance(options *vpcv1.GetInstanceOptions) (*vpcv1.Instance, *core.DetailedResponse, error)
	ListInstances(options *vpcv1.ListInstancesOptions) (*vpcv1.InstanceCollection, *core.DetailedResponse, error)
	GetDedicatedHostByName(dHostName string) (*vpcv1.DedicatedHost, error)
	CreateVPC(options *vpcv1.CreateVPCOptions) (*vpcv1.VPC, *core.DetailedResponse, error)
	DeleteVPC(options *vpcv1.DeleteVPCOptions) (response *core.DetailedResponse, err error)
	ListVpcs(options *vpcv1.ListVpcsOptions) (*vpcv1.VPCCollection, *core.DetailedResponse, error)
	CreateSubnet(options *vpcv1.CreateSubnetOptions) (*vpcv1.Subnet, *core.DetailedResponse, error)
	DeleteSubnet(options *vpcv1.DeleteSubnetOptions) (*core.DetailedResponse, error)
	ListSubnets(options *vpcv1.ListSubnetsOptions) (*vpcv1.SubnetCollection, *core.DetailedResponse, error)
	GetSubnetPublicGateway(options *vpcv1.GetSubnetPublicGatewayOptions) (*vpcv1.PublicGateway, *core.DetailedResponse, error)
	SetSubnetPublicGateway(options *vpcv1.SetSubnetPublicGatewayOptions) (*vpcv1.PublicGateway, *core.DetailedResponse, error)
	UnsetSubnetPublicGateway(options *vpcv1.UnsetSubnetPublicGatewayOptions) (*core.DetailedResponse, error)
	CreatePublicGateway(options *vpcv1.CreatePublicGatewayOptions) (*vpcv1.PublicGateway, *core.DetailedResponse, error)
	DeletePublicGateway(options *vpcv1.DeletePublicGatewayOptions) (*core.DetailedResponse, error)
	ListVPCAddressPrefixes(options *vpcv1.ListVPCAddressPrefixesOptions) (*vpcv1.AddressPrefixCollection, *core.DetailedResponse, error)
	CreateSecurityGroupRule(options *vpcv1.CreateSecurityGroupRuleOptions) (vpcv1.SecurityGroupRuleIntf, *core.DetailedResponse, error)
	CreateLoadBalancer(options *vpcv1.CreateLoadBalancerOptions) (*vpcv1.LoadBalancer, *core.DetailedResponse, error)
	DeleteLoadBalancer(options *vpcv1.DeleteLoadBalancerOptions) (*core.DetailedResponse, error)
	ListLoadBalancers(options *vpcv1.ListLoadBalancersOptions) (*vpcv1.LoadBalancerCollection, *core.DetailedResponse, error)
	GetLoadBalancer(options *vpcv1.GetLoadBalancerOptions) (*vpcv1.LoadBalancer, *core.DetailedResponse, error)
	CreateLoadBalancerPoolMember(options *vpcv1.CreateLoadBalancerPoolMemberOptions) (*vpcv1.LoadBalancerPoolMember, *core.DetailedResponse, error)
	DeleteLoadBalancerPoolMember(options *vpcv1.DeleteLoadBalancerPoolMemberOptions) (*core.DetailedResponse, error)
	ListLoadBalancerPoolMembers(options *vpcv1.ListLoadBalancerPoolMembersOptions) (*vpcv1.LoadBalancerPoolMemberCollection, *core.DetailedResponse, error)
	ListKeys(options *vpcv1.ListKeysOptions) (*vpcv1.KeyCollection, *core.DetailedResponse, error)
	CreateImage(options *vpcv1.CreateImageOptions) (*vpcv1.Image, *core.DetailedResponse, error)
	ListImages(options *vpcv1.ListImagesOptions) (*vpcv1.ImageCollection, *core.DetailedResponse, error)
	GetImage(options *vpcv1.GetImageOptions) (*vpcv1.Image, *core.DetailedResponse, error)
	GetInstanceProfile(options *vpcv1.GetInstanceProfileOptions) (*vpcv1.InstanceProfile, *core.DetailedResponse, error)
	GetVPC(*vpcv1.GetVPCOptions) (*vpcv1.VPC, *core.DetailedResponse, error)
	GetVPCByName(vpcName string) (*vpcv1.VPC, error)
	GetImageByName(imageName string) (*vpcv1.Image, error)
	GetVPCPublicGatewayByName(publicGatewayName string, resourceGroupID string) (*vpcv1.PublicGateway, error)
	GetSubnet(*vpcv1.GetSubnetOptions) (*vpcv1.Subnet, *core.DetailedResponse, error)
	GetVPCSubnetByName(subnetName string) (*vpcv1.Subnet, error)
	GetLoadBalancerPoolByName(loadBalancerID string, poolName string) (*vpcv1.LoadBalancerPool, error)
	GetLoadBalancerByName(loadBalancerName string) (*vpcv1.LoadBalancer, error)
	CreateSecurityGroup(options *vpcv1.CreateSecurityGroupOptions) (*vpcv1.SecurityGroup, *core.DetailedResponse, error)
	DeleteSecurityGroup(options *vpcv1.DeleteSecurityGroupOptions) (*core.DetailedResponse, error)
	ListSecurityGroups(options *vpcv1.ListSecurityGroupsOptions) (*vpcv1.SecurityGroupCollection, *core.DetailedResponse, error)
	GetSecurityGroup(options *vpcv1.GetSecurityGroupOptions) (*vpcv1.SecurityGroup, *core.DetailedResponse, error)
	GetSecurityGroupByName(name string) (*vpcv1.SecurityGroup, error)
	GetSecurityGroupRule(options *vpcv1.GetSecurityGroupRuleOptions) (vpcv1.SecurityGroupRuleIntf, *core.DetailedResponse, error)
	ListSecurityGroupRules(options *vpcv1.ListSecurityGroupRulesOptions) (*vpcv1.SecurityGroupRuleCollection, *core.DetailedResponse, error)
	GetVPCZonesByRegion(region string) ([]string, error)
}
