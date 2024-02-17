package clusterapi

import (
	"context"
	"fmt"

	"google.golang.org/api/compute/v1"
)

// NewComputeService wraps the creation of a gcp compute service creation.
func NewComputeService() (*compute.Service, error) {
	ctx := context.Background()

	service, err := compute.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create compute service: %w", err)
	}

	return service, nil
}
