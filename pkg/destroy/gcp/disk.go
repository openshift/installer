package gcp

import (
	"fmt"

	"github.com/pkg/errors"

	compute "google.golang.org/api/compute/v1"
)


func (o *ClusterUninstaller) listDisks() ([]cloudResource, error) {
	o.Logger.Debug("Listing disks")
	result := []cloudResource{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.Disks.AggregatedList(o.ProjectID).Fields("items/*/disks(name,zone),nextPageToken").Filter(o.clusterIDFilter())
	err := req.Pages(ctx, func(aggregatedList *compute.DiskAggregatedList) error {
		for _, scopedList := range aggregatedList.Items {
			for _, disk := range scopedList.Disks {
				zone := o.getZoneName(disk.Zone)
				result = append(result, cloudResource{
					key:      fmt.Sprintf("%s/%s", zone, disk.Name),
					name:     disk.Name,
					typeName: "disk",
					zone:     zone,
				})
				o.Logger.Debugf("Found disk %s in zone %s", disk.Name, zone)
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch disks")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteDisk(item cloudResource) error {
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	o.Logger.Debugf("Deleting disk %s in zone %s", item.name, item.zone)
	op, err := o.computeSvc.Disks.Delete(o.ProjectID, item.zone, item.name).RequestId(o.requestID(item.typeName, item.zone, item.name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID(item.typeName, item.zone, item.name)
		return errors.Wrapf(err, "failed to delete disk %s in zone %s", item.name, item.zone)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID(item.typeName, item.zone, item.name)
		return errors.Errorf("failed to delete disk %s in zone %s with error: %s", item.name, item.zone, operationErrorMessage(op))
	}
	return nil
}

// destroyDisks searches for disks across all zones that have a name that starts with
// the infra ID prefix. It then deletes each disk found.
func (o *ClusterUninstaller) destroyDisks() error {
	disks, err := o.listDisks()
	if err != nil {
		return err
	}
	errs := []error{}
	found := cloudResources{}
	for _, disk := range disks {
		found.insert(disk)
		err := o.deleteDisk(disk)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deletedItems := o.setPendingItems("disk", found)
	for _, item := range deletedItems {
		o.Logger.Infof("Deleted disk %s", item.name)
	}
	return aggregateError(errs, len(found))
}
