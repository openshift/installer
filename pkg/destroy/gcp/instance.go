package gcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	gcpconsts "github.com/openshift/installer/pkg/constants/gcp"
	"github.com/openshift/installer/pkg/types/gcp"
)

const (
	instanceResourceName     = "instance"
	stopInstanceResourceName = "stopinstance"
)

type instanceFilterFunc func(instance *compute.Instance) bool

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
	instances, err := o.listInstancesWithFilter(ctx, "items/*/instances(name,zone,status,machineType,labels),nextPageToken",
		func(item *compute.Instance) bool {
			if o.isClusterResource(item.Name) {
				return true
			}

			for key, value := range item.Labels {
				if key == fmt.Sprintf(gcpconsts.ClusterIDLabelFmt, o.ClusterID) && value == ownedLabelValue {
					return true
				} else if key == fmt.Sprintf(capgProviderOwnedLabelFmt, o.ClusterID) && value == ownedLabelValue {
					return true
				}
			}
			return false
		})
	if err != nil {
		return nil, err
	}
	return instances, nil
}

// listInstancesWithFilter lists instances in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listInstancesWithFilter(ctx context.Context, fields string, filterFunc instanceFilterFunc) ([]cloudResource, error) {
	o.Logger.Debugf("Listing compute instances")
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	result := []cloudResource{}
	req := o.computeSvc.Instances.AggregatedList(o.ProjectID).Fields(googleapi.Field(fields))

	err := req.Pages(ctx, func(list *compute.InstanceAggregatedList) error {
		for _, scopedList := range list.Items {
			for _, item := range scopedList.Instances {
				if filterFunc(item) {
					zoneName := o.getZoneName(item.Zone)
					o.Logger.Debugf("Found instance: %s in zone %s, status %s", item.Name, zoneName, item.Status)
					result = append(result, cloudResource{
						key:      fmt.Sprintf("%s/%s", zoneName, item.Name),
						name:     item.Name,
						status:   item.Status,
						typeName: instanceResourceName,
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
		return nil, fmt.Errorf("failed to fetch compute instances: %w", err)
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteInstance(ctx context.Context, item cloudResource) error {
	o.Logger.Debugf("Deleting compute instance %s in zone %s", item.name, item.zone)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	op, err := o.computeSvc.Instances.Delete(o.ProjectID, item.zone, item.name).RequestId(o.requestID(item.typeName, item.zone, item.name)).Context(ctx).Do()
	item.scope = zonal
	return o.handleOperation(ctx, op, err, item, "instance")
}

// destroyInstances searches for instances across all zones that have a name that starts with
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyInstances(ctx context.Context) error {
	found, err := o.listInstances(ctx)
	if err != nil {
		return err
	}
	items := o.insertPendingItems(instanceResourceName, found)
	errs := []error{}
	for _, item := range items {
		err := o.deleteInstance(ctx, item)
		if err != nil {
			errs = append(errs, err)
		}
	}
	items = o.getPendingItems(instanceResourceName)
	return aggregateError(errs, len(items))
}

func (o *ClusterUninstaller) stopInstance(ctx context.Context, item cloudResource) error {
	o.Logger.Debugf("Stopping compute instance %s in zone %s", item.name, item.zone)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	op, err := o.computeSvc.Instances.
		Stop(o.ProjectID, item.zone, item.name).
		RequestId(o.requestID(stopInstanceResourceName, item.zone, item.name)).
		DiscardLocalSsd(true).
		Context(ctx).
		Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID(stopInstanceResourceName, item.zone, item.name)
		return errors.Wrapf(err, "failed to stop instance %s in zone %s", item.name, item.zone)
	}
	if op != nil && op.Status == DONE && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID(stopInstanceResourceName, item.zone, item.name)
		return errors.Errorf("failed to stop instance %s in zone %s with error: %s", item.name, item.zone, operationErrorMessage(op))
	}
	if (err != nil && isNoOp(err)) || (op != nil && op.Status == DONE) {
		o.resetRequestID(stopInstanceResourceName, item.name)
		o.deletePendingItems(stopInstanceResourceName, []cloudResource{item})
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
			o.insertPendingItems(stopInstanceResourceName, []cloudResource{item})
		}
	}
	items := o.getPendingItems(stopInstanceResourceName)
	for _, item := range items {
		err := o.stopInstance(ctx, item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}
	if items = o.getPendingItems(stopInstanceResourceName); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
