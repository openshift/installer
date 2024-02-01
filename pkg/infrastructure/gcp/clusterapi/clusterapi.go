package clusterapi

import (
	"context"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"

	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
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

	return nil, nil
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

	// Create the health checks related to the load balancers
	if err := createHealthCheck(ctx, in.InstallConfig, in.InfraID, internalLoadBalancer, internalLoadBalancerPort); err != nil {
		return err
	}
	if err := createHealthCheck(ctx, in.InstallConfig, in.InfraID, externalLoadBalancer, externalLoadBalancerPort); err != nil {
		return err
	}

	// Create the private zone if one does not exist
	if err := createPrivateManagedZone(ctx, in.InstallConfig, in.InfraID); err != nil {
		return err
	}

	// Create the public (optional) and private dns records
	if err := createDNSRecords(ctx, in.InstallConfig, in.InfraID, apiIPAddress, apiIntIPAddress); err != nil {
		return err
	}
	return nil
}
