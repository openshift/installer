package openstack

import (
	"github.com/pkg/errors"

	"github.com/gophercloud/gophercloud/openstack/common/extensions"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	netext "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	"github.com/gophercloud/utils/openstack/clientconfig"

	"github.com/openshift/installer/pkg/types/openstack/validation"
)

type realValidValuesFetcher struct{}

// NewValidValuesFetcher returns a new ValidValuesFetcher.
func NewValidValuesFetcher() validation.ValidValuesFetcher {
	return realValidValuesFetcher{}
}

// GetCloudNames gets the valid cloud names. These are read from clouds.yaml.
func (f realValidValuesFetcher) GetCloudNames() ([]string, error) {
	clouds, err := new(yamlLoadOpts).LoadCloudsYAML()
	if err != nil {
		return nil, err
	}

	i := 0
	cloudNames := make([]string, len(clouds))
	for k := range clouds {
		cloudNames[i] = k
		i++
	}
	return cloudNames, nil
}

// GetNetworkNames gets the valid network names.
func (f realValidValuesFetcher) GetNetworkNames(cloud string) ([]string, error) {
	opts := defaultClientOpts(cloud)

	conn, err := clientconfig.NewServiceClient("network", opts)
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

// GetFlavorNames gets a list of valid flavor names.
func (f realValidValuesFetcher) GetFlavorNames(cloud string) ([]string, error) {
	opts := defaultClientOpts(cloud)

	conn, err := clientconfig.NewServiceClient("compute", opts)
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

func (f realValidValuesFetcher) GetNetworkExtensionsAliases(cloud string) ([]string, error) {
	opts := defaultClientOpts(cloud)

	conn, err := clientconfig.NewServiceClient("network", opts)
	if err != nil {
		return nil, err
	}

	allPages, err := netext.List(conn).AllPages()
	if err != nil {
		return nil, err
	}

	allExts, err := extensions.ExtractExtensions(allPages)
	if err != nil {
		return nil, err
	}

	extAliases := make([]string, len(allExts))
	for i, ext := range allExts {
		extAliases[i] = ext.Alias
	}

	return extAliases, nil
}

func (f realValidValuesFetcher) GetServiceCatalog(cloud string) ([]string, error) {
	opts := defaultClientOpts(cloud)

	conn, err := clientconfig.NewServiceClient("identity", opts)
	if err != nil {
		return nil, err
	}

	authResult := conn.GetAuthResult()
	auth, ok := authResult.(tokens.CreateResult)
	if !ok {
		return nil, errors.New("unable to extract service catalog")
	}

	allServices, err := auth.ExtractServiceCatalog()
	if err != nil {
		return nil, err
	}

	serviceCatalogNames := make([]string, len(allServices.Entries))
	for i, service := range allServices.Entries {
		serviceCatalogNames[i] = service.Name
	}

	return serviceCatalogNames, nil
}

func (f realValidValuesFetcher) GetFloatingIPNames(cloud string, floatingNetworkName string) ([]string, error) {
	opts := defaultClientOpts(cloud)

	conn, err := clientconfig.NewServiceClient("network", opts)
	if err != nil {
		return nil, err
	}

	// floatingips.ListOpts requires an ID so we must get it from the name
	floatingNetworkID, err := networks.IDFromName(conn, floatingNetworkName)
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

func (f realValidValuesFetcher) GetSubnetCIDR(cloud string, subnetID string) (string, error) {
	opts := &clientconfig.ClientOpts{
		Cloud: cloud,
	}

	networkClient, err := clientconfig.NewServiceClient("network", opts)
	if err != nil {
		return "", err
	}

	subnet, err := subnets.Get(networkClient, subnetID).Extract()
	if err != nil {
		return "", err
	}

	return subnet.CIDR, nil
}
