package clusterapi

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/compute/v1"
)

// WaitForOperationGlobal will attempt to wait for a operation to complete where the operational
// resource is globally scoped.
func WaitForOperationGlobal(ctx context.Context, projectID string, operation *compute.Operation) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()

	service, err := NewComputeService()
	if err != nil {
		return err
	}

	g := compute.NewGlobalOperationsService(service)
	if _, err := g.Wait(projectID, operation.Name).Context(ctx).Do(); err != nil {
		return fmt.Errorf("failed to wait for regional operation: %w", err)
	}

	return nil
}

// WaitForOperationRegional will attempt to wait for a operation to complete where the operational
// resource is regionally scoped.
func WaitForOperationRegional(ctx context.Context, projectID, region string, operation *compute.Operation) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()

	service, err := NewComputeService()
	if err != nil {
		return err
	}

	r := compute.NewRegionOperationsService(service)
	if _, err := r.Wait(projectID, region, operation.Name).Context(ctx).Do(); err != nil {
		return fmt.Errorf("failed to wait for regional operation: %w", err)
	}

	return nil
}
