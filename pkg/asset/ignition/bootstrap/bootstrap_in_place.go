package bootstrap

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
)

const (
	singleNodeBootstrapInPlaceIgnFilename = "bootstrap-in-place-for-live-iso.ign"
)

var (
	bootstrapInPlaceEnabledServices = []string{
		"install-to-disk.service",
	}
)

// SingleNodeBootstrapInPlace is an asset that generates the ignition config for single node OpenShift.
type SingleNodeBootstrapInPlace struct {
	Common
}

var _ asset.Asset = (*SingleNodeBootstrapInPlace)(nil)

// Name returns the human-friendly name of the asset.
func (a *SingleNodeBootstrapInPlace) Name() string {
	return "Single node bootstrap In Place Ignition Config"
}

// Generate generates the ignition config for the Bootstrap asset.
func (a *SingleNodeBootstrapInPlace) Generate(dependencies asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	dependencies.Get(installConfig)
	if err := verifyBootstrapInPlace(installConfig.Config); err != nil {
		return err
	}
	templateData := a.getTemplateData(dependencies, true)
	if err := a.generateConfig(dependencies, templateData); err != nil {
		return err
	}
	if err := a.addStorageFiles("/", "bootstrap/bootstrap-in-place/files", templateData); err != nil {
		return err
	}
	if err := a.addSystemdUnits("bootstrap/bootstrap-in-place/systemd/units", templateData, bootstrapInPlaceEnabledServices); err != nil {
		return err
	}
	if err := a.Common.generateFile(singleNodeBootstrapInPlaceIgnFilename); err != nil {
		return err
	}
	return nil
}

// Load returns the bootstrap-in-place ignition from disk.
func (a *SingleNodeBootstrapInPlace) Load(f asset.FileFetcher) (found bool, err error) {
	return a.load(f, singleNodeBootstrapInPlaceIgnFilename)
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
