package machine

import (
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
	mcfgv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
)

var (
	masterMachineConfigFileName = filepath.Join(directory, "99_openshift-installer-ignition_master.yaml")
)

// MasterIgnitionCustomizations is an asset that checks for any customizations a user might
// have made to the pointer ignition configs before creating the cluster. If customizations
// are made, then the updates are reconciled as a MachineConfig file
type MasterIgnitionCustomizations struct {
	asset.DefaultFileWriter
}

var _ asset.WritableAsset = (*MasterIgnitionCustomizations)(nil)

// Dependencies returns the dependencies for MasterIgnitionCustomizations
func (a *MasterIgnitionCustomizations) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
		&tls.RootCA{},
		&Master{},
	}
}

// Generate queries for input from the user.
func (a *MasterIgnitionCustomizations) Generate(dependencies asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	rootCA := &tls.RootCA{}
	master := &Master{}
	dependencies.Get(installConfig, rootCA, master)

	defaultPointerIgnition := pointerIgnitionConfig(installConfig.Config, rootCA.Cert(), "master")
	savedPointerIgnition := master.Config

	savedPointerIgnitionJSON, err := ignition.Marshal(savedPointerIgnition)
	if err != nil {
		return errors.Wrap(err, "failed to Marshal savedPointerIgnition")
	}
	defaultPointerIgnitionJSON, err := ignition.Marshal(defaultPointerIgnition)
	if err != nil {
		return errors.Wrap(err, "failed to Marshal defaultPointerIgnition")
	}
	if string(savedPointerIgnitionJSON) != string(defaultPointerIgnitionJSON) {
		logrus.Infof("Master pointer ignition was modified. Saving contents to a machineconfig")
		mc := &mcfgv1.MachineConfig{}
		mc, err = generatePointerMachineConfig(*savedPointerIgnition, "master")
		if err != nil {
			return errors.Wrap(err, "failed to generate master installer machineconfig")
		}
		configData, err := yaml.Marshal(mc)
		if err != nil {
			return errors.Wrap(err, "failed to marshal master installer machineconfig")
		}
		a.File = &asset.File{
			Filename: masterMachineConfigFileName,
			Data:     configData,
		}
	}

	return nil
}

// Name returns the human-friendly name of the asset.
func (a *MasterIgnitionCustomizations) Name() string {
	return "Master Ignition Customization Check"
}

// Load does nothing, since we consume the ignition-configs
func (a *MasterIgnitionCustomizations) Load(f asset.FileFetcher) (found bool, err error) {
	return false, nil
}
