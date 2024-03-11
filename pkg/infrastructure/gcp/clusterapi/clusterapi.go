package clusterapi

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	capg "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap/gcp"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	"github.com/openshift/installer/pkg/types"
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
	// Create the bucket and presigned url. The url is generated using a known/expected name so that the
	// url can be retrieved from the api by this name.
	ctx, cancel := context.WithTimeout(ctx, time.Minute*2)
	defer cancel()

	bucketName := gcp.GetBootstrapStorageName(in.InfraID)
	bucketHandle, err := gcp.CreateBucketHandle(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to create bucket handle %s: %w", bucketName, err)
	}

	url, err := gcp.ProvisionBootstrapStorage(ctx, in.InstallConfig, bucketHandle, in.InfraID)
	if err != nil {
		return nil, fmt.Errorf("ignition failed to provision storage: %w", err)
	}
	editedIgnitionBytes, err := EditIgnition(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed to edit bootstrap ignition: %w", err)
	}

	ignitionBytes := in.BootstrapIgnData
	if editedIgnitionBytes != nil {
		ignitionBytes = editedIgnitionBytes
	}

	if err := gcp.FillBucket(ctx, bucketHandle, string(ignitionBytes)); err != nil {
		return nil, fmt.Errorf("ignition failed to fill bucket: %w", err)
	}

	ignShim, err := bootstrap.GenerateIgnitionShimWithCertBundleAndProxy(url, in.InstallConfig.Config.AdditionalTrustBundle, in.InstallConfig.Config.Proxy)
	if err != nil {
		return nil, fmt.Errorf("failed to create ignition shim: %w", err)
	}

	return ignShim, nil
}

// InfraReady is called once cluster.Status.InfrastructureReady
// is true, typically after load balancers have been provisioned. It can be used
// to create DNS records.
func (p Provider) InfraReady(ctx context.Context, in clusterapi.InfraReadyInput) error {
	gcpCluster := &capg.GCPCluster{}
	key := client.ObjectKey{
		Name:      in.InfraID,
		Namespace: capiutils.Namespace,
	}
	if err := in.Client.Get(ctx, key, gcpCluster); err != nil {
		return fmt.Errorf("failed to get GCP cluster: %w", err)
	}

	// public load balancer is created by CAPG. The health check for this load balancer is also created by
	// the CAPG.
	apiIPAddress := gcpCluster.Spec.ControlPlaneEndpoint.Host
	if apiIPAddress == "" && in.InstallConfig.Config.Publish == types.ExternalPublishingStrategy {
		logrus.Debugf("publish strategy is set to external but api address is empty")
	}

	// Currently, the internal/private load balancer is not created by CAPG. The load balancer will be created
	// by the installer for now
	// TODO: remove the creation of the LB and health check here when supported by CAPG.
	// https://github.com/kubernetes-sigs/cluster-api-provider-gcp/issues/903
	// Create the public (optional) and private load balancer static ip addresses
	// TODO: Do we then need to setup a subnet for internal load balancing ?
	apiIntIPAddress, err := createInternalLBAddress(ctx, in)
	if err != nil {
		return fmt.Errorf("failed to create internal load balancer address: %w", err)
	}

	if in.InstallConfig.Config.GCP.UserProvisionedDNS != gcptypes.UserProvisionedDNSEnabled {
		// Get the network from the GCP Cluster. The network is used to create the private managed zone.
		if gcpCluster.Status.Network.SelfLink == nil {
			return fmt.Errorf("failed to get GCP network: %w", err)
		}

		// Create the private zone if one does not exist
		if err := createPrivateManagedZone(ctx, in.InstallConfig, in.InfraID, *gcpCluster.Status.Network.SelfLink); err != nil {
			return fmt.Errorf("failed to create the private managed zone: %w", err)
		}

		// Create the public (optional) and private dns records
		if err := createDNSRecords(ctx, in.InstallConfig, in.InfraID, apiIPAddress, apiIntIPAddress); err != nil {
			return fmt.Errorf("failed to create DNS records: %w", err)
		}
	}

	if err := createInternalLoadBalancer(ctx, in.InstallConfig, in.InfraID); err != nil {
		return fmt.Errorf("failed to create internal load balancer")
	}

	return nil
}

// PostProvision should be called to add or update and GCP resources after provisioning has completed.
func (p Provider) PostProvision(ctx context.Context, in clusterapi.PostProvisionInput) error {
	return nil
}
