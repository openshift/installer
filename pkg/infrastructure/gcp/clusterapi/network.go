package clusterapi

import (
	"context"
	"fmt"
)

func getAPIAddressName(infraID string) string {
	return fmt.Sprintf("%s-api-internal", infraID)
}

func getInternalLBAddress(ctx context.Context, project, region, name string) (string, error) {
	service, err := NewComputeService()
	if err != nil {
		return "", err
	}

	addrOutput, err := service.Addresses.Get(project, region, name).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to get compute address %s: %w", name, err)
	}
	return addrOutput.Address, nil
}
