package gcp

import (
	"fmt"

	"github.com/pkg/errors"

	compute "google.golang.org/api/compute/v1"
)

func (o *ClusterUninstaller) listInstanceGroupsWithFilter(filter string) ([]nameAndZone, error) {
	o.Logger.Debugf("Listing instance groups")
	result := []nameAndZone{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.InstanceGroups.AggregatedList(o.ProjectID).Fields("items/*/instanceGroups(name,zone),nextPageToken").Filter(filter)
	err := req.Pages(ctx, func(list *compute.InstanceGroupAggregatedList) error {
		for _, scopedList := range list.Items {
			for _, ig := range scopedList.InstanceGroups {
				zoneName := o.getZoneName(ig.Zone)
				result = append(result, nameAndZone{
					name: ig.Name,
					zone: zoneName,
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

func (o *ClusterUninstaller) listInstanceGroupInstances(ig nameAndZone) ([]nameAndZone, error) {
	o.Logger.Debugf("Listing instance group instances for %v", ig)
	result := []nameAndZone{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.InstanceGroups.ListInstances(o.ProjectID, ig.zone, ig.name, &compute.InstanceGroupsListInstancesRequest{}).Fields("items(instance),nextPageToken")
	err := req.Pages(ctx, func(list *compute.InstanceGroupsListInstances) error {
		for _, item := range list.Items {
			name, zone := o.getInstanceNameAndZone(item.Instance)
			if len(name) == 0 {
				continue
			}
			result = append(result, nameAndZone{
				name: name,
				zone: zone,
			})
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch instance group instances")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteInstanceGroup(ig nameAndZone) error {
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	o.Logger.Debugf("Deleting instance group %s in zone %s", ig.name, ig.zone)
	op, err := o.computeSvc.InstanceGroups.Delete(o.ProjectID, ig.zone, ig.name).RequestId(o.requestID("instancegroup", ig.zone, ig.name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID("instancegroup", ig.zone, ig.name)
		return errors.Wrapf(err, "failed to delete instance group %s in zone %s", ig.name, ig.zone)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID("instancegroup", ig.zone, ig.name)
		return errors.Errorf("failed to delete instance group %s in zone %s with error: %s", ig.name, ig.zone, operationErrorMessage(op))
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
	found := make([]string, 0, len(groups))
	for _, group := range groups {
		found = append(found, fmt.Sprintf("%s/%s", group.zone, group.name))
		err := o.deleteInstanceGroup(group)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deletedItems := o.setPendingItems("computeinstancegroup", found)
	for _, item := range deletedItems {
		o.Logger.Infof("Deleted instance group %s", item)
	}
	return aggregateError(errs, len(found))
}
