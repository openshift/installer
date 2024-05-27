package openstack

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/subnets"
	"github.com/gophercloud/gophercloud/v2/pagination"
	network_utils "github.com/gophercloud/utils/v2/openstack/networking/v2/networks"

	machinev1alpha1 "github.com/openshift/api/machine/v1alpha1"
	types_openstack "github.com/openshift/installer/pkg/types/openstack"
)

type terraformFixedIP struct {
	SubnetID  string `json:"subnet_id"`
	IPAddress string `json:"ip_address"`
}

type terraformPort struct {
	NetworkID string             `json:"network_id"`
	FixedIP   []terraformFixedIP `json:"fixed_ips"`
}

func portTargetToTerraformPort(ctx context.Context, networkClient *gophercloud.ServiceClient, portTarget types_openstack.PortTarget) (terraformPort, error) {
	networkID := portTarget.Network.ID
	if networkID == "" && portTarget.Network.Name != "" {
		var err error
		networkID, err = network_utils.IDFromName(ctx, networkClient, portTarget.Network.Name)
		if err != nil {
			return terraformPort{}, fmt.Errorf("failed to resolve network ID for network name %q: %w", portTarget.Network.Name, err)
		}
	}

	terraformFixedIPs := make([]terraformFixedIP, 0, len(portTarget.FixedIPs))
	for _, fixedIP := range portTarget.FixedIPs {
		subnetFilter := machinev1alpha1.SubnetFilter{
			ID:   fixedIP.Subnet.ID,
			Name: fixedIP.Subnet.Name,
		}
		resolvedSubnetID, resolvedNetworkID, err := resolveSubnetFilter(ctx, networkClient, networkID, subnetFilter)
		if err != nil {
			return terraformPort{}, fmt.Errorf("failed to resolve the subnet filter: %w", err)
		}

		if networkID == "" {
			networkID = resolvedNetworkID
		}

		if networkID != resolvedNetworkID {
			return terraformPort{}, fmt.Errorf("control plane port has ports on multiple networks")
		}

		terraformFixedIPs = append(terraformFixedIPs, terraformFixedIP{
			SubnetID: resolvedSubnetID,
		})
	}

	return terraformPort{
		NetworkID: networkID,
		FixedIP:   terraformFixedIPs,
	}, nil
}

func resolveSubnetFilter(ctx context.Context, networkClient *gophercloud.ServiceClient, networkID string, subnetFilter machinev1alpha1.SubnetFilter) (resolvedSubnetID, resolvedNetworkID string, err error) {
	if subnetFilter.ProjectID != "" {
		subnetFilter.TenantID = ""
	}
	if err = subnets.List(networkClient, subnets.ListOpts{
		NetworkID: networkID,
		Name:      subnetFilter.Name,
		ID:        subnetFilter.ID,
	}).EachPage(ctx, func(ctx context.Context, page pagination.Page) (bool, error) {
		returnedSubnets, err := subnets.ExtractSubnets(page)
		if err != nil {
			return false, err
		}
		for _, subnet := range returnedSubnets {
			if resolvedSubnetID == "" {
				resolvedSubnetID = subnet.ID
				resolvedNetworkID = subnet.NetworkID
			} else {
				return false, fmt.Errorf("more than one subnet found")
			}
		}
		return true, nil
	}); err != nil {
		return "", "", fmt.Errorf("failed to list subnets: %w", err)
	}

	if resolvedSubnetID == "" {
		return "", "", fmt.Errorf("no subnet found")
	}

	return resolvedSubnetID, resolvedNetworkID, err
}
