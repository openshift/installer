package gcp

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"github.com/openshift/installer/pkg/types/gcp"
)

const (
	routerResourceName = "router"
)

func (o *ClusterUninstaller) listRouters(ctx context.Context) ([]cloudResource, error) {
	return o.listRoutersWithFilter(ctx, "items(name),nextPageToken", o.isClusterResource)
}

// listRoutersWithFilter lists routers in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listRoutersWithFilter(ctx context.Context, fields string, filterFunc resourceFilterFunc) ([]cloudResource, error) {
	o.Logger.Debug("Listing routers")
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	result := []cloudResource{}
	req := o.computeSvc.Routers.List(o.ProjectID, o.Region).Fields(googleapi.Field(fields))
	err := req.Pages(ctx, func(list *compute.RouterList) error {
		for _, item := range list.Items {
			if filterFunc(item.Name) {
				o.Logger.Debugf("Found router: %s", item.Name)
				result = append(result, cloudResource{
					key:      item.Name,
					name:     item.Name,
					typeName: routerResourceName,
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
	item.scope = global
	return o.handleOperation(ctx, op, err, item, "router")
}

// destroyRouters removes all router resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyRouters(ctx context.Context) error {
	found, err := o.listRouters(ctx)
	if err != nil {
		return err
	}
	items := o.insertPendingItems(routerResourceName, found)
	for _, item := range items {
		err := o.deleteRouter(ctx, item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}
	if items = o.getPendingItems(routerResourceName); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
