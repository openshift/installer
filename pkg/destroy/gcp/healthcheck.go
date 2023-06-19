package gcp

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"github.com/openshift/installer/pkg/types/gcp"
)

func (o *ClusterUninstaller) listHealthChecks(ctx context.Context) ([]cloudResource, error) {
	return o.listHealthChecksWithFilter(ctx, "items(name),nextPageToken", o.clusterIDFilter(), nil)
}

// listHealthChecksWithFilter lists health checks in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listHealthChecksWithFilter(ctx context.Context, fields string, filter string, filterFunc func(*compute.HealthCheck) bool) ([]cloudResource, error) {
	o.Logger.Debugf("Listing health checks")
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	result := []cloudResource{}
	req := o.computeSvc.HealthChecks.List(o.ProjectID).Fields(googleapi.Field(fields))
	if len(filter) > 0 {
		req = req.Filter(filter)
	}
	err := req.Pages(ctx, func(list *compute.HealthCheckList) error {
		for _, item := range list.Items {
			if filterFunc == nil || filterFunc != nil && filterFunc(item) {
				o.Logger.Debugf("Found health check: %s", item.Name)
				result = append(result, cloudResource{
					key:      item.Name,
					name:     item.Name,
					typeName: "healthcheck",
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
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list health checks")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteHealthCheck(ctx context.Context, item cloudResource) error {
	o.Logger.Debugf("Deleting health check %s", item.name)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	op, err := o.computeSvc.HealthChecks.Delete(o.ProjectID, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID(item.typeName, item.name)
		return errors.Wrapf(err, "failed to delete health check %s", item.name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID(item.typeName, item.name)
		return errors.Errorf("failed to delete health check %s with error: %s", item.name, operationErrorMessage(op))
	}
	if (err != nil && isNoOp(err)) || (op != nil && op.Status == "DONE") {
		o.resetRequestID(item.typeName, item.name)
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted health check %s", item.name)
	}
	return nil
}

// destroyHealthChecks removes all health check resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyHealthChecks(ctx context.Context) error {
	found, err := o.listHealthChecks(ctx)
	if err != nil {
		return err
	}
	items := o.insertPendingItems("healthcheck", found)
	for _, item := range items {
		err := o.deleteHealthCheck(ctx, item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}
	if items = o.getPendingItems("healthcheck"); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
