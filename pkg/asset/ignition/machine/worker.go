package machine

import (
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
)

// worker is an asset that generates the ignition config for worker nodes.
type worker struct {
	installConfig asset.Asset
	rootCA        asset.Asset
}

var _ asset.Asset = (*worker)(nil)

// newWorker creates a new worker asset.
func newWorker(
	installConfigStock installconfig.Stock,
	tlsStock tls.Stock,
) *worker {
	return &worker{
		installConfig: installConfigStock.InstallConfig(),
		rootCA:        tlsStock.RootCA(),
	}
}

// Dependencies returns the assets on which the worker asset depends.
func (a *worker) Dependencies() []asset.Asset {
	return []asset.Asset{
		a.installConfig,
		a.rootCA,
	}
}

// Generate generates the ignition config for the worker asset.
func (a *worker) Generate(dependencies map[asset.Asset]*asset.State) (*asset.State, error) {
	installConfig, err := installconfig.GetInstallConfig(a.installConfig, dependencies)
	if err != nil {
		return nil, err
	}

	return &asset.State{
		Contents: []asset.Content{{
			Name: "worker.ign",
			Data: pointerIgnitionConfig(installConfig, dependencies[a.rootCA].Contents[tls.CertIndex].Data, "worker", ""),
		}},
	}, nil
}

// Name returns the human-friendly name of the asset.
func (a *worker) Name() string {
	return "Worker Ignition Config"
}
