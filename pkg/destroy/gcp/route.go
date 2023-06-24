package gcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"github.com/openshift/installer/pkg/types/gcp"
)

func (o *ClusterUninstaller) listNetworkRoutes(ctx context.Context, networkURL string) ([]cloudResource, error) {
	return o.listRoutesWithFilter(ctx, "items(name),nextPageToken", fmt.Sprintf("network eq %q", networkURL), nil)
}

func (o *ClusterUninstaller) listRoutes(ctx context.Context) ([]cloudResource, error) {
	return o.listRoutesWithFilter(ctx, "items(name),nextPageToken", o.clusterIDFilter(), nil)
}

// listRoutesWithFilter lists routes in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listRoutesWithFilter(ctx context.Context, fields string, filter string, filterFunc func(*compute.Route) bool) ([]cloudResource, error) {
	o.Logger.Debugf("Listing routes")
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	result := []cloudResource{}
	req := o.computeSvc.Routes.List(o.ProjectID).Fields(googleapi.Field(fields))
	if len(filter) > 0 {
		req = req.Filter(filter)
	}
	err := req.Pages(ctx, func(list *compute.RouteList) error {
		for _, item := range list.Items {
			if filterFunc == nil || filterFunc != nil && filterFunc(item) {
				if strings.HasPrefix(item.Name, "default-route-") {
					continue
				}
				o.Logger.Debugf("Found route: %s", item.Name)
				result = append(result, cloudResource{
					key:      item.Name,
					name:     item.Name,
					typeName: "route",
					quota: []gcp.QuotaUsage{{
						Metric: &gcp.Metric{
							Service: gcp.ServiceComputeEngineAPI,
							Limit:   "routes",
						},
						Amount: 1,
					}},
				})
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list routes")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteRoute(ctx context.Context, item cloudResource) error {
	o.Logger.Debugf("Deleting route %s", item.name)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	op, err := o.computeSvc.Routes.Delete(o.ProjectID, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID(item.typeName, item.name)
		return errors.Wrapf(err, "failed to delete route %s", item.name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID(item.typeName, item.name)
		return errors.Errorf("failed to delete route %s with error: %s", item.name, operationErrorMessage(op))
	}
	if (err != nil && isNoOp(err)) || (op != nil && op.Status == "DONE") {
		o.resetRequestID(item.typeName, item.name)
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted route %s", item.name)
	}
	return nil
}

// destroyRutes removes all route resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyRoutes(ctx context.Context) error {
	found, err := o.listRoutes(ctx)
	if err != nil {
		return err
	}
	items := o.insertPendingItems("route", found)
	for _, item := range items {
		err := o.deleteRoute(ctx, item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}
	if items = o.getPendingItems("route"); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
