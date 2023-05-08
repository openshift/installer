package gcp

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"github.com/openshift/installer/pkg/types/gcp"
)

func (o *ClusterUninstaller) listAddresses(ctx context.Context) ([]cloudResource, error) {
	return o.listAddressesWithFilter(ctx, "items(name,region,addressType),nextPageToken", o.clusterIDFilter(), nil)
}

// listAddressesWithFilter lists addresses in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listAddressesWithFilter(ctx context.Context, fields string, filter string, filterFunc func(*compute.Address) bool) ([]cloudResource, error) {
	o.Logger.Debugf("Listing addresses")
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	result := []cloudResource{}
	req := o.computeSvc.Addresses.List(o.ProjectID, o.Region).Fields(googleapi.Field(fields))
	if len(filter) > 0 {
		req = req.Filter(filter)
	}
	err := req.Pages(ctx, func(list *compute.AddressList) error {
		for _, item := range list.Items {
			if filterFunc == nil || filterFunc != nil && filterFunc(item) {
				o.Logger.Debugf("Found address: %s", item.Name)
				var quota []gcp.QuotaUsage
				if item.AddressType == "INTERNAL" {
					quota = []gcp.QuotaUsage{{
						Metric: &gcp.Metric{
							Service: gcp.ServiceComputeEngineAPI,
							Limit:   "internal_addresses",
							Dimensions: map[string]string{
								"region": getNameFromURL("regions", item.Region),
							},
						},
						Amount: 1,
					}}
				}
				result = append(result, cloudResource{
					key:      item.Name,
					name:     item.Name,
					typeName: "address",
					quota:    quota,
				})
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list addresses")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteAddress(ctx context.Context, item cloudResource) error {
	o.Logger.Debugf("Deleting address %s", item.name)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
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
	if (err != nil && isNoOp(err)) || (op != nil && op.Status == "DONE") {
		o.resetRequestID(item.typeName, item.name)
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted address %s", item.name)
	}
	return nil
}

// destroyAddresses removes all address resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyAddresses(ctx context.Context) error {
	found, err := o.listAddresses(ctx)
	if err != nil {
		return err
	}
	items := o.insertPendingItems("address", found)
	for _, item := range items {
		err := o.deleteAddress(ctx, item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}
	if items = o.getPendingItems("address"); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
