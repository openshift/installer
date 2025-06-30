package clusterapi

import (
	"context"
	"fmt"

	configv1 "github.com/openshift/api/config/v1"
	gcpconfig "github.com/openshift/installer/pkg/asset/installconfig/gcp"
)

func getAPIAddressName(infraID string) string {
	return fmt.Sprintf("%s-api-internal", infraID)
}

func getInternalLBAddress(ctx context.Context, project, region, name string, endpoints []configv1.GCPServiceEndpoint) (string, error) {
	service, err := gcpconfig.GetComputeService(ctx, endpoints)
	if err != nil {
		return "", err
	}

	addrOutput, err := service.Addresses.Get(project, region, name).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to get compute address %s: %w", name, err)
	}
	return addrOutput.Address, nil
}
