package installconfig

import (
	"sort"

	"github.com/pkg/errors"
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/types"
)

// Networking is an asset that queries the user for the networking on which to install
// the cluster.
type networking struct {
	types.Networking
}

var _ asset.Asset = (*networking)(nil)

// Dependencies returns no dependencies.
func (a *networking) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate queries for input from the user.
func (a *networking) Generate(asset.Parents) error {
	networkplugin, err := a.queryUserForNetworkPlugin()
	if err != nil {
		return err
	}
	a.Networking.NetworkType = networkplugin
	return nil
}

// Name returns the human-friendly name of the asset.
func (a *networking) Name() string {
	return "Network Plugin Configuration"
}

func (a *networking) queryUserForNetworkPlugin() (networking string, err error) {
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Network Plugin Configuration",
				Options: types.NetworkPluginNames,
				Default: "OpenShiftSDN",
				Help:    "The network plugin on which the cluster will run.  For a full list of networking plugins, including those not supported by this wizard, see https://github.com/openshift/installer",
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				choice := ans.(string)
				i := sort.SearchStrings(types.NetworkPluginNames, choice)
				if i == len(types.NetworkPluginNames) || types.NetworkPluginNames[i] != choice {
					return errors.Errorf("invalid network plugin %q", choice)
				}
				return nil
			}),
		},
	}, &networking)
	return
}

func (a *networking) CurrentName() string {
	return a.Networking.NetworkType
}
