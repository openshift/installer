package gcp

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"github.com/openshift/installer/pkg/types/gcp"
)

func (o *ClusterUninstaller) listBackendServices(ctx context.Context, scope resourceScope) ([]cloudResource, error) {
	return o.listBackendServicesWithFilter(ctx, "items(name),nextPageToken", o.clusterIDFilter(), nil, scope)
}

func createBackendServiceCloudResources(filterFunc func(*compute.BackendService) bool, list *compute.BackendServiceList) []cloudResource {
	result := []cloudResource{}

	for _, item := range list.Items {
		if filterFunc == nil || filterFunc(item) {
			logrus.Debugf("Found backend service: %s", item.Name)
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

	return result
}

// listBackendServicesWithFilter lists backend services in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listBackendServicesWithFilter(ctx context.Context, fields string, filter string, filterFunc func(*compute.BackendService) bool, scope resourceScope) ([]cloudResource, error) {
	o.Logger.Debugf("Listing %s backend services", scope)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	result := []cloudResource{}

	if scope == gcpGlobalResource {
		req := o.computeSvc.BackendServices.List(o.ProjectID).Fields(googleapi.Field(fields))
		if len(filter) > 0 {
			req = req.Filter(filter)
		}
		err := req.Pages(ctx, func(list *compute.BackendServiceList) error {
			result = append(result, createBackendServiceCloudResources(filterFunc, list)...)
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list global backend services: %w", err)
		}
		return result, nil
	}

	// Regional backend services
	req := o.computeSvc.RegionBackendServices.List(o.ProjectID, o.Region).Fields(googleapi.Field(fields))
	if len(filter) > 0 {
		req = req.Filter(filter)
	}
	err := req.Pages(ctx, func(list *compute.BackendServiceList) error {
		result = append(result, createBackendServiceCloudResources(filterFunc, list)...)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list regional backend services: %w", err)
	}

	return result, nil
}

func (o *ClusterUninstaller) deleteBackendService(ctx context.Context, item cloudResource, scope resourceScope) error {
	o.Logger.Debugf("Deleting backend service %s", item.name)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var op *compute.Operation
	var err error
	if scope == gcpGlobalResource {
		op, err = o.computeSvc.BackendServices.Delete(o.ProjectID, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	} else {
		op, err = o.computeSvc.RegionBackendServices.Delete(o.ProjectID, o.Region, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	}

	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID(item.typeName, item.name)
		return fmt.Errorf("failed to delete backend service %s with error: %s: %w", item.name, operationErrorMessage(op), err)
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
	for _, scope := range []resourceScope{gcpGlobalResource, gcpRegionalResource} {
		found, err := o.listBackendServices(ctx, scope)
		if err != nil {
			return fmt.Errorf("failed to list backend services: %w", err)
		}
		items := o.insertPendingItems("backendservice", found)
		for _, item := range items {
			err := o.deleteBackendService(ctx, item, scope)
			if err != nil {
				o.errorTracker.suppressWarning(item.key, err, o.Logger)
			}
		}
		if items = o.getPendingItems("backendservice"); len(items) > 0 {
			for _, item := range items {
				if err := o.deleteBackendService(ctx, item, scope); err != nil {
					return fmt.Errorf("error deleting pending backend service %s: %w", item.name, err)
				}
			}
		}
	}
	return nil
}
