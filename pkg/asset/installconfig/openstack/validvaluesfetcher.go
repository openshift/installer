package openstack

import (
	"github.com/pkg/errors"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/utils/openstack/clientconfig"
	networkutils "github.com/gophercloud/utils/openstack/networking/v2/networks"
)

// getCloudNames gets the valid cloud names. These are read from clouds.yaml.
func getCloudNames() ([]string, error) {
	clouds, err := clientconfig.LoadCloudsYAML()
	if err != nil {
		return nil, err
	}

	cloudNames := []string{}
	for k := range clouds {
		cloudNames = append(cloudNames, k)
	}
	return cloudNames, nil
}

// getNetworkNames gets the valid network names.
func getNetworkNames(cloud string) ([]string, error) {
	conn, err := clientconfig.NewServiceClient("network", &clientconfig.ClientOpts{
		Cloud: cloud,
	})
	if err != nil {
		return nil, err
	}

	listOpts := networks.ListOpts{}
	allPages, err := networks.List(conn, listOpts).AllPages()
	if err != nil {
		return nil, err
	}

	allNetworks, err := networks.ExtractNetworks(allPages)
	if err != nil {
		return nil, err
	}

	networkNames := make([]string, len(allNetworks))
	for x, network := range allNetworks {
		networkNames[x] = network.Name
	}

	return networkNames, nil
}

// getFlavorNames gets a list of valid flavor names.
func getFlavorNames(cloud string) ([]string, error) {
	conn, err := clientconfig.NewServiceClient("compute", &clientconfig.ClientOpts{
		Cloud: cloud,
	})
	if err != nil {
		return nil, err
	}

	listOpts := flavors.ListOpts{}
	allPages, err := flavors.ListDetail(conn, listOpts).AllPages()
	if err != nil {
		return nil, err
	}

	allFlavors, err := flavors.ExtractFlavors(allPages)
	if err != nil {
		return nil, err
	}

	if len(allFlavors) == 0 {
		return nil, errors.New("no OpenStack flavors were found")
	}

	flavorNames := make([]string, len(allFlavors))
	for i, flavor := range allFlavors {
		flavorNames[i] = flavor.Name
	}

	return flavorNames, nil
}
func getFloatingIPNames(cloud string, floatingNetworkName string) ([]string, error) {
	conn, err := clientconfig.NewServiceClient("network", &clientconfig.ClientOpts{
		Cloud: cloud,
	})
	if err != nil {
		return nil, err
	}

	// floatingips.ListOpts requires an ID so we must get it from the name
	floatingNetworkID, err := networkutils.IDFromName(conn, floatingNetworkName)
	if err != nil {
		return nil, err
	}

	// Only show IPs that belong to the network and are not in use
	listOpts := floatingips.ListOpts{
		FloatingNetworkID: floatingNetworkID,
		Status:            "DOWN",
	}

	allPages, err := floatingips.List(conn, listOpts).AllPages()
	if err != nil {
		return nil, err
	}

	allFloatingIPs, err := floatingips.ExtractFloatingIPs(allPages)
	if err != nil {
		return nil, err
	}

	if len(allFloatingIPs) == 0 {
		return nil, errors.New("there are no unassigned floating IP addresses available")
	}

	floatingIPNames := make([]string, len(allFloatingIPs))
	for i, floatingIP := range allFloatingIPs {
		floatingIPNames[i] = floatingIP.FloatingIP
	}

	return floatingIPNames, nil
}
