package gcp

import (
	"fmt"

	"github.com/pkg/errors"

	compute "google.golang.org/api/compute/v1"
)

func (o *ClusterUninstaller) listNetworkRoutes(networkURL string) ([]cloudResource, error) {
	return o.listRoutesWithFilter(fmt.Sprintf("network eq %q", networkURL))
}

func (o *ClusterUninstaller) listRoutes() ([]cloudResource, error) {
	return o.listRoutesWithFilter(o.clusterIDFilter())
}

func (o *ClusterUninstaller) listRoutesWithFilter(filter string) ([]cloudResource, error) {
	o.Logger.Debugf("Listing routes")
	result := []cloudResource{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.Routes.List(o.ProjectID).Fields("items(name),nextPageToken").Filter(filter)
	err := req.Pages(ctx, func(list *compute.RouteList) error {
		for _, route := range list.Items {
			o.Logger.Debugf("Found route: %s", route.Name)
			result = append(result, cloudResource{
				key:      route.Name,
				name:     route.Name,
				typeName: "route",
			})
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
		return errors.Wrapf(err, "failed to delete route %s", item.name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID(item.typeName, item.name)
		return errors.Errorf("failed to delete route %s with error: %s", item.name, operationErrorMessage(op))
	}
	return nil
}

// destroyRutes removes all route resources that have a name prefixed with the
// cluster's infra ID
func (o *ClusterUninstaller) destroyRoutes() error {
	routes, err := o.listRoutes()
	if err != nil {
		return err
	}
	found := cloudResources{}
	errs := []error{}
	for _, route := range routes {
		found.insert(route)
		err := o.deleteRoute(route)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deleted := o.setPendingItems("route", found)
	for _, item := range deleted {
		o.Logger.Infof("Deleted route %s", item.name)
	}
	return aggregateError(errs, len(found))
}
