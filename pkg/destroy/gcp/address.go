package gcp

import (
	"context"
	"fmt"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"github.com/openshift/installer/pkg/types/gcp"
)

const (
	globalAddressResource   = "address"
	regionalAddressResource = "regionaddress"
)

func (o *ClusterUninstaller) listAddresses(ctx context.Context, typeName string) ([]cloudResource, error) {
	return o.listAddressesWithFilter(ctx, typeName, "items(name,region,addressType),nextPageToken", o.clusterIDFilter())
}

// listAddressesWithFilter lists addresses in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results.
func (o *ClusterUninstaller) listAddressesWithFilter(ctx context.Context, typeName, fields, filter string) ([]cloudResource, error) {
	o.Logger.Debugf("Listing addresses")

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var err error
	var list *compute.AddressList
	switch typeName {
	case globalAddressResource:
		list, err = o.computeSvc.GlobalAddresses.List(o.ProjectID).Filter(filter).Fields(googleapi.Field(fields)).Context(ctx).Do()
	case regionalAddressResource:
		list, err = o.computeSvc.Addresses.List(o.ProjectID, o.Region).Filter(filter).Fields(googleapi.Field(fields)).Context(ctx).Do()
	default:
		return nil, fmt.Errorf("invalid address type %q", typeName)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list addresses: %w", err)
	}

	result := []cloudResource{}
	for _, item := range list.Items {
		o.Logger.Debugf("Found address: %s", item.Name)
		if item.AddressType == "INTERNAL" {
			result = append(result, cloudResource{
				key:      item.Name,
				name:     item.Name,
				typeName: typeName,
				quota: []gcp.QuotaUsage{{
					Metric: &gcp.Metric{
						Service: gcp.ServiceComputeEngineAPI,
						Limit:   "internal_addresses",
						Dimensions: map[string]string{
							"region": getNameFromURL("regions", item.Region),
						},
					},
					Amount: 1,
				}},
			})
		}
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteAddress(ctx context.Context, item cloudResource) error {
	o.Logger.Debugf("Deleting address %s", item.name)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var err error
	var op *compute.Operation
	switch item.typeName {
	case globalAddressResource:
		op, err = o.computeSvc.GlobalAddresses.Delete(o.ProjectID, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	case regionalAddressResource:
		op, err = o.computeSvc.Addresses.Delete(o.ProjectID, o.Region, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	default:
		return fmt.Errorf("invalid address type %q", item.typeName)
	}

	if err != nil && !isNoOp(err) {
		o.resetRequestID(item.typeName, item.name)
		return fmt.Errorf("failed to delete address %s: %w", item.name, err)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID(item.typeName, item.name)
		return fmt.Errorf("failed to delete address %s with error: %s", item.name, operationErrorMessage(op))
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
	found, err := o.listAddresses(ctx, globalAddressResource)
	if err != nil {
		return err
	}
	items := o.insertPendingItems(globalAddressResource, found)

	found, err = o.listAddresses(ctx, regionalAddressResource)
	if err != nil {
		return err
	}
	items = append(items, o.insertPendingItems(regionalAddressResource, found)...)

	for _, item := range items {
		err := o.deleteAddress(ctx, item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}

	if items = o.getPendingItems(globalAddressResource); len(items) > 0 {
		return fmt.Errorf("%d global addresses pending", len(items))
	}

	if items = o.getPendingItems(regionalAddressResource); len(items) > 0 {
		return fmt.Errorf("%d region addresses pending", len(items))
	}

	return nil
}
