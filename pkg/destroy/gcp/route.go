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

const (
	routeResourceName = "resource"
)

type routeFilterFunc func(item *compute.Route) bool

func (o *ClusterUninstaller) listNetworkRoutes(ctx context.Context, networkURL string) ([]cloudResource, error) {
	return o.listRoutesWithFilter(ctx, "items(name),nextPageToken", func(item *compute.Route) bool { return item.Network == networkURL })
}

func (o *ClusterUninstaller) listRoutes(ctx context.Context) ([]cloudResource, error) {
	return o.listRoutesWithFilter(ctx, "items(name),nextPageToken", func(item *compute.Route) bool { return o.isClusterResource(item.Name) })
}

// listRoutesWithFilter lists routes in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listRoutesWithFilter(ctx context.Context, fields string, filterFunc routeFilterFunc) ([]cloudResource, error) {
	o.Logger.Debugf("Listing routes")
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	result := []cloudResource{}
	req := o.computeSvc.Routes.List(o.ProjectID).Fields(googleapi.Field(fields))

	err := req.Pages(ctx, func(list *compute.RouteList) error {
		for _, item := range list.Items {
			if filterFunc(item) {
				if strings.HasPrefix(item.Name, "default-route-") {
					continue
				}
				o.Logger.Debugf("Found route: %s", item.Name)
				result = append(result, cloudResource{
					key:      item.Name,
					name:     item.Name,
					typeName: routeResourceName,
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
		return nil, fmt.Errorf("failed to list routes: %w", err)
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteRoute(ctx context.Context, item cloudResource) error {
	o.Logger.Debugf("Deleting route %s", item.name)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	op, err := o.computeSvc.Routes.Delete(o.ProjectID, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	if err = o.handleOperation(op, err, item, "route"); err != nil {
		return err
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
	items := o.insertPendingItems(routeResourceName, found)
	for _, item := range items {
		err := o.deleteRoute(ctx, item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}
	if items = o.getPendingItems(routeResourceName); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
