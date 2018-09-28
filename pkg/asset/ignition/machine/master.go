package machine

import (
	"fmt"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
)

// master is an asset that generates the ignition config for master nodes.
type master struct {
	installConfig asset.Asset
	rootCA        asset.Asset
}

var _ asset.Asset = (*master)(nil)

// newMaster generates a new master asset.
func newMaster(
	installConfigStock installconfig.Stock,
	tlsStock tls.Stock,
) *master {
	return &master{
		installConfig: installConfigStock.InstallConfig(),
		rootCA:        tlsStock.RootCA(),
	}
}

// Dependencies returns the assets on which the master asset depends.
func (a *master) Dependencies() []asset.Asset {
	return []asset.Asset{
		a.installConfig,
		a.rootCA,
	}
}

// Generate generates the ignition config for the master asset.
func (a *master) Generate(dependencies map[asset.Asset]*asset.State) (*asset.State, error) {
	installConfig, err := installconfig.GetInstallConfig(a.installConfig, dependencies)
	if err != nil {
		return nil, err
	}

	state := &asset.State{
		Contents: make([]asset.Content, installConfig.MasterCount()),
	}
	for i := range state.Contents {
		state.Contents[i].Name = fmt.Sprintf("master-%d.ign", i)
		state.Contents[i].Data = pointerIgnitionConfig(installConfig, dependencies[a.rootCA].Contents[tls.CertIndex].Data, "master", fmt.Sprintf("etcd_index=%d", i))
	}

	return state, nil
}

// Name returns the human-friendly name of the asset.
func (a *master) Name() string {
	return "Master Ignition Config(s)"
}
