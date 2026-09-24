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

package scope

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-logr/logr"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"

	"k8s.io/klog/v2"

	"sigs.k8s.io/controller-runtime/pkg/client"

	capiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/patch"

	infrav1beta2 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/utils"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/vpc"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/endpoints"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/record"
)

const subnetSuffix = "-subnet"

// ClusterScopeParams defines the input parameters used to create a new ClusterScope.
type ClusterScopeParams struct {
	IBMVPCClient    vpc.Vpc
	Client          client.Client
	Logger          logr.Logger
	Cluster         *capiv1beta1.Cluster
	IBMVPCCluster   *infrav1beta2.IBMVPCCluster
	ServiceEndpoint []endpoints.ServiceEndpoint
}

// ClusterScope defines a scope defined around a cluster.
type ClusterScope struct {
	logr.Logger
	Client      client.Client
	patchHelper *patch.Helper

	IBMVPCClient    vpc.Vpc
	Cluster         *capiv1beta1.Cluster
	IBMVPCCluster   *infrav1beta2.IBMVPCCluster
	ServiceEndpoint []endpoints.ServiceEndpoint
}

// NewClusterScope creates a new ClusterScope from the supplied parameters.
func NewClusterScope(params ClusterScopeParams) (*ClusterScope, error) {
	if params.Cluster == nil {
		return nil, errors.New("failed to generate new scope from nil Cluster")
	}
	if params.IBMVPCCluster == nil {
		return nil, errors.New("failed to generate new scope from nil IBMVPCCluster")
	}

	if params.Logger == (logr.Logger{}) {
		params.Logger = klog.Background()
	}

	helper, err := patch.NewHelper(params.IBMVPCCluster, params.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to init patch helper: %w", err)
	}

	// Fetch the service endpoint.
	svcEndpoint := endpoints.FetchVPCEndpoint(params.IBMVPCCluster.Spec.Region, params.ServiceEndpoint)

	vpcClient, err := vpc.NewService(svcEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to create IBM VPC session: %w", err)
	}

	if params.Logger.V(DEBUGLEVEL).Enabled() {
		core.SetLoggingLevel(core.LevelDebug)
	}

	return &ClusterScope{
		Logger:        params.Logger,
		Client:        params.Client,
		IBMVPCClient:  vpcClient,
		Cluster:       params.Cluster,
		IBMVPCCluster: params.IBMVPCCluster,
		patchHelper:   helper,
	}, nil
}

// CreateVPC creates a new IBM VPC in specified resource group.
func (s *ClusterScope) CreateVPC() (*vpcv1.VPC, error) {
	vpcReply, err := s.ensureVPCUnique(s.IBMVPCCluster.Spec.VPC)
	if err != nil {
		return nil, err
	} else if vpcReply != nil {
		// TODO need a reasonable wrapped error
		return vpcReply, nil
	}

	options := &vpcv1.CreateVPCOptions{}
	options.SetResourceGroup(&vpcv1.ResourceGroupIdentity{
		ID: &s.IBMVPCCluster.Spec.ResourceGroup,
	})
	options.SetName(s.IBMVPCCluster.Spec.VPC)
	vpc, _, err := s.IBMVPCClient.CreateVPC(options)
	if err != nil {
		record.Warnf(s.IBMVPCCluster, "FailedCreateVPC", "Failed vpc creation - %v", err)
		return nil, err
	} else if err := s.updateDefaultSG(*vpc.DefaultSecurityGroup.ID); err != nil {
		record.Warnf(s.IBMVPCCluster, "FailedUpdateDefaultSecurityGroup", "Failed to update default security group - %v", err)
		return nil, err
	}
	record.Eventf(s.IBMVPCCluster, "SuccessfulCreateVPC", "Created VPC %q", *vpc.Name)
	return vpc, nil
}

