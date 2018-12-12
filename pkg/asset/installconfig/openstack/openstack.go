// Package openstack collects OpenStack-specific configuration.
package openstack

import (
	"sort"
	"strings"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/regions"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/pkg/errors"
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/types/openstack"
)

const (
	defaultVPCCIDR = "10.0.0.0/16"
)

// Read the valid cloud names from the clouds.yaml
func getCloudNames() ([]string, error) {
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
	// Sort cloudNames so we can use sort.SearchStrings
	sort.Strings(cloudNames)
	return cloudNames, nil
}

func getRegionNames(cloud string) ([]string, error) {
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

	sort.Strings(regionNames)
	return regionNames, nil
}

func getImageNames(cloud string) ([]string, error) {
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

	sort.Strings(imageNames)
	return imageNames, nil
}

func getNetworkNames(cloud string) ([]string, error) {
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

	sort.Strings(networkNames)
	return networkNames, nil
}

// Platform collects OpenStack-specific configuration.
func Platform() (*openstack.Platform, error) {
	cloudNames, err := getCloudNames()
	if err != nil {
		return nil, err
	}
	var cloud string
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Cloud",
				Help:    "The OpenStack cloud name from clouds.yaml.",
				Options: cloudNames,
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				value := ans.(string)
				i := sort.SearchStrings(cloudNames, value)
				if i == len(cloudNames) || cloudNames[i] != value {
					return errors.Errorf("invalid cloud name %q, should be one of %+v", value, strings.Join(cloudNames, ", "))
				}
				return nil
			}),
		},
	}, &cloud)
	if err != nil {
		return nil, err
	}

	regionNames, err := getRegionNames(cloud)
	if err != nil {
		return nil, err
	}
	var region string
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Region",
				Help:    "The OpenStack region to be used for installation.",
				Default: "regionOne",
				Options: regionNames,
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				value := ans.(string)
				i := sort.SearchStrings(regionNames, value)
				if i == len(regionNames) || regionNames[i] != value {
					return errors.Errorf("invalid region name %q, should be one of %+v", value, strings.Join(regionNames, ", "))
				}
				return nil
			}),
		},
	}, &region)
	if err != nil {
		return nil, err
	}

	imageNames, err := getImageNames(cloud)
	if err != nil {
		return nil, err
	}
	var image string
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Image",
				Help:    "The OpenStack image name to be used for installation.",
				Default: "rhcos",
				Options: imageNames,
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				value := ans.(string)
				i := sort.SearchStrings(imageNames, value)
				if i == len(imageNames) || imageNames[i] != value {
					return errors.Errorf("invalid image name %q, should be one of %+v", value, strings.Join(imageNames, ", "))
				}
				return nil
			}),
		},
	}, &image)
	if err != nil {
		return nil, err
	}

	networkNames, err := getNetworkNames(cloud)
	if err != nil {
		return nil, err
	}
	var extNet string
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "ExternalNetwork",
				Help:    "The OpenStack external network name to be used for installation.",
				Options: networkNames,
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				value := ans.(string)
				i := sort.SearchStrings(networkNames, value)
				if i == len(networkNames) || networkNames[i] != value {
					return errors.Errorf("invalid network name %q, should be one of %+v", value, strings.Join(networkNames, ", "))
				}
				return nil
			}),
		},
	}, &extNet)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to Marshal %s platform", openstack.Name)
	}

	return &openstack.Platform{
		NetworkCIDRBlock: defaultVPCCIDR,
		Region:           region,
		BaseImage:        image,
		Cloud:            cloud,
		ExternalNetwork:  extNet,
	}, nil
}
