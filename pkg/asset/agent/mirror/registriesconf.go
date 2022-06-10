package mirror

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/asset"
)

var (
	registriesConfFilename = filepath.Join(mirrorConfigDir, "registries.conf")
)

// RegistriesConf generates the registries.conf file.
type RegistriesConf struct {
	asset.DefaultFileWriter
}

var _ asset.WritableAsset = (*RegistriesConf)(nil)

// Name returns a human friendly name for the asset.
func (*RegistriesConf) Name() string {
	return "Mirror Registries Config"
}

// Dependencies returns all of the dependencies directly needed to generate
// the asset.
func (*RegistriesConf) Dependencies() []asset.Asset {
	return []asset.Asset{
		// &installconfig.InstallConfig{},
	}
}

// Generate generates the registries.conf file from install-config.
func (i *RegistriesConf) Generate(dependencies asset.Parents) error {

	// installConfig := &installconfig.InstallConfig{}

	// registries := []sysregistriesv2.Registry{}
	// for _, group := range bootstrap.MergedMirrorSets(installConfig.Config.ImageContentSources) {
	//        if len(group.Mirrors) == 0 {
	//                continue
	//        }

	//        registry := sysregistriesv2.Registry{}
	//        registry.Endpoint.Location = group.Source
	//        registry.MirrorByDigestOnly = true
	//        for _, mirror := range group.Mirrors {
	//                registry.Mirrors = append(registry.Mirrors, sysregistriesv2.Endpoint{Location: mirror})
	//        }
	//        registries = append(registries, registry)
	// }

	// i.File = &asset.File{
	//      Filename: registriesConfFilename,
	//      Data:     registries,
	// }

	return nil
}

// Load returns RegistriesConf asset from the disk.
func (i *RegistriesConf) Load(f asset.FileFetcher) (bool, error) {

	file, err := f.FetchByName(registriesConfFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.Wrap(err, fmt.Sprintf("failed to load %s file", registriesConfFilename))
	}

	i.File = file
	return true, nil
}
