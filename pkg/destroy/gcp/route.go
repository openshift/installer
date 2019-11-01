package gcp

import (
	"fmt"

	"github.com/pkg/errors"

	compute "google.golang.org/api/compute/v1"
)

func (o *ClusterUninstaller) listNetworkRoutes(networkURL string) ([]string, error) {
	return o.listRoutesWithFilter(fmt.Sprintf("network eq %q", networkURL))
}

func (o *ClusterUninstaller) listRoutes() ([]string, error) {
	return o.listRoutesWithFilter(o.clusterIDFilter())
}

func (o *ClusterUninstaller) listRoutesWithFilter(filter string) ([]string, error) {
	o.Logger.Debugf("Listing routes")
	result := []string{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.Routes.List(o.ProjectID).Fields("items(name),nextPageToken").Filter(filter)
	err := req.Pages(ctx, func(list *compute.RouteList) error {
		for _, route := range list.Items {
			o.Logger.Debugf("Found route: %s", route.Name)
			result = append(result, route.Name)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list routes")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteRoute(name string) error {
	o.Logger.Debugf("Deleting route %s", name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	op, err := o.computeSvc.Routes.Delete(o.ProjectID, name).RequestId(o.requestID("route", name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID("route", name)
		return errors.Wrapf(err, "failed to delete route %s", name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID("route", name)
		return errors.Errorf("failed to delete route %s with error: %s", name, operationErrorMessage(op))
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
	found := make([]string, 0, len(routes))
	errs := []error{}
	for _, route := range routes {
		found = append(found, route)
		err := o.deleteRoute(route)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deleted := o.setPendingItems("route", found)
	for _, item := range deleted {
		o.Logger.Infof("Deleted route %s", item)
	}
	return aggregateError(errs, len(found))
}
