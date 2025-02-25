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
	return o.listAddressesWithFilter(ctx, typeName, "items(name,region,addressType),nextPageToken", o.isClusterResource)
}

// listAddressesWithFilter lists addresses in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results.
func (o *ClusterUninstaller) listAddressesWithFilter(ctx context.Context, typeName, fields string, filterFunc resourceFilterFunc) ([]cloudResource, error) {
	o.Logger.Debugf("Listing addresses")
	result := []cloudResource{}

	pagesFunc := func(list *compute.AddressList) error {
		for _, item := range list.Items {
			o.Logger.Debugf("Found address (%s): %s", typeName, item.Name)
			if filterFunc(item.Name) {
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
					typeName: typeName,
					quota:    quota,
				})
			}
		}
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var err error
	switch typeName {
	case globalAddressResource:
		err = o.computeSvc.GlobalAddresses.List(o.ProjectID).Fields(googleapi.Field(fields)).Pages(ctx, pagesFunc)
	case regionalAddressResource:
		err = o.computeSvc.Addresses.List(o.ProjectID, o.Region).Fields(googleapi.Field(fields)).Pages(ctx, pagesFunc)
	default:
		return nil, fmt.Errorf("invalid address type %q", typeName)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list addresses: %w", err)
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
		item.scope = global
	case regionalAddressResource:
		op, err = o.computeSvc.Addresses.Delete(o.ProjectID, o.Region, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
		item.scope = regional
	default:
		return fmt.Errorf("invalid address type %q", item.typeName)
	}

	return o.handleOperation(ctx, op, err, item, "address")
}

// destroyAddresses removes all address resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyAddresses(ctx context.Context) error {
	addressTypes := []string{globalAddressResource, regionalAddressResource}
	for _, addressType := range addressTypes {
		found, err := o.listAddresses(ctx, addressType)
		if err != nil {
			return err
		}
		items := o.insertPendingItems(addressType, found)

		for _, item := range items {
			err := o.deleteAddress(ctx, item)
			if err != nil {
				o.errorTracker.suppressWarning(item.key, err, o.Logger)
			}
		}

		if items := o.getPendingItems(addressType); len(items) > 0 {
			return fmt.Errorf("%d %s resources pending", len(items), addressType)
		}
	}

	return nil
}