// DeleteVPC deletes IBM VPC associated with a VPC id.
func (s *ClusterScope) DeleteVPC() error {
	if s.IBMVPCCluster.Status.VPC.ID == "" {
		return nil
	}

	deleteVpcOptions := &vpcv1.DeleteVPCOptions{}
	deleteVpcOptions.SetID(s.IBMVPCCluster.Status.VPC.ID)
	_, err := s.IBMVPCClient.DeleteVPC(deleteVpcOptions)
	if err != nil {
		record.Warnf(s.IBMVPCCluster, "FailedDeleteVPC", "Failed vpc deletion - %v", err)
	} else {
		record.Eventf(s.IBMVPCCluster, "SuccessfulDeleteVPC", "Deleted VPC %q", s.IBMVPCCluster.Status.VPC.Name)
	}

	return err
}

func (s *ClusterScope) ensureVPCUnique(vpcName string) (*vpcv1.VPC, error) {
	var vpc *vpcv1.VPC
	f := func(start string) (bool, string, error) {
		// check for existing vpcs
		listVpcsOptions := &vpcv1.ListVpcsOptions{}
		if start != "" {
			listVpcsOptions.Start = &start
		}

		vpcsList, _, err := s.IBMVPCClient.ListVpcs(listVpcsOptions)
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

func (s *ClusterScope) updateDefaultSG(sgID string) error {
	options := &vpcv1.CreateSecurityGroupRuleOptions{}
	options.SetSecurityGroupID(sgID)
	options.SetSecurityGroupRulePrototype(&vpcv1.SecurityGroupRulePrototype{
		Direction: core.StringPtr("inbound"),
		Protocol:  core.StringPtr("tcp"),
		IPVersion: core.StringPtr("ipv4"),
		PortMin:   core.Int64Ptr(int64(s.APIServerPort())),
		PortMax:   core.Int64Ptr(int64(s.APIServerPort())),
	})
	_, _, err := s.IBMVPCClient.CreateSecurityGroupRule(options)
	if err != nil {
		record.Warnf(s.IBMVPCCluster, "FailedCreateSecurityGroupRule", "Failed security group rule creation - %v", err)
	}
	return err
}

// CreateSubnet creates a subnet within provided vpc and zone.
func (s *ClusterScope) CreateSubnet() (*vpcv1.Subnet, error) {
	subnetName := s.IBMVPCCluster.Name + subnetSuffix
	subnetReply, err := s.ensureSubnetUnique(subnetName)
	if err != nil {
		return nil, err
	} else if subnetReply != nil {
		// TODO need a reasonable wrapped error
		return subnetReply, nil
	}

	options := &vpcv1.CreateSubnetOptions{}
	cidrBlock, err := s.getSubnetAddrPrefix(s.IBMVPCCluster.Status.VPC.ID, s.IBMVPCCluster.Spec.Zone)
	if err != nil {
		return nil, err
	}
	subnetName = s.IBMVPCCluster.Name + subnetSuffix
	options.SetSubnetPrototype(&vpcv1.SubnetPrototype{
		Ipv4CIDRBlock: &cidrBlock,
		Name:          &subnetName,
		VPC: &vpcv1.VPCIdentity{
			ID: &s.IBMVPCCluster.Status.VPC.ID,
		},
		Zone: &vpcv1.ZoneIdentity{
			Name: &s.IBMVPCCluster.Spec.Zone,
		},
		ResourceGroup: &vpcv1.ResourceGroupIdentity{
			ID: &s.IBMVPCCluster.Spec.ResourceGroup,
		},
	})
	subnet, _, err := s.IBMVPCClient.CreateSubnet(options)
	if err != nil {
		record.Warnf(s.IBMVPCCluster, "FailedCreateSubnet", "Failed subnet creation - %v", err)
	}
	if subnet != nil {
		pgw, err := s.createPublicGateWay(s.IBMVPCCluster.Status.VPC.ID, s.IBMVPCCluster.Spec.Zone, s.IBMVPCCluster.Spec.ResourceGroup)
		if err != nil {
			return subnet, err
		}
		if pgw != nil {
			if _, err := s.attachPublicGateWay(*subnet.ID, *pgw.ID); err != nil {
				return nil, err
			}
		}
	}
	return subnet, err
}

func (s *ClusterScope) getSubnetAddrPrefix(vpcID, zone string) (string, error) {
	var addrPrefix *vpcv1.AddressPrefix
	f := func(start string) (bool, string, error) {
		// check for existing vpcAddressPrefixes
		listVPCAddressPrefixesOptions := &vpcv1.ListVPCAddressPrefixesOptions{
			VPCID: &vpcID,
		}
		if start != "" {
			listVPCAddressPrefixesOptions.Start = &start
		}

		vpcAddressPrefixesList, _, err := s.IBMVPCClient.ListVPCAddressPrefixes(listVPCAddressPrefixesOptions)
		if err != nil {
			return false, "", err
		}

		if vpcAddressPrefixesList == nil {
			return false, "", fmt.Errorf("vpcAddressPrefix list returned is nil")
		}

		for i, addressPrefix := range vpcAddressPrefixesList.AddressPrefixes {
			if (*addressPrefix.Zone.Name) == zone {
				addrPrefix = &vpcAddressPrefixesList.AddressPrefixes[i]
				return true, "", nil
			}
		}

		if vpcAddressPrefixesList.Next != nil && *vpcAddressPrefixesList.Next.Href != "" {
			return false, *vpcAddressPrefixesList.Next.Href, nil
		}
		return true, "", nil
	}

	if err := utils.PagingHelper(f); err != nil {
		return "", err
	}

	if addrPrefix != nil {
		return *addrPrefix.CIDR, nil
	}
	return "", fmt.Errorf("not found a valid CIDR for VPC %s in zone %s", vpcID, zone)
}

func (s *ClusterScope) ensureSubnetUnique(subnetName string) (*vpcv1.Subnet, error) {
	var subnet *vpcv1.Subnet
	f := func(start string) (bool, string, error) {
		// check for existing subnets
		listSubnetsOptions := &vpcv1.ListSubnetsOptions{}
		if start != "" {
			listSubnetsOptions.Start = &start
		}

		subnetsList, _, err := s.IBMVPCClient.ListSubnets(listSubnetsOptions)
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

// DeleteSubnet deletes a subnet associated with subnet id.
func (s *ClusterScope) DeleteSubnet() error {
	if s.IBMVPCCluster.Status.Subnet.ID == nil {
		return nil
	}

	subnetID := *s.IBMVPCCluster.Status.Subnet.ID

	// Lists the subnet available and compare before deleting to avoid any failure(404) later
	found := false
	f := func(start string) (bool, string, error) {
		// check for existing subnets
		listSubnetsOptions := &vpcv1.ListSubnetsOptions{}
		if start != "" {
			listSubnetsOptions.Start = &start
		}

		subnetsList, _, err := s.IBMVPCClient.ListSubnets(listSubnetsOptions)
		if err != nil {
			return false, "", err
		}

		if subnetsList == nil {
			return false, "", fmt.Errorf("subnet list returned is nil")
		}

		for _, s := range subnetsList.Subnets {
			if *s.ID == subnetID {
				found = true
				return true, "", nil
			}
		}

		if subnetsList.Next != nil && *subnetsList.Next.Href != "" {
			return false, *subnetsList.Next.Href, nil
		}
		return true, "", nil
	}

	if err := utils.PagingHelper(f); err != nil {
		return err
	}

	if !found {
		s.Logger.V(3).Info("No subnets found with ID", "Subnet ID", subnetID)
		return nil
	}

	// get the pgw id for given subnet, so we can delete it later
	getPGWOptions := &vpcv1.GetSubnetPublicGatewayOptions{}
	getPGWOptions.SetID(subnetID)
	pgw, _, err := s.IBMVPCClient.GetSubnetPublicGateway(getPGWOptions)
	if err != nil {
		return err
	}
	if pgw != nil { // public gateway found
		// Unset the public gateway for subnet first
		err = s.detachPublicGateway(subnetID, *pgw.ID)
		if err != nil {
			return fmt.Errorf("error when detaching publicgateway for subnet %s: %w", subnetID, err)
		}
	}

	// Delete subnet
	deleteSubnetOption := &vpcv1.DeleteSubnetOptions{}
	deleteSubnetOption.SetID(subnetID)
	_, err = s.IBMVPCClient.DeleteSubnet(deleteSubnetOption)
	if err != nil {
		record.Warnf(s.IBMVPCCluster, "FailedDeleteSubnet", "Failed subnet deletion - %v", err)
		return fmt.Errorf("error when deleting subnet: %w", err)
	}
	return err
}

func (s *ClusterScope) createPublicGateWay(vpcID string, zoneName string, resourceGroupID string) (*vpcv1.PublicGateway, error) {
	options := &vpcv1.CreatePublicGatewayOptions{}
	options.SetVPC(&vpcv1.VPCIdentity{
		ID: &vpcID,
	})
	options.SetZone(&vpcv1.ZoneIdentity{
		Name: &zoneName,
	})
	options.SetResourceGroup(&vpcv1.ResourceGroupIdentity{
		ID: &resourceGroupID,
	})
	publicGateway, _, err := s.IBMVPCClient.CreatePublicGateway(options)
	if err != nil {
		record.Warnf(s.IBMVPCCluster, "FailedCreatePublicGateway", "Failed publicgateway creation - %v", err)
	}
	return publicGateway, err
}

func (s *ClusterScope) attachPublicGateWay(subnetID string, pgwID string) (*vpcv1.PublicGateway, error) {
	options := &vpcv1.SetSubnetPublicGatewayOptions{}
	options.SetID(subnetID)
	options.SetPublicGatewayIdentity(&vpcv1.PublicGatewayIdentity{
		ID: &pgwID,
	})
	publicGateway, _, err := s.IBMVPCClient.SetSubnetPublicGateway(options)
	if err != nil {
		record.Warnf(s.IBMVPCCluster, "FailedAttachPublicGateway", "Failed publicgateway attachment - %v", err)
	}
	return publicGateway, err
}

func (s *ClusterScope) detachPublicGateway(subnetID string, pgwID string) error {
	// Unset the publicgateway first, and then delete it
	unsetPGWOption := &vpcv1.UnsetSubnetPublicGatewayOptions{}
	unsetPGWOption.SetID(subnetID)
	_, err := s.IBMVPCClient.UnsetSubnetPublicGateway(unsetPGWOption)
	if err != nil {
		record.Warnf(s.IBMVPCCluster, "FailedDetachPublicGateway", "Failed publicgateway detachment - %v", err)
		return fmt.Errorf("error when unsetting publicgateway for subnet %s: %w", subnetID, err)
	}

	// Delete the public gateway
	deletePGWOption := &vpcv1.DeletePublicGatewayOptions{}
	deletePGWOption.SetID(pgwID)
	_, err = s.IBMVPCClient.DeletePublicGateway(deletePGWOption)
	if err != nil {
		record.Warnf(s.IBMVPCCluster, "FailedDeletePublicGateway", "Failed publicgateway deletion - %v", err)
		return fmt.Errorf("error when deleting publicgateway for subnet %s: %w", subnetID, err)
	}
	return err
}

// CreateLoadBalancer creates a new IBM VPC load balancer in specified resource group.
func (s *ClusterScope) CreateLoadBalancer() (*vpcv1.LoadBalancer, error) {
	loadBalancerReply, err := s.ensureLoadBalancerUnique(s.IBMVPCCluster.Spec.ControlPlaneLoadBalancer.Name)
	if err != nil {
		return nil, err
	} else if loadBalancerReply != nil {
		// TODO need a reasonable wrapped error
		return loadBalancerReply, nil
	}

	options := &vpcv1.CreateLoadBalancerOptions{}
	options.SetName(s.IBMVPCCluster.Spec.ControlPlaneLoadBalancer.Name)
	options.SetIsPublic(true)
	options.SetResourceGroup(&vpcv1.ResourceGroupIdentity{
		ID: &s.IBMVPCCluster.Spec.ResourceGroup,
	})

	if s.IBMVPCCluster.Status.Subnet.ID != nil {
		subnet := &vpcv1.SubnetIdentity{
			ID: s.IBMVPCCluster.Status.Subnet.ID,
		}
		options.Subnets = append(options.Subnets, subnet)
	} else {
		return nil, fmt.Errorf("error subnet required for load balancer creation")
	}

	options.SetPools([]vpcv1.LoadBalancerPoolPrototype{
		{
			Algorithm:     core.StringPtr("round_robin"),
			HealthMonitor: &vpcv1.LoadBalancerPoolHealthMonitorPrototype{Delay: core.Int64Ptr(5), MaxRetries: core.Int64Ptr(2), Timeout: core.Int64Ptr(2), Type: core.StringPtr("tcp")},
			Name:          core.StringPtr(s.IBMVPCCluster.Spec.ControlPlaneLoadBalancer.Name + "-pool"),
			Protocol:      core.StringPtr("tcp"),
		},
	})

	options.SetListeners([]vpcv1.LoadBalancerListenerPrototypeLoadBalancerContext{
		{
			Protocol: core.StringPtr("tcp"),
			Port:     core.Int64Ptr(int64(s.APIServerPort())),
			DefaultPool: &vpcv1.LoadBalancerPoolIdentityByName{
				Name: core.StringPtr(s.IBMVPCCluster.Spec.ControlPlaneLoadBalancer.Name + "-pool"),
			},
		},
	})

	loadBalancer, _, err := s.IBMVPCClient.CreateLoadBalancer(options)
	if err != nil {
		record.Warnf(s.IBMVPCCluster, "FailedCreateLoadBalancer", "Failed loadBalancer creation - %v", err)
		return nil, err
	}

	record.Eventf(s.IBMVPCCluster, "SuccessfulCreateLoadBalancer", "Created loadBalancer %q", *loadBalancer.Name)
	return loadBalancer, nil
}

// GetLoadBalancerByHostname retrieves a IBM VPC load balancer with specified hostname.
func (s *ClusterScope) GetLoadBalancerByHostname(loadBalancerHostname string) (*vpcv1.LoadBalancer, error) {
	loadBalancer, err := s.getLoadBalancerByHostname(loadBalancerHostname)
	if err != nil {
		return nil, err
	}
	return loadBalancer, nil
}

func (s *ClusterScope) getLoadBalancerByHostname(loadBalancerHostname string) (*vpcv1.LoadBalancer, error) {
	var loadBalancer *vpcv1.LoadBalancer
	f := func(start string) (bool, string, error) {
		// check for existing loadBalancers
		listLoadBalancersOptions := &vpcv1.ListLoadBalancersOptions{}
		if start != "" {
			listLoadBalancersOptions.Start = &start
		}

		loadBalancersList, _, err := s.IBMVPCClient.ListLoadBalancers(listLoadBalancersOptions)
		if err != nil {
			return false, "", err
		}

		if loadBalancersList == nil {
			return false, "", fmt.Errorf("loadBalancer list returned is nil")
		}

		for i, lb := range loadBalancersList.LoadBalancers {
			if (*lb.Hostname) == loadBalancerHostname {
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

func (s *ClusterScope) ensureLoadBalancerUnique(loadBalancerName string) (*vpcv1.LoadBalancer, error) {
	var loadBalancer *vpcv1.LoadBalancer
	f := func(start string) (bool, string, error) {
		// check for existing loadBalancers
		listLoadBalancersOptions := &vpcv1.ListLoadBalancersOptions{}
		if start != "" {
			listLoadBalancersOptions.Start = &start
		}

		loadBalancersList, _, err := s.IBMVPCClient.ListLoadBalancers(listLoadBalancersOptions)
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

// DeleteLoadBalancer deletes IBM VPC load balancer associated with a VPC id.
func (s *ClusterScope) DeleteLoadBalancer() (bool, error) {
	deleted := false
	if lbipID := s.GetLoadBalancerID(); lbipID != "" {
		f := func(start string) (bool, string, error) {
			// check for existing loadBalancers
			listLoadBalancersOptions := &vpcv1.ListLoadBalancersOptions{}
			if start != "" {
				listLoadBalancersOptions.Start = &start
			}

			loadBalancersList, _, err := s.IBMVPCClient.ListLoadBalancers(listLoadBalancersOptions)
			if err != nil {
				return false, "", err
			}

			if loadBalancersList == nil {
				return false, "", fmt.Errorf("loadBalancer list returned is nil")
			}

			for _, lb := range loadBalancersList.LoadBalancers {
				if (*lb.ID) == lbipID {
					deleted = true
					if *lb.ProvisioningStatus != string(infrav1beta2.VPCLoadBalancerStateDeletePending) {
						deleteLoadBalancerOption := &vpcv1.DeleteLoadBalancerOptions{}
						deleteLoadBalancerOption.SetID(lbipID)
						_, err := s.IBMVPCClient.DeleteLoadBalancer(deleteLoadBalancerOption)
						if err != nil {
							record.Warnf(s.IBMVPCCluster, "FailedDeleteLoadBalancer", "Failed loadBalancer deletion - %v", err)
							return false, "", err
						}
					}
				}
			}

			if loadBalancersList.Next != nil && *loadBalancersList.Next.Href != "" {
				return false, *loadBalancersList.Next.Href, nil
			}
			return true, "", nil
		}

		if err := utils.PagingHelper(f); err != nil {
			return false, err
		}
	}
	return deleted, nil
}

// SetReady will set the status as ready for the cluster.
func (s *ClusterScope) SetReady() {
	s.IBMVPCCluster.Status.Ready = true
}

// SetNotReady will set the status as not ready for the cluster.
func (s *ClusterScope) SetNotReady() {
	s.IBMVPCCluster.Status.Ready = false
}

// IsReady will return the status for the cluster.
func (s *ClusterScope) IsReady() bool {
	return s.IBMVPCCluster.Status.Ready
}

// SetLoadBalancerState will set the state for the load balancer.
func (s *ClusterScope) SetLoadBalancerState(status string) {
	s.IBMVPCCluster.Status.ControlPlaneLoadBalancerState = infrav1beta2.VPCLoadBalancerState(status)
}

// GetLoadBalancerState will get the state for the load balancer.
func (s *ClusterScope) GetLoadBalancerState() infrav1beta2.VPCLoadBalancerState {
	return s.IBMVPCCluster.Status.ControlPlaneLoadBalancerState
}

// SetLoadBalancerID will set the id for the load balancer.
func (s *ClusterScope) SetLoadBalancerID(id *string) {
	s.IBMVPCCluster.Status.VPCEndpoint.LBID = id
}

// GetLoadBalancerID will get the id for the load balancer.
func (s *ClusterScope) GetLoadBalancerID() string {
	if s.IBMVPCCluster.Status.VPCEndpoint.LBID == nil {
		return ""
	}

	return *s.IBMVPCCluster.Status.VPCEndpoint.LBID
}

// SetLoadBalancerAddress will set the address for the load balancer.
func (s *ClusterScope) SetLoadBalancerAddress(address *string) {
	s.IBMVPCCluster.Status.VPCEndpoint.Address = address
}

// GetLoadBalancerAddress will get the address for the load balancer.
func (s *ClusterScope) GetLoadBalancerAddress() string {
	if s.IBMVPCCluster.Status.VPCEndpoint.Address == nil {
		return ""
	}

	return *s.IBMVPCCluster.Status.VPCEndpoint.Address
}

// PatchObject persists the cluster configuration and status.
func (s *ClusterScope) PatchObject() error {
	return s.patchHelper.Patch(context.TODO(), s.IBMVPCCluster)
}

// Close closes the current scope persisting the cluster configuration and status.
func (s *ClusterScope) Close() error {
	return s.PatchObject()
}

// APIServerPort returns the APIServerPort to use when creating the ControlPlaneEndpoint.
func (s *ClusterScope) APIServerPort() int32 {
	if s.Cluster.Spec.ClusterNetwork != nil && s.Cluster.Spec.ClusterNetwork.APIServerPort != nil {
		return *s.Cluster.Spec.ClusterNetwork.APIServerPort
	}
	return infrav1beta2.DefaultAPIServerPort
}
