package clusterapi

import (
	"context"

	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	powervstypes "github.com/openshift/installer/pkg/types/powervs"
)

// Provider is the vSphere implementation of the clusterapi InfraProvider.
type Provider struct {
	clusterapi.InfraProvider
}

var _ clusterapi.PreProvider = (*Provider)(nil)
var _ clusterapi.Provider = (*Provider)(nil)

// Name returns the PowerVS provider name.
func (p Provider) Name() string {
	return powervstypes.Name
}

// PreProvision creates the PowerVS objects required prior to running capv.
func (p Provider) PreProvision(ctx context.Context, in clusterapi.PreProvisionInput) error {
	return nil
}
