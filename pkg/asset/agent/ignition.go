package agent

import (
	"path/filepath"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/openshift/installer/pkg/asset"
	am "github.com/openshift/installer/pkg/asset/agent/manifests"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
)

const manifestPathInIso = "/etc/assisted/manifests"

// Ignition is an asset that generates the agent installer ignition file.
type Ignition struct {
	bootstrap.Common
}

// Name returns the human-friendly name of the asset.
func (a *Ignition) Name() string {
	return "Agent Installer Ignition"
}

// Dependencies returns the assets on which the Ignition asset depends.
func (a *Ignition) Dependencies() []asset.Asset {
	return []asset.Asset{
		&am.AgentManifests{},
		// &am.AgentClusterInstall{},
	}
}

// Generate generates the agent installer ignition.
func (a *Ignition) Generate(dependencies asset.Parents) error {
	agentManifests := &am.AgentManifests{}
	// aciAsset := &am.AgentClusterInstall{}
	dependencies.Get(agentManifests)

	a.Config = &igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
		},
	}

	// a.Config.Passwd.Users = append(
	// 	a.Config.Passwd.Users,
	// 	igntypes.PasswdUser{Name: "core", SSHAuthorizedKeys: []igntypes.SSHAuthorizedKey{
	// 		igntypes.SSHAuthorizedKey(aciAsset.Config.Spec.SSHPublicKey),
	// 	}},
	// )

	for _, file := range agentManifests.FileList {
		a.Config.Storage.Files = bootstrap.ReplaceOrAppend(a.Config.Storage.Files,
			ignition.FileFromBytes(filepath.Join(manifestPathInIso, filepath.Base(file.Filename)),
				"root", 0600, file.Data))
	}

	// TODO:
	// Write pull-secret to /root/.docker/config.json
	// Add NMState config
	// Add data/agent/files using Common.addStorageFiles
	// Add data/agent/systemd using Common.addSystemdUnits
	// Template for serviceBaseUrl, pullSecret, apiVIP?, controlPlaneAgents, workerAgents,
	// nodeZeroIP, etc..

	return nil
}
