// Package openstack collects OpenStack-specific configuration.
package openstack

import (
	"github.com/pkg/errors"
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/types/openstack"
)

const (
	defaultVPCCIDR = "10.0.0.0/16"
)

// Platform collects OpenStack-specific configuration.
func Platform() (*openstack.Platform, error) {
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

	var cloudConfig *clientconfig.Cloud
	cloud, err := asset.GenerateUserProvidedAsset(
		"OpenStack Cloud",
		&survey.Question{
			Prompt: &survey.Input{
				Message: "Cloud",
				Help:    "The OpenStack cloud name from clouds.yaml.",
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				clientOpts := new(clientconfig.ClientOpts)
				clientOpts.Cloud = ans.(string)
				cloudConfig, err = clientconfig.GetCloudFromYAML(clientOpts)
				return err
			}),
		},
		"OPENSHIFT_INSTALL_OPENSTACK_CLOUD",
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
		CloudConfig:      cloudConfig,
		ExternalNetwork:  extNet,
	}, nil
}
