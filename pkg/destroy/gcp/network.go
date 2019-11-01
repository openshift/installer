package gcp

import (
	"github.com/pkg/errors"

	compute "google.golang.org/api/compute/v1"
)

func (o *ClusterUninstaller) listNetworks() ([]nameAndURL, error) {
	o.Logger.Debugf("Listing networks")
	result := []nameAndURL{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.Networks.List(o.ProjectID).Fields("items(name,selfLink),nextPageToken").Filter(o.clusterIDFilter())
	err := req.Pages(ctx, func(list *compute.NetworkList) error {
		for _, network := range list.Items {
			o.Logger.Debugf("Found network: %s", network.Name)
			result = append(result, nameAndURL{
				name: network.Name,
				url:  network.SelfLink,
			})
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list networks")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteNetwork(name string) error {
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	o.Logger.Debugf("Deleting network %s", name)
	op, err := o.computeSvc.Networks.Delete(o.ProjectID, name).RequestId(o.requestID("network", name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID("network", name)
		return errors.Wrapf(err, "failed to delete network %s", name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID("network", name)
		return errors.Errorf("failed to delete network %s with error: %s", name, operationErrorMessage(op))
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
	found := make([]string, 0, len(networks))
	errs := []error{}
	for _, network := range networks {
		found = append(found, network.name)
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

		err = o.deleteNetwork(network.name)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deleted := o.setPendingItems("network", found)
	for _, item := range deleted {
		o.Logger.Infof("Deleted network %s", item)
	}
	return aggregateError(errs, len(found))
}
