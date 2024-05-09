package gcp

import (
	"context"
	"fmt"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"github.com/openshift/installer/pkg/types/gcp"
)

// listAddressesWithFilter lists addresses in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listAddressesWithFilter(ctx context.Context, resourceName, typeName, fields, filter string, listFunc addressListFunc) ([]cloudResource, error) {
	o.Logger.Debugf("Listing %s", resourceName)

	if listFunc == nil {
		return nil, fmt.Errorf("%s list func does not exist", resourceName)
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	result := []cloudResource{}
	list, err := listFunc(ctx, filter, fields)
	if err != nil {
		return nil, fmt.Errorf("failed to list %s: %w", resourceName, err)
	}

	for _, item := range list.Items {
		o.Logger.Debugf("Found %s: %s", resourceName, item.Name)
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

type addressListFunc func(ctx context.Context, filter, fields string) (*compute.AddressList, error)

func (o *ClusterUninstaller) addressDelete(ctx context.Context, item cloudResource) (*compute.Operation, error) {
	return o.computeSvc.GlobalAddresses.Delete(o.ProjectID, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
}

func (o *ClusterUninstaller) addressList(ctx context.Context, filter, fields string) (*compute.AddressList, error) {
	return o.computeSvc.GlobalAddresses.List(o.ProjectID).Filter(filter).Fields(googleapi.Field(fields)).Context(ctx).Do()
}

func (o *ClusterUninstaller) regionAddressDelete(ctx context.Context, item cloudResource) (*compute.Operation, error) {
	return o.computeSvc.Addresses.Delete(o.ProjectID, o.Region, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
}

func (o *ClusterUninstaller) regionAddressList(ctx context.Context, filter, fields string) (*compute.AddressList, error) {
	return o.computeSvc.Addresses.List(o.ProjectID, o.Region).Filter(filter).Fields(googleapi.Field(fields)).Context(ctx).Do()
}
