// Package openstack collects OpenStack-specific configuration.
package openstack

import (
	"sort"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/types/openstack"
)

const (
	noExtNet = "<none>"
)

// Platform collects OpenStack-specific configuration.
func Platform() (*openstack.Platform, error) {
	cloudNames, err := getCloudNames()
	if err != nil {
		return nil, err
	}
	// Sort cloudNames so we can use sort.SearchStrings
	sort.Strings(cloudNames)
	var cloud string
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Cloud",
				Help:    "The OpenStack cloud name from clouds.yaml.",
				Options: cloudNames,
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				value := ans.(core.OptionAnswer).Value
				i := sort.SearchStrings(cloudNames, value)
				if i == len(cloudNames) || cloudNames[i] != value {
					return errors.Errorf("invalid cloud name %q, should be one of %+v", value, strings.Join(cloudNames, ", "))
				}
				return nil
			}),
		},
	}, &cloud)
	if err != nil {
		return nil, errors.Wrap(err, "failed UserInput")
	}

	networkNames, err := getExternalNetworkNames(cloud)
	if err != nil {
		return nil, err
	}
	networkNames = append(networkNames, noExtNet)
	sort.Strings(networkNames)
	var extNet string
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "ExternalNetwork",
				Help:    "The OpenStack external network name to be used for installation.",
				Options: networkNames,
				Default: noExtNet,
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				value := ans.(core.OptionAnswer).Value
				i := sort.SearchStrings(networkNames, value)
				if i == len(networkNames) || networkNames[i] != value {
					return errors.Errorf("invalid network name %q, should be one of %+v", value, strings.Join(networkNames, ", "))
				}
				return nil
			}),
		},
	}, &extNet)
	if extNet == noExtNet {
		extNet = ""
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed UserInput")
	}

	var apiFloatingIP string
	if extNet != "" {
		floatingIPs, err := getFloatingIPs(cloud, extNet)
		if err != nil {
			return nil, err
		}
		sort.Sort(floatingIPs)
		err = survey.Ask([]*survey.Question{
			{
				Prompt: &survey.Select{
					Message:     "APIFloatingIPAddress",
					Help:        "The Floating IP address used for external access to the OpenShift API.",
					Options:     floatingIPs.Names(),
					Description: func(_ string, index int) string { return floatingIPs.Description(index) },
				},
				Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
					if value := ans.(core.OptionAnswer).Value; !floatingIPs.Contains(value) {
						return errors.Errorf("invalid floating IP %q, should be one of %+v", value, strings.Join(floatingIPs.Names(), ", "))
					}
					return nil
				}),
			},
		}, &apiFloatingIP)
		if err != nil {
			return nil, errors.Wrap(err, "failed UserInput")
		}
	}

	flavorNames, err := getFlavorNames(cloud)
	if err != nil {
		return nil, err
	}
	sort.Strings(flavorNames)
	var flavor string
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "FlavorName",
				Help:    "The OpenStack flavor to use for control-plane and compute nodes. A flavor with at least 16 GB RAM is recommended.",
				Options: flavorNames,
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				value := ans.(core.OptionAnswer).Value
				i := sort.SearchStrings(flavorNames, value)
				if i == len(flavorNames) || flavorNames[i] != value {
					return errors.Errorf("invalid flavor name %q, should be one of %+v", value, strings.Join(flavorNames, ", "))
				}
				return nil
			}),
		},
	}, &flavor)
	if err != nil {
		return nil, errors.Wrap(err, "failed UserInput")
	}

	return &openstack.Platform{
		APIFloatingIP:   apiFloatingIP,
		Cloud:           cloud,
		ExternalNetwork: extNet,
		DefaultMachinePlatform: &openstack.MachinePool{
			FlavorName: flavor,
		},
	}, nil
}
