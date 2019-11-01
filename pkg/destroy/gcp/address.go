package gcp

import (
	"github.com/pkg/errors"

	compute "google.golang.org/api/compute/v1"
)

func (o *ClusterUninstaller) listAddresses() ([]cloudResource, error) {
	o.Logger.Debugf("Listing addresses")
	result := []cloudResource{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.Addresses.List(o.ProjectID, o.Region).Fields("items(name),nextPageToken").Filter(o.clusterIDFilter())
	err := req.Pages(ctx, func(list *compute.AddressList) error {
		for _, item := range list.Items {
			o.Logger.Debugf("Found address: %s", item.Name)
			result = append(result, cloudResource{
				key:      item.Name,
				name:     item.Name,
				typeName: "address",
			})
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list addresses")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteAddress(item cloudResource, errorOnPending bool) error {
	o.Logger.Debugf("Deleting address %s", item.name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	op, err := o.computeSvc.Addresses.Delete(o.ProjectID, o.Region, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID(item.typeName, item.name)
		return errors.Wrapf(err, "failed to delete address %s", item.name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID(item.typeName, item.name)
		return errors.Errorf("failed to delete address %s with error: %s", item.name, operationErrorMessage(op))
	}
	if errorOnPending && op != nil && op.Status != "DONE" {
		return errors.Errorf("deletion of address %s is pending", item.name)
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
	found := cloudResources{}
	errs := []error{}
	for _, address := range addresses {
		found.insert(address)
		err := o.deleteAddress(address, false)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deleted := o.setPendingItems("address", found)
	for _, item := range deleted {
		o.Logger.Infof("Deleted address %s", item.name)
	}
	return aggregateError(errs, len(found))
}
