package clusterapi

import (
	"context"
	"fmt"

	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	awstypes "github.com/openshift/installer/pkg/types/aws"
)

var _ clusterapi.Provider = (*Provider)(nil)
var _ clusterapi.PreProvider = (*Provider)(nil)

// Provider implements AWS CAPI installation.
type Provider struct{}

// Name gives the name of the provider, AWS.
func (*Provider) Name() string { return awstypes.Name }

func (*Provider) PreProvision(ctx context.Context, in clusterapi.PreProvisionInput) error {
	if err := putIAMRoles(ctx, in.InfraID, in.InstallConfig); err != nil {
		return fmt.Errorf("failed to create IAM roles: %w", err)
	}
	return nil
}
