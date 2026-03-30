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

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/subnets"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/record"
	capoerrors "sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/errors"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/filterconvert"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/names"
)

func (s *Service) ReconcileRouter(openStackCluster *infrav1.OpenStackCluster, clusterResourceName string) error {
	if openStackCluster.Status.Network == nil || openStackCluster.Status.Network.ID == "" {
		s.scope.Logger().V(3).Info("No need to reconcile router since no network exists")
		return nil
	}
	if len(openStackCluster.Status.Network.Subnets) == 0 {
		s.scope.Logger().V(4).Info("No need to reconcile router since no subnet exists")
		return nil
	}
	if openStackCluster.Status.ExternalNetwork == nil || openStackCluster.Status.ExternalNetwork.ID == "" {
		s.scope.Logger().V(3).Info("No need to create router, due to missing ExternalNetworkID")
		return nil
	}

	s.scope.Logger().Info("Reconciling router", "cluster", clusterResourceName)

	router, err := s.getExistingRouter(openStackCluster, clusterResourceName)
	if err != nil {
		return fmt.Errorf("fetching router: %w", err)
	}

	if router == nil {
		if openStackCluster.Spec.Router != nil {
			// Should not happen: will have returned ErrNoMatches above
			return fmt.Errorf("router not found")
		}

		var err error
		router, err = s.createRouter(openStackCluster, clusterResourceName, getRouterName(clusterResourceName))
		if err != nil {
			return err
		}
	}

	routerIPs := []string{}
	for _, ip := range router.GatewayInfo.ExternalFixedIPs {
		routerIPs = append(routerIPs, ip.IPAddress)
	}

	openStackCluster.Status.Router = &infrav1.Router{
		Name: router.Name,
		ID:   router.ID,
		Tags: router.Tags,
		IPs:  routerIPs,
	}

	if len(openStackCluster.Spec.ExternalRouterIPs) > 0 {
		if err := s.setRouterExternalIPs(openStackCluster, router); err != nil {
			return err
		}
	}

	routerInterfaces, err := s.getRouterInterfaces(router.ID)
	if err != nil {
		return err
	}

	// check all router interfaces for an existing port in a cluster subnet.
	for _, subnet := range openStackCluster.Status.Network.Subnets {
		createInterface := true
	INTERFACE_LOOP:
		for _, iface := range routerInterfaces {
			for _, ip := range iface.FixedIPs {
				if ip.SubnetID == subnet.ID {
					createInterface = false
					break INTERFACE_LOOP
				}
			}
		}

		// ... and create a router interface for our subnet.
		if createInterface {
			s.scope.Logger().V(4).Info("Creating RouterInterface", "routerID", router.ID, "subnetID", subnet.ID)
			routerInterface, err := s.client.AddRouterInterface(router.ID, routers.AddInterfaceOpts{
				SubnetID: subnet.ID,
			})
			if err != nil {
				return fmt.Errorf("unable to create router interface: %v", err)
			}
			s.scope.Logger().V(4).Info("Created RouterInterface", "id", routerInterface.ID)
		}
	}
	return nil
}

func (s *Service) getExistingRouter(openStackCluster *infrav1.OpenStackCluster, clusterResourceName string) (*routers.Router, error) {
	// For an externally-managed router we always expect it to exist. We will return an error if it doesn't.
	if openStackCluster.Spec.Router != nil {
		return s.getExternallyManagedRouter(openStackCluster)
	}

	// A managed router may not exist either because we haven't created it
	// or because we deleted it. Swallow NotFound errors and return nil.

	if openStackCluster.Status.Router != nil {
		router, err := s.client.GetRouter(openStackCluster.Status.Router.ID)
		if capoerrors.IsNotFound(err) {
			return nil, nil
		}
		return router, err
	}

	routerName := getRouterName(clusterResourceName)
	listOpts := routers.ListOpts{Name: routerName}
	if openStackCluster.Spec.Router != nil {
		listOpts = filterconvert.RouterFilterToListOpts(openStackCluster.Spec.Router.Filter)
	}

	router, err := s.getRouterByFilter(listOpts)
	// It's ok if our managed router doesn't exist yet, just return nil
	if errors.Is(err, capoerrors.ErrNoMatches) {
		return nil, nil
	}
	return router, err
}

func (s *Service) getExternallyManagedRouter(openStackCluster *infrav1.OpenStackCluster) (*routers.Router, error) {
	if openStackCluster.Spec.Router == nil {
		return nil, fmt.Errorf("getExternallyManagedRouter called with no external router specified")
	}

	// Fetch by ID if we previously resolved it
	if openStackCluster.Status.Router != nil {
		return s.client.GetRouter(openStackCluster.Status.Router.ID)
	}
	return s.GetRouterByParam(openStackCluster.Spec.Router)
}

func (s *Service) GetRouterByParam(routerParam *infrav1.RouterParam) (*routers.Router, error) {
	if routerParam.ID != nil {
		return s.client.GetRouter(*routerParam.ID)
	}

	if routerParam.Filter == nil {
		// Should have been caught by validation
		return nil, errors.New("invalid router param, either ID or Filter must be set")
	}

	listOpts := filterconvert.RouterFilterToListOpts(routerParam.Filter)
	return s.getRouterByFilter(listOpts)
}

