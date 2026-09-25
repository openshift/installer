package clusterapi

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/compute/v1"
)

// operationStatusDone is the terminal status of a compute operation.
const operationStatusDone = "DONE"

// operationPollInterval bounds how often an operation is re-polled. The compute
// Wait methods normally block server-side for up to two minutes, but are
// documented to return after zero seconds when the server is overloaded, so
// polling needs its own floor. Overridden in tests.
var operationPollInterval = 5 * time.Second

// waitForOperation polls a compute operation until it reports DONE.
//
// A single call to the compute Wait methods is not enough: they return once the
// operation is DONE *or* once the request approaches its two minute deadline,
// whichever comes first, so a long-running operation such as image creation is
// still in progress when the first call returns. Callers bound the total wait
// through ctx.
func waitForOperation(ctx context.Context, name string, wait func() (*compute.Operation, error)) error {
	for {
		op, err := wait()
		if err != nil {
			return fmt.Errorf("failed to wait for operation %s: %w", name, err)
		}

		if op != nil && op.Status == operationStatusDone {
			// A DONE operation still reports failures through Error, so
			// returning without checking it would swallow them.
			if op.Error != nil && len(op.Error.Errors) > 0 {
				return fmt.Errorf("operation %s failed: %s: %s", name, op.Error.Errors[0].Code, op.Error.Errors[0].Message)
			}
			return nil
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting for operation %s to complete: %w", name, ctx.Err())
		case <-time.After(operationPollInterval):
		}
	}
}

// WaitForOperationGlobal will attempt to wait for a operation to complete where the operational
// resource is globally scoped.
func WaitForOperationGlobal(ctx context.Context, svc *compute.Service, projectID string, operation *compute.Operation) error {
	g := compute.NewGlobalOperationsService(svc)
	return waitForOperation(ctx, operation.Name, func() (*compute.Operation, error) {
		return g.Wait(projectID, operation.Name).Context(ctx).Do()
	})
}

// WaitForOperationRegional will attempt to wait for a operation to complete where the operational
// resource is regionally scoped.
func WaitForOperationRegional(ctx context.Context, svc *compute.Service, projectID, region string, operation *compute.Operation) error {
	r := compute.NewRegionOperationsService(svc)
	return waitForOperation(ctx, operation.Name, func() (*compute.Operation, error) {
		return r.Wait(projectID, region, operation.Name).Context(ctx).Do()
	})
}
