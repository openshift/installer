package clusterapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	capg "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"

	ignutil "github.com/coreos/ignition/v2/config/util"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"

	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

type Provider struct {
}

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

	if err := FillBucket(ctx, BootstrapIgnitionBucket, BootstrapIgnitionBucket, string(in.BootstrapIgnData)); err != nil {
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
	// Create the public (optional) and private load balancer static ip addresses
	apiIntIPAddress, err := createLoadBalancerAddress(ctx, in.InstallConfig, in.InfraID, "", internalLoadBalancer)
	if err != nil {
		return nil
	}
	apiIPAddress, err := createLoadBalancerAddress(ctx, in.InstallConfig, in.InfraID, "", externalLoadBalancer)
	if err != nil {
		return nil
	}

	// Create the health checks related to the internal load balancer.
	if err := createHealthCheck(ctx, in.InstallConfig, in.InfraID, internalLoadBalancer, internalLoadBalancerPort); err != nil {
		return err
	}
	// Create the health checks related to the external load balancer.
	if err := createHealthCheck(ctx, in.InstallConfig, in.InfraID, externalLoadBalancer, externalLoadBalancerPort); err != nil {
		return err
	}

	gcpCluster := &capg.GCPCluster{}
	key := client.ObjectKey{
		Name:      in.InfraID,
		Namespace: capiutils.Namespace,
	}
	if err := in.Client.Get(ctx, key, gcpCluster); err != nil {
		return fmt.Errorf("failed to get GCP cluster: %w", err)
	}

	if gcpCluster.Status.Network.SelfLink == nil {
		return fmt.Errorf("failed to get GCP network: %w", err)
	}

	// Create the private zone if one does not exist
	if err := createPrivateManagedZone(ctx, in.InstallConfig, in.InfraID, *gcpCluster.Status.Network.SelfLink); err != nil {
		return err
	}

	// Create the public (optional) and private dns records
	if err := createDNSRecords(ctx, in.InstallConfig, in.InfraID, apiIPAddress, apiIntIPAddress); err != nil {
		return err
	}
	return nil
}
