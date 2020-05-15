// Package openstack collects OpenStack-specific configuration.
package openstack

import (
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/types/openstack"
	openstackvalidation "github.com/openshift/installer/pkg/types/openstack/validation"
)

// Platform collects OpenStack-specific configuration.
func Platform() (*openstack.Platform, error) {
	validValuesFetcher := openstackvalidation.NewValidValuesFetcher()

	cloudNames, err := validValuesFetcher.GetCloudNames()
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

	networkNames, err := validValuesFetcher.GetNetworkNames(cloud)
	if err != nil {
		return nil, err
	}
	sort.Strings(networkNames)
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
		return nil, err
	}

	floatingIPNames, err := validValuesFetcher.GetFloatingIPNames(cloud, extNet)
	if err != nil {
		return nil, err
	}
	sort.Strings(floatingIPNames)
	var lbFloatingIP string
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "APIFloatingIPAddress",
				Help:    "The Floating IP address used for external access to the OpenShift API.",
				Options: floatingIPNames,
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				value := ans.(string)
				i := sort.SearchStrings(floatingIPNames, value)
				if i == len(floatingIPNames) || floatingIPNames[i] != value {
					return errors.Errorf("invalid floating IP %q, should be one of %+v", value, strings.Join(floatingIPNames, ", "))
				}
				return nil
			}),
		},
	}, &lbFloatingIP)
	if err != nil {
		return nil, err
	}

	flavorNames, err := validValuesFetcher.GetFlavorNames(cloud)
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
				value := ans.(string)
				i := sort.SearchStrings(flavorNames, value)
				if i == len(flavorNames) || flavorNames[i] != value {
					return errors.Errorf("invalid flavor name %q, should be one of %+v", value, strings.Join(flavorNames, ", "))
				}
				return nil
			}),
		},
	}, &flavor)
	if err != nil {
		return nil, err
	}

	trunkSupport := "0"
	var i int
	netExts, err := validValuesFetcher.GetNetworkExtensionsAliases(cloud)
	if err != nil {
		logrus.Warning("Could not retrieve networking extension aliases. Assuming trunk ports are not supported.")
	} else {
		sort.Strings(netExts)
		i = sort.SearchStrings(netExts, "trunk")
		if i != len(netExts) && netExts[i] == "trunk" {
			trunkSupport = "1"
		}
	}

	octaviaSupport := "0"
	serviceCatalog, err := validValuesFetcher.GetServiceCatalog(cloud)
	if err != nil {
		logrus.Warning("Could not retrieve service catalog. Assuming there is no Octavia load balancer service available.")
	} else {
		sort.Strings(serviceCatalog)
		i = sort.SearchStrings(serviceCatalog, "octavia")
		if i != len(serviceCatalog) && serviceCatalog[i] == "octavia" {
			octaviaSupport = "1"
		}
	}

	return &openstack.Platform{
		Cloud:           cloud,
		ExternalNetwork: extNet,
		FlavorName:      flavor,
		LbFloatingIP:    lbFloatingIP,
		TrunkSupport:    trunkSupport,
		OctaviaSupport:  octaviaSupport,
	}, nil
}
