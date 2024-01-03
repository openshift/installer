package clusterapi

import (
	"fmt"

	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
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
	if err := createDNSRecords(in.Cluster, in.InstallConfig); err != nil {
		return fmt.Errorf("failed to create DNS records: %w", err)
	}
	return nil
}
