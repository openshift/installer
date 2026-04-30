package machine

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/types"
)

// CustomPool is an asset that generates pointer ignition configs for all
// user-defined custom compute pools. Each stub points to /config/<pool-name>
// on the MCS so that MCO serves the correct rendered MachineConfig for the pool.
type CustomPool struct {
	// FilesByPool maps pool name to its pointer ignition file.
	FilesByPool map[string]*asset.File
}

var _ asset.WritableAsset = (*CustomPool)(nil)

// Dependencies returns the assets required to generate the custom pool ignitions.
func (a *CustomPool) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
		&tls.RootCA{},
	}
}

// Generate generates pointer ignition configs for every custom pool in the
// install config. It is a no-op when no custom pools are defined.
func (a *CustomPool) Generate(_ context.Context, dependencies asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	rootCA := &tls.RootCA{}
	dependencies.Get(installConfig, rootCA)

	logrus.Debugf("[custom-pool] CustomPool ignition: scanning %d compute pool(s)", len(installConfig.Config.Compute))

	a.FilesByPool = make(map[string]*asset.File)
	for i := range installConfig.Config.Compute {
		name := installConfig.Config.Compute[i].Name
		if !types.IsCustomPool(name) {
			logrus.Debugf("[custom-pool] CustomPool ignition: pool[%d] name=%q — not custom, skipping", i, name)
			continue
		}

		logrus.Infof("[custom-pool] CustomPool ignition: generating pointer ignition for pool %q (MCS path: /config/%s)", name, name)
		cfg := pointerIgnitionConfig(installConfig.Config, rootCA.Cert(), name)
		data, err := ignition.Marshal(cfg)
		if err != nil {
			return fmt.Errorf("failed to marshal ignition config for custom pool %q: %w", name, err)
		}
		filename := fmt.Sprintf("%s.ign", name)
		logrus.Infof("[custom-pool] CustomPool ignition: wrote %q (%d bytes)", filename, len(data))
		a.FilesByPool[name] = &asset.File{
			Filename: filename,
			Data:     data,
		}
	}

	if len(a.FilesByPool) == 0 {
		logrus.Debugf("[custom-pool] CustomPool ignition: no custom pools found, skipping ignition generation")
	}
	return nil
}

// Name returns the human-friendly name of the asset.
func (a *CustomPool) Name() string {
	return "Custom Pool Ignition Config"
}

// Files returns the ignition files generated for all custom pools.
func (a *CustomPool) Files() []*asset.File {
	files := make([]*asset.File, 0, len(a.FilesByPool))
	for _, f := range a.FilesByPool {
		files = append(files, f)
	}
	return files
}

// Load returns false because ignition filenames depend on pool names which are
// not known until the install config is available. Generate() is always used.
func (a *CustomPool) Load(_ asset.FileFetcher) (bool, error) {
	return false, nil
}
