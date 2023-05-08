package gcp

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"github.com/openshift/installer/pkg/types/gcp"
)

func (o *ClusterUninstaller) listBackendServices(ctx context.Context) ([]cloudResource, error) {
	return o.listBackendServicesWithFilter(ctx, "items(name),nextPageToken", o.clusterIDFilter(), nil)
}

// listBackendServicesWithFilter lists backend services in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listBackendServicesWithFilter(ctx context.Context, fields string, filter string, filterFunc func(*compute.BackendService) bool) ([]cloudResource, error) {
	o.Logger.Debugf("Listing backend services")
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	result := []cloudResource{}
	req := o.computeSvc.RegionBackendServices.List(o.ProjectID, o.Region).Fields(googleapi.Field(fields))
	if len(filter) > 0 {
		req = req.Filter(filter)
	}
	err := req.Pages(ctx, func(list *compute.BackendServiceList) error {
		for _, item := range list.Items {
			if filterFunc == nil || filterFunc != nil && filterFunc(item) {
				o.Logger.Debugf("Found backend service: %s", item.Name)
				result = append(result, cloudResource{
					key:      item.Name,
					name:     item.Name,
					typeName: "backendservice",
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
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list backend services")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteBackendService(ctx context.Context, item cloudResource) error {
	o.Logger.Debugf("Deleting backend service %s", item.name)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	op, err := o.computeSvc.RegionBackendServices.Delete(o.ProjectID, o.Region, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID(item.typeName, item.name)
		return errors.Wrapf(err, "failed to delete backend service %s", item.name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID(item.typeName, item.name)
		return errors.Errorf("failed to delete backend service %s with error: %s", item.name, operationErrorMessage(op))
	}
	if (err != nil && isNoOp(err)) || (op != nil && op.Status == "DONE") {
		o.resetRequestID(item.typeName, item.name)
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted backend service %s", item.name)
	}
	return nil
}

// destroyBackendServices removes backend services with a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyBackendServices(ctx context.Context) error {
	found, err := o.listBackendServices(ctx)
	if err != nil {
		return err
	}
	items := o.insertPendingItems("backendservice", found)
	for _, item := range items {
		err := o.deleteBackendService(ctx, item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}
	if items = o.getPendingItems("backendservice"); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
