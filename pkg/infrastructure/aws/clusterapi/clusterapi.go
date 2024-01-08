package clusterapi

import (
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	"github.com/sirupsen/logrus"
)

type InfraHelper struct {
	clusterapi.CAPIInfraHelper
}

func (a InfraHelper) PreProvision(in clusterapi.PreProvisionInput) error {
	logrus.Infoln("Calling AWS PreProvision override")
	return nil
}

func (a InfraHelper) ControlPlaneAvailable(in clusterapi.ControlPlaneAvailableInput) error {
	logrus.Infoln("Calling AWS ControlPlaneAvailable")
	return nil
}
