package validation

import (
	"github.com/gophercloud/gophercloud/openstack/common/extensions"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/regions"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	netext "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/utils/openstack/clientconfig"
)

type realValidValuesFetcher struct{}

// NewValidValuesFetcher returns a new ValidValuesFetcher.
func NewValidValuesFetcher() ValidValuesFetcher {
	return realValidValuesFetcher{}
}

// GetCloudNames gets the valid cloud names. These are read from clouds.yaml.
func (f realValidValuesFetcher) GetCloudNames() ([]string, error) {
	clouds, err := clientconfig.LoadCloudsYAML()
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

// GetRegionNames gets the valid region names.
func (f realValidValuesFetcher) GetRegionNames(cloud string) ([]string, error) {
	opts := &clientconfig.ClientOpts{
		Cloud: cloud,
	}

	conn, err := clientconfig.NewServiceClient("identity", opts)
	if err != nil {
		return nil, err
	}

	listOpts := regions.ListOpts{}
	allPages, err := regions.List(conn, listOpts).AllPages()
	if err != nil {
		return nil, err
	}

	allRegions, err := regions.ExtractRegions(allPages)
	if err != nil {
		return nil, err
	}

	regionNames := make([]string, len(allRegions))
	for x, region := range allRegions {
		regionNames[x] = region.ID
	}

	return regionNames, nil
}

// GetImageNames gets the valid image names.
func (f realValidValuesFetcher) GetImageNames(cloud string) ([]string, error) {
	opts := &clientconfig.ClientOpts{
		Cloud: cloud,
	}

	conn, err := clientconfig.NewServiceClient("image", opts)
	if err != nil {
		return nil, err
	}

	listOpts := images.ListOpts{}
	allPages, err := images.List(conn, listOpts).AllPages()
	if err != nil {
		return nil, err
	}

	allImages, err := images.ExtractImages(allPages)
	if err != nil {
		return nil, err
	}

	imageNames := make([]string, len(allImages))
	for x, image := range allImages {
		imageNames[x] = image.Name
	}

	return imageNames, nil
}

// GetNetworkNames gets the valid network names.
func (f realValidValuesFetcher) GetNetworkNames(cloud string) ([]string, error) {
	opts := &clientconfig.ClientOpts{
		Cloud: cloud,
	}

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
	opts := &clientconfig.ClientOpts{
		Cloud: cloud,
	}

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

	flavorNames := make([]string, len(allFlavors))
	for i, flavor := range allFlavors {
		flavorNames[i] = flavor.Name
	}

	return flavorNames, nil
}

func (f realValidValuesFetcher) GetNetworkExtensionsAliases(cloud string) ([]string, error) {
	opts := &clientconfig.ClientOpts{
		Cloud: cloud,
	}

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

	return extAliases, err
}
