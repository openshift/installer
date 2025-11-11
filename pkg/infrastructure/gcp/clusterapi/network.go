package clusterapi

import (
	"context"
	"fmt"

	"google.golang.org/api/option"

	gcpconfig "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

func getAPIAddressName(infraID string) string {
	return fmt.Sprintf("%s-api-internal", infraID)
}

func getInternalLBAddress(ctx context.Context, project, region, name string, endpoint *gcptypes.PSCEndpoint) (string, error) {
	opts := []option.ClientOption{}
	if gcptypes.ShouldUseEndpointForInstaller(endpoint) {
		opts = append(opts, gcpconfig.CreateEndpointOption(endpoint.Name, gcpconfig.ServiceNameGCPCompute))
	}
	service, err := gcpconfig.GetComputeService(ctx, opts...)
	if err != nil {
		return "", err
	}

	addrOutput, err := service.Addresses.Get(project, region, name).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to get compute address %s: %w", name, err)
	}
	return addrOutput.Address, nil
}
