package gcp

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"github.com/openshift/installer/pkg/types/gcp"
)

func (o *ClusterUninstaller) listInstanceGroups(ctx context.Context) ([]cloudResource, error) {
	return o.listInstanceGroupsWithFilter(ctx, "items/*/instanceGroups(name,selfLink,zone),nextPageToken", o.clusterIDFilter(), nil)
}

// listInstanceGroupsWithFilter lists addresses in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listInstanceGroupsWithFilter(ctx context.Context, fields string, filter string, filterFunc func(*compute.InstanceGroup) bool) ([]cloudResource, error) {
	o.Logger.Debugf("Listing instance groups")
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	result := []cloudResource{}
	req := o.computeSvc.InstanceGroups.AggregatedList(o.ProjectID).Fields(googleapi.Field(fields))
	if len(filter) > 0 {
		req = req.Filter(filter)
	}
	err := req.Pages(ctx, func(list *compute.InstanceGroupAggregatedList) error {
		for _, scopedList := range list.Items {
			for _, item := range scopedList.InstanceGroups {
				if filterFunc == nil || filterFunc != nil && filterFunc(item) {
					zoneName := o.getZoneName(item.Zone)
					o.Logger.Debugf("Found instance group: %s in zone %s", item.Name, zoneName)
					result = append(result, cloudResource{
						key:      fmt.Sprintf("%s/%s", zoneName, item.Name),
						name:     item.Name,
						typeName: "instancegroup",
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
	if err != nil && !isNoOp(err) {
		o.resetRequestID(item.typeName, item.zone, item.name)
		return errors.Wrapf(err, "failed to delete instance group %s in zone %s", item.name, item.zone)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID("instancegroup", item.zone, item.name)
		return errors.Errorf("failed to delete instance group %s in zone %s with error: %s", item.name, item.zone, operationErrorMessage(op))
	}
	if (err != nil && isNoOp(err)) || (op != nil && op.Status == "DONE") {
		o.resetRequestID(item.typeName, item.name)
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted instance group %s", item.name)
	}
	return nil
}

// destroyInstanceGroups searches for instance groups that have a name prefixed
// with the cluster's infra ID, and then deletes each instance found.
func (o *ClusterUninstaller) destroyInstanceGroups(ctx context.Context) error {
	found, err := o.listInstanceGroups(ctx)
	if err != nil {
		return err
	}
	items := o.insertPendingItems("instancegroup", found)
	for _, item := range items {
		err := o.deleteInstanceGroup(ctx, item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}
	if items = o.getPendingItems("instancegroup"); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
