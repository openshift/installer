package gcp

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	compute "google.golang.org/api/compute/v1"
)

// getInstanceNameAndZone extracts an instance and zone name from an instance URL in the form:
// https://www.googleapis.com/compute/v1/projects/project-id/zones/us-central1-a/instances/instance-name
// After trimming the service's base path, you get:
// project-id/zones/us-central1-a/instances/instance-name
func (o *ClusterUninstaller) getInstanceNameAndZone(instanceURL string) (string, string) {
	path := strings.TrimLeft(instanceURL, o.computeSvc.BasePath)
	parts := strings.Split(path, "/")
	if len(parts) >= 5 {
		return parts[4], parts[2]
	}
	return "", ""
}

func (o *ClusterUninstaller) listComputeInstances() ([]nameAndZone, error) {
	o.Logger.Debugf("Listing compute instances")
	result := []nameAndZone{}
	req := o.computeSvc.Instances.AggregatedList(o.ProjectID).Filter(o.clusterIDFilter()).Fields("items/*/instances(name,zone,status),nextPageToken")
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	err := req.Pages(ctx, func(list *compute.InstanceAggregatedList) error {
		for _, scopedList := range list.Items {
			for _, instance := range scopedList.Instances {
				zoneName := o.getZoneName(instance.Zone)
				result = append(result, nameAndZone{
					name:   instance.Name,
					zone:   zoneName,
					status: instance.Status,
				})
				o.Logger.Debugf("Found instance %s in zone %s, status %s", instance.Name, zoneName, instance.Status)
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch compute instances")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteComputeInstance(instance nameAndZone) error {
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	o.Logger.Debugf("Deleting compute instance %s in zone %s", instance.name, instance.zone)
	op, err := o.computeSvc.Instances.Delete(o.ProjectID, instance.zone, instance.name).RequestId(o.requestID("instance", instance.zone, instance.name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID("instance", instance.zone, instance.name)
		return errors.Wrapf(err, "failed to delete instance %s in zone %s", instance.name, instance.zone)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID("instance", instance.zone, instance.name)
		return errors.Errorf("failed to delete instance %s in zone %s with error: %s", instance.name, instance.zone, operationErrorMessage(op))
	}
	return nil
}

// destroyComputeInstances searches for instances across all zones that have a name that starts with
// the infra ID prefix and are not yet terminated. It then deletes each instance found.
func (o *ClusterUninstaller) destroyComputeInstances() error {
	instances, err := o.listComputeInstances()
	if err != nil {
		return err
	}
	errs := []error{}
	found := make([]string, 0, len(instances))
	for _, instance := range instances {
		found = append(found, fmt.Sprintf("%s/%s", instance.zone, instance.name))
		err := o.deleteComputeInstance(instance)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deletedItems := o.setPendingItems("computeinstance", found)
	for _, item := range deletedItems {
		o.Logger.Infof("Deleted instance %s", item)
	}
	return aggregateError(errs, len(found))
}
