package gcp

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"github.com/openshift/installer/pkg/types/gcp"
)

func (o *ClusterUninstaller) listHTTPHealthChecks(ctx context.Context) ([]cloudResource, error) {
	return o.listHTTPHealthChecksWithFilter(ctx, "items(name),nextPageToken", o.clusterIDFilter(), nil)
}

// listHTTPHealthChecksWithFilter lists HTTP Health Checks in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listHTTPHealthChecksWithFilter(ctx context.Context, fields string, filter string, filterFunc func(*compute.HttpHealthCheck) bool) ([]cloudResource, error) {
	o.Logger.Debugf("Listing HTTP health checks")
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	result := []cloudResource{}
	req := o.computeSvc.HttpHealthChecks.List(o.ProjectID).Fields(googleapi.Field(fields))
	if len(filter) > 0 {
		req = req.Filter(filter)
	}
	err := req.Pages(ctx, func(list *compute.HttpHealthCheckList) error {
		for _, item := range list.Items {
			if filterFunc == nil || filterFunc != nil && filterFunc(item) {
				o.Logger.Debugf("Found HTTP health check: %s", item.Name)
				result = append(result, cloudResource{
					key:      item.Name,
					name:     item.Name,
					typeName: "httphealthcheck",
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
		return nil, errors.Wrapf(err, "failed to list HTTP health checks")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteHTTPHealthCheck(ctx context.Context, item cloudResource) error {
	o.Logger.Debugf("Deleting HTTP health check %s", item.name)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	op, err := o.computeSvc.HttpHealthChecks.Delete(o.ProjectID, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID(item.typeName, item.name)
		return errors.Wrapf(err, "failed to delete HTTP health check %s", item.name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID(item.typeName, item.name)
		return errors.Errorf("failed to delete HTTP health check %s with error: %s", item.name, operationErrorMessage(op))
	}
	if (err != nil && isNoOp(err)) || (op != nil && op.Status == "DONE") {
		o.resetRequestID(item.typeName, item.name)
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted HTTP health check %s", item.name)
	}
	return nil
}

// destroyHTTPHealthChecks removes all HTTP health check resources that have a name prefixed
// with the cluster's infra ID
func (o *ClusterUninstaller) destroyHTTPHealthChecks(ctx context.Context) error {
	found, err := o.listHTTPHealthChecks(ctx)
	if err != nil {
		return err
	}
	items := o.insertPendingItems("httphealthcheck", found)
	for _, item := range items {
		err := o.deleteHTTPHealthCheck(ctx, item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}
	if items = o.getPendingItems("httphealthcheck"); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
