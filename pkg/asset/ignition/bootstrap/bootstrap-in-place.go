package bootstrap

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/types"
)

// verifyBootstrapInPlace validate the number of control plane replica is one
func verifyBootstrapInPlace(installConfig *types.InstallConfig) error {
	if *installConfig.ControlPlane.Replicas != 1 {
		return errors.Errorf("bootstrap in place require single control plane replica, current value: %d", *installConfig.ControlPlane.Replicas)
	}
	logrus.Warnf("Creating single node bootstrap in place configuration")
	return nil
}
