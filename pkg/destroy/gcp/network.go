package gcp

import (
	"github.com/pkg/errors"

	compute "google.golang.org/api/compute/v1"
)

func (o *ClusterUninstaller) listNetworks() ([]cloudResource, error) {
	o.Logger.Debugf("Listing networks")
	result := []cloudResource{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.Networks.List(o.ProjectID).Fields("items(name,selfLink),nextPageToken").Filter(o.clusterIDFilter())
	err := req.Pages(ctx, func(list *compute.NetworkList) error {
		for _, network := range list.Items {
			o.Logger.Debugf("Found network: %s", network.Name)
			result = append(result, cloudResource{
				key:      network.Name,
				name:     network.Name,
				typeName: "network",
				url:      network.SelfLink,
			})
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list networks")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteNetwork(item cloudResource) error {
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	o.Logger.Debugf("Deleting network %s", item.name)
	op, err := o.computeSvc.Networks.Delete(o.ProjectID, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID(item.typeName, item.name)
		return errors.Wrapf(err, "failed to delete network %s", item.name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID(item.typeName, item.name)
		return errors.Errorf("failed to delete network %s with error: %s", item.name, operationErrorMessage(op))
	}
	return nil
}

// destroyNetworks removes all vpc network resources prefixed with the
// cluster's infra ID
func (o *ClusterUninstaller) destroyNetworks() error {
	networks, err := o.listNetworks()
	if err != nil {
		return err
	}
	found := cloudResources{}
	errs := []error{}
	for _, network := range networks {
		found.insert(network)
		// destroy any network routes that are not named with the infra ID
		routes, err := o.listNetworkRoutes(network.url)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		for _, route := range routes {
			err := o.deleteRoute(route)
			if err != nil {
				o.Logger.Debugf("Failed to delete route %s: %v", route, err)
			}
		}

		err = o.deleteNetwork(network)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deleted := o.setPendingItems("network", found)
	for _, item := range deleted {
		o.Logger.Infof("Deleted network %s", item.name)
	}
	return aggregateError(errs, len(found))
}
