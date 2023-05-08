package gcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"github.com/openshift/installer/pkg/types/gcp"
)

// getInstanceNameAndZone extracts an instance and zone name from an instance URL in the form:
// https://www.googleapis.com/compute/v1/projects/project-id/zones/us-central1-a/instances/instance-name
// After splitting the service's base path with the work `/projects/`, you get:
// project-id/zones/us-central1-a/instances/instance-name
// TODO: Find a better way to get the instance name and zone to account for changes in base path
func (o *ClusterUninstaller) getInstanceNameAndZone(instanceURL string) (string, string) {
	path := strings.Split(instanceURL, "/projects/")[1]
	parts := strings.Split(path, "/")
	if len(parts) >= 5 {
		return parts[4], parts[2]
	}
	return "", ""
}

func (o *ClusterUninstaller) listInstances(ctx context.Context) ([]cloudResource, error) {
	byName, err := o.listInstancesWithFilter(ctx, "items/*/instances(name,zone,status,machineType),nextPageToken", o.clusterIDFilter(), nil)
	if err != nil {
		return nil, err
	}

	byLabel, err := o.listInstancesWithFilter(ctx, "items/*/instances(name,zone,status,machineType),nextPageToken", o.clusterLabelFilter(), nil)
	if err != nil {
		return nil, err
	}
	return append(byName, byLabel...), nil
}

// listInstancesWithFilter lists instances in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listInstancesWithFilter(ctx context.Context, fields string, filter string, filterFunc func(*compute.Instance) bool) ([]cloudResource, error) {
	o.Logger.Debugf("Listing compute instances")
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	result := []cloudResource{}
	req := o.computeSvc.Instances.AggregatedList(o.ProjectID).Fields(googleapi.Field(fields))
	if len(filter) > 0 {
		req = req.Filter(filter)
	}
	err := req.Pages(ctx, func(list *compute.InstanceAggregatedList) error {
		for _, scopedList := range list.Items {
			for _, item := range scopedList.Instances {
				if filterFunc == nil || filterFunc != nil && filterFunc(item) {
					zoneName := o.getZoneName(item.Zone)
					o.Logger.Debugf("Found instance: %s in zone %s, status %s", item.Name, zoneName, item.Status)
					result = append(result, cloudResource{
						key:      fmt.Sprintf("%s/%s", zoneName, item.Name),
						name:     item.Name,
						status:   item.Status,
						typeName: "instance",
						zone:     zoneName,
						quota: []gcp.QuotaUsage{{
							Metric: &gcp.Metric{
								Service: gcp.ServiceComputeEngineAPI,
								Limit:   "cpus",
								Dimensions: map[string]string{
									"region": getRegionFromZone(zoneName),
								},
							},
							Amount: o.cpusByMachineType[getNameFromURL("machineTypes", item.MachineType)],
						}},
					})
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch compute instances")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteInstance(ctx context.Context, item cloudResource) error {
	o.Logger.Debugf("Deleting compute instance %s in zone %s", item.name, item.zone)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	op, err := o.computeSvc.Instances.Delete(o.ProjectID, item.zone, item.name).RequestId(o.requestID(item.typeName, item.zone, item.name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID(item.typeName, item.zone, item.name)
		return errors.Wrapf(err, "failed to delete instance %s in zone %s", item.name, item.zone)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID(item.typeName, item.zone, item.name)
		return errors.Errorf("failed to delete instance %s in zone %s with error: %s", item.name, item.zone, operationErrorMessage(op))
	}
	if (err != nil && isNoOp(err)) || (op != nil && op.Status == "DONE") {
		o.resetRequestID(item.typeName, item.name)
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted instance %s", item.name)
	}
	return nil
}

// destroyInstances searches for instances across all zones that have a name that starts with
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyInstances(ctx context.Context) error {
	found, err := o.listInstances(ctx)
	if err != nil {
		return err
	}
	items := o.insertPendingItems("instance", found)
	errs := []error{}
	for _, item := range items {
		err := o.deleteInstance(ctx, item)
		if err != nil {
			errs = append(errs, err)
		}
	}
	items = o.getPendingItems("instance")
	return aggregateError(errs, len(items))
}

func (o *ClusterUninstaller) stopInstance(ctx context.Context, item cloudResource) error {
	o.Logger.Debugf("Stopping compute instance %s in zone %s", item.name, item.zone)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	op, err := o.computeSvc.Instances.Stop(o.ProjectID, item.zone, item.name).RequestId(o.requestID("stopinstance", item.zone, item.name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID("stopinstance", item.zone, item.name)
		return errors.Wrapf(err, "failed to stop instance %s in zone %s", item.name, item.zone)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID("stopinstance", item.zone, item.name)
		return errors.Errorf("failed to stop instance %s in zone %s with error: %s", item.name, item.zone, operationErrorMessage(op))
	}
	if (err != nil && isNoOp(err)) || (op != nil && op.Status == "DONE") {
		o.resetRequestID("stopinstance", item.name)
		o.deletePendingItems("stopinstance", []cloudResource{item})
		o.Logger.Infof("Stopped instance %s", item.name)
	}
	return nil
}

// stopComputeInstances searches for instances across all zones that have a name that starts with
// the infra ID prefix and are not yet stopped. It then stops each instance found.
func (o *ClusterUninstaller) stopInstances(ctx context.Context) error {
	found, err := o.listInstances(ctx)
	if err != nil {
		return err
	}
	for _, item := range found {
		if item.status != "TERMINATED" {
			// we record instance quota when we delete the instance, not when we terminate it
			item.quota = nil
			o.insertPendingItems("stopinstance", []cloudResource{item})
		}
	}
	items := o.getPendingItems("stopinstance")
	for _, item := range items {
		err := o.stopInstance(ctx, item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}
	if items = o.getPendingItems("stopinstance"); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
