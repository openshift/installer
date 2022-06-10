package mirror

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/openshift/installer/pkg/asset"
	"github.com/pkg/errors"
)

var (
	caBundleFilename = filepath.Join(mirrorConfigDir, "ca-bundle.crt")
)

// CaBundle generates the cetificate file for disconnected mirrors.
type CaBundle struct {
	asset.DefaultFileWriter
}

var _ asset.WritableAsset = (*CaBundle)(nil)

// Name returns a human friendly name for the asset.
func (*CaBundle) Name() string {
	return "Mirror Registries Certificate File"
}

// Dependencies returns all of the dependencies directly needed to generate
// the asset.
func (*CaBundle) Dependencies() []asset.Asset {
	return []asset.Asset{
		// &installconfig.InstallConfig{},
	}
}

// Generate generates the Mirror Registries certificate file from install-config.
func (i *CaBundle) Generate(dependencies asset.Parents) error {

	// installConfig := &installconfig.InstallConfig{}
	// dependencies.Get(installConfig)

	// if installConfig.Config.AdditionalTrustBundle == "" {
	//        return nil
	// }
	// data, err := parseCertificates(installConfig.Config.AdditionalTrustBundle)

	//if err != nil {
	//        return err
	//}

	// i.File = &asset.File{
	//      Filename: caBundleFilename,
	//      Data:     data,
	// }

	return nil
}

// Load returns the Mirror Registries certificate file from the disk.
func (i *CaBundle) Load(f asset.FileFetcher) (bool, error) {

	file, err := f.FetchByName(caBundleFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.Wrap(err, fmt.Sprintf("failed to load %s file", caBundleFilename))
	}

	i.File = file
	return true, nil
}
