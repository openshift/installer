package gcp

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"github.com/openshift/installer/pkg/types/gcp"
)

const (
	httpHealthCheckResourceName = "httphealthcheck"
)

func (o *ClusterUninstaller) listHTTPHealthChecks(ctx context.Context) ([]cloudResource, error) {
	return o.listHTTPHealthChecksWithFilter(ctx, "items(name),nextPageToken", o.isClusterResource)
}

// listHTTPHealthChecksWithFilter lists HTTP Health Checks in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listHTTPHealthChecksWithFilter(ctx context.Context, fields string, filterFunc resourceFilterFunc) ([]cloudResource, error) {
	o.Logger.Debugf("Listing HTTP health checks")
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	result := []cloudResource{}
	req := o.computeSvc.HttpHealthChecks.List(o.ProjectID).Fields(googleapi.Field(fields))

	err := req.Pages(ctx, func(list *compute.HttpHealthCheckList) error {
		for _, item := range list.Items {
			if filterFunc(item.Name) {
				o.Logger.Debugf("Found HTTP health check: %s", item.Name)
				result = append(result, cloudResource{
					key:      item.Name,
					name:     item.Name,
					typeName: httpHealthCheckResourceName,
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
	item.scope = global
	return o.handleOperation(ctx, op, err, item, "HTTP health check")
}

// destroyHTTPHealthChecks removes all HTTP health check resources that have a name prefixed
// with the cluster's infra ID
func (o *ClusterUninstaller) destroyHTTPHealthChecks(ctx context.Context) error {
	found, err := o.listHTTPHealthChecks(ctx)
	if err != nil {
		return err
	}
	items := o.insertPendingItems(httpHealthCheckResourceName, found)
	for _, item := range items {
		err := o.deleteHTTPHealthCheck(ctx, item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}
	if items = o.getPendingItems(httpHealthCheckResourceName); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
