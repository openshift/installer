package gcp

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"github.com/openshift/installer/pkg/types/gcp"
)

func (o *ClusterUninstaller) listAddresses(ctx context.Context, scope resourceScope) ([]cloudResource, error) {
	return o.listAddressesWithFilter(ctx, "items(name,region,addressType),nextPageToken", o.clusterIDFilter(), nil, scope)
}

func createAddressCloudResources(filterFunc func(address *compute.Address) bool, list *compute.AddressList) []cloudResource {
	result := []cloudResource{}

	for _, item := range list.Items {
		if filterFunc == nil || filterFunc(item) {
			logrus.Debugf("Found address: %s", item.Name)
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

	return result
}

// listAddressesWithFilter lists addresses in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listAddressesWithFilter(ctx context.Context, fields string, filter string, filterFunc func(*compute.Address) bool, scope resourceScope) ([]cloudResource, error) {
	o.Logger.Debugf("Listing %s addresses", scope)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	result := []cloudResource{}

	if scope == gcpGlobalResource {
		req := o.computeSvc.GlobalAddresses.List(o.ProjectID).Fields(googleapi.Field(fields))
		if len(filter) > 0 {
			req = req.Filter(filter)
		}
		err := req.Pages(ctx, func(list *compute.AddressList) error {
			result = append(result, createAddressCloudResources(filterFunc, list)...)
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list global addresses: %w", err)
		}
		return result, nil
	}

	// Regional addresses
	req := o.computeSvc.Addresses.List(o.ProjectID, o.Region).Fields(googleapi.Field(fields))
	if len(filter) > 0 {
		req = req.Filter(filter)
	}
	err := req.Pages(ctx, func(list *compute.AddressList) error {
		result = append(result, createAddressCloudResources(filterFunc, list)...)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list regional addresses: %w", err)
	}

	return result, nil
}

func (o *ClusterUninstaller) deleteAddress(ctx context.Context, item cloudResource, scope resourceScope) error {
	o.Logger.Debugf("Deleting address %s", item.name)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var op *compute.Operation
	var err error
	if scope == gcpGlobalResource {
		op, err = o.computeSvc.GlobalAddresses.Delete(o.ProjectID, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	} else {
		op, err = o.computeSvc.Addresses.Delete(o.ProjectID, o.Region, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	}

	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID(item.typeName, item.name)
		return fmt.Errorf("failed to delete address %s with error: %s: %w", item.name, operationErrorMessage(op), err)
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
	for _, scope := range []resourceScope{gcpGlobalResource, gcpRegionalResource} {
		found, err := o.listAddresses(ctx, scope)
		if err != nil {
			return fmt.Errorf("failed to list %s addresses: %w", scope, err)
		}
		items := o.insertPendingItems("address", found)
		for _, item := range items {
			err := o.deleteAddress(ctx, item, scope)
			if err != nil {
				o.errorTracker.suppressWarning(item.key, err, o.Logger)
			}
		}
		for _, item := range o.getPendingItems("address") {
			if err := o.deleteAddress(ctx, item, scope); err != nil {
				return fmt.Errorf("error deleting pending address %s: %w", item.name, err)
			}
		}
	}
	return nil
}
