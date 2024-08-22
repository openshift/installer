package gcp

import (
	"context"
	"fmt"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openshift/installer/pkg/types/gcp"
)

const (
	globalBackendServiceResource = "backendservice"
	regionBackendServiceResource = "regionbackendservice"
)

func (o *ClusterUninstaller) listBackendServices(ctx context.Context, typeName string) ([]cloudResource, error) {
	return o.listBackendServicesWithFilter(ctx, typeName, "items(name),nextPageToken", o.clusterIDFilter(), nil)
}

func backendServiceBelongsToInstanceGroup(item *compute.BackendService, igURLs sets.Set[string]) bool {
	if igURLs == nil {
		return true
	}

	if len(item.Backends) == 0 {
		return false
	}
	for _, backend := range item.Backends {
		if !igURLs.Has(backend.Group) {
			return false
		}
	}
	return true
}

// listBackendServicesWithFilter lists backend services in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results.
func (o *ClusterUninstaller) listBackendServicesWithFilter(ctx context.Context, typeName, fields, filter string, urls sets.Set[string]) ([]cloudResource, error) {
	o.Logger.Debugf("Listing backend services")

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var err error
	var list *compute.BackendServiceList
	switch typeName {
	case globalBackendServiceResource:
		list, err = o.computeSvc.BackendServices.List(o.ProjectID).Filter(filter).Fields(googleapi.Field(fields)).Context(ctx).Do()
	case regionBackendServiceResource:
		list, err = o.computeSvc.RegionBackendServices.List(o.ProjectID, o.Region).Filter(filter).Fields(googleapi.Field(fields)).Context(ctx).Do()
	default:
		return nil, fmt.Errorf("invalid backend service type %q", typeName)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list backend services: %w", err)
	}

	result := []cloudResource{}
	for _, item := range list.Items {
		o.Logger.Debugf("Found backend service: %s", item.Name)
		if !backendServiceBelongsToInstanceGroup(item, urls) {
			o.Logger.Debug("No matching instance group for backend service: %s", item.Name)
			continue
		}
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
	case regionBackendServiceResource:
		op, err = o.computeSvc.RegionBackendServices.Delete(o.ProjectID, o.Region, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	default:
		return fmt.Errorf("invalid backend service type %q", item.typeName)
	}

	if err != nil && !isNoOp(err) {
		o.resetRequestID(item.typeName, item.name)
		return fmt.Errorf("failed to delete backend service %s: %w", item.name, err)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID(item.typeName, item.name)
		return fmt.Errorf("failed to delete backend service %s with error: %s", item.name, operationErrorMessage(op))
	}
	if (err != nil && isNoOp(err)) || (op != nil && op.Status == "DONE") {
		o.resetRequestID(item.typeName, item.name)
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted backend service %s", item.name)
	}
	return nil
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
