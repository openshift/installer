package gcp

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"github.com/openshift/installer/pkg/types/gcp"
)

func (o *ClusterUninstaller) listRouters(ctx context.Context) ([]cloudResource, error) {
	return o.listRoutersWithFilter(ctx, "items(name),nextPageToken", o.clusterIDFilter(), nil)
}

// listRoutersWithFilter lists routers in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listRoutersWithFilter(ctx context.Context, fields string, filter string, filterFunc func(*compute.Router) bool) ([]cloudResource, error) {
	o.Logger.Debug("Listing routers")
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	result := []cloudResource{}
	req := o.computeSvc.Routers.List(o.ProjectID, o.Region).Fields(googleapi.Field(fields))
	if len(filter) > 0 {
		req = req.Filter(filter)
	}
	err := req.Pages(ctx, func(list *compute.RouterList) error {
		for _, item := range list.Items {
			if filterFunc == nil || filterFunc != nil && filterFunc(item) {
				o.Logger.Debugf("Found router: %s", item.Name)
				result = append(result, cloudResource{
					key:      item.Name,
					name:     item.Name,
					typeName: "router",
					quota: []gcp.QuotaUsage{{
						Metric: &gcp.Metric{
							Service: gcp.ServiceComputeEngineAPI,
							Limit:   "routers",
						},
						Amount: 1,
					}},
				})
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to list routers")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteRouter(ctx context.Context, item cloudResource) error {
	o.Logger.Debugf("Deleting router %s", item.name)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	op, err := o.computeSvc.Routers.Delete(o.ProjectID, o.Region, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID(item.typeName, item.name)
		return errors.Wrapf(err, "failed to delete router %s", item.name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID(item.typeName, item.name)
		return errors.Errorf("failed to delete router %s with error: %s", item.name, operationErrorMessage(op))
	}
	if (err != nil && isNoOp(err)) || (op != nil && op.Status == "DONE") {
		o.resetRequestID(item.typeName, item.name)
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted router %s", item.name)
	}
	return nil
}

// destroyRouters removes all router resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyRouters(ctx context.Context) error {
	found, err := o.listRouters(ctx)
	if err != nil {
		return err
	}
	items := o.insertPendingItems("router", found)
	for _, item := range items {
		err := o.deleteRouter(ctx, item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}
	if items = o.getPendingItems("router"); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
