package clusterapi

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	ignutil "github.com/coreos/ignition/v2/config/util"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	capg "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

type Provider struct {
}

var _ clusterapi.PreProvider = (*Provider)(nil)
var _ clusterapi.IgnitionProvider = (*Provider)(nil)
var _ clusterapi.InfraReadyProvider = (*Provider)(nil)
var _ clusterapi.PostProvider = (*Provider)(nil)

func (p Provider) Name() string {
	return gcptypes.Name
}

func (p Provider) PreProvision(ctx context.Context, in clusterapi.PreProvisionInput) error {
	return nil
}

func (p Provider) Ignition(ctx context.Context, in clusterapi.IgnitionInput) ([]byte, error) {
	// Create the bucket and presigned url. The url is generated using a known/expected name so that the
	// url can be retrieved from the api by this name.

	ctx, cancel := context.WithTimeout(ctx, time.Minute*2)
	defer cancel()

	url, err := ProvisionBootstrapStorage(ctx, in.InstallConfig, in.InfraID)
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

	if err := FillBucket(ctx, BootstrapIgnitionBucket, BootstrapIgnitionBucket, string(ignitionBytes)); err != nil {
		return nil, fmt.Errorf("ignition failed to fill bucket: %w", err)
	}

	// Generate an ignition stub where the URL is stored
	ign := igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
			Config: igntypes.IgnitionConfig{
				Replace: igntypes.Resource{
					Source: ignutil.StrToPtr(url),
				},
			},
		},
	}

	ignShimBytes, err := json.Marshal(ign)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ignition shim: %w", err)
	}

	return ignShimBytes, nil
}

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
	apiIPAddress := *gcpCluster.Status.Network.APIServerAddress

	// Currently, the internal/private load balancer is not created by CAPG. The load balancer will be created
	// by the installer for now
	// TODO: remove the creation of the LB and health check here when supported by CAPG.
	// https://github.com/kubernetes-sigs/cluster-api-provider-gcp/issues/903
	// Create the public (optional) and private load balancer static ip addresses
	// TODO: Do we then need to setup a subnet for internal load balancing ?
	apiIntIPAddress, err := getInternalLBAddress(ctx, in.InstallConfig.Config.GCP.ProjectID, in.InstallConfig.Config.GCP.Region, getApiAddressName(in.InfraID))
	if err != nil {
		return fmt.Errorf("failed to create the internal load balancer address: %w", err)
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

	return nil
}

func (p Provider) PostProvision(ctx context.Context, in clusterapi.PostProvisionInput) error {
	// Create the backend service for the internal load balancer
	// TODO: remove when CAPG supports private installs and load balancer creation
	if err := createInternalLBBackendService(ctx, in); err != nil {
		return fmt.Errorf("failed to create the internal load balancer backend service: %w", err)
	}

	return nil
}
