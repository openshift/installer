package gcp

import (
	"fmt"

	"github.com/pkg/errors"

	compute "google.golang.org/api/compute/v1"
)

func (o *ClusterUninstaller) listInstanceGroupsWithFilter(filter string) ([]cloudResource, error) {
	o.Logger.Debugf("Listing instance groups")
	result := []cloudResource{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.InstanceGroups.AggregatedList(o.ProjectID).Fields("items/*/instanceGroups(name,zone),nextPageToken").Filter(filter)
	err := req.Pages(ctx, func(list *compute.InstanceGroupAggregatedList) error {
		for _, scopedList := range list.Items {
			for _, ig := range scopedList.InstanceGroups {
				zoneName := o.getZoneName(ig.Zone)
				result = append(result, cloudResource{
					key:      fmt.Sprintf("%s/%s", zoneName, ig.Name),
					name:     ig.Name,
					typeName: "instancegroup",
					zone:     zoneName,
				})
				o.Logger.Debugf("Found instance group %s in zone %s", ig.Name, zoneName)
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch compute instance groups")
	}
	return result, nil
}

func (o *ClusterUninstaller) listInstanceGroupInstances(ig cloudResource) ([]cloudResource, error) {
	o.Logger.Debugf("Listing instance group instances for %v", ig)
	result := []cloudResource{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.InstanceGroups.ListInstances(o.ProjectID, ig.zone, ig.name, &compute.InstanceGroupsListInstancesRequest{}).Fields("items(instance),nextPageToken")
	err := req.Pages(ctx, func(list *compute.InstanceGroupsListInstances) error {
		for _, item := range list.Items {
			name, zone := o.getInstanceNameAndZone(item.Instance)
			if len(name) == 0 {
				continue
			}
			result = append(result, cloudResource{
				key:      fmt.Sprintf("%s/%s", zone, name),
				name:     name,
				typeName: "instance",
				zone:     zone,
			})
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch instance group instances")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteInstanceGroup(item cloudResource) error {
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	o.Logger.Debugf("Deleting instance group %s in zone %s", item.name, item.zone)
	op, err := o.computeSvc.InstanceGroups.Delete(o.ProjectID, item.zone, item.name).RequestId(o.requestID(item.typeName, item.zone, item.name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID(item.typeName, item.zone, item.name)
		return errors.Wrapf(err, "failed to delete instance group %s in zone %s", item.name, item.zone)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID("instancegroup", item.zone, item.name)
		return errors.Errorf("failed to delete instance group %s in zone %s with error: %s", item.name, item.zone, operationErrorMessage(op))
	}
	return nil
}

// destroyInstanceGroups searches for instance groups across all zones that have a name that starts with
// the infra ID prefix, and then deletes each instance found.
func (o *ClusterUninstaller) destroyInstanceGroups() error {
	groups, err := o.listInstanceGroupsWithFilter(o.clusterIDFilter())
	if err != nil {
		return err
	}
	errs := []error{}
	found := cloudResources{}
	for _, group := range groups {
		found.insert(group)
		err := o.deleteInstanceGroup(group)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deletedItems := o.setPendingItems("computeinstancegroup", found)
	for _, item := range deletedItems {
		o.Logger.Infof("Deleted instance group %s", item.name)
	}
	return aggregateError(errs, len(found))
}