func (s *Service) createRouter(openStackCluster *infrav1.OpenStackCluster, clusterResourceName, name string) (*routers.Router, error) {
	opts := routers.CreateOpts{
		Description: names.GetDescription(clusterResourceName),
		Name:        name,
	}
	// only set the GatewayInfo right now when no externalIPs
	// should be configured because at least in our environment
	// we can only set the routerIP via gateway update not during create
	// That's also the same way terraform provider OpenStack does it
	if len(openStackCluster.Spec.ExternalRouterIPs) == 0 {
		opts.GatewayInfo = &routers.GatewayInfo{
			NetworkID: openStackCluster.Status.ExternalNetwork.ID,
		}
	}

	router, err := s.client.CreateRouter(opts)
	if err != nil {
		record.Warnf(openStackCluster, "FailedCreateRouter", "Failed to create router %s: %v", name, err)
		return nil, err
	}
	record.Eventf(openStackCluster, "SuccessfulCreateRouter", "Created router %s with id %s", name, router.ID)

	if len(openStackCluster.Spec.Tags) > 0 {
		_, err = s.client.ReplaceAllAttributesTags("routers", router.ID, attributestags.ReplaceAllOpts{
			Tags: openStackCluster.Spec.Tags,
		})
		if err != nil {
			return nil, err
		}
	}

	return router, nil
}

func (s *Service) setRouterExternalIPs(openStackCluster *infrav1.OpenStackCluster, router *routers.Router) error {
	updateOpts := routers.UpdateOpts{
		GatewayInfo: &routers.GatewayInfo{
			NetworkID: openStackCluster.Status.ExternalNetwork.ID,
		},
	}

	for i := range openStackCluster.Spec.ExternalRouterIPs {
		externalRouterIP := openStackCluster.Spec.ExternalRouterIPs[i]
		subnetID, err := s.GetSubnetIDByParam(&externalRouterIP.Subnet)
		if err != nil {
			return fmt.Errorf("failed to get subnet for external router: %w", err)
		}
		updateOpts.GatewayInfo.ExternalFixedIPs = append(updateOpts.GatewayInfo.ExternalFixedIPs, routers.ExternalFixedIP{
			IPAddress: externalRouterIP.FixedIP,
			SubnetID:  subnetID,
		})
	}

	_, err := s.client.UpdateRouter(router.ID, updateOpts)
	if err != nil {
		record.Warnf(openStackCluster, "FailedUpdateRouter", "Failed to update router %s with id %s: %v", router.Name, router.ID, err)
		return err
	}

	record.Eventf(openStackCluster, "SuccessfulUpdateRouter", "Updated router %s with id %s", router.Name, router.ID)
	return nil
}

func (s *Service) DeleteRouter(openStackCluster *infrav1.OpenStackCluster, clusterResourceName string) error {
	router, err := s.getExistingRouter(openStackCluster, clusterResourceName)
	if err != nil {
		return err
	}

	if router == nil {
		return nil
	}

	subnetName := getSubnetName(clusterResourceName)
	subnet, err := s.getSubnetByName(subnetName)
	if err != nil {
		return err
	}

	if subnet.ID != "" {
		_, err = s.client.RemoveRouterInterface(router.ID, routers.RemoveInterfaceOpts{
			SubnetID: subnet.ID,
		})
		if err != nil {
			if !capoerrors.IsNotFound(err) {
				return fmt.Errorf("unable to remove router interface: %v", err)
			}
			s.scope.Logger().V(4).Info("Router interface already removed, nothing to do", "id", router.ID)
		} else {
			s.scope.Logger().V(4).Info("Removed RouterInterface of router", "id", router.ID)
		}
	}

	if openStackCluster.Spec.Router != nil {
		s.scope.Logger().V(4).Info("Not deleting pre-existing router", "name", router.Name)
		return nil
	}

	err = s.client.DeleteRouter(router.ID)
	if err != nil {
		record.Warnf(openStackCluster, "FailedDeleteRouter", "Failed to delete router %s with id %s: %v", router.Name, router.ID, err)
		return err
	}

	record.Eventf(openStackCluster, "SuccessfulDeleteRouter", "Deleted router %s with id %s", router.Name, router.ID)
	return nil
}

func (s *Service) getRouterInterfaces(routerID string) ([]ports.Port, error) {
	return s.client.ListPort(ports.ListOpts{
		DeviceID: routerID,
	})
}

func (s *Service) getRouterByFilter(opts routers.ListOpts) (*routers.Router, error) {
	routerList, err := s.client.ListRouter(opts)
	if err != nil {
		return nil, err
	}

	switch len(routerList) {
	case 0:
		return nil, capoerrors.ErrNoMatches
	case 1:
		return &routerList[0], nil
	}
	return nil, capoerrors.ErrMultipleMatches
}

func (s *Service) getSubnetByName(subnetName string) (subnets.Subnet, error) {
	opts := subnets.ListOpts{
		Name: subnetName,
	}

	subnetList, err := s.client.ListSubnet(opts)
	if err != nil {
		return subnets.Subnet{}, err
	}

	switch len(subnetList) {
	case 0:
		return subnets.Subnet{}, nil
	case 1:
		return subnetList[0], nil
	}
	return subnets.Subnet{}, fmt.Errorf("found %d subnets with the name %s, which should not happen", len(subnetList), subnetName)
}

func getRouterName(clusterResourceName string) string {
	return fmt.Sprintf("%s-cluster-%s", networkPrefix, clusterResourceName)
}
