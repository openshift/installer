package ibmcloud

import (
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

const imageTypeName = "image"

// listImages lists images in the vpc
func (o *ClusterUninstaller) listImages() (cloudResources, error) {
	o.Logger.Debugf("Listing images")
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	options := o.vpcSvc.NewListImagesOptions()
	resources, _, err := o.vpcSvc.ListImagesWithContext(ctx, options)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list images")
	}

	result := []cloudResource{}
	for _, image := range resources.Images {
		if strings.Contains(*image.Name, o.InfraID) {
			result = append(result, cloudResource{
				key:      *image.ID,
				name:     *image.Name,
				status:   *image.Status,
				typeName: imageTypeName,
				id:       *image.ID,
			})
		}
	}

	return cloudResources{}.insert(result...), nil
}

func (o *ClusterUninstaller) deleteImage(item cloudResource) error {
	if item.status == vpcv1.ImageStatusDeletingConst {
		o.Logger.Debugf("Waiting for image %q to delete", item.name)
		return nil
	}

	o.Logger.Debugf("Deleting image %q", item.name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	options := o.vpcSvc.NewDeleteImageOptions(item.id)
	details, err := o.vpcSvc.DeleteImageWithContext(ctx, options)

	if err != nil && details != nil && details.StatusCode == http.StatusNotFound {
		// The resource is gone
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted image %q", item.name)
		return nil
	}

	if err != nil && details != nil && details.StatusCode != http.StatusNotFound {
		return errors.Wrapf(err, "Failed to delete image %s", item.name)
	}

	return nil
}

// destroyImages removes all image resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyImages() error {
	found, err := o.listImages()
	if err != nil {
		return err
	}

	items := o.insertPendingItems(imageTypeName, found.list())
	for _, item := range items {
		if _, ok := found[item.key]; !ok {
			// This item has finished deletion.
			o.deletePendingItems(item.typeName, []cloudResource{item})
			o.Logger.Infof("Deleted image %q", item.name)
			continue
		}
		err := o.deleteImage(item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}

	if items = o.getPendingItems(imageTypeName); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
