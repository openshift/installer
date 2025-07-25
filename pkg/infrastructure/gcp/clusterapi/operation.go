package clusterapi

import (
	"context"
	"fmt"

	"google.golang.org/api/compute/v1"
)

// WaitForOperationGlobal will attempt to wait for a operation to complete where the operational
// resource is globally scoped.
func WaitForOperationGlobal(ctx context.Context, svc *compute.Service, projectID string, operation *compute.Operation) error {
	g := compute.NewGlobalOperationsService(svc)
	if _, err := g.Wait(projectID, operation.Name).Context(ctx).Do(); err != nil {
		return fmt.Errorf("failed to wait for global operation: %w", err)
	}

	return nil
}

// WaitForOperationRegional will attempt to wait for a operation to complete where the operational
// resource is regionally scoped.
func WaitForOperationRegional(ctx context.Context, svc *compute.Service, projectID, region string, operation *compute.Operation) error {
	r := compute.NewRegionOperationsService(svc)
	if _, err := r.Wait(projectID, region, operation.Name).Context(ctx).Do(); err != nil {
		return fmt.Errorf("failed to wait for regional operation: %w", err)
	}

	return nil
}
