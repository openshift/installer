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
	instanceGroupResourceName = "instancegroup"
)

func (o *ClusterUninstaller) listInstanceGroups(ctx context.Context) ([]cloudResource, error) {
	return o.listInstanceGroupsWithFilter(ctx, "items/*/instanceGroups(name,selfLink,zone),nextPageToken", o.isClusterResource)
}

// listInstanceGroupsWithFilter lists addresses in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listInstanceGroupsWithFilter(ctx context.Context, fields string, filterFunc resourceFilterFunc) ([]cloudResource, error) {
	o.Logger.Debugf("Listing instance groups")
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	result := []cloudResource{}
	req := o.computeSvc.InstanceGroups.AggregatedList(o.ProjectID).Fields(googleapi.Field(fields))
	err := req.Pages(ctx, func(list *compute.InstanceGroupAggregatedList) error {
		for _, scopedList := range list.Items {
			for _, item := range scopedList.InstanceGroups {
				if filterFunc(item.Name) {
					zoneName := o.getZoneName(item.Zone)
					o.Logger.Debugf("Found instance group: %s in zone %s", item.Name, zoneName)
					result = append(result, cloudResource{
						key:      fmt.Sprintf("%s/%s", zoneName, item.Name),
						name:     item.Name,
						typeName: instanceGroupResourceName,
						zone:     zoneName,
						url:      item.SelfLink,
						quota: []gcp.QuotaUsage{{
							Metric: &gcp.Metric{
								Service: gcp.ServiceComputeEngineAPI,
								Limit:   "instance_groups",
								Dimensions: map[string]string{
									"region": getRegionFromZone(zoneName),
								},
							},
							Amount: 1,
						}},
					})
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch compute instance groups")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteInstanceGroup(ctx context.Context, item cloudResource) error {
	o.Logger.Debugf("Deleting instance group %s in zone %s", item.name, item.zone)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	op, err := o.computeSvc.InstanceGroups.Delete(o.ProjectID, item.zone, item.name).RequestId(o.requestID(item.typeName, item.zone, item.name)).Context(ctx).Do()
	item.scope = zonal
	return o.handleOperation(ctx, op, err, item, "instance group")
}

// destroyInstanceGroups searches for instance groups that have a name prefixed
// with the cluster's infra ID, and then deletes each instance found.
func (o *ClusterUninstaller) destroyInstanceGroups(ctx context.Context) error {
	found, err := o.listInstanceGroups(ctx)
	if err != nil {
		return err
	}
	items := o.insertPendingItems(instanceGroupResourceName, found)
	for _, item := range items {
		err := o.deleteInstanceGroup(ctx, item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}
	if items = o.getPendingItems(instanceGroupResourceName); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
