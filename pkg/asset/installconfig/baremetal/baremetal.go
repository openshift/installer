// Package baremetal collects bare metal specific configuration.
package baremetal

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types/baremetal"
	baremetaldefaults "github.com/openshift/installer/pkg/types/baremetal/defaults"
)

// Platform collects bare metal specific configuration.
func Platform() (*baremetal.Platform, error) {
	var provisioningNetworkCIDR, externalBridge, provisioningBridge, provisioningNetwork, provisioningNetworkInterface string
	var parsedCIDR *ipnet.IPNet
	var hosts []*baremetal.Host

	survey.AskOne(&survey.Select{
		Message: "Provisioning Network",
		Help:    "Select whether the provisioning network will be managed, unmanaged, or disabled. In managed mode, the cluster deploys DHCP and TFTP services for PXE provisioning.",
		Options: []string{"Managed", "Unmanaged", "Disabled"},
		Default: "Managed",
	}, &provisioningNetwork, nil)

	if provisioningNetwork != string(baremetal.DisabledProvisioningNetwork) {
		if err := survey.Ask([]*survey.Question{
			{
				Prompt: &survey.Input{
					Message: "Provisioning Network CIDR",
					Help:    "The network used for provisioning.",
					Default: "172.22.0.0/24",
				},
				Validate: survey.ComposeValidators(survey.Required, ipNetValidator),
			},
		}, &provisioningNetworkCIDR); err != nil {
			return nil, err
		}
		provNetCIDR, err := ipnet.ParseCIDR(provisioningNetworkCIDR)
		if err != nil {
			return nil, err
		}
		parsedCIDR = provNetCIDR

		if err := survey.Ask([]*survey.Question{
			{
				Prompt: &survey.Input{
					Message: "Provisioning bridge",
					Help:    "Provisioning bridge is used to provision machines by the bootstrap virtual machine.",
					Default: baremetaldefaults.ProvisioningBridge,
				},
			},
		}, &provisioningBridge); err != nil {
			return nil, err
		}

		if err := survey.Ask([]*survey.Question{
			{
				Prompt: &survey.Input{
					Message: "Provisioning Network Interface",
					Help:    "The name of the network interface on a control plane host connected to the provisioning network.",
				},
			},
		}, &provisioningNetworkInterface); err != nil {
			return nil, err
		}
	}

	if err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "External bridge",
				Help:    "External bridge is used for external communication by the bootstrap virtual machine.",
				Default: baremetaldefaults.ExternalBridge,
			},
		},
	}, &externalBridge); err != nil {
		return nil, err
	}

	// Keep prompting for hosts
	for {
		var hostRole string
		survey.AskOne(&survey.Select{
			Message: "Add a Host:",
			Options: []string{"control plane", "worker"},
		}, &hostRole, nil)

		var host *baremetal.Host
		var err error
		host, err = Host()
		// Check for keyboard interrupt or else we'll loop forever
		if errors.Is(err, terminal.InterruptErr) {
			fmt.Println("interrupted - hosts were not added")
			break
		} else if err != nil {
			fmt.Printf("invalid host - please try again")
			continue
		}
		if hostRole == "control plane" {
			host.Role = "master"
		} else {
			host.Role = hostRole
		}
		hosts = append(hosts, host)

		more := false
		survey.AskOne(&survey.Confirm{
			Message: "Add another host?",
		}, &more, nil)
		if !more {
			break
		}
	}

	return &baremetal.Platform{
		ExternalBridge:               externalBridge,
		ProvisioningBridge:           provisioningBridge,
		ProvisioningNetworkCIDR:      parsedCIDR,
		ProvisioningNetworkInterface: provisioningNetworkInterface,
		Hosts:                        hosts,
	}, nil
}

// ipNetValidator validates for a valid IP
func ipNetValidator(ans interface{}) error {
	_, err := ipnet.ParseCIDR(ans.(string))
	return err
}
