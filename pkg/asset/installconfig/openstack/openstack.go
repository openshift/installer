// Package openstack collects OpenStack-specific configuration.
package openstack

import (
	"sort"
	"strings"

	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/pkg/errors"
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/asset"
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

// Platform collects OpenStack-specific configuration.
func Platform() (*openstack.Platform, error) {
	cloudNames, err := getCloudNames()
	if err != nil {
		return nil, err
	}
	cloud, err := asset.GenerateUserProvidedAsset(
		"OpenStack Cloud",
		&survey.Question{
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
		"OPENSHIFT_INSTALL_OPENSTACK_CLOUD",
	)
	if err != nil {
		return nil, err
	}

	region, err := asset.GenerateUserProvidedAsset(
		"OpenStack Region",
		&survey.Question{
			Prompt: &survey.Input{
				Message: "Region",
				Help:    "The OpenStack region to be used for installation.",
				Default: "regionOne",
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				//value := ans.(string)
				//FIXME(shardy) add some validation here
				return nil
			}),
		},
		"OPENSHIFT_INSTALL_OPENSTACK_REGION",
	)
	if err != nil {
		return nil, err
	}

	image, err := asset.GenerateUserProvidedAsset(
		"OpenStack Image",
		&survey.Question{
			Prompt: &survey.Input{
				Message: "Image",
				Help:    "The OpenStack image to be used for installation.",
				Default: "rhcos",
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				//value := ans.(string)
				//FIXME(shardy) add some validation here
				return nil
			}),
		},
		"OPENSHIFT_INSTALL_OPENSTACK_IMAGE",
	)
	if err != nil {
		return nil, err
	}

	extNet, err := asset.GenerateUserProvidedAsset(
		"OpenStack External Network",
		&survey.Question{
			Prompt: &survey.Input{
				Message: "ExternalNetwork",
				Help:    "The OpenStack external network to be used for installation.",
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				//value := ans.(string)
				//FIXME(shadower) add some validation here
				return nil
			}),
		},
		"OPENSHIFT_INSTALL_OPENSTACK_EXTERNAL_NETWORK",
	)
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
