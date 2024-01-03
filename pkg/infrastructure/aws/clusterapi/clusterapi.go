package clusterapi

import (
	"fmt"

	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	"github.com/sirupsen/logrus"
)

type InfraHelper struct {
	clusterapi.CAPIInfraHelper
}

func (a InfraHelper) PreProvision(in clusterapi.PreProvisionInput) error {
	// TODO(padillon): skip if users bring their own roles
	if err := putIAMRoles(in.ClusterID, in.InstallConfig); err != nil {
		return fmt.Errorf("failed to create IAM roles: %w", err)
	}
	return nil
}

func (a InfraHelper) ControlPlaneAvailable(in clusterapi.ControlPlaneAvailableInput) error {
	logrus.Infoln("Calling AWS ControlPlaneAvailable")
	return nil
}
