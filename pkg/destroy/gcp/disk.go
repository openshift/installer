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
	maxGCEPDNameLength    = 63
	estimatedPVNameLength = 40
	// Removing an extra value (-1) for the "-" separated between the storage name and the pv name
	storageNameLength = maxGCEPDNameLength - estimatedPVNameLength - 1
)

// formatClusterIDForStorage will format the Cluster ID as it will be used for destroying
// GCE PDs. The maximum length is 63 characters, and can end with "-dynamic".
// https://github.com/kubernetes/kubernetes/blob/master/pkg/volume/util/util.go, GenerateVolumeName()
func (o *ClusterUninstaller) formatClusterIDForStorage() string {
	storageName := o.ClusterID + "-dynamic"
	slicedLength := storageNameLength
	if len(storageName) < slicedLength {
		slicedLength = len(storageName)
	}
	return storageName[:slicedLength]
}

func (o *ClusterUninstaller) storageIDFilter() string {
	return fmt.Sprintf("name : \"%s-*\"", o.formatClusterIDForStorage())
}

func (o *ClusterUninstaller) storageLabelFilter() string {
	return fmt.Sprintf("labels.kubernetes-io-cluster-%s = \"owned\"", o.formatClusterIDForStorage())
}

// storageLabelOrClusterIDFilter will perform the search for resources with the ClusterID, but
// it will also search for specific disk name formats.
func (o *ClusterUninstaller) storageLabelOrClusterIDFilter() string {
	return fmt.Sprintf("%s OR (%s) OR (%s)", o.clusterLabelOrClusterIDFilter(), o.storageIDFilter(), o.storageLabelFilter())
}

func (o *ClusterUninstaller) listDisks(ctx context.Context) ([]cloudResource, error) {
	return o.listDisksWithFilter(ctx, "items/*/disks(name,zone,type,sizeGb),nextPageToken", o.storageLabelOrClusterIDFilter(), nil)
}

// listDisksWithFilter lists disks in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listDisksWithFilter(ctx context.Context, fields string, filter string, filterFunc func(*compute.Disk) bool) ([]cloudResource, error) {
	o.Logger.Debug("Listing disks")
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
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

func (o *ClusterUninstaller) deleteDisk(ctx context.Context, item cloudResource) error {
	o.Logger.Debugf("Deleting disk %s in zone %s", item.name, item.zone)
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
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
func (o *ClusterUninstaller) destroyDisks(ctx context.Context) error {
	found, err := o.listDisks(ctx)
	if err != nil {
		return err
	}
	items := o.insertPendingItems("disk", found)
	for _, item := range items {
		err := o.deleteDisk(ctx, item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}
	if items = o.getPendingItems("disk"); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
