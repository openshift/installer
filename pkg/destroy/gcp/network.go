package gcp

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

const (
	networkResourceName = "network"
)

func (o *ClusterUninstaller) listNetworks(ctx context.Context) ([]cloudResource, error) {
	return o.listNetworksWithFilter(ctx, "items(name,selfLink),nextPageToken")
}

// listNetworksWithFilter lists addresses in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listNetworksWithFilter(ctx context.Context, fields string) ([]cloudResource, error) {
	o.Logger.Debugf("Listing networks")
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	result := []cloudResource{}
	req := o.computeSvc.Networks.List(o.ProjectID).Fields(googleapi.Field(fields))

	err := req.Pages(ctx, func(list *compute.NetworkList) error {
		for _, item := range list.Items {
			if o.isClusterResource(item.Name) {
				o.Logger.Debugf("Found network: %s", item.Name)
				result = append(result, cloudResource{
					key:      item.Name,
					name:     item.Name,
					typeName: networkResourceName,
					url:      item.SelfLink,
				})
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list networks: %w", err)
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteNetwork(ctx context.Context, item cloudResource) error {
	o.Logger.Debugf("Deleting network %s", item.name)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	op, err := o.computeSvc.Networks.Delete(o.ProjectID, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	if err = o.handleOperation(op, err, item, "network"); err != nil {
		return err
	}
	return nil
}

// destroyNetworks removes all vpc network resources prefixed
// with the cluster's infra ID, and deletes all of the routes.
func (o *ClusterUninstaller) destroyNetworks(ctx context.Context) error {
	found, err := o.listNetworks(ctx)
	if err != nil {
		return err
	}
	items := o.insertPendingItems(networkResourceName, found)
	for _, item := range items {
		foundRoutes, err := o.listNetworkRoutes(ctx, item.url)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
			continue
		}
		routes := o.insertPendingItems(routeResourceName, foundRoutes)
		for _, route := range routes {
			err := o.deleteRoute(ctx, route)
			if err != nil {
				o.errorTracker.suppressWarning(route.key, err, o.Logger)
			}
		}
		err = o.deleteNetwork(ctx, item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}
	if items = o.getPendingItems(networkResourceName); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
