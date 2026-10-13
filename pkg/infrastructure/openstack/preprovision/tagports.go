package preprovision

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/subnets"
	"github.com/gophercloud/gophercloud/v2/pagination"
	network_utils "github.com/gophercloud/utils/v2/openstack/networking/v2/networks"

	machinev1alpha1 "github.com/openshift/api/machine/v1alpha1"
	"github.com/openshift/installer/pkg/asset/installconfig"
	types_openstack "github.com/openshift/installer/pkg/types/openstack"
	openstackdefaults "github.com/openshift/installer/pkg/types/openstack/defaults"
)

// TagVIPPorts tags the VIP Ports pre-created by the user.
func TagVIPPorts(ctx context.Context, installConfig *installconfig.InstallConfig, infraID string) error {
	networkClient, err := openstackdefaults.NewServiceClient(ctx, "network", openstackdefaults.DefaultClientOpts(installConfig.Config.Platform.OpenStack.Cloud))
	if err != nil {
		return fmt.Errorf("failed to build an OpenStack service client: %w", err)
	}

	if controlPlanePort := installConfig.Config.Platform.OpenStack.ControlPlanePort; controlPlanePort != nil {
		for _, vips := range [][]string{installConfig.Config.OpenStack.APIVIPs, installConfig.Config.OpenStack.IngressVIPs} {
			// Tagging the API and Ingress ports if pre-created by the user.
			if len(vips) > 0 {
				networkID, err := getNetworkIDByPortTarget(ctx, networkClient, *controlPlanePort)
				if err != nil {
					return err
				}

				// Assuming the VIP addresses are on the same port.
				port, err := getPortByNetworkIDAndIP(ctx, networkClient, networkID, vips[0])
				if err != nil {
					return err
				}

				if port != nil {
					err := attributestags.Add(ctx, networkClient, "ports", port.ID, infraID+openstackdefaults.DualStackVIPsPortTag).ExtractErr()
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func getNetworkIDByPortTarget(ctx context.Context, networkClient *gophercloud.ServiceClient, portTarget types_openstack.PortTarget) (string, error) {
	networkID := portTarget.Network.ID
	if networkID == "" && portTarget.Network.Name != "" {
		var err error
		networkID, err = network_utils.IDFromName(ctx, networkClient, portTarget.Network.Name)
		if err != nil {
			return "", fmt.Errorf("failed to resolve network ID for network name %q: %w", portTarget.Network.Name, err)
		}
	}

	for _, fixedIP := range portTarget.FixedIPs {
		subnetFilter := machinev1alpha1.SubnetFilter{
			ID:   fixedIP.Subnet.ID,
			Name: fixedIP.Subnet.Name,
		}
		_, resolvedNetworkID, err := resolveSubnetFilter(ctx, networkClient, networkID, subnetFilter)
		if err != nil {
			return "", fmt.Errorf("failed to resolve the subnet filter: %w", err)
		}

		if networkID == "" {
			networkID = resolvedNetworkID
		}

		if networkID != resolvedNetworkID {
			return "", fmt.Errorf("control plane port has ports on multiple networks")
		}
	}
	return networkID, nil
}

func getPortByNetworkIDAndIP(ctx context.Context, networkClient *gophercloud.ServiceClient, networkID, ipAddress string) (*ports.Port, error) {
	allPagesPort, err := ports.List(networkClient, ports.ListOpts{
		NetworkID: networkID,
		FixedIPs:  []ports.FixedIPOpts{{IPAddress: ipAddress}},
	}).AllPages(ctx)
	if err != nil {
		return nil, err
	}

	allPorts, err := ports.ExtractPorts(allPagesPort)
	if err != nil {
		return nil, err
	}

	switch len(allPorts) {
	case 0:
		return nil, nil
	case 1:
		return &allPorts[0], nil
	default:
		return nil, fmt.Errorf("found multiple ports matching network ID %s and IP %s, cannot proceed", networkID, ipAddress)
	}
}

func resolveSubnetFilter(ctx context.Context, networkClient *gophercloud.ServiceClient, networkID string, subnetFilter machinev1alpha1.SubnetFilter) (resolvedSubnetID, resolvedNetworkID string, err error) {
	if subnetFilter.ProjectID != "" {
		subnetFilter.TenantID = ""
	}
	if err = subnets.List(networkClient, subnets.ListOpts{
		NetworkID: networkID,
		Name:      subnetFilter.Name,
		ID:        subnetFilter.ID,
	}).EachPage(ctx, func(_ context.Context, page pagination.Page) (bool, error) {
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
