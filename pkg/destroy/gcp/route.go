package gcp

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

func (o *ClusterUninstaller) listNetworkRoutes(networkURL string) ([]cloudResource, error) {
	return o.listRoutesWithFilter("items(name),nextPageToken", fmt.Sprintf("network eq %q", networkURL), nil)
}

func (o *ClusterUninstaller) listRoutes() ([]cloudResource, error) {
	return o.listRoutesWithFilter("items(name),nextPageToken", o.clusterIDFilter(), nil)
}

// listRoutesWithFilter lists routes in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listRoutesWithFilter(fields string, filter string, filterFunc func(*compute.Route) bool) ([]cloudResource, error) {
	o.Logger.Debugf("Listing routes")
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	result := []cloudResource{}
	req := o.computeSvc.Routes.List(o.ProjectID).Fields(googleapi.Field(fields))
	if len(filter) > 0 {
		req = req.Filter(filter)
	}
	err := req.Pages(ctx, func(list *compute.RouteList) error {
		for _, item := range list.Items {
			if filterFunc == nil || filterFunc != nil && filterFunc(item) {
				o.Logger.Debugf("Found route: %s", item.Name)
				result = append(result, cloudResource{
					key:      item.Name,
					name:     item.Name,
					typeName: "route",
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

func (o *ClusterUninstaller) deleteRoute(item cloudResource) error {
	o.Logger.Debugf("Deleting route %s", item.name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	op, err := o.computeSvc.Routes.Delete(o.ProjectID, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID(item.typeName, item.name)
		if strings.HasPrefix(item.name, "default-route") {
			return errors.New("this looks like a default route, which cannot be deleted manually but will be deleted with the corresponding network")
		}
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
func (o *ClusterUninstaller) destroyRoutes() error {
	found, err := o.listRoutes()
	if err != nil {
		return err
	}
	items := o.insertPendingItems("route", found)
	errs := []error{}
	for _, item := range items {
		err := o.deleteRoute(item)
		if err != nil {
			errs = append(errs, err)
		}
	}
	items = o.getPendingItems("route")
	return aggregateError(errs, len(items))
}
