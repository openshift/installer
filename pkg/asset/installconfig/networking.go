package installconfig

import (
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/kubevirt"
)

type networking struct {
	machineNetwork []types.MachineNetworkEntry
}

var _ asset.Asset = (*networking)(nil)

// Dependencies returns no dependencies.
func (a *networking) Dependencies() []asset.Asset {
	return []asset.Asset{
		&platform{},
	}
}

// Generate queries for the networking from the user.
func (a *networking) Generate(parents asset.Parents) error {
	platform := &platform{}
	parents.Get(platform)

	switch platform.CurrentName() {
	case kubevirt.Name:
		selectedCIDR, err := selectMachineNetworkCIDR()
		if err != nil {
			return err
		}
		CIDR, err := ipnet.ParseCIDR(selectedCIDR)
		if err != nil {
			return err
		}
		a.machineNetwork = []types.MachineNetworkEntry{
			{CIDR: *CIDR},
		}
	}
	return nil
}

func selectMachineNetworkCIDR() (string, error) {
	var selectedCIDR string

	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Machine Network CIDR",
				Help:    "The IP address pool for machines.",
			},
		},
	}, &selectedCIDR)

	return selectedCIDR, err
}

// Name returns the human-friendly name of the asset.
func (a *networking) Name() string {
	return "Networking"
}
