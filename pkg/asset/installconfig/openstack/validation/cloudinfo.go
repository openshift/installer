package validation

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	"github.com/gophercloud/utils/openstack/clientconfig"
	flavorutils "github.com/gophercloud/utils/openstack/compute/v2/flavors"
	networkutils "github.com/gophercloud/utils/openstack/networking/v2/networks"
	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/types"
)

// CloudInfo caches data fetched from the user's openstack cloud
type CloudInfo struct {
	ExternalNetwork *networks.Network
	Flavors         map[string]*flavors.Flavor
	MachinesSubnet  *subnets.Subnet

	clients *clients
}

type clients struct {
	networkClient *gophercloud.ServiceClient
	computeClient *gophercloud.ServiceClient
}

// GetCloudInfo fetches and caches metadata from openstack
func GetCloudInfo(ic *types.InstallConfig) (*CloudInfo, error) {
	var err error
	ci := CloudInfo{
		clients: &clients{},
		Flavors: map[string]*flavors.Flavor{},
	}

	opts := &clientconfig.ClientOpts{Cloud: ic.OpenStack.Cloud}

	ci.clients.networkClient, err = clientconfig.NewServiceClient("network", opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a network client")
	}

	ci.clients.computeClient, err = clientconfig.NewServiceClient("compute", opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a compute client")
	}

	err = ci.collectInfo(ic)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate OpenStack cloud info")
	}

	return &ci, nil
}

func (ci *CloudInfo) collectInfo(ic *types.InstallConfig) error {
	var err error

	ci.ExternalNetwork, err = ci.getNetwork(ic.OpenStack.ExternalNetwork)
	if err != nil {
		return errors.Wrap(err, "failed to fetch external network info")
	}

	ci.Flavors[ic.OpenStack.FlavorName], err = ci.getFlavor(ic.OpenStack.FlavorName)
	if err != nil {
		return errors.Wrap(err, "failed to fetch platform flavor info")
	}

	if ic.ControlPlane != nil && ic.ControlPlane.Platform.OpenStack != nil {
		crtlPlaneFlavor := ic.ControlPlane.Platform.OpenStack.FlavorName
		if crtlPlaneFlavor != "" {
			ci.Flavors[crtlPlaneFlavor], err = ci.getFlavor(crtlPlaneFlavor)
			if err != nil {
				return err
			}
		}
	}

	for _, machine := range ic.Compute {
		if machine.Platform.OpenStack != nil {
			flavorName := machine.Platform.OpenStack.FlavorName
			if flavorName != "" {
				if _, seen := ci.Flavors[flavorName]; !seen {
					ci.Flavors[flavorName], err = ci.getFlavor(flavorName)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	ci.MachinesSubnet, err = ci.getSubnet(ic.OpenStack.MachinesSubnet)
	if err != nil {
		return errors.Wrap(err, "failed to fetch machine subnet info")
	}

	return nil
}

func (ci *CloudInfo) getSubnet(subnetID string) (*subnets.Subnet, error) {
	subnet, err := subnets.Get(ci.clients.networkClient, subnetID).Extract()
	if err != nil {
		var gerr *gophercloud.ErrResourceNotFound
		if errors.As(err, &gerr) {
			return nil, nil
		}
		return nil, err
	}

	return subnet, nil
}

func (ci *CloudInfo) getFlavor(flavorName string) (*flavors.Flavor, error) {
	flavorID, err := flavorutils.IDFromName(ci.clients.computeClient, flavorName)
	if err != nil {
		var gerr *gophercloud.ErrResourceNotFound
		if errors.As(err, &gerr) {
			return nil, nil
		}
		return nil, err
	}

	flavor, err := flavors.Get(ci.clients.computeClient, flavorID).Extract()
	if err != nil {
		return nil, err
	}

	return flavor, nil
}

func (ci *CloudInfo) getNetwork(networkName string) (*networks.Network, error) {
	if networkName == "" {
		return nil, nil
	}
	networkID, err := networkutils.IDFromName(ci.clients.networkClient, networkName)
	if err != nil {
		var gerr *gophercloud.ErrResourceNotFound
		if errors.As(err, &gerr) {
			return nil, nil
		}
		return nil, err
	}

	network, err := networks.Get(ci.clients.networkClient, networkID).Extract()
	if err != nil {
		return nil, err
	}

	return network, nil
}
