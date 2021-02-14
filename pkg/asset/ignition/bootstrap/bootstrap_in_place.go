package bootstrap

import (
	"encoding/json"
	"os"

	igntypes "github.com/coreos/ignition/v2/config/v3_1/types"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap/baremetal"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/manifests"
	"github.com/openshift/installer/pkg/asset/releaseimage"
	"github.com/openshift/installer/pkg/asset/rhcos"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/types"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

const (
	singleNodeBootstrapInPlaceIgnFilename = "bootstrap-in-place-for-live-iso.ign"
)

// SingleNodeBootstrapInPlace is an asset that generates the ignition config for single node OpenShift.
type SingleNodeBootstrapInPlace struct {
	Bootstrap *Bootstrap
	File      *asset.File
	Config    *igntypes.Config
}

var _ asset.Asset = (*SingleNodeBootstrapInPlace)(nil)

// Dependencies returns no dependencies.
func (a *SingleNodeBootstrapInPlace) Dependencies() []asset.Asset {
	a.Bootstrap.bootstrapInPlace = true
	return a.Bootstrap.Dependencies()
}

// Name returns the human-friendly name of the asset.
func (a *SingleNodeBootstrapInPlace) Name() string {
	return "Single node bootstrap In Place Ignition Config"
}

// Files returns the password file.
func (a *SingleNodeBootstrapInPlace) Files() []*asset.File {
	if a.File != nil {
		return []*asset.File{a.File}
	}
	return []*asset.File{}
}

// Generate generates the ignition config for the Bootstrap asset.
func (a *SingleNodeBootstrapInPlace) Generate(dependencies asset.Parents) error {
	a.Bootstrap.bootstrapInPlace = true
	installConfig := &installconfig.InstallConfig{}
	proxy := &manifests.Proxy{}
	releaseImage := &releaseimage.Image{}
	rhcosImage := new(rhcos.Image)
	bootstrapSSHKeyPair := &tls.BootstrapSSHKeyPair{}
	ironicCreds := &baremetal.IronicCreds{}
	dependencies.Get(installConfig, proxy, releaseImage, rhcosImage, bootstrapSSHKeyPair, ironicCreds)
	if err := verifyBootstrapInPlace(installConfig.Config); err != nil {
		return err
	}
	templateData, err := a.Bootstrap.getTemplateData(installConfig.Config, releaseImage.PullSpec, installConfig.Config.ImageContentSources, proxy.Config, rhcosImage, ironicCreds)

	if err = a.Bootstrap.Generate(dependencies); err != nil {
		return err
	}
	err = a.Bootstrap.addStorageFiles("/", "bootstrap/bootstrap-in-place/files", templateData)
	if err != nil {
		return err
	}
	err = a.Bootstrap.addSystemdUnits("bootstrap/bootstrap-in-place/systemd/units", templateData)
	if err != nil {
		return err
	}

	a.Config = a.Bootstrap.Config
	data, err := ignition.Marshal(a.Config)
	if err != nil {
		return errors.Wrap(err, "failed to Marshal Ignition config")
	}

	a.File = &asset.File{
		Filename: singleNodeBootstrapInPlaceIgnFilename,
		Data:     data,
	}
	return nil
}

// Load returns the bootstrap-in-place ignition from disk.
func (a *SingleNodeBootstrapInPlace) Load(f asset.FileFetcher) (found bool, err error) {
	file, err := f.FetchByName(singleNodeBootstrapInPlaceIgnFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	config := &igntypes.Config{}
	if err := json.Unmarshal(file.Data, config); err != nil {
		return false, errors.Wrapf(err, "failed to unmarshal %s", bootstrapIgnFilename)
	}

	a.File, a.Config = file, config
	bootstrapFile := &asset.File{
		Filename: bootstrapIgnFilename,
		Data:     a.File.Data,
	}
	a.Bootstrap = &Bootstrap{Config: config, File: bootstrapFile}
	a.Bootstrap.bootstrapInPlace = true
	warnIfCertificatesExpired(a.Config)
	return true, nil
}

// verifyBootstrapInPlace validate the number of control plane replica is one and that installation disk is set
func verifyBootstrapInPlace(installConfig *types.InstallConfig) error {
	errorList := field.ErrorList{}
	if installConfig.ControlPlane.Replicas == nil {
		errorList = append(errorList, field.Invalid(field.NewPath("controlPlane", "replicas"), installConfig.ControlPlane.Replicas,
			"bootstrap in place requires ControlPlane.Replicas configuration"))
	}
	if *installConfig.ControlPlane.Replicas != 1 {
		errorList = append(errorList, field.Invalid(field.NewPath("controlPlane", "replicas"), installConfig.ControlPlane.Replicas,
			"bootstrap in place requires a single ControlPlane replica"))
	}
	if installConfig.BootstrapInPlace == nil {
		errorList = append(errorList, field.Required(field.NewPath("bootstrapInPlace"), "bootstrapInPlace is required when creating a single node bootstrap-in-place ignition"))
	} else if installConfig.BootstrapInPlace.InstallationDisk == "" {
		errorList = append(errorList, field.Required(field.NewPath("bootstrapInPlace", "installationDisk"),
			"installationDisk must be set the target disk drive for the installation"))
	}
	return errorList.ToAggregate()
}
