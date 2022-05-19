package image

import (
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"

	"github.com/openshift/installer/pkg/agent/imagebuilder"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/manifests"
)

// Ignition is an asset that generates the agent installer ignition file.
type Ignition struct {
	Config *igntypes.Config
}

// Name returns the human-friendly name of the asset.
func (a *Ignition) Name() string {
	return "Agent Installer Ignition"
}

// Dependencies returns the assets on which the Ignition asset depends.
func (a *Ignition) Dependencies() []asset.Asset {
	return []asset.Asset{
		&manifests.AgentManifests{},
	}
}

// Generate generates the agent installer ignition.
func (a *Ignition) Generate(dependencies asset.Parents) error {

	agentManifests := &manifests.AgentManifests{}
	dependencies.Get(agentManifests)

	configBuilder, err := imagebuilder.New(*agentManifests)
	if err != nil {
		return err
	}

	ignition, err := configBuilder.Ignition()
	if err != nil {
		return err
	}

	a.Config = &ignition
	return nil
}
