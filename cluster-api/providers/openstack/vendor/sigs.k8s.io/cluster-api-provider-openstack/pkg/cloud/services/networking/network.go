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
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/external"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	"k8s.io/utils/pointer"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/metrics"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/record"
	capoerrors "sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/errors"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/filterconvert"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/names"
)

var (
	ErrFilterMatch     = fmt.Errorf("filter match error")
	ErrMultipleMatches = multipleMatchesError{}
	ErrNoMatches       = noMatchesError{}
)

type (
	multipleMatchesError struct{}
	noMatchesError       struct{}
)

func (e multipleMatchesError) Error() string {
	return "filter matched more than one resource"
}

func (e multipleMatchesError) Is(err error) bool {
	return err == ErrFilterMatch
}

func (e noMatchesError) Error() string {
	return "filter matched no resources"
}

func (e noMatchesError) Is(err error) bool {
	return err == ErrFilterMatch
}

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
	var listOpts external.ListOptsExt

	if pointer.BoolDeref(openStackCluster.Spec.DisableExternalNetwork, false) {
		s.scope.Logger().Info("External network is disabled - proceeding with internal network only")
		openStackCluster.Status.ExternalNetwork = nil
		return nil
	}

	listOpts = external.ListOptsExt{
		ListOptsBuilder: filterconvert.NetworkFilterToListOpts(openStackCluster.Spec.ExternalNetwork),
		External:        pointer.Bool(true),
	}
	networkList, err := s.client.ListNetwork(listOpts)
	if err != nil {
		return err
	}

	switch len(networkList) {
	case 0:
		if openStackCluster.Spec.ExternalNetwork == nil {
			// Not finding an external network is fine if ExternalNetwork is not set
			openStackCluster.Status.ExternalNetwork = nil
			s.scope.Logger().Info("No external network found - proceeding with internal network only")
			return nil
		}
		return fmt.Errorf("no external network found")
	case 1:
		openStackCluster.Status.ExternalNetwork = &infrav1.NetworkStatus{
			ID:   networkList[0].ID,
			Name: networkList[0].Name,
			Tags: networkList[0].Tags,
		}
		s.scope.Logger().Info("External network found", "id", networkList[0].ID)
		return nil
	}
	return fmt.Errorf("found %d external networks, which should not happen", len(networkList))
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

	if pointer.BoolDeref(openStackCluster.Spec.DisablePortSecurity, false) {
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

// GetNetworksByFilter retrieves networks by querying openstack with filters.
func (s *Service) GetNetworksByFilter(opts networks.ListOptsBuilder) ([]networks.Network, error) {
	if opts == nil {
		return nil, fmt.Errorf("no Filters were passed")
	}
	networkList, err := s.client.ListNetwork(opts)
	if err != nil {
		return nil, err
	}
	if len(networkList) == 0 {
		return nil, fmt.Errorf("no networks could be found with the filters provided")
	}
	return networkList, nil
}

// GetNetworkIDsByFilter retrieves network ids by querying openstack with filters.
func (s *Service) GetNetworkIDsByFilter(opts networks.ListOptsBuilder) ([]string, error) {
	nets, err := s.GetNetworksByFilter(opts)
	if err != nil {
		return nil, err
	}
	ids := []string{}
	for _, network := range nets {
		ids = append(ids, network.ID)
	}
	return ids, nil
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
		return nil, ErrNoMatches
	}
	return subnetList, nil
}

// GetSubnetByFilter gets a single subnet specified by the given SubnetFilter.
// It returns an ErrFilterMatch if no or multiple subnets are found.
func (s *Service) GetSubnetByFilter(filter *infrav1.SubnetFilter) (*subnets.Subnet, error) {
	listOpts := filterconvert.SubnetFilterToListOpts(filter)
	return s.getSubnetByFilter(listOpts)
}

// GetNetworkSubnetByFilter gets a single subnet of the given network, specified by the given SubnetFilter.
// It returns an ErrFilterMatch if no or multiple subnets are found.
func (s *Service) GetNetworkSubnetByFilter(networkID string, filter *infrav1.SubnetFilter) (*subnets.Subnet, error) {
	listOpts := filterconvert.SubnetFilterToListOpts(filter)
	listOpts.NetworkID = networkID

	return s.getSubnetByFilter(listOpts)
}

// getSubnetByFilter gets a single subnet specified by the given gophercloud ListOpts.
// It returns an ErrFilterMatch if no or multiple subnets are found.
func (s *Service) getSubnetByFilter(listOpts subnets.ListOpts) (*subnets.Subnet, error) {
	// If the ID is set, we can just get the subnet by ID.
	if listOpts.ID != "" {
		subnet, err := s.client.GetSubnet(listOpts.ID)
		if capoerrors.IsNotFound(err) {
			return nil, ErrNoMatches
		}
		return subnet, err
	}

	subnets, err := s.GetSubnetsByFilter(listOpts)
	if err != nil {
		return nil, err
	}
	if len(subnets) == 0 {
		return nil, ErrNoMatches
	}
	if len(subnets) > 1 {
		return nil, ErrMultipleMatches
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
