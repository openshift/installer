package bootstrap

import (
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/types"
)

// SingleNodeBootstrapInPlace is the data to use to replace values in bootstrap template files.
type SingleNodeBootstrapInPlace struct {
	BootstrapInPlace    bool
	CoreosInstallerArgs string
}

// GetSingleNodeBootstrapInPlaceConfig generates the config for the BootstrapInPlace.
func GetSingleNodeBootstrapInPlaceConfig() *SingleNodeBootstrapInPlace {
	return &SingleNodeBootstrapInPlace{
		BootstrapInPlace:    true,
		CoreosInstallerArgs: getCoreosInstallerArgs(),
	}
}

// verifyBootstrapInPlace validate the number of control plane replica is one
func verifyBootstrapInPlace(installConfig *types.InstallConfig) error {
	if *installConfig.ControlPlane.Replicas != 1 {
		return errors.Errorf("bootstrap in place require single control plane replica, current value: %d", *installConfig.ControlPlane.Replicas)
	}
	logrus.Warnf("Creating single node bootstrap in place configuration")
	return nil
}

// getCoreosInstallerArgs checks for bootstrap in place coreos installer args env
func getCoreosInstallerArgs() string {
	coreosInstallerEnv := os.Getenv("OPENSHIFT_INSTALL_EXPERIMENTAL_BOOTSTRAP_IN_PLACE_COREOS_INSTALLER_ARGS")
	if coreosInstallerEnv != "" {
		logrus.Warnf("Setting coreos-installer args: %s", coreosInstallerEnv)
	}
	return coreosInstallerEnv
}
