package machine

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
)

const (
	arbiterIgnFilename = "arbiter.ign"
)

// Arbiter is an asset that generates the ignition config for arbiter nodes.
type Arbiter struct {
	Config *igntypes.Config
	File   *asset.File
}

var _ asset.WritableAsset = (*Arbiter)(nil)

// Dependencies returns the assets on which the Arbiter asset depends.
func (a *Arbiter) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
		&tls.RootCA{},
	}
}

// Generate generates the ignition config for the Arbiter asset.
func (a *Arbiter) Generate(_ context.Context, dependencies asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	rootCA := &tls.RootCA{}
	dependencies.Get(installConfig, rootCA)

	// Avoid creating ignition files when not an arbiter deployment.
	if !installConfig.Config.IsArbiterEnabled() {
		return nil
	}

	a.Config = pointerIgnitionConfig(installConfig.Config, rootCA.Cert(), "arbiter")

	data, err := ignition.Marshal(a.Config)
	if err != nil {
		return fmt.Errorf("failed to marshal Ignition config: %w", err)
	}
	a.File = &asset.File{
		Filename: arbiterIgnFilename,
		Data:     data,
	}

	return nil
}

// Name returns the human-friendly name of the asset.
func (a *Arbiter) Name() string {
	return "Arbiter Ignition Config"
}

// Files returns the files generated by the asset.
func (a *Arbiter) Files() []*asset.File {
	if a.File != nil {
		return []*asset.File{a.File}
	}
	return []*asset.File{}
}

// Load returns the arbiter ignitions from disk.
func (a *Arbiter) Load(f asset.FileFetcher) (found bool, err error) {
	file, err := f.FetchByName(arbiterIgnFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	config := &igntypes.Config{}
	if err := json.Unmarshal(file.Data, config); err != nil {
		return false, fmt.Errorf("failed to unmarshal %s: %w", arbiterIgnFilename, err)
	}

	a.File, a.Config = file, config
	return true, nil
}
