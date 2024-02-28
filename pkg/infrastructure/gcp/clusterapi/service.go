package clusterapi

import (
	"context"
	"fmt"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"

	"github.com/openshift/installer/pkg/asset/installconfig/gcp"
)

// NewComputeService wraps the creation of a gcp compute service creation.
func NewComputeService() (*compute.Service, error) {
	ctx := context.Background()

	ssn, err := gcp.GetSession(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	service, err := compute.NewService(ctx, option.WithCredentials(ssn.Credentials))
	if err != nil {
		return nil, fmt.Errorf("failed to create compute service: %w", err)
	}

	return service, nil
}
