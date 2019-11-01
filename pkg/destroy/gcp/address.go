package gcp

import (
	"github.com/pkg/errors"

	compute "google.golang.org/api/compute/v1"
)

func (o *ClusterUninstaller) listAddresses() ([]string, error) {
	o.Logger.Debugf("Listing addresses")
	result := []string{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.Addresses.List(o.ProjectID, o.Region).Fields("items(name),nextPageToken").Filter(o.clusterIDFilter())
	err := req.Pages(ctx, func(list *compute.AddressList) error {
		for _, address := range list.Items {
			o.Logger.Debugf("Found address: %s", address.Name)
			result = append(result, address.Name)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list addresses")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteAddress(name string, errorOnPending bool) error {
	o.Logger.Debugf("Deleting address %s", name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	op, err := o.computeSvc.Addresses.Delete(o.ProjectID, o.Region, name).RequestId(o.requestID("address", name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID("address", name)
		return errors.Wrapf(err, "failed to delete address %s", name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID("address", name)
		return errors.Errorf("failed to delete address %s with error: %s", name, operationErrorMessage(op))
	}
	if errorOnPending && op != nil && op.Status != "DONE" {
		return errors.Errorf("deletion of address %s is pending", name)
	}
	return nil
}

// destroyAddresses removes all address resources that have a name prefixed with the
// cluster's infra ID
func (o *ClusterUninstaller) destroyAddresses() error {
	addresses, err := o.listAddresses()
	if err != nil {
		return err
	}
	found := make([]string, 0, len(addresses))
	errs := []error{}
	for _, address := range addresses {
		found = append(found, address)
		err := o.deleteAddress(address, false)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deleted := o.setPendingItems("address", found)
	for _, item := range deleted {
		o.Logger.Infof("Deleted address %s", item)
	}
	return aggregateError(errs, len(found))
}
