package ovirt

import (
	"fmt"
	"sort"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	ovirtsdk4 "github.com/ovirt/go-ovirt"
	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/defaults"
	"github.com/openshift/installer/pkg/types/ovirt"
	"github.com/openshift/installer/pkg/types/validation"
	"github.com/openshift/installer/pkg/validate"
)

func askNetwork(c *ovirtsdk4.Connection, p *ovirt.Platform) error {
	var networkName string
	var networkByNames = make(map[string]*ovirtsdk4.Network)
	var networkNames []string
	systemService := c.SystemService()
	networksResponse, err := systemService.ClustersService().ClusterService(p.ClusterID).NetworksService().List().Send()
	if err != nil {
		return err
	}
	networks, ok := networksResponse.Networks()
	if !ok {
		return fmt.Errorf("there are no available networks for cluster %s", p.ClusterID)
	}

	for _, network := range networks.Slice() {
		networkByNames[network.MustName()] = network
		networkNames = append(networkNames, network.MustName())
	}
	if err := survey.AskOne(
		&survey.Select{
			Message: "Network",
			Help:    "The Engine network of the deployed VMs. 'ovirtmgmt' is the default network. It is recommended to use a dedicated network for each OpenShift cluster.",
			Options: networkNames,
		},
		&networkName,
		survey.WithValidator(func(ans interface{}) error {
			choice := ans.(core.OptionAnswer).Value
			sort.Strings(networkNames)
			i := sort.SearchStrings(networkNames, choice)
			if i == len(networkNames) || networkNames[i] != choice {
				return fmt.Errorf("invalid network %s", choice)
			}
			network, ok := networkByNames[choice]
			if !ok {
				return fmt.Errorf("cannot find a network by name %s", choice)
			}
			p.NetworkName = network.MustName()
			return nil
		}),
	); err != nil {
		return errors.Wrap(err, "failed UserInput")
	}
	return nil
}

func askVNICProfileID(c *ovirtsdk4.Connection, p *ovirt.Platform) error {
	var profileID string
	var profilesByNames = make(map[string]*ovirtsdk4.VnicProfile)
	var profileNames []string
	profiles, err := FetchVNICProfileByClusterNetwork(c, p.ClusterID, p.NetworkName)
	if err != nil {
		return err
	}

	for _, profile := range profiles {
		profilesByNames[profile.MustName()] = profile
		profileNames = append(profileNames, profile.MustName())
	}

	if len(profilesByNames) == 1 {
		p.VNICProfileID = profilesByNames[profileNames[0]].MustId()
		return nil
	}

	// we have multiple vnic profile for the selected network
	if err := survey.AskOne(
		&survey.Select{
			Message: "VNIC Profile",
			Help:    "The Engine VNIC profile of the VMs.",
			Options: profileNames,
		},
		&profileID,
		survey.WithValidator(func(ans interface{}) error {
			choice := ans.(core.OptionAnswer).Value
			sort.Strings(profileNames)
			i := sort.SearchStrings(profileNames, choice)
			if i == len(profileNames) || profileNames[i] != choice {
				return fmt.Errorf("invalid VNIC profile %s", choice)
			}
			profile, ok := profilesByNames[choice]
			if !ok {
				return fmt.Errorf("cannot find a VNIC profile id by the name %s", choice)
			}
			p.VNICProfileID = profile.MustId()
			return nil
		}),
	); err != nil {
		return errors.Wrap(err, "failed UserInput")
	}
	return nil
}

func askVIPs(p *ovirt.Platform) error {
	//TODO: Add support to specify multiple VIPs (-> dual-stack)
	var apiVIP, ingressVIP string

	defaultMachineNetwork := &types.Networking{
		MachineNetwork: []types.MachineNetworkEntry{
			{
				CIDR: *defaults.DefaultMachineCIDR,
			},
		},
	}

	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Internal API virtual IP",
				Help:    "This is the virtual IP address that will be used to address the OpenShift control plane. Make sure the IP address is not in use.",
				Default: "",
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				err := validate.IP((ans).(string))
				if err != nil {
					return err
				}
				return validation.ValidateIPinMachineCIDR((ans).(string), defaultMachineNetwork)
			}),
		},
	}, &apiVIP)
	if err != nil {
		return errors.Wrap(err, "failed UserInput")
	}
	p.APIVIPs = []string{apiVIP}

	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Ingress virtual IP",
				Help:    "This is the virtual IP address that will be used to address the OpenShift ingress routers. Make sure the IP address is not in use.",
				Default: "",
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				if apiVIP == (ans.(string)) {
					return fmt.Errorf("%q should not be equal to the Virtual IP address for the API", ans.(string))
				}
				err := validate.IP((ans).(string))
				if err != nil {
					return err
				}
				return validation.ValidateIPinMachineCIDR((ans).(string), defaultMachineNetwork)
			}),
		},
	}, &ingressVIP)
	if err != nil {
		return errors.Wrap(err, "failed UserInput")
	}
	p.IngressVIPs = []string{ingressVIP}

	return nil
}
