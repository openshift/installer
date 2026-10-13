package bootstrap

import (
	"context"

	"github.com/openshift/installer/pkg/asset"
)

const (
	bootstrapIgnFilename = "bootstrap.ign"
)

// Bootstrap is an asset that generates the ignition config for bootstrap nodes.
type Bootstrap struct {
	Common
}

var _ asset.WritableAsset = (*Bootstrap)(nil)

// Generate generates the ignition config for the Bootstrap asset.
func (a *Bootstrap) Generate(_ context.Context, dependencies asset.Parents) error {
	templateData := a.getTemplateData(dependencies, false)
	if err := a.generateConfig(dependencies, templateData); err != nil {
		return err
	}

	if err := a.generateFile(bootstrapIgnFilename); err != nil {
		return err
	}
	return nil
}

// Name returns the human-friendly name of the asset.
func (a *Bootstrap) Name() string {
	return "Bootstrap Ignition Config"
}

// Load returns the bootstrap ignition from disk.
func (a *Bootstrap) Load(f asset.FileFetcher) (found bool, err error) {
	return a.load(f, bootstrapIgnFilename)
}
