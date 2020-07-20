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
	PlatformFlavor  *flavors.Flavor
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
	}

	opts := &clientconfig.ClientOpts{Cloud: ic.OpenStack.Cloud}

	ci.clients.networkClient, err = clientconfig.NewServiceClient("network", opts)
	if err != nil {
		return nil, err
	}

	ci.clients.computeClient, err = clientconfig.NewServiceClient("compute", opts)
	if err != nil {
		return nil, err
	}

	err = ci.collectInfo(ic)
	if err != nil {
		return nil, err
	}

	return &ci, nil
}

func (ci *CloudInfo) collectInfo(ic *types.InstallConfig) error {
	var err error

	ci.ExternalNetwork, err = ci.getNetwork(ic.OpenStack.ExternalNetwork)
	if err != nil {
		return err
	}

	ci.PlatformFlavor, err = ci.getFlavor(ic.OpenStack.FlavorName)
	if err != nil {
		return err
	}

	ci.MachinesSubnet, err = ci.getSubnet(ic.OpenStack.MachinesSubnet)
	if err != nil {
		return err
	}

	return nil
}

func (ci *CloudInfo) getSubnet(subnetID string) (*subnets.Subnet, error) {
	subnet, err := subnets.Get(ci.clients.networkClient, subnetID).Extract()
	if err != nil {
		if errors.Is(err, gophercloud.ErrResourceNotFound{}) {
			return nil, nil
		}
		return nil, err
	}

	return subnet, nil
}

func (ci *CloudInfo) getFlavor(flavorName string) (*flavors.Flavor, error) {
	flavorID, err := flavorutils.IDFromName(ci.clients.computeClient, flavorName)
	if err != nil {
		if errors.Is(err, gophercloud.ErrResourceNotFound{}) {
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
		return &networks.Network{}, nil
	}
	networkID, err := networkutils.IDFromName(ci.clients.networkClient, networkName)
	if err != nil {
		if errors.Is(err, gophercloud.ErrResourceNotFound{}) {
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
