package gcp

import (
	"fmt"

	"github.com/pkg/errors"

	compute "google.golang.org/api/compute/v1"
)


func (o *ClusterUninstaller) listDisks() ([]nameAndZone, error) {
	o.Logger.Debug("Listing disks")
	result := []nameAndZone{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.Disks.AggregatedList(o.ProjectID).Fields("items/*/disks(name,zone),nextPageToken").Filter(o.clusterIDFilter())
	err := req.Pages(ctx, func(aggregatedList *compute.DiskAggregatedList) error {
		for _, scopedList := range aggregatedList.Items {
			for _, disk := range scopedList.Disks {
				zone := o.getZoneName(disk.Zone)
				result = append(result, nameAndZone{name: disk.Name, zone: zone})
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

func (o *ClusterUninstaller) deleteDisk(disk nameAndZone) error {
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	o.Logger.Debugf("Deleting disk %s in zone %s", disk.name, disk.zone)
	op, err := o.computeSvc.Disks.Delete(o.ProjectID, disk.zone, disk.name).RequestId(o.requestID("disk", disk.zone, disk.name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID("disk", disk.zone, disk.name)
		return errors.Wrapf(err, "failed to delete disk %s in zone %s", disk.name, disk.zone)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID("disk", disk.zone, disk.name)
		return errors.Errorf("failed to delete disk %s in zone %s with error: %s", disk.name, disk.zone, operationErrorMessage(op))
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
	found := make([]string, 0, len(disks))
	for _, disk := range disks {
		found = append(found, fmt.Sprintf("%s/%s", disk.zone, disk.name))
		err := o.deleteDisk(disk)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deletedItems := o.setPendingItems("disk", found)
	for _, item := range deletedItems {
		o.Logger.Infof("Deleted disk %s", item)
	}
	return aggregateError(errs, len(found))
}
