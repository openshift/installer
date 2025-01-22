package gcp

import (
	"context"
	"fmt"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"github.com/openshift/installer/pkg/types/gcp"
)

const (
	globalBackendServiceResource = "backendservice"
	regionBackendServiceResource = "regionbackendservice"
)

func (o *ClusterUninstaller) listBackendServices(ctx context.Context, typeName string) ([]cloudResource, error) {
	return o.listBackendServicesWithFilter(ctx, typeName, "items(name),nextPageToken",
		func(item *compute.BackendService) bool {
			return o.isClusterResource(item.Name)
		})
}

// listBackendServicesWithFilter lists backend services in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result.
func (o *ClusterUninstaller) listBackendServicesWithFilter(ctx context.Context, typeName, fields string, filterFunc func(item *compute.BackendService) bool) ([]cloudResource, error) {
	o.Logger.Debugf("Listing backend services")
	result := []cloudResource{}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	pagesFunc := func(list *compute.BackendServiceList) error {
		for _, item := range list.Items {
			o.Logger.Debugf("Found backend service: %s", item.Name)
			if filterFunc(item) {
				result = append(result, cloudResource{
					key:      item.Name,
					name:     item.Name,
					typeName: typeName,
					quota: []gcp.QuotaUsage{{
						Metric: &gcp.Metric{
							Service: gcp.ServiceComputeEngineAPI,
							Limit:   "backend_services",
						},
						Amount: 1,
					}},
				})
			}
		}
		return nil
	}

	var err error
	switch typeName {
	case globalBackendServiceResource:
		err = o.computeSvc.BackendServices.List(o.ProjectID).Fields(googleapi.Field(fields)).Pages(ctx, pagesFunc)
	case regionBackendServiceResource:
		err = o.computeSvc.RegionBackendServices.List(o.ProjectID, o.Region).Fields(googleapi.Field(fields)).Pages(ctx, pagesFunc)
	default:
		return nil, fmt.Errorf("invalid backend service type %q", typeName)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list backend services: %w", err)
	}

	return result, nil
}

func (o *ClusterUninstaller) deleteBackendService(ctx context.Context, item cloudResource) error {
	o.Logger.Debugf("Deleting backend service %s", item.name)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var err error
	var op *compute.Operation
	switch item.typeName {
	case globalBackendServiceResource:
		op, err = o.computeSvc.BackendServices.Delete(o.ProjectID, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
		item.scope = global
	case regionBackendServiceResource:
		op, err = o.computeSvc.RegionBackendServices.Delete(o.ProjectID, o.Region, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
		item.scope = regional
	default:
		return fmt.Errorf("invalid backend service type %q", item.typeName)
	}

	return o.handleOperation(ctx, op, err, item, "backend service")
}

// destroyBackendServices removes all backend services resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyBackendServices(ctx context.Context) error {
	found, err := o.listBackendServices(ctx, globalBackendServiceResource)
	if err != nil {
		return err
	}
	items := o.insertPendingItems(globalBackendServiceResource, found)

	found, err = o.listBackendServices(ctx, regionBackendServiceResource)
	if err != nil {
		return err
	}
	items = append(items, o.insertPendingItems(regionBackendServiceResource, found)...)

	for _, item := range items {
		err := o.deleteBackendService(ctx, item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}

	if items = o.getPendingItems(globalBackendServiceResource); len(items) > 0 {
		return fmt.Errorf("%d global backend service pending", len(items))
	}

	if items = o.getPendingItems(regionBackendServiceResource); len(items) > 0 {
		return fmt.Errorf("%d region backend service pending", len(items))
	}

	return nil
}
