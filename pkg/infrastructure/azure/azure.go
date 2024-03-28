package azure

import (
	"context"

	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
)

var _ clusterapi.Provider = (*Provider)(nil)

// Provider implements Azure CAPI installation.
type Provider struct{}

// Name gives the name of the provider, Azure.
func (*Provider) Name() string { return azuretypes.Name }

// InfraReady sets the DNS currently after the ignition is done.
func (p *Provider) InfraReady(ctx context.Context, in clusterapi.InfraReadyInput) error {
	return createDNSEntries(ctx, in)
}
