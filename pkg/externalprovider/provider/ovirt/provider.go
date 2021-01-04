package ovirt

import (
	"github.com/openshift/installer/pkg/externalprovider/provider"
	"github.com/openshift/installer/pkg/externalprovider/provider/defaultprovider"
	ovirttypes "github.com/openshift/installer/pkg/types/ovirt"
)

// NewOvirtProvider creates a new oVirt provider.
func NewOvirtProvider() provider.ExternalProvider {
	return &ovirtProvider{}
}

// OvirtProvider
type ovirtProvider struct {
	defaultprovider.DefaultProvider
}

// Name returns the name of the provider
func (ovirt *ovirtProvider) Name() string {
	return ovirttypes.Name
}
