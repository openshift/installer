package gcp

import (
	"fmt"

	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/pkg/errors"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

func (o *ClusterUninstaller) listDisks() ([]cloudResource, error) {
	return o.listDisksWithFilter("items/*/disks(name,zone,type,sizeGb),nextPageToken", o.clusterLabelOrClusterIDFilter(), nil)
}

// listDisksWithFilter lists disks in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listDisksWithFilter(fields string, filter string, filterFunc func(*compute.Disk) bool) ([]cloudResource, error) {
	o.Logger.Debug("Listing disks")
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	result := []cloudResource{}
	req := o.computeSvc.Disks.AggregatedList(o.ProjectID).Fields(googleapi.Field(fields))
	if len(filter) > 0 {
		req = req.Filter(filter)
	}
	err := req.Pages(ctx, func(list *compute.DiskAggregatedList) error {
		for _, scopedList := range list.Items {
			for _, item := range scopedList.Disks {
				if filterFunc == nil || filterFunc != nil && filterFunc(item) {
					zone := o.getZoneName(item.Zone)
					o.Logger.Debugf("Found disk: %s in zone %s", item.Name, zone)
					result = append(result, cloudResource{
						key:      fmt.Sprintf("%s/%s", zone, item.Name),
						name:     item.Name,
						typeName: "disk",
						zone:     zone,
						quota: []gcp.QuotaUsage{{
							Metric: &gcp.Metric{
								Service: gcp.ServiceComputeEngineAPI,
								Limit:   getDiskLimit(item.Type),
								Dimensions: map[string]string{
									"region": getRegionFromZone(zone),
								},
							},
							Amount: item.SizeGb,
						}},
					})
				}
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
	o.Logger.Debugf("Deleting disk %s in zone %s", item.name, item.zone)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	op, err := o.computeSvc.Disks.Delete(o.ProjectID, item.zone, item.name).RequestId(o.requestID(item.typeName, item.zone, item.name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID(item.typeName, item.zone, item.name)
		return errors.Wrapf(err, "failed to delete disk %s in zone %s", item.name, item.zone)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID(item.typeName, item.zone, item.name)
		return errors.Errorf("failed to delete disk %s in zone %s with error: %s", item.name, item.zone, operationErrorMessage(op))
	}
	if (err != nil && isNoOp(err)) || (op != nil && op.Status == "DONE") {
		o.resetRequestID(item.typeName, item.name)
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted disk %s", item.name)
	}
	return nil
}

// destroyDisks removes all disk resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyDisks() error {
	found, err := o.listDisks()
	if err != nil {
		return err
	}
	items := o.insertPendingItems("disk", found)
	for _, item := range items {
		err := o.deleteDisk(item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}
	if items = o.getPendingItems("disk"); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
