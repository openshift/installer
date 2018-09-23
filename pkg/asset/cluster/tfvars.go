package cluster

import (
	"fmt"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types/config"
)

const (
	tfvarsFilename  = "terraform.tfvars"
	tfvarsAssetName = "Terraform Variables"
)

// TerraformVariables depends on InstallConfig and
// Ignition to generate the terrafor.tfvars.
type TerraformVariables struct {
	// The Assets that this tfvars file depends.
	installConfig     asset.Asset
	bootstrapIgnition asset.Asset
	masterIgnition    asset.Asset
	workerIgnition    asset.Asset
}

var _ asset.Asset = (*TerraformVariables)(nil)

// Name returns the human-friendly name of the asset.
func (t *TerraformVariables) Name() string {
	return tfvarsAssetName
}

// Dependencies returns the dependency of the TerraformVariable
func (t *TerraformVariables) Dependencies() []asset.Asset {
	return []asset.Asset{t.installConfig, t.bootstrapIgnition, t.masterIgnition, t.workerIgnition}
}

// Generate generates the terraform.tfvars file.
func (t *TerraformVariables) Generate(parents map[asset.Asset]*asset.State) (*asset.State, error) {
	installCfg, err := installconfig.GetInstallConfig(t.installConfig, parents)
	if err != nil {
		return nil, fmt.Errorf("failed to get install config state in the parent asset states")
	}

	contents := map[asset.Asset][]string{}

	for _, ign := range []asset.Asset{
		t.bootstrapIgnition,
		t.masterIgnition,
		t.workerIgnition,
	} {
		state, ok := parents[ign]
		if !ok {
			return nil, fmt.Errorf("failed to get the ignition state for %v in the parent asset states", ign)
		}

		for _, content := range state.Contents {
			contents[ign] = append(contents[ign], string(content.Data))
		}
	}

	cluster, err := config.ConvertInstallConfigToTFVars(installCfg, contents[t.bootstrapIgnition][0], contents[t.masterIgnition], contents[t.workerIgnition][0])
	if err != nil {
		return nil, err
	}

	if cluster.Platform == config.PlatformLibvirt {
		if err := cluster.Libvirt.UseCachedImage(); err != nil {
			return nil, err
		}
	}

	data, err := cluster.TFVars()
	if err != nil {
		return nil, err
	}

	return &asset.State{
		Contents: []asset.Content{
			{
				Name: tfvarsFilename,
				Data: []byte(data),
			},
		},
	}, nil
}
