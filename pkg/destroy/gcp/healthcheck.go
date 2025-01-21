package gcp

import (
	"context"
	"fmt"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"github.com/openshift/installer/pkg/types/gcp"
)

const (
	globalHealthCheckResource = "healthcheck"
	regionHealthCheckResource = "regionHealthCheck"
)

func (o *ClusterUninstaller) listHealthChecks(ctx context.Context, typeName string) ([]cloudResource, error) {
	return o.listHealthChecksWithFilter(ctx, typeName, "items(name),nextPageToken", o.isClusterResource)
}

// listHealthChecksWithFilter lists health checks in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listHealthChecksWithFilter(ctx context.Context, typeName, fields string, filterFunc resourceFilterFunc) ([]cloudResource, error) {
	o.Logger.Debugf("Listing health checks")
	result := []cloudResource{}

	pagesFunc := func(list *compute.HealthCheckList) error {
		for _, item := range list.Items {
			if filterFunc(item.Name) {
				o.Logger.Debugf("Found health check: %s", item.Name)
				result = append(result, cloudResource{
					key:      item.Name,
					name:     item.Name,
					typeName: typeName,
					quota: []gcp.QuotaUsage{{
						Metric: &gcp.Metric{
							Service: gcp.ServiceComputeEngineAPI,
							Limit:   "health_checks",
						},
						Amount: 1,
					}},
				})
			}
		}
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var err error
	switch typeName {
	case globalHealthCheckResource:
		err = o.computeSvc.HealthChecks.List(o.ProjectID).Fields(googleapi.Field(fields)).Pages(ctx, pagesFunc)
	case regionHealthCheckResource:
		err = o.computeSvc.RegionHealthChecks.List(o.ProjectID, o.Region).Fields(googleapi.Field(fields)).Pages(ctx, pagesFunc)
	default:
		return nil, fmt.Errorf("invalid health check type %q", typeName)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list health checks: %w", err)
	}

	return result, nil
}

func (o *ClusterUninstaller) deleteHealthCheck(ctx context.Context, item cloudResource) error {
	o.Logger.Debugf("Deleting health check %s", item.name)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var err error
	var op *compute.Operation
	switch item.typeName {
	case globalHealthCheckResource:
		op, err = o.computeSvc.HealthChecks.Delete(o.ProjectID, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	case regionHealthCheckResource:
		op, err = o.computeSvc.RegionHealthChecks.Delete(o.ProjectID, o.Region, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	default:
		return fmt.Errorf("invalid health check type %q", item.typeName)
	}

	if err = o.handleOperation(op, err, item, "health check"); err != nil {
		return err
	}
	return nil
}

// destroyHealthChecks removes all health check resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyHealthChecks(ctx context.Context) error {
	healthCheckTypes := []string{globalHealthCheckResource, regionHealthCheckResource}
	for _, hct := range healthCheckTypes {
		found, err := o.listHealthChecks(ctx, hct)
		if err != nil {
			return err
		}
		items := o.insertPendingItems(hct, found)

		for _, item := range items {
			err := o.deleteHealthCheck(ctx, item)
			if err != nil {
				o.errorTracker.suppressWarning(item.key, err, o.Logger)
			}
		}

		if items := o.getPendingItems(hct); len(items) > 0 {
			return fmt.Errorf("%d %s resources pending", len(items), hct)
		}
	}

	return nil
}
