package openstack

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/external"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	"github.com/gophercloud/utils/v2/openstack/clientconfig"
	networkutils "github.com/gophercloud/utils/v2/openstack/networking/v2/networks"

	openstackdefaults "github.com/openshift/installer/pkg/types/openstack/defaults"
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

// getExternalNetworkNames interrogates OpenStack to get the external network
// names.
func getExternalNetworkNames(ctx context.Context, cloud string) ([]string, error) {
	conn, err := openstackdefaults.NewServiceClient(ctx, "network", openstackdefaults.DefaultClientOpts(cloud))
	if err != nil {
		return nil, err
	}

	iTrue := true
	listOpts := external.ListOptsExt{
		ListOptsBuilder: networks.ListOpts{},
		External:        &iTrue,
	}

	allPages, err := networks.List(conn, listOpts).AllPages(ctx)
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
func getFlavorNames(ctx context.Context, cloud string) ([]string, error) {
	conn, err := openstackdefaults.NewServiceClient(ctx, "compute", openstackdefaults.DefaultClientOpts(cloud))
	if err != nil {
		return nil, err
	}

	listOpts := flavors.ListOpts{}
	allPages, err := flavors.ListDetail(conn, listOpts).AllPages(ctx)
	if err != nil {
		return nil, err
	}

	allFlavors, err := flavors.ExtractFlavors(allPages)
	if err != nil {
		return nil, err
	}

	if len(allFlavors) == 0 {
		return nil, fmt.Errorf("no OpenStack flavors were found")
	}

	flavorNames := make([]string, len(allFlavors))
	for i, flavor := range allFlavors {
		flavorNames[i] = flavor.Name
	}

	return flavorNames, nil
}

type sortableFloatingIPCollection []floatingips.FloatingIP

func (fips sortableFloatingIPCollection) Len() int { return len(fips) }
func (fips sortableFloatingIPCollection) Less(i, j int) bool {
	return fips[i].FloatingIP < fips[j].FloatingIP
}
func (fips sortableFloatingIPCollection) Swap(i, j int) {
	fips[i], fips[j] = fips[j], fips[i]
}

func (fips sortableFloatingIPCollection) Names() []string {
	names := make([]string, len(fips))
	for i := range fips {
		names[i] = fips[i].FloatingIP
	}
	return names
}

func (fips sortableFloatingIPCollection) Description(index int) string {
	return fips[index].Description
}

func (fips sortableFloatingIPCollection) Contains(value string) bool {
	for i := range fips {
		if value == fips[i].FloatingIP {
			return true
		}
	}
	return false
}

func getFloatingIPs(ctx context.Context, cloud string, floatingNetworkName string) (sortableFloatingIPCollection, error) {
	conn, err := openstackdefaults.NewServiceClient(ctx, "network", openstackdefaults.DefaultClientOpts(cloud))
	if err != nil {
		return nil, err
	}

	// floatingips.ListOpts requires an ID so we must get it from the name
	floatingNetworkID, err := networkutils.IDFromName(ctx, conn, floatingNetworkName)
	if err != nil {
		return nil, err
	}

	// Only show IPs that belong to the network and are not in use
	listOpts := floatingips.ListOpts{
		FloatingNetworkID: floatingNetworkID,
		Status:            "DOWN",
	}

	allPages, err := floatingips.List(conn, listOpts).AllPages(ctx)
	if err != nil {
		return nil, err
	}

	allFloatingIPs, err := floatingips.ExtractFloatingIPs(allPages)
	if err != nil {
		return nil, err
	}

	if len(allFloatingIPs) == 0 {
		return nil, fmt.Errorf("there are no unassigned floating IP addresses available")
	}

	return allFloatingIPs, nil
}
