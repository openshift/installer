package clusterapi

import (
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
)

type Provider struct {
	clusterapi.DefaultCAPIProvider
}

func (p Provider) PreProvision(in clusterapi.PreProvisionInput) error {

	return nil
}

func (p Provider) ControlPlaneAvailable(in clusterapi.ControlPlaneAvailableInput) error {
	return nil
}
