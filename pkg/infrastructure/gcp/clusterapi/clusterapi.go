package clusterapi

import (
	"context"

	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

// Provider implements gcp infrastructure in conjunction with the
// GCP CAPI provider.
type Provider struct {
}

var _ clusterapi.PreProvider = (*Provider)(nil)
var _ clusterapi.IgnitionProvider = (*Provider)(nil)
var _ clusterapi.InfraReadyProvider = (*Provider)(nil)
var _ clusterapi.PostProvider = (*Provider)(nil)

// Name returns the name for the platform.
func (p Provider) Name() string {
	return gcptypes.Name
}

// PreProvision is called before provisioning using CAPI controllers has initiated.
// GCP resources that are not created by CAPG (and are required for other stages of the install) are
// created here using the gcp sdk.
func (p Provider) PreProvision(ctx context.Context, in clusterapi.PreProvisionInput) error {
	return nil
}

// Ignition provisions the GCP bucket and url that points to the bucket. Bootstrap ignition data cannot
// populate the metadata field of the bootstrap instance as the data can be too large. Instead, the data is
// added to a bucket. A signed url is generated to point to the bucket and the ignition data will be
// updated to point to the url. This is also allows for bootstrap data to be edited after its initial creation.
func (p Provider) Ignition(ctx context.Context, in clusterapi.IgnitionInput) ([]byte, error) {
	return nil, nil
}

// InfraReady is called once cluster.Status.InfrastructureReady
// is true, typically after load balancers have been provisioned. It can be used
// to create DNS records.
func (p Provider) InfraReady(ctx context.Context, in clusterapi.InfraReadyInput) error {
	return nil
}

// PostProvision should be called to add or update and GCP resources after provisioning has completed.
func (p Provider) PostProvision(ctx context.Context, in clusterapi.PostProvisionInput) error {
	return nil
}
