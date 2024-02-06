package ibmcloud

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

func (o *ClusterUninstaller) listDisks() ([]cloudResource, error) {
	o.Logger.Infof("Listing disks")

	result := []cloudResource{}
	clusterOwnedTag := o.clusterLabelFilter()
	options := o.vpcSvc.NewListVolumesOptions()

	for {
		ctx, cancel := o.contextWithTimeout()
		defer cancel()
		options.SetLimit(100)
		resources, _, err := o.vpcSvc.ListVolumesWithContext(ctx, options)
		if err != nil {
			return nil, errors.Wrap(err, "Listing disks failed")
		}

		for _, volume := range resources.Volumes {
			userTags := strings.Join(volume.UserTags, ",")
			if strings.Contains(userTags, clusterOwnedTag) {
				o.Logger.Debugf("Found disk: %s", *volume.ID)
				result = append(result, cloudResource{
					key:      *volume.ID,
					name:     *volume.Name,
					status:   *volume.Status,
					typeName: "disk",
					id:       *volume.ID,
				})
			}
		}

		//This was the last page, please exit the loop.
		if resources.Next == nil {
			o.Logger.Debugf("All disks fetched")
			break
		}

		//Set the start for the next page.
		start, _ := resources.GetNextStart()
		o.Logger.Debugf("Listing next page %s", *start)
		options.SetStart(*start)
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteDisk(item cloudResource) error {
	o.Logger.Infof("Deleting disk %s", item.id)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	options := o.vpcSvc.NewDeleteVolumeOptions(item.id)
	details, err := o.vpcSvc.DeleteVolumeWithContext(ctx, options)

	if err != nil {
		if details == nil || details.StatusCode != http.StatusNotFound {
			return fmt.Errorf("failed to delete disk name=%s, id=%s.If this error continues to persist for more than 20 minutes then please try to manually cleanup the volume using - ibmcloud is vold %s: %w", item.name, item.id, item.id, err)
		}

		if details.StatusCode == http.StatusNotFound {
			// The resource is gone
			o.deletePendingItems(item.typeName, []cloudResource{item})
			o.Logger.Infof("Deleted disk %s", item.id)
		}
	}

	return nil
}

func (o *ClusterUninstaller) waitForDiskDeletion(item cloudResource) error {
	o.Logger.Infof("Waiting for disk %s to be deleted", item.id)
	var skip = false

	err := o.Retry(func() (error, bool) {
		ctx, cancel := o.contextWithTimeout()
		defer cancel()
		volumeOptions := o.vpcSvc.NewGetVolumeOptions(item.id)
		_, response, err := o.vpcSvc.GetVolumeWithContext(ctx, volumeOptions)
		// Keep retry, until GetVolume returns volume not found
		if err != nil {
			if response != nil && response.StatusCode == http.StatusNotFound {
				skip = true
				return nil, skip
			}
		}
		return err, false // continue retry as we are not seeing error which means volume is available
	})

	if err == nil && skip {
		// The resource is gone
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted disk %s", item.id)
	} else {
		return errors.Wrapf(err, "Failed to delete disk name=%s, id=%s.If this error continues to persist for more than 20 minutes then please try to manually cleanup the volume using - ibmcloud is vold %s", item.name, item.id, item.id)
	}

	return err
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

	for _, item := range items {
		err := o.waitForDiskDeletion(item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}

	if items = o.getPendingItems("disk"); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
