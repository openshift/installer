package gcp

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"github.com/openshift/installer/pkg/types/gcp"
)

const (
	subnetworkResourceName = "subnetwork"
)

func (o *ClusterUninstaller) deleteSubnetworkByName(ctx context.Context, resourceName string) error {
	items, err := o.listSubnetworksWithFilter(ctx, "items(name),nextPageToken", func(item string) bool { return item == resourceName })
	if err != nil {
		return fmt.Errorf("failed to list subnetworks by name: %w", err)
	}
	for _, item := range items {
		if err := o.deleteSubnetwork(ctx, item); err != nil {
			return fmt.Errorf("failed to delete subnetwork by name: %w", err)
		}
	}
	return nil
}

func (o *ClusterUninstaller) listSubnetworks(ctx context.Context) ([]cloudResource, error) {
	return o.listSubnetworksWithFilter(ctx, "items(name,network),nextPageToken", o.isClusterResource)
}

// listSubnetworksWithFilter lists subnetworks in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listSubnetworksWithFilter(ctx context.Context, fields string, filterFunc resourceFilterFunc) ([]cloudResource, error) {
	o.Logger.Debugf("Listing subnetworks")
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	result := []cloudResource{}
	req := o.computeSvc.Subnetworks.List(o.ProjectID, o.Region).Fields(googleapi.Field(fields))
	err := req.Pages(ctx, func(list *compute.SubnetworkList) error {
		for _, item := range list.Items {
			if filterFunc(item.Name) {
				o.Logger.Debugf("Found subnetwork: %s", item.Name)
				result = append(result, cloudResource{
					key:      item.Name,
					name:     item.Name,
					typeName: subnetworkResourceName,
					quota: []gcp.QuotaUsage{{
						Metric: &gcp.Metric{
							Service: gcp.ServiceComputeEngineAPI,
							Limit:   "subnetworks",
						},
						Amount: 1,
					}, {
						Metric: &gcp.Metric{
							Service: gcp.ServiceComputeEngineAPI,
							Limit:   "subnet_ranges_per_vpc_network",
							Dimensions: map[string]string{
								"network_id": getNameFromURL("networks", item.Network),
							},
						},
						Amount: 1,
					}},
				})
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list subnetworks")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteSubnetwork(ctx context.Context, item cloudResource) error {
	o.Logger.Debugf("Deleting subnetwork %s", item.name)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	op, err := o.computeSvc.Subnetworks.Delete(o.ProjectID, o.Region, item.name).RequestId(o.requestID(item.typeName, item.name)).Context(ctx).Do()
	item.scope = regional
	return o.handleOperation(ctx, op, err, item, "subnetwork")
}

// destroySubNetworks removes all subnetwork resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroySubnetworks(ctx context.Context) error {
	found, err := o.listSubnetworks(ctx)
	if err != nil {
		return err
	}
	items := o.insertPendingItems(subnetworkResourceName, found)
	for _, item := range items {
		err := o.deleteSubnetwork(ctx, item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}
	if items = o.getPendingItems(subnetworkResourceName); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
