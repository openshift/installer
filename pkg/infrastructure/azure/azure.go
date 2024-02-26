package azure

import (
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
)

var _ clusterapi.Provider = (*Provider)(nil)

// Provider implements Azure CAPI installation.
type Provider struct{}

// Name gives the name of the provider, Azure.
func (*Provider) Name() string { return azuretypes.Name }
