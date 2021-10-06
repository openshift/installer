package bootstrap

import (
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition"
)

const (
	bootstrapIgnFilename = "bootstrap.ign"
)

// Bootstrap is an asset that generates the ignition config for bootstrap nodes.
type Bootstrap struct {
	Common
}

var _ asset.WritableAsset = (*Bootstrap)(nil)

// Dependencies returns the assets on which the Bootstrap asset depends.
func (a *Bootstrap) Dependencies() []asset.Asset {
	return append(
		a.Common.Dependencies(),
		&CVOIgnore{},
	)
}

// Generate generates the ignition config for the Bootstrap asset.
func (a *Bootstrap) Generate(dependencies asset.Parents) error {
	templateData := a.getTemplateData(dependencies, false)
	if err := a.generateConfig(dependencies, templateData); err != nil {
		return err
	}

	// replace cvo-overrides.yaml so that CVO does not try to create resources for the manifests that the installer
	// is supplying to cluster-bootstrap
	cvoIgnore := &CVOIgnore{}
	dependencies.Get(cvoIgnore)
	for _, file := range ignition.FilesFromAsset(rootDir, "root", 0644, cvoIgnore) {
		a.Config.Storage.Files = replaceOrAppend(a.Config.Storage.Files, file)
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
