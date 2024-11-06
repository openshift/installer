/*
Copyright 2018 The Kubernetes Authors.

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

package networking

import (
	"errors"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/external"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/subnets"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/metrics"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/record"
	capoerrors "sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/errors"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/filterconvert"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/names"
)

type createOpts struct {
	AdminStateUp        *bool  `json:"admin_state_up,omitempty"`
	Name                string `json:"name,omitempty"`
	PortSecurityEnabled *bool  `json:"port_security_enabled,omitempty"`
	MTU                 *int   `json:"mtu,omitempty"`
}

func (c createOpts) ToNetworkCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(c, "network")
}

// ReconcileExternalNetwork will try to find an external network and set it in the cluster status.
// The external network can be specified in the cluster spec or will be searched for if not specified.
// OpenStackCluster.Status.ExternalNetwork will be set to nil if one of these conditions are met:
// - no external network was given in the cluster spec and no external network was found
// - the user has set OpenStackCluster.Spec.DisableExternalNetwork to true.
func (s *Service) ReconcileExternalNetwork(openStackCluster *infrav1.OpenStackCluster) error {
	if ptr.Deref(openStackCluster.Spec.DisableExternalNetwork, false) {
		s.scope.Logger().Info("External network is disabled - proceeding with internal network only")
		openStackCluster.Status.ExternalNetwork = nil
		return nil
	}

	var network *networks.Network
	if openStackCluster.Spec.ExternalNetwork == nil {
		// No external network specified in the cluster spec. Default behaviour: query all external networks.
		// * If there's only one, use that.
		// * If there's none don't use an external network.
		// * If there's more than one it's an error.

		// Empty NetworkFilter will query all networks
		var err error
		network, err = s.getNetworkByFilter(&infrav1.NetworkFilter{}, ExternalNetworksOnly)
		if errors.Is(err, capoerrors.ErrNoMatches) {
			openStackCluster.Status.ExternalNetwork = nil
			s.scope.Logger().Info("No external network found - proceeding with internal network only")
			return nil
		}
		if err != nil {
			return fmt.Errorf("failed to get external network: %w", err)
		}
	} else {
		var err error
		network, err = s.GetNetworkByParam(openStackCluster.Spec.ExternalNetwork, ExternalNetworksOnly)
		if err != nil {
			return fmt.Errorf("failed to get external network: %w", err)
		}
	}

	openStackCluster.Status.ExternalNetwork = &infrav1.NetworkStatus{
		ID:   network.ID,
		Name: network.Name,
		Tags: network.Tags,
	}
	s.scope.Logger().Info("External network found", "id", network.ID)
	return nil
}

func (s *Service) ReconcileNetwork(openStackCluster *infrav1.OpenStackCluster, clusterResourceName string) error {
	networkName := getNetworkName(clusterResourceName)
	s.scope.Logger().Info("Reconciling network", "name", networkName)

	res, err := s.getNetworkByName(networkName)
	if err != nil {
		return err
	}

	if res.ID != "" {
		// Network exists
		openStackCluster.Status.Network = &infrav1.NetworkStatusWithSubnets{}
		openStackCluster.Status.Network.ID = res.ID
		openStackCluster.Status.Network.Name = res.Name
		openStackCluster.Status.Network.Tags = res.Tags
		s.scope.Logger().V(6).Info("Reusing existing network", "name", res.Name, "id", res.ID)
		return nil
	}

	opts := createOpts{
		AdminStateUp: gophercloud.Enabled,
		Name:         networkName,
	}

	if ptr.Deref(openStackCluster.Spec.DisablePortSecurity, false) {
		opts.PortSecurityEnabled = gophercloud.Disabled
	}

	if openStackCluster.Spec.NetworkMTU != nil {
		opts.MTU = openStackCluster.Spec.NetworkMTU
	}

	network, err := s.client.CreateNetwork(opts)
	if err != nil {
		record.Warnf(openStackCluster, "FailedCreateNetwork", "Failed to create network %s: %v", networkName, err)
		return err
	}
	record.Eventf(openStackCluster, "SuccessfulCreateNetwork", "Created network %s with id %s", networkName, network.ID)

	if len(openStackCluster.Spec.Tags) > 0 {
		_, err = s.client.ReplaceAllAttributesTags("networks", network.ID, attributestags.ReplaceAllOpts{
			Tags: openStackCluster.Spec.Tags,
		})
		if err != nil {
			return err
		}
	}

	openStackCluster.Status.Network = &infrav1.NetworkStatusWithSubnets{}
	openStackCluster.Status.Network.ID = network.ID
	openStackCluster.Status.Network.Name = network.Name
	openStackCluster.Status.Network.Tags = openStackCluster.Spec.Tags
	return nil
}

func (s *Service) DeleteNetwork(openStackCluster *infrav1.OpenStackCluster, clusterResourceName string) error {
	networkName := getNetworkName(clusterResourceName)
	network, err := s.getNetworkByName(networkName)
	if err != nil {
		return err
	}
	if network.ID == "" {
		return nil
	}

	err = s.client.DeleteNetwork(network.ID)
	if err != nil {
		record.Warnf(openStackCluster, "FailedDeleteNetwork", "Failed to delete network %s with id %s: %v", network.Name, network.ID, err)
		return err
	}

	record.Eventf(openStackCluster, "SuccessfulDeleteNetwork", "Deleted network %s with id %s", network.Name, network.ID)
	return nil
}

func (s *Service) ReconcileSubnet(openStackCluster *infrav1.OpenStackCluster, clusterResourceName string) error {
	if openStackCluster.Status.Network == nil || openStackCluster.Status.Network.ID == "" {
		s.scope.Logger().V(4).Info("No need to reconcile network components since no network exists")
		return nil
	}

	subnetName := getSubnetName(clusterResourceName)
	s.scope.Logger().Info("Reconciling subnet", "name", subnetName)

	subnetList, err := s.client.ListSubnet(subnets.ListOpts{
		NetworkID: openStackCluster.Status.Network.ID,
		// Currently we only support 1 SubnetSpec.
		CIDR: openStackCluster.Spec.ManagedSubnets[0].CIDR,
	})
	if err != nil {
		return err
	}

	if len(subnetList) > 1 {
		return fmt.Errorf("found %d subnets with the CIDR %s and network %s, which should not happen",
			len(subnetList), openStackCluster.Spec.ManagedSubnets[0], openStackCluster.Status.Network.ID)
	}

	var subnet *subnets.Subnet
	if len(subnetList) == 0 {
		var err error
		subnet, err = s.createSubnet(openStackCluster, clusterResourceName, subnetName)
		if err != nil {
			return err
		}
	} else if len(subnetList) == 1 {
		subnet = &subnetList[0]
		s.scope.Logger().V(6).Info("Reusing existing subnet", "name", subnet.Name, "id", subnet.ID)
	}

	openStackCluster.Status.Network.Subnets = []infrav1.Subnet{
		{
			ID:   subnet.ID,
			Name: subnet.Name,
			CIDR: subnet.CIDR,
			Tags: subnet.Tags,
		},
	}
	return nil
}

func (s *Service) createSubnet(openStackCluster *infrav1.OpenStackCluster, clusterResourceName string, name string) (*subnets.Subnet, error) {
	opts := subnets.CreateOpts{
		NetworkID:      openStackCluster.Status.Network.ID,
		Name:           name,
		IPVersion:      4,
		CIDR:           openStackCluster.Spec.ManagedSubnets[0].CIDR,
		DNSNameservers: openStackCluster.Spec.ManagedSubnets[0].DNSNameservers,
		Description:    names.GetDescription(clusterResourceName),
	}

	for _, pool := range openStackCluster.Spec.ManagedSubnets[0].AllocationPools {
		opts.AllocationPools = append(opts.AllocationPools, subnets.AllocationPool{Start: pool.Start, End: pool.End})
	}

	subnet, err := s.client.CreateSubnet(opts)
	if err != nil {
		record.Warnf(openStackCluster, "FailedCreateSubnet", "Failed to create subnet %s: %v", name, err)
		return nil, err
	}
	record.Eventf(openStackCluster, "SuccessfulCreateSubnet", "Created subnet %s with id %s", name, subnet.ID)

	if len(openStackCluster.Spec.Tags) > 0 {
		mc := metrics.NewMetricPrometheusContext("subnet", "update")
		_, err = s.client.ReplaceAllAttributesTags("subnets", subnet.ID, attributestags.ReplaceAllOpts{
			Tags: openStackCluster.Spec.Tags,
		})
		if mc.ObserveRequest(err) != nil {
			return nil, err
		}
	}

	return subnet, nil
}

func (s *Service) getNetworkByName(networkName string) (networks.Network, error) {
	opts := networks.ListOpts{
		Name: networkName,
	}

	networkList, err := s.client.ListNetwork(opts)
	if err != nil {
		return networks.Network{}, err
	}

	switch len(networkList) {
	case 0:
		return networks.Network{}, nil
	case 1:
		return networkList[0], nil
	}
	return networks.Network{}, fmt.Errorf("found %d networks with the name %s, which should not happen", len(networkList), networkName)
}

type GetNetworkOpts func(networks.ListOptsBuilder) networks.ListOptsBuilder

func ExternalNetworksOnly(opts networks.ListOptsBuilder) networks.ListOptsBuilder {
	return &external.ListOptsExt{
		ListOptsBuilder: opts,
		External:        ptr.To(true),
	}
}

// getNetworksByFilter retrieves networks by querying openstack with filters.
func (s *Service) getNetworkByFilter(filter *infrav1.NetworkFilter, opts ...GetNetworkOpts) (*networks.Network, error) {
	var listOpts networks.ListOptsBuilder
	listOpts = filterconvert.NetworkFilterToListOpts(filter)
	for _, opt := range opts {
		listOpts = opt(listOpts)
	}

	networks, err := s.client.ListNetwork(listOpts)
	if err != nil {
		return nil, err
	}
	if len(networks) == 0 {
		return nil, capoerrors.ErrNoMatches
	}
	if len(networks) > 1 {
		return nil, capoerrors.ErrMultipleMatches
	}
	return &networks[0], nil
}

// GetNetworkByParam gets the network specified by the given NetworkParam.
func (s *Service) GetNetworkByParam(param *infrav1.NetworkParam, opts ...GetNetworkOpts) (*networks.Network, error) {
	if param.ID != nil {
		return s.GetNetworkByID(*param.ID)
	}

	if param.Filter == nil {
		return nil, errors.New("no filter or ID provided")
	}

	return s.getNetworkByFilter(param.Filter, opts...)
}

// GetNetworkIDByParam returns the ID of the network specified by the given
// NetworkParam. It does not make an OpenStack call if the network is specified
// by ID.
func (s *Service) GetNetworkIDByParam(param *infrav1.NetworkParam, opts ...GetNetworkOpts) (string, error) {
	if param.ID != nil {
		return *param.ID, nil
	}

	if param.Filter == nil {
		return "", errors.New("no filter or ID provided")
	}

	network, err := s.getNetworkByFilter(param.Filter, opts...)
	if err != nil {
		return "", err
	}
	return network.ID, nil
}

// GetSubnetsByFilter gets the id of a subnet by querying openstack with filters.
func (s *Service) GetSubnetsByFilter(opts subnets.ListOptsBuilder) ([]subnets.Subnet, error) {
	if opts == nil {
		return []subnets.Subnet{}, fmt.Errorf("no Filters were passed")
	}
	subnetList, err := s.client.ListSubnet(opts)
	if err != nil {
		return []subnets.Subnet{}, err
	}
	if len(subnetList) == 0 {
		return nil, capoerrors.ErrNoMatches
	}
	return subnetList, nil
}

// GetSubnetIDByParam gets the id of a subnet from the given SubnetParam. It
// does not make any OpenStack API calls if the subnet is specified by ID.
func (s *Service) GetSubnetIDByParam(param *infrav1.SubnetParam) (string, error) {
	if param.ID != nil {
		return *param.ID, nil
	}
	subnet, err := s.GetSubnetByParam(param)
	if err != nil {
		return "", err
	}
	return subnet.ID, nil
}

// GetSubnetByParam gets a single subnet specified by the given SubnetParam
// It returns an ErrFilterMatch if no or multiple subnets are found.
func (s *Service) GetSubnetByParam(param *infrav1.SubnetParam) (*subnets.Subnet, error) {
	return s.GetNetworkSubnetByParam("", param)
}

// GetNetworkSubnetByParam gets a single subnet of the given network, specified by the given SubnetParam.
// It returns an ErrFilterMatch if no or multiple subnets are found.
func (s *Service) GetNetworkSubnetByParam(networkID string, param *infrav1.SubnetParam) (*subnets.Subnet, error) {
	if param.ID != nil {
		subnet, err := s.client.GetSubnet(*param.ID)
		if capoerrors.IsNotFound(err) {
			return nil, capoerrors.ErrNoMatches
		}

		if networkID != "" && subnet.NetworkID != networkID {
			s.scope.Logger().V(4).Info("Subnet specified by ID does not belong to the given network", "subnetID", subnet.ID, "networkID", networkID)
			return nil, capoerrors.ErrNoMatches
		}
		return subnet, err
	}

	if param.Filter == nil {
		// Should have been caught by validation
		return nil, errors.New("subnet filter: both id and filter are nil")
	}

	listOpts := filterconvert.SubnetFilterToListOpts(param.Filter)
	if networkID != "" {
		listOpts.NetworkID = networkID
	}

	subnets, err := s.GetSubnetsByFilter(listOpts)
	if err != nil {
		return nil, err
	}
	if len(subnets) == 0 {
		return nil, capoerrors.ErrNoMatches
	}
	if len(subnets) > 1 {
		return nil, capoerrors.ErrMultipleMatches
	}
	return &subnets[0], nil
}

func getSubnetName(clusterResourceName string) string {
	return fmt.Sprintf("%s-cluster-%s", networkPrefix, clusterResourceName)
}

func getNetworkName(clusterResourceName string) string {
	return fmt.Sprintf("%s-cluster-%s", networkPrefix, clusterResourceName)
}

// GetNetworkByID retrieves network by the ID.
func (s *Service) GetNetworkByID(networkID string) (*networks.Network, error) {
	network, err := s.client.GetNetwork(networkID)
	if err != nil {
		return &networks.Network{}, err
	}
	return network, nil
}
