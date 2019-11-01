package gcp

import (
	"github.com/pkg/errors"

	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

func (o *ClusterUninstaller) listBackendServices() ([]string, error) {
	return o.listBackendServicesWithFilter("items(name),nextPageToken", o.clusterIDFilter(), nil)
}

// listBackendServicesWithFilter lists backend services in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listBackendServicesWithFilter(fields string, filter string, filterFunc func(*compute.BackendService) bool) ([]string, error) {
	o.Logger.Debugf("Listing backend services")
	result := []string{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.RegionBackendServices.List(o.ProjectID, o.Region).Fields(googleapi.Field(fields))
	if len(filter) > 0 {
		req = req.Filter(filter)
	}
	err := req.Pages(ctx, func(list *compute.BackendServiceList) error {
		for _, backendService := range list.Items {
			if filterFunc == nil || filterFunc != nil && filterFunc(backendService) {
				o.Logger.Debugf("Found backend service: %s", backendService.Name)
				result = append(result, backendService.Name)
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list backend services")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteBackendService(name string) error {
	o.Logger.Debugf("Deleting backend service %s", name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	op, err := o.computeSvc.RegionBackendServices.Delete(o.ProjectID, o.Region, name).RequestId(o.requestID("backendservice", name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID("backendservice", name)
		return errors.Wrapf(err, "failed to delete backend service %s", name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID("backendservice", name)
		return errors.Errorf("failed to delete backend service %s with error: %s", name, operationErrorMessage(op))
	}
	return nil
}

// destroyBackendServices removes backend services with a name prefixed by the
// cluster's infra ID.
func (o *ClusterUninstaller) destroyBackendServices() error {
	backendServices, err := o.listBackendServices()
	if err != nil {
		return err
	}
	found := make([]string, 0, len(backendServices))
	errs := []error{}
	for _, backendService := range backendServices {
		found = append(found, backendService)
		err := o.deleteBackendService(backendService)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deleted := o.setPendingItems("backendservice", found)
	for _, item := range deleted {
		o.Logger.Infof("Deleted backend service %s", item)
	}
	return aggregateError(errs, len(found))
}
