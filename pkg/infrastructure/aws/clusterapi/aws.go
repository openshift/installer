package clusterapi

import (
	"context"
	"fmt"

	"k8s.io/utils/ptr"
	capa "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"

	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	awstypes "github.com/openshift/installer/pkg/types/aws"
)

var _ clusterapi.Provider = (*Provider)(nil)
var _ clusterapi.PreProvider = (*Provider)(nil)

// Provider implements AWS CAPI installation.
type Provider struct{}

// Name gives the name of the provider, AWS.
func (*Provider) Name() string { return awstypes.Name }

// PreProvision creates the IAM roles used by all nodes in the cluster.
func (*Provider) PreProvision(ctx context.Context, in clusterapi.PreProvisionInput) error {
	if err := createIAMRoles(ctx, in.InfraID, in.InstallConfig); err != nil {
		return fmt.Errorf("failed to create IAM roles: %w", err)
	}

	amiID, err := copyAMIToRegion(ctx, in.InstallConfig, in.InfraID, string(*in.RhcosImage))
	if err != nil {
		return fmt.Errorf("failed to copy AMI: %w", err)
	}
	for i := range in.MachineManifests {
		if awsMachine, ok := in.MachineManifests[i].(*capa.AWSMachine); ok {
			awsMachine.Spec.AMI = capa.AMIReference{ID: ptr.To(amiID)}
		}
	}
	return nil
}
